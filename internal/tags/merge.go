// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tags

// RemoveDefaultTags removes default tags from a tag map, returning only resource-specific tags.
// This is used to separate user-defined tags from provider default_tags.
func RemoveDefaultTags(allTags, defaultTags map[string]*string) map[string]*string {
	resourceTags := make(map[string]*string)

	// Copy all tags first
	if allTags != nil {
		for k, v := range allTags {
			resourceTags[k] = v
		}
	}

	// Remove any tags that match defaults
	if defaultTags != nil {
		for k := range defaultTags {
			delete(resourceTags, k)
		}
	}

	return resourceTags
}
