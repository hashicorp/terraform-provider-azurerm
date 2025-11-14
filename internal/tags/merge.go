// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tags

// MergeDefaultTags merges provider default tags with resource-specific tags.
// Resource tags take precedence over default tags in case of conflicts.
// Returns a map[string]interface{} suitable for Terraform operations.
func MergeDefaultTags(defaultTags map[string]*string, resourceTags map[string]interface{}) map[string]interface{} {
	mergedTags := make(map[string]interface{})

	// Add default tags first
	if defaultTags != nil {
		for k, v := range defaultTags {
			mergedTags[k] = *v
		}
	}

	// Add/override with resource-specific tags
	for k, v := range resourceTags {
		mergedTags[k] = v
	}

	return mergedTags
}

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
