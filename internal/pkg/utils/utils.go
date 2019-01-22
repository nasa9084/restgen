package utils

import (
	"sort"
	"strings"

	openapi "github.com/nasa9084/go-openapi"
)

func SortedPaths(paths openapi.Paths) []string {
	var ret []string
	for path := range paths {
		ret = append(ret, path)
	}
	sort.Strings(ret)
	return ret
}

func SortedMethods(pathItem *openapi.PathItem) []string {
	var ret []string
	for method := range pathItem.Operations() {
		ret = append(ret, method)
	}
	sort.Strings(ret)
	return ret
}

func NameFromRef(ref string) string {
	s := strings.Split(ref, "/")
	return s[len(s)-1]
}
