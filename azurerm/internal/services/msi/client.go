package msi

import "github.com/Azure/azure-sdk-for-go/services/preview/msi/mgmt/2015-08-31-preview/msi"

type Client struct {
	UserAssignedIdentitiesClient msi.UserAssignedIdentitiesClient
}
