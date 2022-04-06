package log

import "sync"

// =====================================================

var (
	std = New(NewOptions())
	mu sync.Mutex
)

// Init 使用给定的 Options 初始化 logger
func Init(opts *Options) {
	mu.Lock()
	defer mu.Unlock()
	std = New(opts)
}