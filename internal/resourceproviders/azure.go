package resourceproviders

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/profiles/2017-03-09/resources/mgmt/resources"
)

func availableResourceProviders(ctx context.Context, client *resources.ProvidersClient) ([]resources.Provider, error) {
	providers, err := client.List(ctx, nil, "")
	if err != nil {
		return nil, fmt.Errorf("listing Resource Providers: %+v", err)
	}
	return providers.Values(), nil
}
