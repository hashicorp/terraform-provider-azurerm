package azurerm

import (
	"context"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/go-azure-helpers/resourceproviders"
)

func TestAccAzureRMEnsureRequiredResourceProvidersAreRegistered(t *testing.T) {
	config := testGetAzureConfig(t)
	if config == nil {
		return
	}

	builder := armClientBuilder{
		authConfig:                  config,
		terraformVersion:            "0.0.0",
		partnerId:                   "",
		disableCorrelationRequestID: true,
		disableTerraformPartnerID:   false,
		// this test intentionally checks all the RP's are registered - so this is intentional
		skipProviderRegistration: true,
	}
	armClient, err := getArmClient(context.Background(), builder)
	if err != nil {
		t.Fatalf("Error building ARM Client: %+v", err)
	}

	client := armClient.Resource.ProvidersClient
	ctx := testAccProvider.StopContext()
	providerList, err := client.List(ctx, nil, "")
	if err != nil {
		t.Fatalf("Unable to list provider registration status, it is possible that this is due to invalid "+
			"credentials or the service principal does not have permission to use the Resource Manager API, Azure "+
			"error: %s", err)
	}

	availableResourceProviders := providerList.Values()
	requiredResourceProviders := requiredResourceProviders()
	err = ensureResourceProvidersAreRegistered(ctx, *client, availableResourceProviders, requiredResourceProviders)
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
