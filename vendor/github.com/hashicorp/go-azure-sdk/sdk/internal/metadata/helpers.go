package metadata

import "strings"

func normalizeResourceId(resourceId string) string {
	return strings.TrimRight(resourceId, "/")
}
