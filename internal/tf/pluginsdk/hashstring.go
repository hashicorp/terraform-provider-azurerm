package pluginsdk

import (
	"hash/crc32"
)

// HashString hashes strings. If you want a Set of strings, this is the
// SchemaSetFunc you want.
func HashString(input interface{}) int {
	v := int(crc32.ChecksumIEEE([]byte(input.(string))))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}
