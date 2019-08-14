package conf

import "os"

var (
	// APIURL is the base url of Github APIS.
	APIURL = "https://api.github.com"
	// OAuth is githubs auth token.
	OAuth = os.Getenv("OATH")
)
