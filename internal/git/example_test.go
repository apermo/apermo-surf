package git

import "fmt"

func ExampleTicket() {
	ticket, _ := Ticket("feature/PROJ-123-add-login", `PROJ-\d+`)
	fmt.Println(ticket)
	// Output: PROJ-123
}

func ExampleTicket_noMatch() {
	ticket, _ := Ticket("main", `PROJ-\d+`)
	fmt.Println(ticket)
	// Output:
}
