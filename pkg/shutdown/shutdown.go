package shutdown

import "sync"

// ShutdownManager 是被 ShutdownManagers 实现的一个接口
// 1. GetName
//    返回 ShutdownManager 的名字
// 2. Start、ShutdownStart、ShutdownFinish
//    ShutdownManagers 在 Start 函数中开始监听 shutdown 请求。
// 当调用 GSInterface 中的 StartShutdown 时，
// 首先调用 ShutdownStart()；
// 然后所有的 ShutdownCallbacks 将被执行；
// 一旦所有的 ShutdownCallbacks 返回，ShutdownFinish() 将被调用。
type ShutdownManager interface {
	GetName() string
	Start(gs GSInterface) error
	ShutdownStart() error
	ShutdownFinish() error
}

// GSInterface 是被 GracefulShutdown 实现的一个接口
// 当接收到 Shutdown 请求时，通过 ShutdownManager 去调用 StartShutdown
type GSInterface interface {
	StartShutdown(sm ShutdownManager)
	ReportError(err error)
	AddShutdownCallback(sc ShutdownCallback)
}

// ShutdownCallback 是一个接口，所有的 callback 必须实现这个接口。
// 当收到 shutdown 请求时，OnShutdown 将被调用。
// OnShutdown 参数是请求 shutdown 的 ShutdownManager 的 name
type ShutdownCallback interface {
	OnShutdown(string) error
}

// ShutdownFunc 是一个帮助类型，你可以简单的提供匿名函数作为 ShutdownCallback
type ShutdownFunc func(string) error

func (f ShutdownFunc) OnShutdown(shutdownManager string) error {
	return f(shutdownManager)
}

// ErrorHandler 是一个接口，可以通过 SetErrorHandler 去处理异步错误
type ErrorHandler interface {
	OnError(err error)
}

// ErrorFunc is a helper type, so you can easily provide anonymous functions
// as ErrorHandlers.
type ErrorFunc func(err error)

// OnError defines the action needed to run when error occurred.
func (f ErrorFunc) OnError(err error) {
	f(err)
}

// GracefulShutdown
// 处理 ShutdownManagers 和 ShutdownCallbacks 的主要结构体
// 使用 New() 进行初始化
type GracefulShutdown struct {
	managers     []ShutdownManager
	callbacks    []ShutdownCallback
	errorHandler ErrorHandler
}

// New 初始化 GracefulShutdown
func New() *GracefulShutdown {
	return &GracefulShutdown{
		managers: make([]ShutdownManager, 0, 3),
		callbacks: make([]ShutdownCallback, 0, 10),
	}
}

// AddShutdownManager 添加一个侦听 shutdown 请求的 ShutdownManager
func (gs *GracefulShutdown) AddShutdownManager(manager ShutdownManager)  {
	gs.managers = append(gs.managers, manager)
}

// AddShutdownCallback 添加一个 ShutdownCallback（当收到 shutdown 请求时被调用）
// 你可以提供任何实现了 ShutdownCallback 接口的函数，
// 或者按照下面的格式实现函数：
//   AddShutdownCallback(shutdown.ShutdownFunc(func() error {
//       // callback code
//       return nil
//   }))
func (gs *GracefulShutdown) AddShutdownCallback(sc ShutdownCallback)  {
	gs.callbacks = append(gs.callbacks, sc)
}

// SetErrorHandler 设置一个 ErrorHandler
// 当 ShutdownCallback 或者 ShutdownManager 发生错误的时候被调用。
// 你可以提供任何实现了 ErrorHandler 接口的函数，
// 或者按照下面的格式实现函数：
//   SetErrorHandler(shutdown.ErrorFunc(func(err error) {
//       // handle error
//   }))
func (gs *GracefulShutdown) SetErrorHandler(errorHandler ErrorHandler) {
	gs.errorHandler = errorHandler
}

// ReportError is a function that can be used to report errors to
// ErrorHandler. It is used in ShutdownManagers.
func (gs *GracefulShutdown) ReportError(err error) {
	if err != nil && gs.errorHandler != nil {
		gs.errorHandler.OnError(err)
	}
}

// Start 调用所有添加的 ShutdownManager 的 Start() 方法
// ShutdownManager 开始监听 shutdown 请求，并返回错误（如果 ShutdownManager 返回一个错误）
func (gs *GracefulShutdown) Start() error {
	for _, manager := range gs.managers {
		if err := manager.Start(gs); err != nil {
			return err
		}
	}
	return nil
}

// StartShutdown is called from a ShutdownManager and will initiate shutdown.
// first call ShutdownStart on Shutdownmanager,
// call all ShutdownCallbacks, wait for callbacks to finish and
// call ShutdownFinish on ShutdownManager.
func (gs *GracefulShutdown) StartShutdown(sm ShutdownManager) {
	gs.ReportError(sm.ShutdownStart())

	var wg sync.WaitGroup
	for _, shutdownCallback := range gs.callbacks {
		wg.Add(1)
		go func(shutdownCallback ShutdownCallback) {
			defer wg.Done()
			gs.ReportError(shutdownCallback.OnShutdown(sm.GetName()))
		}(shutdownCallback)
	}
	wg.Wait()

	gs.ReportError(sm.ShutdownFinish())
}