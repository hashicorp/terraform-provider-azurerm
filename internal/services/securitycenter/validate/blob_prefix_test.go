// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"strings"
	"testing"
)

func TestBlobPrefix(t *testing.T) {
	validPrefixes := []string{
		"container",
		"co",
		"container/",
		"container/blob",
		"con/blob",
		strings.Repeat("a", 63),
		fmt.Sprintf("container/%s", strings.Repeat("a", 479)),
	}
	for _, v := range validPrefixes {
		_, errors := BlobPrefix(v, "prefix")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid Blob Prefix: %q", v, errors)
		}
	}

	invalidPrefixes := []string{
		"",
		fmt.Sprintf("container/%s", strings.Repeat("a", 480)),
		"/blob",
		"container?",
		"container/blob\\",
		"Container",
		strings.Repeat("a", 64),
		"co/blob",
		"-container",
	}
	for _, v := range invalidPrefixes {
		if _, errors := BlobPrefix(v, "prefix"); len(errors) == 0 {
			t.Fatalf("%q should be an invalid Blob Prefix", v)
		}
	}
}
