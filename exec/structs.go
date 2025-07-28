package exec

type CommandFinishedMsg struct {
	Err     error
	Command string
	Output  string
	NewDir  string
}
