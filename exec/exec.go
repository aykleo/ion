package exec

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func ExecSysCommand(command string, args []string) tea.Cmd {
	fullCommand := command
	if len(args) > 0 {
		fullCommand = command + " " + strings.Join(args, " ")
	}

	return func() tea.Msg {
		if strings.HasPrefix(strings.ToLower(strings.TrimSpace(fullCommand)), "cd ") {
			return handleCDCommand(fullCommand)
		}

		cmd := exec.Command("powershell", "-NoProfile", "-Command", fullCommand)
		cmd.Dir = currentDir

		output, err := cmd.CombinedOutput()

		return CommandFinishedMsg{
			Err:     err,
			Command: fullCommand,
			Output:  string(output),
			NewDir:  currentDir,
		}
	}
}

func handleCDCommand(fullCommand string) CommandFinishedMsg {
	parts := strings.Fields(fullCommand)
	if len(parts) < 2 {
		if homeDir, err := os.UserHomeDir(); err == nil {
			currentDir = homeDir
			return CommandFinishedMsg{
				Command: fullCommand,
				// Output:  "moved to " + currentDir,
				NewDir: currentDir,
			}
		}
	}

	targetPath := parts[1]

	if !filepath.IsAbs(targetPath) {
		targetPath = filepath.Join(currentDir, targetPath)
	}

	targetPath = filepath.Clean(targetPath)

	if info, err := os.Stat(targetPath); err != nil {
		return CommandFinishedMsg{
			Err:     err,
			Command: fullCommand,
			Output:  "directory not found " + targetPath,
			NewDir:  currentDir,
		}
	} else if !info.IsDir() {
		return CommandFinishedMsg{
			Err:     nil,
			Command: fullCommand,
			Output:  "not a directory " + targetPath,
			NewDir:  currentDir,
		}
	}

	currentDir = targetPath

	return CommandFinishedMsg{
		Command: fullCommand,
		// Output:  "moved to " + currentDir,
		NewDir: currentDir,
	}
}
