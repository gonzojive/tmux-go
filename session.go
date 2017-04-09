package tmux

import (
	"fmt"
	"github.com/gonzojive/tmux-go/cmd"
)

type Session struct {
	Name string
}

func (s Session) String() string {
	return s.Name
}

func (s Session) Kill() error {
	return cmd.KillSession(s.Name)
}

func (s *Session) Rename(name string) {
	err := cmd.RenameSession(s.Name, name)

	if err != nil {
		return
	}

	s.Name = name
}

func (s Session) Exists() bool {
	for _, session := range ListSessions() {
		if session.Name == s.Name {
			return true
		}
	}

	return false
}

func (s *Session) SendKeys(windowName string, keys ...string) error {
	return cmd.SendKeys(fmt.Sprintf("%s:%s.0", s.Name, windowName), keys...)
}

func (s *Session) EnsureWindowExists(name string) error {
	exists, err := s.windowExists(name)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return cmd.NewWindow(s.Name, name)
}

func (s *Session) windowExists(name string) (bool, error) {
	wins, err := cmd.ListWindows(s.Name)
	if err != nil {
		return false, err
	}
	for _, w := range wins {
		if w.Name == name {
			return true, nil
		}
	}
	return false, nil
}

type NewSessionOptions struct {
	// WindowName is the name of the initial window
	WindowName string
}

func (o *NewSessionOptions) args() []string {
	if o == nil {
		return nil
	}
	var args []string
	if o.WindowName != "" {
		args = append(args, "-n", o.WindowName)
	}
	return args
}

func NewSession(name string, opts *NewSessionOptions) (*Session, error) {
	err := cmd.NewSession(name, opts.args()...)

	return &Session{Name: name}, err
}

func SessionExists(name string) bool {
	for _, session := range ListSessions() {
		if session.Name == name {
			return true
		}
	}

	return false
}

func ListSessions() []*Session {
	names := cmd.ListSessions()

	var sessions []*Session

	for _, name := range names {
		sessions = append(sessions, &Session{Name: name})
	}

	return sessions
}
