package version

import (
	"context"
	"testing"

	"github.com/urfave/cli/v3"
)

func TestCommand(t *testing.T) {
	// 验证命令基本属性
	if Command.Name != "version" {
		t.Errorf("Command.Name = %q, want %q", Command.Name, "version")
	}

	if Command.Usage == "" {
		t.Error("Command.Usage should not be empty")
	}

	// 验证 flags 存在
	if len(Command.Flags) != 2 {
		t.Errorf("Command should have 2 flags, got %d", len(Command.Flags))
	}

	// 验证 short flag
	var hasShort, hasJSON bool
	for _, flag := range Command.Flags {
		if bf, ok := flag.(*cli.BoolFlag); ok {
			switch bf.Name {
			case "short":
				hasShort = true
				if len(bf.Aliases) == 0 || bf.Aliases[0] != "s" {
					t.Error("short flag should have alias 's'")
				}
			case "json":
				hasJSON = true
				if len(bf.Aliases) == 0 || bf.Aliases[0] != "j" {
					t.Error("json flag should have alias 'j'")
				}
			}
		}
	}

	if !hasShort {
		t.Error("Command should have 'short' flag")
	}
	if !hasJSON {
		t.Error("Command should have 'json' flag")
	}
}

func TestAction(t *testing.T) {
	restore := saveAndRestore()
	defer restore()

	// 设置测试数据
	AppRawName = "testapp"
	AppProject = "test-project"
	AppVersion = "v1.0.0"
	GitCommit = "abc1234"
	BuildTime = "2024-12-07 18:30:00 CST"
	Developer = "http://github.com/testuser"

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "default action",
			args: []string{"version"},
		},
		{
			name: "short action",
			args: []string{"version", "--short"},
		},
		{
			name: "json action",
			args: []string{"version", "--json"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 捕获输出但不验证内容，只验证不报错
			_ = captureStdout(func() {
				err := Command.Run(context.Background(), tt.args)
				if err != nil {
					t.Errorf("Command.Run() error = %v", err)
				}
			})
		})
	}
}
