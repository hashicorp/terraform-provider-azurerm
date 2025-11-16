// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tags

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func TestMergeDefaultTags_NoDefaults(t *testing.T) {
	resourceTags := map[string]interface{}{
		"env":  "dev",
		"team": "backend",
	}

	result := MergeDefaultTags(nil, resourceTags)

	if len(result) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(result))
	}
	if result["env"] != "dev" {
		t.Errorf("Expected env=dev, got %v", result["env"])
	}
	if result["team"] != "backend" {
		t.Errorf("Expected team=backend, got %v", result["team"])
	}
}

func TestMergeDefaultTags_NoResourceTags(t *testing.T) {
	defaultTags := map[string]*string{
		"managed_by": utils.String("terraform"),
		"owner":      utils.String("platform"),
	}

	result := MergeDefaultTags(defaultTags, map[string]interface{}{})

	if len(result) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(result))
	}
	if result["managed_by"] != "terraform" {
		t.Errorf("Expected managed_by=terraform, got %v", result["managed_by"])
	}
	if result["owner"] != "platform" {
		t.Errorf("Expected owner=platform, got %v", result["owner"])
	}
}

func TestMergeDefaultTags_BothPresent(t *testing.T) {
	defaultTags := map[string]*string{
		"managed_by": utils.String("terraform"),
		"owner":      utils.String("platform"),
	}
	resourceTags := map[string]interface{}{
		"env":  "prod",
		"team": "backend",
	}

	result := MergeDefaultTags(defaultTags, resourceTags)

	if len(result) != 4 {
		t.Errorf("Expected 4 tags, got %d", len(result))
	}
	if result["managed_by"] != "terraform" {
		t.Errorf("Expected managed_by=terraform, got %v", result["managed_by"])
	}
	if result["owner"] != "platform" {
		t.Errorf("Expected owner=platform, got %v", result["owner"])
	}
	if result["env"] != "prod" {
		t.Errorf("Expected env=prod, got %v", result["env"])
	}
	if result["team"] != "backend" {
		t.Errorf("Expected team=backend, got %v", result["team"])
	}
}

func TestMergeDefaultTags_Conflict(t *testing.T) {
	defaultTags := map[string]*string{
		"env":        utils.String("dev"),
		"managed_by": utils.String("terraform"),
	}
	resourceTags := map[string]interface{}{
		"env":  "prod", // This should override the default
		"team": "backend",
	}

	result := MergeDefaultTags(defaultTags, resourceTags)

	if len(result) != 3 {
		t.Errorf("Expected 3 tags, got %d", len(result))
	}
	// Resource tag should take precedence
	if result["env"] != "prod" {
		t.Errorf("Expected env=prod (resource override), got %v", result["env"])
	}
	if result["managed_by"] != "terraform" {
		t.Errorf("Expected managed_by=terraform, got %v", result["managed_by"])
	}
	if result["team"] != "backend" {
		t.Errorf("Expected team=backend, got %v", result["team"])
	}
}

func TestMergeDefaultTags_EmptyMaps(t *testing.T) {
	result := MergeDefaultTags(map[string]*string{}, map[string]interface{}{})

	if len(result) != 0 {
		t.Errorf("Expected 0 tags, got %d", len(result))
	}
}

func TestRemoveDefaultTags_Basic(t *testing.T) {
	allTags := map[string]*string{
		"managed_by": utils.String("terraform"),
		"owner":      utils.String("platform"),
		"env":        utils.String("prod"),
		"team":       utils.String("backend"),
	}
	defaultTags := map[string]*string{
		"managed_by": utils.String("terraform"),
		"owner":      utils.String("platform"),
	}

	result := RemoveDefaultTags(allTags, defaultTags)

	if len(result) != 2 {
		t.Errorf("Expected 2 tags remaining, got %d", len(result))
	}
	if _, exists := result["managed_by"]; exists {
		t.Errorf("Expected managed_by to be removed")
	}
	if _, exists := result["owner"]; exists {
		t.Errorf("Expected owner to be removed")
	}
	if *result["env"] != "prod" {
		t.Errorf("Expected env=prod, got %v", result["env"])
	}
	if *result["team"] != "backend" {
		t.Errorf("Expected team=backend, got %v", result["team"])
	}
}

