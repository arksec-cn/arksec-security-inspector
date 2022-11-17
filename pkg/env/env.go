package env

import "os"

const (
	OpenSearchAddress  = "OPEN_SEARCH_ADDRESS"
	OpenSearchUsername = "OPEN_SEARCH_USERNAME"
	OpenSearchPassword = "OPEN_SEARCH_PASSWORD"
	OpenSearchIndex    = "OPEN_SEARCH_INDEX"
)

func GetEnv(key string, arg ...string) string {
	s := os.Getenv(key)
	if s != "" {
		return s
	}

	if len(arg) == 1 {
		s = arg[0]
	}

	return s
}
