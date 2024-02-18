package util

import "regexp"

func RegexpMatchWithGroup(text string, exp *regexp.Regexp) string {

	matches := exp.FindStringSubmatch(text)
	if len(matches) < 2 {
		return ""
	}
	return matches[len(matches)-1]
}
