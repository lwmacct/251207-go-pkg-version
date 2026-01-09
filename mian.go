package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/lwmacct/251207-go-pkg-version/pkg/version"
	"github.com/urfave/cli/v3"
)

// Command 客户端命令
var Command = &cli.Command{
	Name:      "cli examples",
	Usage:     "一个 CLI 应用程序示例",
	UsageText: `演示如何使用包 version 来管理应用程序版本信息。`,
	Action:    action,
	Version:   version.GetVersion(),
	Commands:  []*cli.Command{version.Command},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "config, c",
			Usage: "指定配置文件路径",
			Value: "config.yaml",
		},
	},
}

func action(ctx context.Context, cmd *cli.Command) error {
	// 默认行为：显示帮助
	return cli.ShowAppHelp(cmd)
}

func main() {
	if err := Command.Run(context.Background(), os.Args); err != nil {
		slog.Error("应用程序运行失败", "error", err)
		os.Exit(1)
	}
}
