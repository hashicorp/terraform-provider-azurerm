package authentication

import (
	"testing"

	"github.com/Azure/go-autorest/autorest/azure/cli"
)

func TestAzureCliProfile_populateSubscriptionIdMissing(t *testing.T) {
	cliProfile := azureCLIProfile{
		profile: cli.Profile{
			Subscriptions: []cli.Subscription{},
		},
	}

	err := cliProfile.populateSubscriptionID()
	if err == nil {
		t.Fatalf("Expected an error to be returned - but didn't get one")
	}
}

func TestAzureCliProfile_populateSubscriptionIdNoDefault(t *testing.T) {
	cliProfile := azureCLIProfile{
		profile: cli.Profile{
			Subscriptions: []cli.Subscription{
				{
					IsDefault: false,
					ID:        "abc123",
				},
			},
		},
	}

	err := cliProfile.populateSubscriptionID()
	if err == nil {
		t.Fatalf("Expected an error to be returned - but didn't get one")
	}
}

func TestAzureCliProfile_populateSubscriptionIdValid(t *testing.T) {
	subscriptionId := "abc123"
	cliProfile := azureCLIProfile{
		profile: cli.Profile{
			Subscriptions: []cli.Subscription{
				{
					IsDefault: true,
					ID:        subscriptionId,
				},
			},
		},
	}

	err := cliProfile.populateSubscriptionID()
	if err != nil {
		t.Fatalf("Expected no error to be returned - but got: %+v", err)
	}

	if cliProfile.subscriptionId != subscriptionId {
		t.Fatalf("Expected the Subscription ID to be %q but got %q", subscriptionId, cliProfile.subscriptionId)
	}
}

func TestAzureCliProfile_populateTenantIdEmpty(t *testing.T) {
	cliProfile := azureCLIProfile{
		profile: cli.Profile{
			Subscriptions: []cli.Subscription{},
		},
	}

	err := cliProfile.populateEnvironment()
	if err == nil {
		t.Fatalf("Expected an error to be returned - but didn't get one")
	}
}

func TestAzureCliProfile_populateTenantIdMissingSubscription(t *testing.T) {
	cliProfile := azureCLIProfile{
		subscriptionId: "bcd234",
		profile: cli.Profile{
			Subscriptions: []cli.Subscription{
				{
					IsDefault: false,
					ID:        "abc123",
				},
			},
		},
	}

	err := cliProfile.populateTenantID()
	if err == nil {
		t.Fatalf("Expected an error to be returned - but didn't get one")
	}
}

func TestAzureCliProfile_populateTenantIdValid(t *testing.T) {
	cliProfile := azureCLIProfile{
		subscriptionId: "abc123",
		profile: cli.Profile{
			Subscriptions: []cli.Subscription{
				{
					IsDefault: false,
					ID:        "abc123",
					TenantID:  "bcd234",
				},
			},
		},
	}

	err := cliProfile.populateTenantID()
	if err != nil {
		t.Fatalf("Expected no error to be returned - but got: %+v", err)
	}

	if cliProfile.subscriptionId != "abc123" {
		t.Fatalf("Expected Subscription ID to be 'abc123' - got %q", cliProfile.subscriptionId)
	}

	if cliProfile.tenantId != "bcd234" {
		t.Fatalf("Expected Tenant ID to be 'bcd234' - got %q", cliProfile.tenantId)
	}
}
