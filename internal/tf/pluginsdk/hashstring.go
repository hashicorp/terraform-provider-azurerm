// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package pluginsdk

import (
	"hash/crc32"
	"strings"
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

// HashStringInsensitively provides case-insensitive hashing for TypeSet elements
func HashStringInsensitively(v interface{}) int {
	return HashString(strings.ToLower(v.(string)))
}
