package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dynatrace/2023-04-27/monitors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	MonitorClient *monitors.MonitorsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	monitorClient, err := monitors.NewMonitorsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Dynatrace Monitor client: %+v", err)
	}
	o.Configure(monitorClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		MonitorClient: monitorClient,
	}, nil
}
