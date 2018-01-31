package authentication

import (
	"strings"

	"fmt"

	"github.com/Azure/go-autorest/autorest/azure/cli"
)

type AzureCLIProfile struct {
	cli.Profile
}

func (a AzureCLIProfile) FindDefaultSubscriptionId() (string, error) {
	for _, subscription := range a.Subscriptions {
		if subscription.IsDefault {
			return subscription.ID, nil
		}
	}

	return "", fmt.Errorf("No Subscription was Marked as Default in the Azure Profile.")
}

func (a AzureCLIProfile) FindSubscription(subscriptionId string) (*cli.Subscription, error) {
	for _, subscription := range a.Subscriptions {
		if strings.EqualFold(subscription.ID, subscriptionId) {
			return &subscription, nil
		}
	}

	return nil, fmt.Errorf("Subscription %q was not found in your Azure CLI credentials. Please verify it exists in `az account list`.", subscriptionId)
}
