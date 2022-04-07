package term

import (
	"fmt"
	"io"
	"github.com/moby/term"
)

// TerminalSize 返回用户终端的当前宽度和高度。
// 如果它不是终端，则返回 nil。 出错时，宽度和高度返回零值。
// 通常 w 必须是进程的标准输出。 标准错误将不起作用。
func TerminalSize(w io.Writer) (int, int, error) {
	outFd, isTerminal := term.GetFdInfo(w)
	if !isTerminal {
		return 0, 0, fmt.Errorf("given writer is no terminal")
	}
	winsize, err := term.GetWinsize(outFd)
	if err != nil {
		return 0, 0, err
	}
	return int(winsize.Width), int(winsize.Height), nil
}
