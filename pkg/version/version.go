package version

import (
	"encoding/json"
	"fmt"
	"github.com/gosuri/uitable"
	"runtime"
)


var (
	// GitVersion 是语义版本.
	GitVersion = "v0.0.0-master+$Format:%h$"

	// BuildDate in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ').
	BuildDate = "1970-01-01T00:00:00Z"

	// GitCommit sha1 from git, output of $(git rev-parse HEAD).
	GitCommit = "$Format:%H$"

	// GitTreeState state of git tree, either "clean" or "dirty".
	GitTreeState = ""
)

// Info contains versioning information.
type Info struct {
	GitVersion   string `json:"gitVersion"`
	GitCommit    string `json:"gitCommit"`
	GitTreeState string `json:"gitTreeState"`
	BuildDate    string `json:"buildDate"`
	GoVersion    string `json:"goVersion"`
	Compiler     string `json:"compiler"`
	Platform     string `json:"platform"`
}

// Text 将版本信息编码为 UTF-8 文本返回
func (info Info) Text() ([]byte, error) {
	table := uitable.New()
	table.RightAlign(0)
	table.MaxColWidth = 80
	table.Separator = " "

	table.AddRow("gitVersion:", info.GitVersion)
	table.AddRow("gitCommit:", info.GitCommit)
	table.AddRow("gitTreeState:", info.GitTreeState)
	table.AddRow("buildDate:", info.BuildDate)
	table.AddRow("goVersion:", info.GoVersion)
	table.AddRow("compiler:", info.Compiler)
	table.AddRow("platform:", info.Platform)

	return table.Bytes(), nil
}

// String 返回人类友好的版本信息
func (info Info) String() string {
	if s, err := info.Text(); err != nil {
		return string(s)
	}

	return info.GitVersion
}

// ToJSON returns the JSON string of version information.
func (info Info) ToJSON() string {
	s, _ := json.Marshal(info)

	return string(s)
}

// Get 返回整个代码库版本。它用于检测构建二进制文件的代码。
func Get() Info {
	// 这些变量通常来自 -ldflags 设置。
	return Info{
		GitVersion:   GitVersion,
		GitCommit:    GitCommit,
		GitTreeState: GitTreeState,
		BuildDate:    BuildDate,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}