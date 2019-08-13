package search

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/prakashpandey/go-github/conf"
)

// UserSearchResponse is used to map the response returned when a users are searched based on the filter query.
type UserSearchResponse struct {
	TotalCount       int          `json:"total_count"`
	IncompleteResult bool         `json:"incomplete_results"`
	Items            []UserSearch `json:"items"`
}

// UserSearch is used to map the response of user search API.
type UserSearch struct {
	Login     string  `json:"login"`
	ID        int     `json:"id"`
	NodeID    string  `json:"node_id"`
	AvatarURL string  `json:"avatar_url"`
	URL       string  `json:"url"`
	HTMLURL   string  `json:"html_url"`
	Type      string  `json:"type"`
	SiteAdmin bool    `json:"site_admin"`
	Score     float64 `json:"score"`
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

// NewFilter returns new filter object.
func NewFilter(lang, loc string, perPage, page int) *UserFilter {
	f := &UserFilter{Lang: lang, Location: loc, Page: page}
	if perPage == 0 {
		f.PerPage = defaultPerPageUserSearchLimit
	}
	return f
}

// Users searches github users based on the filter criteria.
func Users(filter *UserFilter) []UserSearch {
	url := fmt.Sprintf("%s%s", conf.APIURL, "/search/users?q=")
	addQualifier := func(qualifier string) {
		if strings.HasSuffix(url, "q=") {
			url = fmt.Sprintf("%s%s", url, qualifier)
		} else {
			url = fmt.Sprintf("%s+%s", url, qualifier)
		}
	}
	addParam := func(param string) {
		url = fmt.Sprintf("%s&%s", url, param)
	}
	if filter.Lang != "" {
		addQualifier(fmt.Sprintf("language:%s", filter.Lang))
	}
	if filter.Location != "" {
		addQualifier(fmt.Sprintf("location:%s", fmt.Sprintf("%s", filter.Location)))
	}
	if filter.PerPage > 0 {
		addQualifier(fmt.Sprintf("&per_page=%s", fmt.Sprintf("%d", filter.PerPage)))
	} else {
		filter.PerPage = defaultPerPageUserSearchLimit
		addParam(fmt.Sprintf("per_page=%s", fmt.Sprintf("%d", filter.PerPage)))
	}
	getUsers := func(page int) ([]UserSearch, int) {
		addParam(fmt.Sprintf("page=%s", fmt.Sprintf("%d", filter.Page)))
		var users []UserSearch
		fmt.Printf("url: %s\n", url)
		resp, err := http.Get(url)
		if err != nil {
			return users, 0
		}
		searchResp := &UserSearchResponse{}
		if err := json.NewDecoder(resp.Body).Decode(searchResp); err != nil {
			fmt.Printf("error decoding response: %s", err)
			return users, 0
		}
		for _, usr := range searchResp.Items {
			users = append(users, usr)
		}
		_ = resp.Body.Close()
		return users, searchResp.TotalCount
	}
	var users []UserSearch
	if filter.Page > 1 {
		users, _ = getUsers(filter.Page)
	} else {
		// Find all users.
		page := 1 // initial page
		for i := 1; i <= page; i++ {
			usrs, total := getUsers(page)
			// In the first iteration, calculate the page size.
			if i == 1 {
				if total%filter.PerPage == 0 {
					page = total / filter.PerPage
				} else {
					page = total/filter.PerPage + 1
				}
			}
			for _, usr := range usrs {
				users = append(users, usr)
			}
		}
	}
	return users
}