func TestRemoveDefaultTags_NoMatches(t *testing.T) {
	allTags := map[string]*string{
		"env":  utils.String("prod"),
		"team": utils.String("backend"),
	}
	defaultTags := map[string]*string{
		"managed_by": utils.String("terraform"),
		"owner":      utils.String("platform"),
	}

	result := RemoveDefaultTags(allTags, defaultTags)

	if len(result) != 2 {
		t.Errorf("Expected 2 tags remaining, got %d", len(result))
	}
	if *result["env"] != "prod" {
		t.Errorf("Expected env=prod, got %v", result["env"])
	}
	if *result["team"] != "backend" {
		t.Errorf("Expected team=backend, got %v", result["team"])
	}
}

func TestRemoveDefaultTags_NilAllTags(t *testing.T) {
	defaultTags := map[string]*string{
		"managed_by": utils.String("terraform"),
	}

	result := RemoveDefaultTags(nil, defaultTags)

	if len(result) != 0 {
		t.Errorf("Expected 0 tags, got %d", len(result))
	}
}

func TestRemoveDefaultTags_PartialValueOverride(t *testing.T) {
	allTags := map[string]*string{
		"managed_by": utils.String("custom"),   // Different value from default
		"owner":      utils.String("platform"), // Same value as default
		"env":        utils.String("prod"),
	}
	defaultTags := map[string]*string{
		"managed_by": utils.String("terraform"),
		"owner":      utils.String("platform"),
	}

	result := RemoveDefaultTags(allTags, defaultTags)

	if len(result) != 2 {
		t.Errorf("Expected 2 tags remaining (managed_by override + env), got %d", len(result))
	}
	// managed_by should be kept because value is different from default
	if _, exists := result["managed_by"]; !exists {
		t.Errorf("Expected managed_by to be kept (resource override with different value)")
	}
	if *result["managed_by"] != "custom" {
		t.Errorf("Expected managed_by=custom, got %v", *result["managed_by"])
	}
	// owner should be removed because value matches default exactly
	if _, exists := result["owner"]; exists {
		t.Errorf("Expected owner to be removed (same key and value)")
	}
	if *result["env"] != "prod" {
		t.Errorf("Expected env=prod, got %v", *result["env"])
	}
}

// ===== Unit Tests for Ensure* Helper Functions =====

func TestEnsureTagsAllSet_NoDefaults(t *testing.T) {
	client := &clients.Client{
		DefaultTags: map[string]*string{},
	}

	// Create a minimal ResourceData with tags
	resourceSchema := map[string]*schema.Schema{
		"tags": {
			Type: schema.TypeMap,
			Elem: &schema.Schema{Type: schema.TypeString},
		},
		"tags_all": {
			Type: schema.TypeMap,
			Elem: &schema.Schema{Type: schema.TypeString},
		},
	}
	d := schema.TestResourceDataRaw(t, resourceSchema, map[string]interface{}{
		"tags": map[string]interface{}{
			"env":  "dev",
			"team": "backend",
		},
	})

	result := EnsureTagsAllSet(d, client)

	if len(result) != 2 {
		t.Errorf("Expected 2 tags in result, got %d", len(result))
	}
	if result["env"] != "dev" {
		t.Errorf("Expected env=dev, got %v", result["env"])
	}
	if result["team"] != "backend" {
		t.Errorf("Expected team=backend, got %v", result["team"])
	}

	// Verify tags_all was set in state
	tagsAll := d.Get("tags_all").(map[string]interface{})
	if len(tagsAll) != 2 {
		t.Errorf("Expected 2 tags in tags_all, got %d", len(tagsAll))
	}
}

func TestEnsureTagsAllSet_WithDefaults(t *testing.T) {
	client := &clients.Client{
		DefaultTags: map[string]*string{
			"managed_by": utils.String("terraform"),
			"owner":      utils.String("platform"),
		},
	}

	resourceSchema := map[string]*schema.Schema{
		"tags": {
			Type: schema.TypeMap,
			Elem: &schema.Schema{Type: schema.TypeString},
		},
		"tags_all": {
			Type: schema.TypeMap,
			Elem: &schema.Schema{Type: schema.TypeString},
		},
	}
	d := schema.TestResourceDataRaw(t, resourceSchema, map[string]interface{}{
		"tags": map[string]interface{}{
			"env": "prod",
		},
	})

	result := EnsureTagsAllSet(d, client)

	if len(result) != 3 {
		t.Errorf("Expected 3 tags in result (2 defaults + 1 resource), got %d", len(result))
	}
	if result["managed_by"] != "terraform" {
		t.Errorf("Expected managed_by=terraform, got %v", result["managed_by"])
	}
	if result["owner"] != "platform" {
		t.Errorf("Expected owner=platform, got %v", result["owner"])
	}
	if result["env"] != "prod" {
		t.Errorf("Expected env=prod, got %v", result["env"])
	}

	// Verify tags_all was set in state
	tagsAll := d.Get("tags_all").(map[string]interface{})
	if len(tagsAll) != 3 {
		t.Errorf("Expected 3 tags in tags_all, got %d", len(tagsAll))
	}
}

