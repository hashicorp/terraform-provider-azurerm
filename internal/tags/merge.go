// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tags

import (
	"context"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

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

// SetTagsDiff computes tags_all by merging provider default_tags with resource-specific tags.
// This function allows tags_all to be computed at plan time rather than showing as (known after apply).
// It should be used with pluginsdk.CustomizeDiffShim:
//
//	CustomizeDiff: pluginsdk.CustomizeDiffShim(tags.SetTagsDiff),
func SetTagsDiff(ctx context.Context, d *pluginsdk.ResourceDiff, meta interface{}) error {
	client := meta.(*clients.Client)

	// Compute tags_all by merging default_tags and tags
	resourceTags := make(map[string]interface{})
	if tags, ok := d.GetOk("tags"); ok {
		resourceTags = tags.(map[string]interface{})
	}
	mergedTags := MergeDefaultTags(client.DefaultTags, resourceTags)

	d.SetNew("tags_all", mergedTags)
	return nil
}
