package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/prakashpandey/go-github/conf"
)

// User is a github user's data structure.
type User struct {
	Login     string  `json:"login"`
	ID        int     `json:"id"`
	NodeID    string  `json:"node_id"`
	AvatarURL string  `json:"avatar_url"`
	URL       string  `json:"url"`
	HTMLURL   string  `json:"html_url"`
	Type      string  `json:"type"`
	SiteAdmin bool    `json:"site_admin"`
	Score     float64 `json:"score"`
	Name      string  `json:"name"`
	Company   string  `json:"company"`
	Email     string  `json:"email"`
	Location  string  `json:"location"`
	Hireable  bool    `json:"hireable"`
}

// GetUser returns user profile data.
func GetUser(username string) *User {
	url := fmt.Sprintf("%s%s%s", conf.APIURL, "/users", fmt.Sprintf("/%s", username))
	fmt.Printf("user profile url: %s\n", url)
	var user = &User{}
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("\nerror finding user: %s", err)
		return nil
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(user); err != nil {
		fmt.Printf("\nerror decoding response: %s", err)
	}
	return user
}
