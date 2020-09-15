package acceptance

import (
	"context"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/go-azure-helpers/resourceproviders"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	rmResourceProviders "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceproviders"
)

// since this depends on GetAuthConfig which lives in this package
// unfortunately this has to live in a different package to the other func

func TestAccEnsureRequiredResourceProvidersAreRegistered(t *testing.T) {
	config := GetAuthConfig(t)
	if config == nil {
		return
	}

	builder := clients.ClientBuilder{
		AuthConfig:                  config,
		TerraformVersion:            "0.0.0",
		PartnerId:                   "",
		DisableCorrelationRequestID: true,
		DisableTerraformPartnerID:   false,
		// this test intentionally checks all the RP's are registered - so this is intentional
		SkipProviderRegistration: true,
	}
	armClient, err := clients.Build(context.Background(), builder)
	if err != nil {
		t.Fatalf("Error building ARM Client: %+v", err)
	}

	client := armClient.Resource.ProvidersClient
	ctx := AzureProvider.StopContext()
	providerList, err := client.List(ctx, nil, "")
	if err != nil {
		t.Fatalf("Unable to list provider registration status, it is possible that this is due to invalid "+
			"credentials or the service principal does not have permission to use the Resource Manager API, Azure "+
			"error: %s", err)
	}

	availableResourceProviders := providerList.Values()
	requiredResourceProviders := rmResourceProviders.Required()
	err = rmResourceProviders.EnsureRegistered(ctx, *client, availableResourceProviders, requiredResourceProviders)
	if err != nil {
		t.Fatalf("Error registering Resource Providers: %+v", err)
	}

	// refresh the list now things have been re-registered
	providerList, err = client.List(ctx, nil, "")
	if err != nil {
		t.Fatalf("Unable to list provider registration status, it is possible that this is due to invalid "+
			"credentials or the service principal does not have permission to use the Resource Manager API, Azure "+
			"error: %s", err)
	}

	stillRequiringRegistration := resourceproviders.DetermineResourceProvidersRequiringRegistration(providerList.Values(), requiredResourceProviders)
	if len(stillRequiringRegistration) > 0 {
		t.Fatalf("'%d' Resource Providers are still Pending Registration: %s", len(stillRequiringRegistration), spew.Sprint(stillRequiringRegistration))
	}
}
