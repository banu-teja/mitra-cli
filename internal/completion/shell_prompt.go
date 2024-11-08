package completion

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/process"
)

// If multiple steps are required, combine them using && or ;.

// Provide only {{.Shell}} commands for {{.OS}} without any description.
// If there is a lack of details, provide the most logical solution.
// Ensure the output is a valid shell command.
// Output only plain text without any formatting.
// Do not use backticks, markdown, or any other special formatting.
// Use appropriate flags and options for {{.Shell}} commands.
// Prefer long-form options (e.g., --help instead of -h) for clarity.
// Use variables when appropriate, following {{.Shell}} conventions.
// Include error handling where necessary (e.g., set -e for bash).
// For file operations, use quotation marks to handle spaces in filenames.

const shellRoleTemplate = `You are a {{.OS}} terminal running {{.Shell}}. Provide only {{.Shell}} commands for {{.OS}} without any description. If there is a lack of details, provide the most logical solution. Ensure the output is a valid shell command. Output only plain text without any formatting. Do not use backticks, markdown, or any other special formatting. Use appropriate flags and options for {{.Shell}} commands. Prefer long-form options (e.g., --help instead of -h) for clarity. Use variables when appropriate, following {{.Shell}} conventions. Include error handling where necessary. For file operations, use quotation marks to handle spaces in filenames.

`

type ShellPrompt struct {
	Shell string
	OS    string
}

type SystemInfo struct {
	OS              string
	Hostname        string
	Platform        string
	PlatformFamily  string
	PlatformVersion string
	CurrentShell    string
	OperatingSystem string
	Architecture    string
}

func GetSystemInfo() SystemInfo {
	info, err := host.Info()
	sysInfo := SystemInfo{
		OperatingSystem: runtime.GOOS,
		Architecture:    runtime.GOARCH,
		CurrentShell:    DetectShell(),
	}

	if err == nil {
		sysInfo.OS = info.OS
		sysInfo.Hostname = info.Hostname
		sysInfo.Platform = info.Platform
		sysInfo.PlatformFamily = info.PlatformFamily
		sysInfo.PlatformVersion = info.PlatformVersion
	}

	return sysInfo
}

func FormatSystemInfo(info SystemInfo) string {
	var buffer bytes.Buffer

	buffer.WriteString("=== System Information ===\n")
	buffer.WriteString(fmt.Sprintf("%-20s: %s\n", "Operating System", info.OperatingSystem))
	buffer.WriteString(fmt.Sprintf("%-20s: %s\n", "Architecture", info.Architecture))
	buffer.WriteString(fmt.Sprintf("%-20s: %s\n", "Current Shell", info.CurrentShell))
	if info.OS != "" {
		buffer.WriteString(fmt.Sprintf("%-20s: %s\n", "OS", info.OS))
		buffer.WriteString(fmt.Sprintf("%-20s: %s\n", "Hostname", info.Hostname))
		buffer.WriteString(fmt.Sprintf("%-20s: %s\n", "Platform", info.Platform))
		if info.PlatformFamily != "" {
			buffer.WriteString(fmt.Sprintf("%-20s: %s\n", "Platform Family", info.PlatformFamily))
		}
		if info.PlatformVersion != "" {
			buffer.WriteString(fmt.Sprintf("%-20s: %s\n", "Platform Version", info.PlatformVersion))
		}
	}
	buffer.WriteString("==========================\n")

	return buffer.String()
}

func GetCurrentShell() (string, error) {
	pid := os.Getppid()
	p, err := process.NewProcess(int32(pid))
	if err != nil {
		return "", err
	}
	name, err := p.Name()
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(name, filepath.Ext(name)), nil
}

func normalizeShellName(shell string) string {
	shell = strings.ToLower(shell)
	switch {
	case strings.Contains(shell, "bash"):
		return "bash"
	case strings.Contains(shell, "zsh"):
		return "zsh"
	case strings.Contains(shell, "fish"):
		return "fish"
	case strings.Contains(shell, "tcsh"):
		return "tcsh"
	case strings.Contains(shell, "csh"):
		return "csh"
	case strings.Contains(shell, "powershell"):
		return "powershell"
	case strings.Contains(shell, "cmd"):
		return "cmd"
	default:
		return shell
	}
}

func DetectShell() string {
	if runtime.GOOS == "windows" {
		return detectWindowsShell()
	}
	return detectUnixShell()
}

func detectWindowsShell() string {
	shell, err := GetCurrentShell()
	if err == nil && shell != "" {
		return normalizeShellName(shell)
	}

	if _, exists := os.LookupEnv("PSModulePath"); exists {
		return "powershell"
	}

	if shell := os.Getenv("ComSpec"); shell != "" {
		return normalizeShellName(filepath.Base(shell))
	}

	return "unknown"
}

func detectUnixShell() string {
	shell, err := GetCurrentShell()
	if err == nil && shell != "" {
		return normalizeShellName(shell)
	}

	if shell := os.Getenv("SHELL"); shell != "" {
		return normalizeShellName(filepath.Base(shell))
	}

	if runtime.GOOS == "darwin" {
		out, err := exec.Command("dscl", ".", "-read", "/Users/"+os.Getenv("USER"), "UserShell").Output()
		if err == nil {
			fields := strings.Fields(string(out))
			if len(fields) > 0 {
				return normalizeShellName(filepath.Base(fields[len(fields)-1]))
			}
		}
	}

	return "unknown"
}

func GetShellCommand(shell, os string) (string, error) {
	tmpl, err := template.New("shellRole").Parse(shellRoleTemplate)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %w", err)
	}

	prompt := ShellPrompt{
		Shell: shell,
		OS:    os,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, prompt); err != nil {
		return "", fmt.Errorf("error executing template: %w", err)
	}

	return buf.String(), nil
}
