package client

import (
	paloalto_2022_08_29 "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	*paloalto_2022_08_29.Client
}

func NewClient(_ *common.ClientOptions) (*Client, error) {

	return &Client{}, nil
}
