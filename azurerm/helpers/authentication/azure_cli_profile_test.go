package authentication

import (
	"testing"

	"github.com/Azure/go-autorest/autorest/azure/cli"
)

func TestAzureCLIProfileFindDefaultSubscription(t *testing.T) {
	cases := []struct {
		Description            string
		Subscriptions          []cli.Subscription
		ExpectedSubscriptionId string
		ExpectError            bool
	}{
		{
			Description:   "Empty Subscriptions",
			Subscriptions: []cli.Subscription{},
			ExpectError:   true,
		},
		{
			Description: "Single Subscription",
			Subscriptions: []cli.Subscription{
				{
					ID:        "7f68fe06-9404-4db8-a5c7-29639dc4b299",
					IsDefault: true,
				},
			},
			ExpectError:            false,
			ExpectedSubscriptionId: "7f68fe06-9404-4db8-a5c7-29639dc4b299",
		},
		{
			Description: "Multiple Subscriptions with First as the Default",
			Subscriptions: []cli.Subscription{
				{
					ID:        "7f68fe06-9404-4db8-a5c7-29639dc4b299",
					IsDefault: true,
				},
				{
					ID:        "f36508bb-53b9-4aad-a2ac-2df86acf0c31",
					IsDefault: false,
				},
			},
			ExpectError:            false,
			ExpectedSubscriptionId: "7f68fe06-9404-4db8-a5c7-29639dc4b299",
		},
		{
			Description: "Multiple Subscriptions with Second as the Default",
			Subscriptions: []cli.Subscription{
				{
					ID:        "7f68fe06-9404-4db8-a5c7-29639dc4b299",
					IsDefault: false,
				},
				{
					ID:        "f36508bb-53b9-4aad-a2ac-2df86acf0c31",
					IsDefault: true,
				},
			},
			ExpectError:            false,
			ExpectedSubscriptionId: "f36508bb-53b9-4aad-a2ac-2df86acf0c31",
		},
		{
			Description: "Multiple Subscriptions with None as the Default",
			Subscriptions: []cli.Subscription{
				{
					ID:        "7f68fe06-9404-4db8-a5c7-29639dc4b299",
					IsDefault: false,
				},
				{
					ID:        "f36508bb-53b9-4aad-a2ac-2df86acf0c31",
					IsDefault: false,
				},
			},
			ExpectError: true,
		},
	}

	for _, v := range cases {
		profile := azureCLIProfile{
			profile: cli.Profile{
				Subscriptions: v.Subscriptions,
			},
		}
		actualSubscriptionId, err := profile.findDefaultSubscriptionId()

		if v.ExpectError && err == nil {
			t.Fatalf("Expected an error for %q: didn't get one", v.Description)
		}

		if !v.ExpectError && err != nil {
			t.Fatalf("Expected there to be no error for %q - but got: %v", v.Description, err)
		}

		if actualSubscriptionId != v.ExpectedSubscriptionId {
			t.Fatalf("Expected Subscription ID to be %q - got %q", v.ExpectedSubscriptionId, actualSubscriptionId)
		}
	}
}

func TestAzureCLIProfileFindSubscription(t *testing.T) {
	cases := []struct {
		Description               string
		Subscriptions             []cli.Subscription
		SubscriptionIdToSearchFor string
		ExpectError               bool
	}{
		{
			Description:               "Empty Subscriptions",
			Subscriptions:             []cli.Subscription{},
			SubscriptionIdToSearchFor: "7f68fe06-9404-4db8-a5c7-29639dc4b299",
			ExpectError:               true,
		},
		{
			Description:               "Single Subscription",
			SubscriptionIdToSearchFor: "7f68fe06-9404-4db8-a5c7-29639dc4b299",
			Subscriptions: []cli.Subscription{
				{
					ID:        "7f68fe06-9404-4db8-a5c7-29639dc4b299",
					IsDefault: true,
				},
			},
			ExpectError: false,
		},
		{
			Description:               "Finding the default subscription",
			SubscriptionIdToSearchFor: "7f68fe06-9404-4db8-a5c7-29639dc4b299",
			Subscriptions: []cli.Subscription{
				{
					ID:        "7f68fe06-9404-4db8-a5c7-29639dc4b299",
					IsDefault: true,
				},
				{
					ID:        "f36508bb-53b9-4aad-a2ac-2df86acf0c31",
					IsDefault: false,
				},
			},
			ExpectError: false,
		},
		{
			Description:               "Finding a non default Subscription",
			SubscriptionIdToSearchFor: "7f68fe06-9404-4db8-a5c7-29639dc4b299",
			Subscriptions: []cli.Subscription{
				{
					ID:        "7f68fe06-9404-4db8-a5c7-29639dc4b299",
					IsDefault: false,
				},
				{
					ID:        "f36508bb-53b9-4aad-a2ac-2df86acf0c31",
					IsDefault: true,
				},
			},
			ExpectError: false,
		},
		{
			Description:               "Multiple Subscriptions with None as the Default",
			SubscriptionIdToSearchFor: "224f4ca6-117f-4928-bc0f-3df018feba3e",
			Subscriptions: []cli.Subscription{
				{
					ID:        "7f68fe06-9404-4db8-a5c7-29639dc4b299",
					IsDefault: false,
				},
				{
					ID:        "f36508bb-53b9-4aad-a2ac-2df86acf0c31",
					IsDefault: false,
				},
			},
			ExpectError: true,
		},
	}

	for _, v := range cases {
		profile := azureCLIProfile{
			profile: cli.Profile{
				Subscriptions: v.Subscriptions,
			},
		}

		subscription, err := profile.findSubscription(v.SubscriptionIdToSearchFor)

		if v.ExpectError && err == nil {
			t.Fatalf("Expected an error for %q: didn't get one", v.Description)
		}

		if !v.ExpectError && err != nil {
			t.Fatalf("Expected there to be no error for %q - but got: %v", v.Description, err)
		}

		if subscription != nil && subscription.ID != v.SubscriptionIdToSearchFor {
			t.Fatalf("Expected to find Subscription ID %q - got %q", subscription.ID, v.SubscriptionIdToSearchFor)
		}
	}
}
