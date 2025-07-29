package exec

type CommandFinishedMsg struct {
	Err     error
	Neutral bool
	Command string
	Output  string
	NewDir  string
}
