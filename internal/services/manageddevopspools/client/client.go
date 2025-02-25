package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
}

func NewClient(o *common.ClientOptions) *Client {
	return &Client{}
}