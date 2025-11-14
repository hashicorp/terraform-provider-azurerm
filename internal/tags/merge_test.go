// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tags

import (
	"testing"
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
		"managed_by": stringPtr("terraform"),
		"owner":      stringPtr("platform"),
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
		"managed_by": stringPtr("terraform"),
		"owner":      stringPtr("platform"),
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
		"env":        stringPtr("dev"),
		"managed_by": stringPtr("terraform"),
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
	result := MergeDefaultTags(nil, map[string]interface{}{})

	if len(result) != 0 {
		t.Errorf("Expected 0 tags, got %d", len(result))
	}
}

func TestMergeDefaultTags_EmptyDefaultsEmptyResource(t *testing.T) {
	result := MergeDefaultTags(map[string]*string{}, map[string]interface{}{})

	if len(result) != 0 {
		t.Errorf("Expected 0 tags, got %d", len(result))
	}
}

func TestRemoveDefaultTags_Basic(t *testing.T) {
	allTags := map[string]*string{
		"managed_by": stringPtr("terraform"),
		"owner":      stringPtr("platform"),
		"env":        stringPtr("prod"),
		"team":       stringPtr("backend"),
	}
	defaultTags := map[string]*string{
		"managed_by": stringPtr("terraform"),
		"owner":      stringPtr("platform"),
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
		"env":  stringPtr("prod"),
		"team": stringPtr("backend"),
	}
	defaultTags := map[string]*string{
		"managed_by": stringPtr("terraform"),
		"owner":      stringPtr("platform"),
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
		"managed_by": stringPtr("terraform"),
	}

	result := RemoveDefaultTags(nil, defaultTags)

	if len(result) != 0 {
		t.Errorf("Expected 0 tags, got %d", len(result))
	}
}

// Note: SetTagsDiff is tested through acceptance tests in resource_group_resource_test.go
// because it requires a full pluginsdk.ResourceDiff implementation which is complex to mock.
// Unit tests for MergeDefaultTags (which SetTagsDiff calls) are provided above.

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
