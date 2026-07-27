package git

import "fmt"

// CommitMessage represents a generated Conventional Commit message.
type CommitMessage struct {
	Type    string
	Scope   string
	Subject string
}

// String formats the message as a Conventional Commit subject.
func (c CommitMessage) String() string {
	if c.Type == "" {
		return c.Subject
	}

	if c.Scope == "" {
		return fmt.Sprintf("%s: %s", c.Type, c.Subject)
	}

	return fmt.Sprintf("%s(%s): %s", c.Type, c.Scope, c.Subject)
}