func TestEnsureTagsAllSet_WithOverride(t *testing.T) {
	client := &clients.Client{
		DefaultTags: map[string]*string{
			"env": utils.String("dev"),
		},
	}

	resourceSchema := map[string]*schema.Schema{
		"tags": {
			Type: schema.TypeMap,
			Elem: &schema.Schema{Type: schema.TypeString},
		},
		"tags_all": {
			Type: schema.TypeMap,
			Elem: &schema.Schema{Type: schema.TypeString},
		},
	}
	d := schema.TestResourceDataRaw(t, resourceSchema, map[string]interface{}{
		"tags": map[string]interface{}{
			"env": "prod", // Resource overrides default
		},
	})

	result := EnsureTagsAllSet(d, client)

	if len(result) != 1 {
		t.Errorf("Expected 1 tag in result, got %d", len(result))
	}
	// Resource tag should override default
	if result["env"] != "prod" {
		t.Errorf("Expected env=prod (resource override), got %v", result["env"])
	}
}

func TestEnsureTagsAllReadSet_Basic(t *testing.T) {
	client := &clients.Client{
		DefaultTags: map[string]*string{
			"managed_by": utils.String("terraform"),
		},
	}

	resourceSchema := map[string]*schema.Schema{
		"tags": {
			Type: schema.TypeMap,
			Elem: &schema.Schema{Type: schema.TypeString},
		},
		"tags_all": {
			Type: schema.TypeMap,
			Elem: &schema.Schema{Type: schema.TypeString},
		},
	}
	d := schema.TestResourceDataRaw(t, resourceSchema, map[string]interface{}{})

	// API returns all tags (including defaults)
	apiTags := map[string]*string{
		"managed_by": utils.String("terraform"),
		"env":        utils.String("prod"),
	}

	err := EnsureTagsAllReadSet(d, apiTags, client)
	if err != nil {
		t.Fatalf("EnsureTagsAllReadSet failed: %v", err)
	}

	// tags_all should have all tags
	tagsAll := d.Get("tags_all").(map[string]interface{})
	if len(tagsAll) != 2 {
		t.Errorf("Expected 2 tags in tags_all, got %d", len(tagsAll))
	}

	// tags should only have resource-specific tags (managed_by removed as it matches default)
	tags := d.Get("tags").(map[string]interface{})
	if len(tags) != 1 {
		t.Errorf("Expected 1 tag in tags (only env), got %d", len(tags))
	}
	if tags["env"] != "prod" {
		t.Errorf("Expected env=prod in tags, got %v", tags["env"])
	}
	if _, exists := tags["managed_by"]; exists {
		t.Errorf("Expected managed_by to be removed from tags (it's a default)")
	}
}

func TestEnsureTagsAllReadSetFromStringMap_Basic(t *testing.T) {
	client := &clients.Client{
		DefaultTags: map[string]*string{
			"owner": utils.String("platform"),
		},
	}

	resourceSchema := map[string]*schema.Schema{
		"tags": {
			Type: schema.TypeMap,
			Elem: &schema.Schema{Type: schema.TypeString},
		},
		"tags_all": {
			Type: schema.TypeMap,
			Elem: &schema.Schema{Type: schema.TypeString},
		},
	}
	d := schema.TestResourceDataRaw(t, resourceSchema, map[string]interface{}{})

	// API returns map[string]string (e.g., Storage Account)
	apiTags := map[string]string{
		"owner": "platform",
		"env":   "staging",
	}

	err := EnsureTagsAllReadSetFromStringMap(d, apiTags, client)
	if err != nil {
		t.Fatalf("EnsureTagsAllReadSetFromStringMap failed: %v", err)
	}

	// tags_all should have all tags
	tagsAll := d.Get("tags_all").(map[string]interface{})
	if len(tagsAll) != 2 {
		t.Errorf("Expected 2 tags in tags_all, got %d", len(tagsAll))
	}

	// tags should only have resource-specific tags
	tags := d.Get("tags").(map[string]interface{})
	if len(tags) != 1 {
		t.Errorf("Expected 1 tag in tags (only env), got %d", len(tags))
	}
	if tags["env"] != "staging" {
		t.Errorf("Expected env=staging in tags, got %v", tags["env"])
	}
}

// Note: SetTagsDiff is tested through acceptance tests in resource_group_resource_test.go
// because it requires a full pluginsdk.ResourceDiff implementation which is complex to mock.
