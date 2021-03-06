package tmux

import "testing"

func kill(s *Session) {
	s.Kill()
}

func TestNewSession(t *testing.T) {
	if SessionExists("foo") {
		t.Fatal("Expected to not have a session(foo)")
	}

	session, err := NewSession("foo", nil)
	defer kill(session)

	if err != nil {
		t.Fatalf("Cannot create a session: %s", err)
	}

	if !SessionExists("foo") {
		t.Fatalf("Expected to have a session(%s)", session)
	}
}

func TestListSessions(t *testing.T) {
	session, _ := NewSession("foo", nil)
	defer kill(session)

	sessions := ListSessions()

	if len(sessions) < 1 {
		t.Fatal("Expected to have at least one session")
	}

	founded := false
	for _, s := range sessions {
		if s.Name == session.Name {
			founded = true
			break
		}
	}

	if !founded {
		t.Fatalf("Expected to have a session(%s)", session)
	}
}

func TestSessionExists(t *testing.T) {
	if SessionExists("foo") {
		t.Fatal("Expect to not exists 'foo' session")
	}

	session, _ := NewSession("foo", nil)
	defer kill(session)

	if !SessionExists("foo") {
		t.Fatal("Expect to exists 'foo' session")
	}
}

func TestSession(t *testing.T) {
	session, _ := NewSession("foo", nil)
	defer kill(session)

	// Test session.SessionExists()
	if !session.Exists() {
		t.Fatal("Expect to exists 'foo' session")
	}

	// Test session.Rename()
	session.Rename("bar")

	if SessionExists("foo") {
		t.Fatal("Expect to not exists 'foo' session")
	}

	if !SessionExists("bar") {
		t.Fatal("Expect to exists 'bar' session")
	}
}
