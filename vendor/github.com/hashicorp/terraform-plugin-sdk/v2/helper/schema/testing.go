// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"context"

	testing "github.com/mitchellh/go-testing-interface"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TestResourceDataRaw creates a ResourceData from a raw configuration map.
func TestResourceDataRaw(t testing.T, schema map[string]*Schema, raw map[string]interface{}) *ResourceData {
	t.Helper()

	c := terraform.NewResourceConfigRaw(raw)

	sm := schemaMap(schema)
	diff, err := sm.Diff(context.Background(), nil, c, nil, nil, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	result, err := sm.Data(nil, diff)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	return result
}

// TestResourceDataWithIdentityRaw creates a ResourceData with an identity from a raw identity map.
func TestResourceDataWithIdentityRaw(t testing.T, schema map[string]*Schema, identitySchema map[string]*Schema, raw map[string]string) *ResourceData {
	t.Helper()

	sm := schemaMapWithIdentity{schema, identitySchema}
	state := terraform.InstanceState{
		Identity: raw,
	}

	result, err := sm.Data(&state, nil)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	return result
}
