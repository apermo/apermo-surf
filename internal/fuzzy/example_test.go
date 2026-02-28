package fuzzy

import "fmt"

func ExampleBestMatch_exact() {
	match, _ := BestMatch("production", []string{"production", "staging", "local"})
	fmt.Println(match)
	// Output: production
}

func ExampleBestMatch_fuzzy() {
	match, _ := BestMatch("prod", []string{"production", "staging", "local"})
	fmt.Println(match)
	// Output: production
}

func ExampleBestMatch_noMatch() {
	match, _ := BestMatch("zzzzzzz", []string{"production", "staging"})
	fmt.Println(match == "")
	// Output: true
}
