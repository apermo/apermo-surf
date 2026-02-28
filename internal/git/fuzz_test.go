package git

import "testing"

func FuzzTicket(f *testing.F) {
	f.Add("feature/PROJ-123", `PROJ-\d+`)
	f.Add("main", `PROJ-\d+`)
	f.Add("", "")
	f.Add("feature/PROJ-123-add-login", `PROJ-\d+`)
	f.Add("bugfix/123", `\d+`)
	f.Add("release/v1.0.0", `v\d+\.\d+\.\d+`)
	f.Fuzz(func(t *testing.T, branch, pattern string) {
		// Must never panic, even with adversarial inputs
		Ticket(branch, pattern)
	})
}
