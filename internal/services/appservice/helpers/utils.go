// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import "strings"

func trimURLScheme(input string) string {
	schemes := []string{
		"https://",
		"http://",
	}
	for _, v := range schemes {
		if strings.HasPrefix(strings.ToLower(input), v) {
			return input[len(v):]
		}
	}

	return input
}

// trimDockerFxPrefix removes the `DOCKER|` prefix from a Linux or Windows FX version string.
// App Service can return the prefix in lowercase (e.g. after the stack is saved in the portal), so the match ignores case.
func trimDockerFxPrefix(input string) string {
	prefix := "docker|"
	if strings.HasPrefix(strings.ToLower(input), prefix) {
		return input[len(prefix):]
	}

	return input
}
