// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package statecheck

import (
	"context"
	"strings"
	"testing"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestExpectStateContainsIdentityValueAtPath_CaseSensitivity(t *testing.T) {
	var version uint64 = 1
	req := statecheck.CheckStateRequest{
		State: &tfjson.State{
			Values: &tfjson.StateValues{
				RootModule: &tfjson.StateModule{
					Resources: []*tfjson.StateResource{
						{
							Address:               "azurerm_test.test",
							IdentitySchemaVersion: &version,
							IdentityValues: map[string]interface{}{
								"name": "acctest-MC12345",
							},
							AttributeValues: map[string]interface{}{
								"parent_id": "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Maintenance/maintenanceConfigurations/acctest-mc12345",
							},
						},
					},
				},
			},
		},
	}

	// Case-sensitive should fail
	sensitiveCheck := ExpectStateContainsIdentityValueAtPath(
		"azurerm_test.test",
		tfjsonpath.New("name"),
		tfjsonpath.New("parent_id"),
	)
	var sensitiveResp statecheck.CheckStateResponse
	sensitiveCheck.CheckState(context.Background(), req, &sensitiveResp)
	if sensitiveResp.Error == nil {
		t.Fatalf("expected case-sensitive check to fail, but it passed")
	}
	if !strings.Contains(sensitiveResp.Error.Error(), "expected state (azurerm_test.test.parent_id) to contain identity value (azurerm_test.test.name)") {
		t.Fatalf("unexpected error message: %v", sensitiveResp.Error)
	}

	// Case-insensitive should pass
	insensitiveCheck := ExpectStateContainsIdentityValueAtPathCaseInsensitive(
		"azurerm_test.test",
		tfjsonpath.New("name"),
		tfjsonpath.New("parent_id"),
	)
	var insensitiveResp statecheck.CheckStateResponse
	insensitiveCheck.CheckState(context.Background(), req, &insensitiveResp)
	if insensitiveResp.Error != nil {
		t.Fatalf("expected case-insensitive check to succeed, but got error: %v", insensitiveResp.Error)
	}
}
