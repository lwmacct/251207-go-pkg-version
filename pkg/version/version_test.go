package version

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"
)

// captureStdout 捕获标准输出
func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

// saveAndRestore 保存并恢复全局变量
func saveAndRestore() func() {
	saved := struct {
		appRawName, appProject, appVersion string
		gitCommit, buildTime, developer    string
	}{
		AppRawName, AppProject, AppVersion,
		GitCommit, BuildTime, Developer,
	}
	return func() {
		AppRawName = saved.appRawName
		AppProject = saved.appProject
		AppVersion = saved.appVersion
		GitCommit = saved.gitCommit
		BuildTime = saved.buildTime
		Developer = saved.developer
	}
}

func TestFormatBuildTime(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid RFC3339 time",
			input:    "2024-12-07T10:30:00Z",
			expected: "2024-12-07 18:30:00 CST",
		},
		{
			name:     "valid RFC3339 with timezone",
			input:    "2024-01-15T08:00:00+00:00",
			expected: "2024-01-15 16:00:00 CST",
		},
		{
			name:     "invalid time format",
			input:    "not-a-time",
			expected: "not-a-time",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatBuildTime(tt.input)
			if result != tt.expected {
				t.Errorf("formatBuildTime(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetVersion(t *testing.T) {
	restore := saveAndRestore()
	defer restore()

	tests := []struct {
		name       string
		appVersion string
		gitCommit  string
		expected   string
	}{
		{
			name:       "with valid version",
			appVersion: "v1.2.3",
			gitCommit:  "abc1234",
			expected:   "v1.2.3",
		},
		{
			name:       "unknown version with git commit",
			appVersion: "Unknown",
			gitCommit:  "abc1234",
			expected:   "dev-abc1234",
		},
		{
			name:       "unknown version with dirty commit",
			appVersion: "Unknown",
			gitCommit:  "abc1234-dirty",
			expected:   "dev-abc1234-dirty",
		},
		{
			name:       "unknown version and unknown commit",
			appVersion: "Unknown",
			gitCommit:  "Unknown",
			expected:   "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AppVersion = tt.appVersion
			GitCommit = tt.gitCommit

			result := GetVersion()
			if result != tt.expected {
				t.Errorf("GetVersion() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestGetBuildInfo(t *testing.T) {
	restore := saveAndRestore()
	defer restore()

	AppVersion = "v1.0.0"
	GitCommit = "abc1234"
	BuildTime = "2024-12-07 18:30:00 CST"

	result := GetBuildInfo()
	expected := "版本: v1.0.0, 提交: abc1234, 构建时间: 2024-12-07 18:30:00 CST"

	if result != expected {
		t.Errorf("GetBuildInfo() = %q, want %q", result, expected)
	}
}

func TestPrintBuildInfo(t *testing.T) {
	restore := saveAndRestore()
	defer restore()

	AppRawName = "myapp"
	AppProject = "251207-myapp"
	AppVersion = "v1.0.0"
	GitCommit = "abc1234"
	BuildTime = "2024-12-07 18:30:00 CST"
	Developer = "http://github.com/testuser"

	output := captureStdout(func() {
		PrintBuildInfo()
	})

	expectedLines := []string{
		"AppRawName:   myapp",
		"AppVersion:   v1.0.0",
		"Git Commit:   abc1234",
		"Build Time:   2024-12-07 18:30:00 CST",
		"AppProject:   251207-myapp",
		"Developer :   http://github.com/testuser",
	}

	for _, line := range expectedLines {
		if !strings.Contains(output, line) {
			t.Errorf("PrintBuildInfo() output missing %q", line)
		}
	}
}

func TestPrintVersionJSON(t *testing.T) {
	restore := saveAndRestore()
	defer restore()

	AppRawName = "myapp"
	AppProject = "251207-myapp"
	AppVersion = "v1.0.0"
	GitCommit = "abc1234"
	BuildTime = "2024-12-07 18:30:00 CST"
	Developer = "http://github.com/testuser"

	output := captureStdout(func() {
		PrintVersionJSON()
	})

	// 验证输出是有效的 JSON
	var result map[string]string
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("PrintVersionJSON() output is not valid JSON: %v\nOutput: %s", err, output)
	}

	// 验证字段值
	expectedFields := map[string]string{
		"appRawName": "myapp",
		"appProject": "251207-myapp",
		"appVersion": "v1.0.0",
		"gitCommit":  "abc1234",
		"buildTime":  "2024-12-07 18:30:00 CST",
		"developer":  "http://github.com/testuser",
	}

	for key, expected := range expectedFields {
		if result[key] != expected {
			t.Errorf("PrintVersionJSON() field %q = %q, want %q", key, result[key], expected)
		}
	}
}

func TestDatePrefixRegex(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "with date prefix",
			input:    "251207-myapp",
			expected: "myapp",
		},
		{
			name:     "without date prefix",
			input:    "myapp",
			expected: "myapp",
		},
		{
			name:     "only date prefix",
			input:    "251207-",
			expected: "",
		},
		{
			name:     "partial match",
			input:    "25120-myapp",
			expected: "25120-myapp",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := datePrefix.ReplaceAllString(tt.input, "")
			if result != tt.expected {
				t.Errorf("datePrefix.ReplaceAllString(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
