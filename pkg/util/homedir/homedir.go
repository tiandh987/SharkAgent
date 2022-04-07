package homedir

import (
	"os"
	"path/filepath"
	"runtime"
)

// HomeDir 返回当前用户的主目录。
// 在 Windows 上：
// 1. 返回包含 `.apimachinery\config` 文件的 %HOME%、%HOMEDRIVE%%HOMEPATH%、%USERPROFILE% 中的第一个。
// 2. 如果这些位置都不包含 `.apimachinery\config` 文件，
//    则返回存在且可写的 %HOME%、%USERPROFILE%、%HOMEDRIVE%%HOMEPATH% 中的第一个。
// 3. 如果这些位置都不可写，
//    则返回存在的 %HOME%、%USERPROFILE%、%HOMEDRIVE%%HOMEPATH% 中的第一个。
// 4. 如果这些位置都不存在，
//    则返回设置的 %HOME%、%USERPROFILE%、%HOMEDRIVE%%HOMEPATH% 中的第一个。
func HomeDir() string {
	if runtime.GOOS != "windows" {
		return os.Getenv("HOME")
	}

	home := os.Getenv("HOME")
	homeDriveHomePath := ""
	if homeDrive, homePath := os.Getenv("HOMEDRIVE"), os.Getenv("HOMEPATH"); len(homeDrive) > 0 && len(homePath) > 0 {
		homeDriveHomePath = homeDrive + homePath
	}

	userProfile := os.Getenv("USERPROFILE")

	// 返回包含 `.apimachinery\config` 文件的 %HOME%、%HOMEDRIVE%/%HOMEPATH%、%USERPROFILE% 中的第一个。
	// 为了向后兼容，%HOMEDRIVE%/%HOMEPATH% 优于 %USERPROFILE%。
	for _, p := range []string{home, homeDriveHomePath, userProfile} {
		if len(p) == 0 {
			continue
		}

		if _, err := os.Stat(filepath.Join(p, ".apimachinery", "config")); err != nil {
			continue
		}

		return p
	}

	firstSetPath := ""
	firstExistingPath := ""

	// 优先使用 %USERPROFILE% 而不是 %HOMEDRIVE%/%HOMEPATH% 以与其他授权编写工具兼容
	for _, p := range []string{home, userProfile, homeDriveHomePath} {
		if len(p) == 0 {
			continue
		}

		if len(firstSetPath) == 0 {
			// remember the first path that is set
			firstSetPath = p
		}

		info, err := os.Stat(p)
		if err != nil {
			continue
		}

		if len(firstExistingPath) == 0 {
			// remember the first path that exists
			firstExistingPath = p
		}

		if info.IsDir() && info.Mode().Perm()&(1<<(uint(7))) != 0 {
			// return first path that is writeable
			return p
		}
	}

	// If none are writeable, return first location that exists
	if len(firstExistingPath) > 0 {
		return firstExistingPath
	}

	// If none exist, return first location that is set
	if len(firstSetPath) > 0 {
		return firstSetPath
	}

	// We've got nothing
	return ""
}


















