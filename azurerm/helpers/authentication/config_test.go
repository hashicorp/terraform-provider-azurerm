package authentication

import (
	"testing"

	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure/cli"
)

func TestAzurePopulateSubscriptionFromCLIProfile_Missing(t *testing.T) {
	config := Config{}
	profiles := AzureCLIProfile{
		Profile: cli.Profile{
			Subscriptions: []cli.Subscription{},
		},
	}

	err := config.populateSubscriptionFromCLIProfile(profiles)
	if err == nil {
		t.Fatalf("Expected an error to be returned - but didn't get one")
	}
}

func TestAzurePopulateSubscriptionFromCLIProfile_NoDefault(t *testing.T) {
	config := Config{}
	profiles := AzureCLIProfile{
		Profile: cli.Profile{
			Subscriptions: []cli.Subscription{
				{
					IsDefault: false,
					ID:        "abc123",
				},
			},
		},
	}

	err := config.populateSubscriptionFromCLIProfile(profiles)
	if err == nil {
		t.Fatalf("Expected an error to be returned - but didn't get one")
	}
}

func TestAzurePopulateSubscriptionFromCLIProfile_Default(t *testing.T) {
	subscriptionId := "abc123"
	config := Config{}
	profiles := AzureCLIProfile{
		Profile: cli.Profile{
			Subscriptions: []cli.Subscription{
				{
					IsDefault: true,
					ID:        subscriptionId,
				},
			},
		},
	}

	err := config.populateSubscriptionFromCLIProfile(profiles)
	if err != nil {
		t.Fatalf("Expected no error to be returned - but got: %+v", err)
	}

	if config.SubscriptionID != subscriptionId {
		t.Fatalf("Expected the Subscription ID to be %q but got %q", subscriptionId, config.SubscriptionID)
	}
}

func TestAzurePopulateTenantFromCLIProfile_Empty(t *testing.T) {
	config := Config{}
	profiles := AzureCLIProfile{
		Profile: cli.Profile{
			Subscriptions: []cli.Subscription{},
		},
	}

	err := config.populateEnvironmentFromCLIProfile(profiles)
	if err == nil {
		t.Fatalf("Expected an error to be returned - but didn't get one")
	}
}

func TestAzurePopulateTenantFromCLIProfile_MissingSubscription(t *testing.T) {
	config := Config{
		SubscriptionID: "bcd234",
	}
	profiles := AzureCLIProfile{
		Profile: cli.Profile{
			Subscriptions: []cli.Subscription{
				{
					IsDefault: false,
					ID:        "abc123",
				},
			},
		},
	}

	err := config.populateTenantFromCLIProfile(profiles)
	if err == nil {
		t.Fatalf("Expected an error to be returned - but didn't get one")
	}
}

func TestAzurePopulateTenantFromCLIProfile_PopulateTenantId(t *testing.T) {
	config := Config{
		SubscriptionID: "abc123",
	}
	profiles := AzureCLIProfile{
		Profile: cli.Profile{
			Subscriptions: []cli.Subscription{
				{
					IsDefault: false,
					ID:        "abc123",
					TenantID:  "bcd234",
				},
			},
		},
	}

	err := config.populateTenantFromCLIProfile(profiles)
	if err != nil {
		t.Fatalf("Expected no error to be returned - but got: %+v", err)
	}

	if config.SubscriptionID != "abc123" {
		t.Fatalf("Expected Subscription ID to be 'abc123' - got %q", config.SubscriptionID)
	}

	if config.TenantID != "bcd234" {
		t.Fatalf("Expected Tenant ID to be 'bcd234' - got %q", config.TenantID)
	}
}

func TestAzurePopulateTenantFromCLIProfile_Complete(t *testing.T) {
	config := Config{
		SubscriptionID: "abc123",
		TenantID:       "bcd234",
	}
	profiles := AzureCLIProfile{
		Profile: cli.Profile{
			Subscriptions: []cli.Subscription{
				{
					IsDefault: false,
					ID:        "abc123",
				},
			},
		},
	}

	err := config.populateTenantFromCLIProfile(profiles)
	if err != nil {
		t.Fatalf("Expected no error to be returned - but got: %+v", err)
	}

	if config.SubscriptionID != "abc123" {
		t.Fatalf("Expected Subscription ID to be 'abc123' - got %q", config.SubscriptionID)
	}

	if config.TenantID != "bcd234" {
		t.Fatalf("Expected Tenant ID to be 'bcd234' - got %q", config.TenantID)
	}
}

func TestAzurePopulateFromAccessToken_Missing(t *testing.T) {
	config := Config{}

	successful, err := config.populateFromAccessToken(nil)
	if err == nil {
		t.Fatalf("Expected an error but didn't get one")
	}

	if successful {
		t.Fatalf("Expected the population of a null token to be false, got true")
	}
}

func TestAzurePopulateFromAccessToken_Exists(t *testing.T) {
	config := Config{}

	token := AccessToken{
		AccessToken: &adal.Token{
			AccessToken: "abc123",
		},
		ClientID:     "bcd234",
		IsCloudShell: true,
	}

	successful, err := config.populateFromAccessToken(&token)
	if err != nil {
		t.Fatalf("Expected no error but got: %+v", err)
	}

	if !successful {
		t.Fatalf("Expected the population of an existing token to be successful, it wasn't")
	}

	if config.usingCloudShell != token.IsCloudShell {
		t.Fatalf("Expected `usingCloudShell` to be %t, got %t", token.IsCloudShell, config.usingCloudShell)
	}

	if config.ClientID != token.ClientID {
		t.Fatalf("Expected `ClientID` to be %q, got %q", token.ClientID, config.ClientID)
	}

	if config.accessToken != token.AccessToken {
		t.Fatalf("Expected `accessToken` to be %+v, got %+v", token.AccessToken, config.accessToken)
	}
}
