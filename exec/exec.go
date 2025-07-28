package exec

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aykleo/ion/config"
	"github.com/aykleo/ion/data"
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

func ExecIonCommand(args []string, dataRef data.IData) tea.Cmd {
	configPath := config.GetConfigPath()

	if args[0] == "ionize" && len(args) == 1 {
		return func() tea.Msg {
			return CommandFinishedMsg{
				Command: "ionize",
				Output:  "ionize",
			}
		}
	}
	if len(args) < 3 {
		return func() tea.Msg {
			return CommandFinishedMsg{
				Err:     errors.New("ion [category] [action] <args>"),
				Command: "ion",
				Output:  "usage: ion [category] [action] <args>",
				NewDir:  currentDir,
			}
		}
	}

	category := args[1]
	action := args[2]
	args = args[3:]

	switch category {
	case "user":
		switch action {
		case "set":
			return changeUsername(args, configPath, dataRef)
		default:
			return func() tea.Msg {
				return CommandFinishedMsg{
					Command: "ion",
					Output:  "command not found",
					NewDir:  currentDir,
				}
			}
		}
	case "secret":
		switch action {
		case "add":
			return addSecret(args, configPath, dataRef)
		case "update":
			return updateSecret(args, configPath, dataRef)
		default:
			return func() tea.Msg {
				return CommandFinishedMsg{
					Command: "ion",
					Output:  "command not found",
					NewDir:  currentDir,
				}
			}
		}
	default:
		return func() tea.Msg {
			return CommandFinishedMsg{
				Command: "ion",
				Output:  "command not found",
				NewDir:  currentDir,
				Err:     errors.New("command not found"),
			}
		}
	}
}
