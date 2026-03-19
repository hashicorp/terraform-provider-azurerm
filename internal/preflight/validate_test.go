package preflight_test

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/preflight"
	preflightvalidation "github.com/hashicorp/terraform-provider-azurerm/internal/preflight/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/preflight/testdata"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

var testLocation = "westeurope"

func TestValidateResource(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	p := provider.AzureProvider()
	_ = p.Configure(ctx, terraform.NewResourceConfigRaw(nil))
	client := p.Meta().(*clients.Client)

	cases := []struct {
		Name        string
		Request     preflight.ValidationRequest
		ExpectError bool
	}{
		{
			Name: "valid_example",
			Request: preflight.ValidationRequest{
				Location:   &testLocation,
				Provider:   "Microsoft.Example",
				ResourceId: pointer.To(commonids.NewResourceGroupID(client.Account.SubscriptionId, "sampleResourceGroup")),
				Type:       "examples",
				Resource: preflightvalidation.ResourceValidationRequestResource{
					ApiVersion: "2026-01-14",
					Name:       "exampleResource1",
					Type:       "Microsoft.Example/examples",
					Properties: map[string]interface{}{
						"exampleProperty": "exampleValue1",
					},
				},
			},
			ExpectError: false,
		},
		{
			Name: "valid_webapp",
			Request: preflight.ValidationRequest{
				Location:   pointer.To("westeurope"),
				Provider:   "Microsoft.Web",
				ResourceId: pointer.To(commonids.NewAppServiceID(client.Account.SubscriptionId, "testResourceGroup", "testSiteName")),
				Type:       "sites",
				Resource: preflightvalidation.ResourceValidationRequestResource{
					ApiVersion: "2024-11-01",
					Name:       "testWebAppValid",
					Type:       "Microsoft.Web/sites",
					Properties: testdata.WebAppExample(client.Account.SubscriptionId),
				},
			},
			ExpectError: false,
		},
		{
			Name: "valid_vnet",
			Request: preflight.ValidationRequest{
				Location:   pointer.To("westeurope"),
				Provider:   "Microsoft.Network",
				ResourceId: pointer.To(commonids.NewVirtualNetworkID(client.Account.SubscriptionId, "testResourceGroup", "testSiteName")),
				Type:       "virtualNetworks",
				Resource: preflightvalidation.ResourceValidationRequestResource{
					ApiVersion: "2024-01-01",
					Name:       "testVirtualValid",
					Type:       "Microsoft.Network/virtualNetworks",
					Properties: testdata.VirtualNetworkExampleMissingRequiredProperty(client.Account.SubscriptionId),
				},
			},
			ExpectError: true,
		},
	}

	metadata := sdk.ResourceMetaData{
		Client: client,
	}

	for _, tc := range cases {
		err := tc.Request.ValidateResource(ctx, metadata)
		if err != nil && !tc.ExpectError {
			t.Fatalf("expected no error, got %s", err)
		}

		if err == nil && tc.ExpectError {
			t.Fatalf("expected error, but didn't get one: %s", tc.Name)
		}
	}
}
