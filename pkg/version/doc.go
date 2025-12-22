// Package version 提供应用程序的版本信息管理和显示功能。
//
// # 版本信息来源
//
// 本包支持两种版本信息来源，按优先级排序：
//
//  1. 构建时注入 (ldflags)：通过 go build -ldflags 在编译时注入
//  2. 运行时读取 (BuildInfo)：自动从 runtime/debug.BuildInfo 读取（Go 1.18+）
//
// # 构建时注入示例
//
// 使用 ldflags 注入版本信息（注意：路径必须是本库的实际导入路径）：
//
//	PKG=github.com/lwmacct/251207-go-pkg-version/pkg/version
//	go build -ldflags "\
//	  -X ${PKG}.AppVersion=v1.0.0 \
//	  -X ${PKG}.GitCommit=$(git rev-parse --short HEAD) \
//	  -X ${PKG}.BuildTime=$(date -u '+%Y-%m-%d %H:%M:%S UTC')"
//
// # 输出格式
//
// 本包提供三种版本信息输出格式：
//
//   - 默认格式 ([PrintBuildInfo])：详细的构建信息，包含应用名称、版本、Git 提交、构建时间等
//   - 短格式 ([GetVersion])：仅返回版本号字符串
//   - JSON 格式 ([PrintVersionJSON])：以 JSON 结构输出完整版本信息
//
// # CLI 集成
//
// 本包提供预定义的 [Command] 变量，可直接集成到 urfave/cli/v3 应用中：
//
//	import "github.com/lwmacct/251207-go-pkg-version/pkg/version"
//
//	app := &cli.Command{
//	    Commands: []*cli.Command{
//	        version.Command,
//	    },
//	}
//
// 支持的命令行选项：
//
//	app version          # 显示详细版本信息
//	app version -s       # 显示简短版本号
//	app version --json   # 以 JSON 格式输出
//
// # 导出变量
//
// 以下包级变量可在构建时通过 ldflags 注入，或由 init() 自动从 BuildInfo 读取：
//
//   - [AppRawName]：应用原始名称（去除日期前缀）
//   - [AppProject]：项目名称（通常为 Git 仓库名）
//   - [AppVersion]：应用版本号（语义化版本）
//   - [GitCommit]：Git 提交短哈希（7 位，dirty 时带 -dirty 后缀）
//   - [BuildTime]：构建时间（UTC+8 格式）
//   - [Developer]：开发者/维护者 URL
package version
