// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package preflight_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachineextensions"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/preflight"
	"github.com/hashicorp/terraform-provider-azurerm/internal/preflight/testdata"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func TestNewValidationRequest_requestShape(t *testing.T) {
	cases := []struct {
		Name             string
		ID               resourceids.ResourceId
		TypeOverride     string
		ExpectedScope    string
		ExpectedProvider string
		ExpectedType     string
		ExpectedName     string
	}{
		{
			Name:             "top_level",
			ID:               pointer.To(commonids.NewAppServiceID("00000000-0000-0000-0000-000000000000", "rg", "site")),
			ExpectedScope:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg",
			ExpectedProvider: "Microsoft.Web",
			ExpectedType:     "sites",
			ExpectedName:     "site",
		},
		{
			Name:             "child",
			ID:               pointer.To(virtualmachineextensions.NewExtensionID("00000000-0000-0000-0000-000000000000", "rg", "vm", "ext")),
			ExpectedScope:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg",
			ExpectedProvider: "Microsoft.Compute",
			ExpectedType:     "virtualMachines/extensions",
			ExpectedName:     "vm/ext",
		},
		{
			Name:             "type_override",
			ID:               pointer.To(commonids.NewAppServiceID("00000000-0000-0000-0000-000000000000", "rg", "site")),
			TypeOverride:     "override",
			ExpectedScope:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg",
			ExpectedProvider: "Microsoft.Web",
			ExpectedType:     "override",
			ExpectedName:     "site",
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			var request preflight.ValidationRequest
			var err error
			if tc.TypeOverride != "" {
				request, err = preflight.NewValidationRequestWithTypeOverride(pointer.To("westeurope"), tc.ID, tc.TypeOverride, "2024-01-01", nil)
			} else {
				request, err = preflight.NewValidationRequest(pointer.To("westeurope"), tc.ID, "2024-01-01", nil)
			}
			if err != nil {
				t.Fatalf("building validation request: %s", err)
			}

			if request.Scope != tc.ExpectedScope {
				t.Errorf("expected scope %q, got %q", tc.ExpectedScope, request.Scope)
			}
			if request.Provider != tc.ExpectedProvider {
				t.Errorf("expected provider %q, got %q", tc.ExpectedProvider, request.Provider)
			}
			if request.Type != tc.ExpectedType {
				t.Errorf("expected type %q, got %q", tc.ExpectedType, request.Type)
			}
			if expected := tc.ExpectedProvider + "/" + tc.ExpectedType; request.Resource.Type != expected {
				t.Errorf("expected resource type %q, got %q", expected, request.Resource.Type)
			}
			if request.Resource.Name != tc.ExpectedName {
				t.Errorf("expected resource name %q, got %q", tc.ExpectedName, request.Resource.Name)
			}
		})
	}
}

func TestAccPreflight_validateResource(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skipf("Acceptance tests skipped unless env 'TF_ACC' set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	p := provider.AzureProvider()
	d := p.Configure(ctx, terraform.NewResourceConfigRaw(nil))
	if d.HasError() {
		t.Fatalf("diags had error!")
	}
	client := p.Meta().(*clients.Client)

	cases := []struct {
		Name        string
		Location    string
		ID          resourceids.ResourceId
		APIVersion  string
		Properties  any
		ExpectError bool
	}{
		{
			Name:        "valid_webapp",
			Location:    "westeurope",
			ID:          pointer.To(commonids.NewAppServiceID(client.Account.SubscriptionId, "testResourceGroup", "testSiteName")),
			APIVersion:  "2024-11-01",
			Properties:  testdata.WebAppExample(client.Account.SubscriptionId),
			ExpectError: false,
		},
		{
			Name:        "invalid_vnet_missing_required_property",
			Location:    "westeurope",
			ID:          pointer.To(commonids.NewVirtualNetworkID(client.Account.SubscriptionId, "testResourceGroup", "testVirtualValid")),
			APIVersion:  "2024-01-01",
			Properties:  testdata.VirtualNetworkExampleMissingRequiredProperty(client.Account.SubscriptionId),
			ExpectError: true,
		},
	}

	metadata := sdk.ResourceMetaData{
		Client: client,
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			request, err := preflight.NewValidationRequest(pointer.To(tc.Location), tc.ID, tc.APIVersion, tc.Properties)
			if err != nil {
				t.Fatalf("building validation request: %s", err)
			}

			err = request.ValidateResource(ctx, metadata)
			if err != nil && !tc.ExpectError {
				t.Fatalf("expected no error, got: %s", err)
			}
			if err == nil && tc.ExpectError {
				t.Fatalf("expected an error but got none")
			}
		})
	}
}
