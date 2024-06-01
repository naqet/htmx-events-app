package chttp

import "fmt"

func resolvePattern(basePath, path, method string) string {
	pattern := path

	if len(basePath) > 0 && basePath != "/" {
		if pattern == "/{$}" {
			pattern = basePath
		} else {
			pattern = basePath + path
		}
	}

	return fmt.Sprintf("%s %s", method, pattern)
}
