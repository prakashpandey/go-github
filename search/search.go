package search

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// UserSearchResponse is used to map the response returned when a users are searched based on the filter query.
type UserSearchResponse struct {
	TotalCount       int          `json:"total_count"`
	IncompleteResult bool         `json:"incomplete_results"`
	Items            []UserSearch `json:"items"`
}

// UserSearch is used to map the response of user search API.
type UserSearch struct {
	Login     string `json:"login"`
	ID        string `json:"id"`
	NodeID    string `json:"node_id"`
	AvatarURL string `json:"avatar_url"`
	URL       string `json:"url"`
	HTMLURL   string `json:"html_url"`
	Type      string `json:"type"`
	SiteAdmin string `json:"site_admin"`
	Score     int    `json:"score"`
}

const defaultPerPageUserSearchLimit = 100

// UserFilter defines the filter criteria use to search a user.
type UserFilter struct {
	Lang     string // programming language, not to be confused with human language.
	Location string
	// If nil, default is 100.
	PerPage int
	// Page is used for pagination.
	// Page < 1 means do not use this as query param.
	Page int
}

// Users searches github users based on the filter criteria.
func Users(filter *UserFilter) []UserSearch {
	api := "/search/users?"
	addQualifier := func(qualifier string) {
		if strings.HasSuffix(api, "?") {
			api = fmt.Sprintf("%s%s", api, qualifier)
		} else {
			api = fmt.Sprintf("%s+%s", api, qualifier)
		}
	}
	if filter.Lang != "" {
		addQualifier(fmt.Sprintf("language:%s", filter.Lang))
	}
	if filter.Location != "" {
		addQualifier(fmt.Sprintf("location:%s", fmt.Sprintf("%s", filter.Location)))
	}
	if filter.PerPage != 0 {
		addQualifier(fmt.Sprintf("&per_page=%s", fmt.Sprintf("%d", filter.PerPage)))
	} else {
		addQualifier(fmt.Sprintf("&per_page=%s", fmt.Sprintf("%d", defaultPerPageUserSearchLimit)))
	}
	getUsers := func(page int) (*http.Response, error) {
		if page > 0 {
			addQualifier(fmt.Sprintf("&page=%s", fmt.Sprintf("%d", filter.Page)))
		}
		return http.Get(api)
	}
	// upperBound is used to run loop at least onces.
	var upperBound int
	if filter.Page < 1 {
		upperBound = 1
	} else {
		upperBound = filter.Page
	}
	var users []UserSearch
	for i := 1; i <= upperBound; i++ {
		if filter.Page < 1 {
			upperBound = filter.Page
		}
		var resp *http.Response
		resp, err := getUsers(upperBound)
		if err != nil {
			break
		}
		searchResp := &UserSearchResponse{}
		if err := json.NewDecoder(resp.Body).Decode(searchResp); err != nil {
			for _, usr := range searchResp.Items {
				users = append(users, usr)
			}
		}
		_ = resp.Body.Close()
	}
	return users
}
