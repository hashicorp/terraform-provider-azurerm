package client

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/iotoperations/armiotoperations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	InstanceClient *armiotoperations.InstanceClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	// Build a TokenCredential using azidentity (uses environment/default chain).
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("building Azure credential: %+v", err)
	}

	client, err := armiotoperations.NewInstanceClient(o.SubscriptionId, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("building IoTOperations InstanceClient: %+v", err)
	}

	return &Client{
		InstanceClient: client,
	}, nil
}
