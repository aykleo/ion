package exec

type CommandFinishedMsg struct {
	Err         error
	IsSystemCmd bool
	Command     string
	Output      string
	NewDir      string
}

type CommandClearMsg struct{}
