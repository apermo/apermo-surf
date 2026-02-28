package resolve

import (
	"fmt"
	"os"

	"github.com/apermo/apermo-surf/internal/config"
)

func ExampleResolve_simple() {
	link := config.Link{URL: "https://example.com"}
	result := Resolve(link, os.TempDir(), "")
	fmt.Println(result.URL)
	// Output: https://example.com
}

func ExampleResolve_withTicket() {
	link := config.Link{
		URL:     "https://jira.example.com/browse/{ticket}",
		Pattern: `PROJ-\d+`,
	}
	result := Resolve(link, os.TempDir(), "PROJ-456")
	fmt.Println(result.URL)
	// Output: https://jira.example.com/browse/PROJ-456
}
