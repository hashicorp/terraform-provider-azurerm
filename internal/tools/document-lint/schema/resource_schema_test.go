// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package schema_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/schema"
)

func TestResourceFile(t *testing.T) {
	p := automation.WatcherResource{}
	t.Log(schema.FileForResource(p.Read().Func))

	// inspect schema
	r := schema.NewResourceByTyped(p)
	if r.ResourceType != p.ResourceType() {
		t.Fatalf("resource type not equal: want: %s, got: %s", p.ResourceType(), r.ResourceType)
	}
}
