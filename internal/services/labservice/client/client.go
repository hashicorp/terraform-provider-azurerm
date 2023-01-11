package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/labservices/2022-08-01/lab"
	"github.com/hashicorp/go-azure-sdk/resource-manager/labservices/2022-08-01/labplan"
	"github.com/hashicorp/go-azure-sdk/resource-manager/labservices/2022-08-01/user"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	LabClient     *lab.LabClient
	LabPlanClient *labplan.LabPlanClient
	UserClient    *user.UserClient
}

func NewClient(o *common.ClientOptions) *Client {
	LabClient := lab.NewLabClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&LabClient.Client, o.ResourceManagerAuthorizer)

	LabPlanClient := labplan.NewLabPlanClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&LabPlanClient.Client, o.ResourceManagerAuthorizer)

	UserClient := user.NewUserClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&UserClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		LabClient:     &LabClient,
		LabPlanClient: &LabPlanClient,
		UserClient:    &UserClient,
	}
}
