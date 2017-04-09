package cmd

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const (
	binary string = "tmux"
)

func command(commandName string, args ...string) ([]string, error) {
	fullArgs := []string{commandName}

	if len(args) > 0 {
		fullArgs = append(fullArgs, args...)
	}

	output, err := exec.Command(binary, fullArgs...).CombinedOutput()

	splittedOutput := make([]string, 0)

	for _, row := range strings.Split(string(output), "\n") {
		if len(row) > 0 {
			splittedOutput = append(splittedOutput, row)
		}
	}

	return splittedOutput, err
}

// NewSession executes the `tmux new-session` command and returns an error if a
// non-zero exit occurs. The command executed is `tmux new-session -d -s name
// [args...]`.
func NewSession(name string, args ...string) error {
	/*
		  Documentation from new-session man page:

			new-session [-AdDEP] [-c start-directory] [-F format] [-n window-name] [-s session-name] [-t target-session] [-x width] [-y height] [shell-command]
			(alias: new)
			Create a new session with name session-name.

			The new session is attached to the current terminal unless -d is given.  window-name and shell-command are the name of and shell command to execute in the initial
			window.  If -d is used, -x and -y specify the size of the initial window (80 by 24 if not given).

			If run from a terminal, any termios(3) special characters are saved and used for new windows in the new session.

			The -A flag makes new-session behave like attach-session if session-name already exists; in this case, -D behaves like -d to attach-session.

			If -t is given, the new session is grouped with target-session.  This means they share the same set of windows - all windows from target-session are linked to the new
			session and any subsequent new windows or windows being closed are applied to both sessions.  The current and previous window and any session options remain indepen‐
			dent and either session may be killed without affecting the other.  Giving -n or shell-command are invalid if -t is used.

			The -P option prints information about the new session after it has been created.  By default, it uses the format ‘#{session_name}:’ but a different format may be
			specified with -F.

			If -E is used, update-environment option will not be applied.  update-environment.
	*/

	fullArgs := []string{"-s", name, "-d"}

	if len(args) > 0 {
		fullArgs = append(fullArgs, args...)
	}

	_, err := command("new-session", fullArgs...)

	return err
}

func RenameSession(oldSession, newSession string) error {
	_, err := command("rename-session", "-t", oldSession, newSession)

	return err
}

func ListSessions() []string {
	names, err := command("list-sessions", "-F", "#{session_name}")

	if err != nil {
		return make([]string, 0)
	}

	return names
}

func KillSession(name string) error {
	_, err := command("kill-session", "-t", name)

	return err
}

func NewWindow(sessionName, windowName string) error {
	_, err := command("new-window", "-t", sessionName, "-n", windowName)
	return err
}

var listWindowsFormatRE = regexp.MustCompile("^([\\d]+) (.*)$")

type ListWindowsResult struct {
	Index int
	Name  string
}

func ListWindows(sessionName string) ([]*ListWindowsResult, error) {
	lines, err := command("list-windows", "-t", sessionName, "-F", "#{window_index} #{window_name}")
	if err != nil {
		return nil, err
	}
	var wins []*ListWindowsResult
	for _, l := range lines {
		m := listWindowsFormatRE.FindStringSubmatch(l)
		if len(m) != 3 || m[2] == "" {
			return nil, fmt.Errorf("Unexpected format in list-windows result: %s", l)
		}
		idx, err := strconv.Atoi(m[1])
		if err != nil {
			return nil, err
		}
		wins = append(wins, &ListWindowsResult{Index: idx, Name: m[2]})
	}
	return wins, nil
}

func SendKeys(targetPane string, keys ...string) error {
	// target-pane may be a pane ID or takes a similar form to target-window but
	// with the optional addition of a period followed by a pane index or pane ID,
	// for example: ‘mysession:mywindow.1’
	_, err := command("send-keys", append([]string{"-t", targetPane}, keys...)...)
	return err
}
