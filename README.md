# 251207-go-pkg-version

[![License](https://img.shields.io/github/license/lwmacct/251207-go-pkg-version)](LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/lwmacct/251207-go-pkg-version.svg)](https://pkg.go.dev/github.com/lwmacct/251207-go-pkg-version)
[![Go CI](https://github.com/lwmacct/251207-go-pkg-version/actions/workflows/go-ci.yml/badge.svg)](https://github.com/lwmacct/251207-go-pkg-version/actions/workflows/go-ci.yml)
[![codecov](https://codecov.io/gh/lwmacct/251207-go-pkg-version/branch/main/graph/badge.svg)](https://codecov.io/gh/lwmacct/251207-go-pkg-version)
[![Go Report Card](https://goreportcard.com/badge/github.com/lwmacct/251207-go-pkg-version)](https://goreportcard.com/report/github.com/lwmacct/251207-go-pkg-version)
[![GitHub Tag](https://img.shields.io/github/v/tag/lwmacct/251207-go-pkg-version?sort=semver)](https://github.com/lwmacct/251207-go-pkg-version/tags)

Go 应用程序版本信息管理库，支持 ldflags 注入和 BuildInfo 自动读取。

## 安装

```bash
go get github.com/lwmacct/251207-go-pkg-version
```

## 快速开始

```go
import "github.com/lwmacct/251207-go-pkg-version/pkg/version"

// 集成到 urfave/cli/v3
app := &cli.Command{
    Commands: []*cli.Command{version.Command},
}
```

构建时注入版本信息：

```bash
PKG=github.com/lwmacct/251207-go-pkg-version/pkg/version
go build -ldflags "-X ${PKG}.AppVersion=v1.0.0 -X ${PKG}.GitCommit=$(git rev-parse --short HEAD)"
```

## 输出示例

```
AppRawName:   go-pkg-version
AppVersion:   v1.0.0
Go Version:   go1.21.0
Git Commit:   abc1234
Build Time:   2024-12-07 18:30:00 CST
AppProject:   251207-go-pkg-version
Developer :   http://github.com/lwmacct
```

## License

MIT
