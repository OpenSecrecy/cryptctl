package utils

import (
	"regexp"
)

func RemoveStatus(yaml []byte) []byte {
	// just a simple hack to remove status field from the yaml
	re := regexp.MustCompile(`status:[\s\S]*?(?:---|$)`)
	result := re.ReplaceAll(yaml, nil)
	// finally, return the result
	return result
}
