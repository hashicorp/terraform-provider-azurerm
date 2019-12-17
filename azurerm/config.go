package azurerm

import (
	"context"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

type ArmClient struct {
	// inherit the fields from the parent, so that we should be able to set/access these at either level
	clients.Client
}

func getArmClient(ctx context.Context, builder clients.ClientBuilder) (*ArmClient, error) {
	// NOTE: this only needs to exist until ArmClient is removed
	//		 at this point this method can disappear along
	//	  	 with the rest of this file

	innerClient, err := clients.Build(ctx, builder)
	if err != nil {
		return nil, err
	}

	client := ArmClient{
		Client: *innerClient,
	}
	return &client, nil
}
