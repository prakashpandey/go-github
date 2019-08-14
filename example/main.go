package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/prakashpandey/go-github/search/user"

	"github.com/prakashpandey/go-github/search"
)

func searchUser() []search.UserSearch {
	lang := flag.String("lang", "Golang", "programming language")
	loc := flag.String("loc", "India", "country or city")
	flag.Parse()
	return search.Users(search.NewFilter(*lang, *loc, 100, 0))
}

func getUserProfile(username string) *user.User {
	return user.GetUser(username)
}
func main() {
	users := searchUser()
	//fmt.Printf("Total users: %d\n", len(users))
	fmt.Printf("Login, Name, Email, Location, Hirable\n")
	for _, user := range users {
		userProfile := getUserProfile(user.Login)
		userProfile.Location = strings.ReplaceAll(userProfile.Location, ",", " ")
		// fmt.Printf("username: %s name: %s email: %s location: %s hireable: %t \n", userProfile.Login, userProfile.Name, userProfile.Email, userProfile.Location, userProfile.Hireable)
		fmt.Printf("%s, %s, %s, %s, %t \n", userProfile.Login, userProfile.Name, userProfile.Email, userProfile.Location, userProfile.Hireable)
	}
}
