package client

import (
	"strings"
)

func MultiUrlParse(url string) (isMultiUrl bool, urlList []string) {
	if strings.Contains(url, ",") {
		return true, strings.Split(url, ",")
	}
	return false, []string{}
}
