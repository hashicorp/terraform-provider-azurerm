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
// Only removes tags that have both the same key AND value as defaults.
// If a resource tag overrides a default (same key, different value), it's kept as a resource-specific tag.
func RemoveDefaultTags(allTags, defaultTags map[string]*string) map[string]*string {
	resourceTags := make(map[string]*string)

	// Copy all tags first
	if allTags != nil {
		for k, v := range allTags {
			resourceTags[k] = v
		}
	}

	// Remove only tags that exactly match defaults (same key AND value)
	if defaultTags != nil {
		for k, defaultValue := range defaultTags {
			if resourceValue, exists := resourceTags[k]; exists {
				// Only remove if the value matches the default
				if resourceValue != nil && defaultValue != nil && *resourceValue == *defaultValue {
					delete(resourceTags, k)
				}
			}
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

// EnsureTagsAllSet merges provider default tags with resource-specific tags and sets tags_all in state.
// This is a convenience helper for Create/Update operations to avoid boilerplate code.
// Usage: mergedTags := tags.EnsureTagsAllSet(d, client)
func EnsureTagsAllSet(d *pluginsdk.ResourceData, client *clients.Client) map[string]interface{} {
	resourceTags := make(map[string]interface{})
	if tagsRaw := d.Get("tags"); tagsRaw != nil {
		resourceTags = tagsRaw.(map[string]interface{})
	}
	mergedTags := MergeDefaultTags(client.DefaultTags, resourceTags)

	// Set tags_all early so it shows in plan output
	d.Set("tags_all", mergedTags)

	return mergedTags
}

// EnsureTagsAllReadSet handles tag flattening and separation for Read operations.
// It sets tags_all to all tags from Azure (includes defaults) and tags to only resource-specific tags.
// This is a convenience helper for Read operations for APIs that return map[string]*string tags.
// Usage: return tags.EnsureTagsAllReadSet(d, resp.Tags, client)
func EnsureTagsAllReadSet(d *pluginsdk.ResourceData, apiTags map[string]*string, client *clients.Client) error {
	// Set tags_all to all tags from Azure (includes defaults)
	if err := d.Set("tags_all", Flatten(apiTags)); err != nil {
		return err
	}

	// Set tags to only resource-specific tags (remove defaults)
	resourceTags := RemoveDefaultTags(apiTags, client.DefaultTags)
	return FlattenAndSet(d, resourceTags)
}

// EnsureTagsAllReadSetFromStringMap handles tag flattening and separation for Read operations.
// It sets tags_all to all tags from Azure (includes defaults) and tags to only resource-specific tags.
// This is a convenience helper for Read operations for APIs that return map[string]string tags (e.g., Storage Account).
// Usage: return tags.EnsureTagsAllReadSetFromStringMap(d, resp.Tags, client)
func EnsureTagsAllReadSetFromStringMap(d *pluginsdk.ResourceData, apiTags map[string]string, client *clients.Client) error {
	// Convert map[string]string to map[string]*string for internal handling
	convertedTags := make(map[string]*string)
	if apiTags != nil {
		for k, v := range apiTags {
			v := v // capture loop variable
			convertedTags[k] = &v
		}
	}

	// Now use the standard handler
	return EnsureTagsAllReadSet(d, convertedTags, client)
}
