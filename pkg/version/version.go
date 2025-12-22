package version

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

// 构建时注入的版本信息变量。
// 通过 go build -ldflags "-X package.Variable=value" 注入。
var (
	AppRawName string = "Unknown" // 应用原始名称
	AppProject string = "Unknown" // 项目名称（通常为 Git 仓库名）
	AppVersion string = "Unknown" // 应用版本号（语义化版本）
	GitCommit  string = "Unknown" // Git 提交哈希
	BuildTime  string = "Unknown" // 构建时间
	Developer  string = "Unknown" // 开发者/维护者
)

// datePrefix 用于匹配项目名称前缀中的日期格式（如 "251203-"）
var datePrefix = regexp.MustCompile(`^[0-9-]{7}`)

func init() {
	initFromBuildInfo()
}

// initFromBuildInfo 从 runtime/debug.BuildInfo 读取版本信息作为后备。
// 当 ldflags 未注入时，自动从 Go 模块信息和 VCS 设置中提取。
func initFromBuildInfo() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	// 从模块路径提取项目名称和开发者
	if info.Main.Path != "" {
		parts := strings.Split(info.Main.Path, "/")
		if AppProject == "Unknown" {
			AppProject = path.Base(info.Main.Path)
		}
		// 从 <domain>/<user>/<repo> 格式中提取开发者
		if Developer == "Unknown" && len(parts) >= 2 {
			Developer = "http://" + parts[0] + "/" + parts[1] // 第二部分是用户名/组织名
		}
	}

	// 从模块版本提取应用版本
	if AppVersion == "Unknown" && info.Main.Version != "" && info.Main.Version != "(devel)" {
		AppVersion = info.Main.Version
	}

	// 从 VCS 设置中提取 Git 信息
	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			if GitCommit == "Unknown" && setting.Value != "" {
				// 使用短 hash（7位），与 git log -n 1 --format=%h 一致
				if len(setting.Value) > 7 {
					GitCommit = setting.Value[:7]
				} else {
					GitCommit = setting.Value
				}
			}
		case "vcs.time":
			if BuildTime == "Unknown" && setting.Value != "" {
				BuildTime = formatBuildTime(setting.Value)
			}
		case "vcs.modified":
			// 如果有未提交的修改，标记为 dirty
			if setting.Value == "true" && GitCommit != "Unknown" && !strings.HasSuffix(GitCommit, "-dirty") {
				GitCommit += "-dirty"
			}
		}
	}

	// 从 AppProject 中提取 AppRawName（去除日期前缀）
	if AppRawName == "Unknown" && AppProject != "Unknown" {
		AppRawName = datePrefix.ReplaceAllString(AppProject, "")
	}
}

// formatBuildTime 将 VCS 时间（RFC3339 格式）转换为 UTC+8 时区格式。
func formatBuildTime(vcsTime string) string {
	t, err := time.Parse(time.RFC3339, vcsTime)
	if err != nil {
		return vcsTime // 解析失败则返回原始值
	}
	// 使用固定偏移量 UTC+8，避免依赖系统时区数据库
	cst := time.FixedZone("CST", 8*60*60)
	return t.In(cst).Format("2006-01-02 15:04:05 MST")
}

// PrintBuildInfo 打印详细的构建信息（包括版本号、Git提交、构建时间等）
// 这是 version 命令的默认输出格式
func PrintBuildInfo() {
	_, _ = fmt.Fprintf(os.Stdout, "AppRawName:   %s\n", AppRawName)
	_, _ = fmt.Fprintf(os.Stdout, "AppVersion:   %s\n", AppVersion)
	_, _ = fmt.Fprintf(os.Stdout, "Go Version:   %s\n", runtime.Version())
	_, _ = fmt.Fprintf(os.Stdout, "Git Commit:   %s\n", GitCommit)
	_, _ = fmt.Fprintf(os.Stdout, "Build Time:   %s\n", BuildTime)
	_, _ = fmt.Fprintf(os.Stdout, "AppProject:   %s\n", AppProject)
	_, _ = fmt.Fprintf(os.Stdout, "Developer :   %s\n", Developer)
}

// PrintVersionJSON 以 JSON 格式打印完整的构建信息。
// 输出字段包括：appRawName、appProject、appVersion、gitCommit、buildTime、developer。
func PrintVersionJSON() {
	_, _ = fmt.Fprintf(os.Stdout, `{
  "appRawName": "%s",
  "appProject": "%s",
  "appVersion": "%s",
  "gitCommit": "%s",
  "buildTime": "%s",
  "developer": "%s"
}
`, AppRawName, AppProject, AppVersion, GitCommit, BuildTime, Developer)
}

// GetVersion 返回应用版本号。
// 若 AppVersion 为 "Unknown" 但 GitCommit 已知，则返回 "dev-<commit>" 格式。
func GetVersion() string {
	if AppVersion == "Unknown" && GitCommit != "Unknown" {
		return "dev-" + GitCommit
	}
	return AppVersion
}

// GetBuildInfo 返回构建相关信息 (用于健康检查等)
func GetBuildInfo() string {
	return fmt.Sprintf("版本: %s, 提交: %s, 构建时间: %s", AppVersion, GitCommit, BuildTime)
}
