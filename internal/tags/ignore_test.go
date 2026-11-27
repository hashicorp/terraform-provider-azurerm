// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package tags

import (
	"testing"

	rmtags "github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
)

func TestExpandHonoursIgnoreConfig(t *testing.T) {
	t.Cleanup(func() { rmtags.SetIgnore(nil) })
	rmtags.SetIgnore(&rmtags.IgnoreConfig{Keys: []string{"createdBy"}, KeyPrefixes: []string{"azure-policy-"}})

	result := Expand(map[string]interface{}{
		"environment":     "prod",
		"createdBy":       "policy",
		"azure-policy-id": "abc",
	})

	if len(result) != 1 {
		t.Fatalf("expected ignored keys to be excluded on expand, got %v", result)
	}
	if result["environment"] == nil || *result["environment"] != "prod" {
		t.Fatalf("expected environment=prod to be retained, got %v", result["environment"])
	}
}

func TestExpandWithoutIgnoreConfigIsUnchanged(t *testing.T) {
	t.Cleanup(func() { rmtags.SetIgnore(nil) })
	rmtags.SetIgnore(nil)

	result := Expand(map[string]interface{}{
		"environment": "prod",
		"createdBy":   "policy",
	})

	if len(result) != 2 {
		t.Fatalf("expected no ignore filtering when unset, got %v", result)
	}
}

func TestFlattenHonoursIgnoreConfig(t *testing.T) {
	t.Cleanup(func() { rmtags.SetIgnore(nil) })
	rmtags.SetIgnore(&rmtags.IgnoreConfig{KeyPrefixes: []string{"internal:"}})

	prod := "prod"
	owner := "x"
	result := Flatten(map[string]*string{
		"environment":    &prod,
		"internal:owner": &owner,
	})

	if len(result) != 1 {
		t.Fatalf("expected ignored prefix key to be scrubbed on flatten, got %v", result)
	}
	if result["environment"] != "prod" {
		t.Fatalf("expected environment=prod, got %v", result)
	}
}
