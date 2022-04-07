package posixsignal

import (
	"github.com/tiandh987/SharkAgent/pkg/shutdown"
	"os"
	"os/signal"
	"syscall"
)

// Name 定义 shutdown manager 的名字
const Name = "PosixSignalManager"

// PosixSignalManager 实现 ShutdownManager 接口，被添加到 GracefaulShutdown
// 使用 NewPosixSignalManager 进行初始化
type  PosixSignalManager struct {
	signals []os.Signal
}

// NewPosixSignalManager 初始化 PosixSignalManager
// 参数为要监听的信号,如果未提供,默认监听 SIGINT 和 SIGTERM
func NewPosixSignalManager(sig ...os.Signal) *PosixSignalManager {
	if len(sig) == 0 {
		sig = make([]os.Signal, 2)
		sig[0] = os.Interrupt
		sig[1] = syscall.SIGTERM
	}

	return &PosixSignalManager{
		signals: sig,
	}
}

// GetName 返回 ShutdownManager 的 Name
func (psm *PosixSignalManager) GetName() string {
	return Name
}

// Start 开始监听 posix 信号
func (psm *PosixSignalManager) Start(gs shutdown.GSInterface) error {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, psm.signals...)

		// 阻塞,直到接收到信号
		<-c

		gs.StartShutdown(psm)
	}()

	return nil
}

//
func (psm *PosixSignalManager) ShutdownStart() error {
	return nil
}

// ShutdownFinish 使用 os.Exit(0) 退出
func (psm *PosixSignalManager) ShutdownFinish() error {
	os.Exit(0)
	return nil
}