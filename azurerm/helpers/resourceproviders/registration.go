package resourceproviders

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2017-05-10/resources"
)

func RegisterWithSubscription(ctx context.Context, providerName string, client resources.ProvidersClient) error {
	_, err := client.Register(ctx, providerName)
	if err != nil {
		return fmt.Errorf("Cannot register provider %s with Azure Resource Manager: %s.", providerName, err)
	}

	return nil
}
