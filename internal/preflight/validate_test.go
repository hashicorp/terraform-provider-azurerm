package preflight_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/preflight"
	"github.com/hashicorp/terraform-provider-azurerm/internal/preflight/testdata"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func TestValidateResource(t *testing.T) {
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
