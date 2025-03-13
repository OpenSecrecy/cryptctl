package utils

import (
	"regexp"
)

var statusRegexp = regexp.MustCompile(`status:[\s\S]*?(?:---|$)`)

func RemoveStatus(yaml []byte) []byte {
	// just a simple hack to remove status field from the yaml
	result := statusRegexp.ReplaceAll(yaml, nil)
	// finally, return the result
	return result
}
