package main

import (
	"fmt"

	"github.com/prakashpandey/go-github/search"
)

func searchUser() []search.UserSearch {
	return search.Users(search.NewFilter("golang", "delhi", 100, 0))
}

func main() {
	totalUsers := searchUser()
	fmt.Printf("Total users: %d\n", len(totalUsers))
}
