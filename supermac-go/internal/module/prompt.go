package module

// PromptInterface handles interactive user input. Mockable for tests.
type PromptInterface interface {
	Confirm(msg string) (bool, error)
	Input(msg string) (string, error)
	Select(msg string, opts []string) (int, error)
}

// TerminalPrompt is the real implementation using stdin/stdout.
type TerminalPrompt struct{}

func (TerminalPrompt) Confirm(msg string) (bool, error) {
	// TODO: implement with bufio.Scanner
	return false, nil
}

func (TerminalPrompt) Input(msg string) (string, error) {
	// TODO: implement with bufio.Scanner
	return "", nil
}

func (TerminalPrompt) Select(msg string, opts []string) (int, error) {
	// TODO: implement with bufio.Scanner
	return -1, nil
}

// AutoYesPrompt always returns true/empty/first — used when --yes flag is set.
type AutoYesPrompt struct{}

func (AutoYesPrompt) Confirm(msg string) (bool, error)  { return true, nil }
func (AutoYesPrompt) Input(msg string) (string, error)   { return "", nil }
func (AutoYesPrompt) Select(msg string, opts []string) (int, error) { return 0, nil }
