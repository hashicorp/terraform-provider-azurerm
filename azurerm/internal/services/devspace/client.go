package devspace

import "github.com/Azure/azure-sdk-for-go/services/preview/devspaces/mgmt/2018-06-01-preview/devspaces"

type Client struct {
	ControllersClient devspaces.ControllersClient
}
