package git

import "testing"

func TestTicket_MatchesPattern(t *testing.T) {
	ticket, err := Ticket("feature/PROJ-123-add-login", `PROJ-\d+`)
	if err != nil {
		t.Fatal(err)
	}
	if ticket != "PROJ-123" {
		t.Errorf("got %q, want PROJ-123", ticket)
	}
}

func TestTicket_NoMatch(t *testing.T) {
	ticket, err := Ticket("main", `PROJ-\d+`)
	if err != nil {
		t.Fatal(err)
	}
	if ticket != "" {
		t.Errorf("got %q, want empty", ticket)
	}
}

func TestTicket_EmptyInputs(t *testing.T) {
	ticket, err := Ticket("", `PROJ-\d+`)
	if err != nil {
		t.Fatal(err)
	}
	if ticket != "" {
		t.Errorf("got %q for empty branch, want empty", ticket)
	}

	ticket, err = Ticket("feature/PROJ-123", "")
	if err != nil {
		t.Fatal(err)
	}
	if ticket != "" {
		t.Errorf("got %q for empty pattern, want empty", ticket)
	}
}

func TestRepoNameFromURL_SSH(t *testing.T) {
	got := repoNameFromURL("git@github.com:apermo/apermo-surf.git")
	if got != "apermo-surf" {
		t.Errorf("got %q, want apermo-surf", got)
	}
}

func TestRepoNameFromURL_HTTPS(t *testing.T) {
	got := repoNameFromURL("https://github.com/apermo/apermo-surf.git")
	if got != "apermo-surf" {
		t.Errorf("got %q, want apermo-surf", got)
	}
}

func TestRepoNameFromURL_NoSuffix(t *testing.T) {
	got := repoNameFromURL("https://github.com/apermo/apermo-surf")
	if got != "apermo-surf" {
		t.Errorf("got %q, want apermo-surf", got)
	}
}
