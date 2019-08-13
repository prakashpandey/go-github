package main

import (
	"fmt"

	"github.com/prakashpandey/go-github/search/user"

	"github.com/prakashpandey/go-github/search"
)

func searchUser() []search.UserSearch {
	return search.Users(search.NewFilter("golang", "delhi", 100, 0))
}

func getUserProfile(username string) *user.User {
	return user.GetUser(username)
}
func main() {
	users := searchUser()
	fmt.Printf("Total users: %d\n", len(users))
	for _, user := range users {
		userProfile := getUserProfile(user.Login)
		fmt.Printf("username: %s name: %s email: %s location: %s hireable: %t ", userProfile.Login, userProfile.Name, userProfile.Email, userProfile.Location, userProfile.Hireable)
	}
}
