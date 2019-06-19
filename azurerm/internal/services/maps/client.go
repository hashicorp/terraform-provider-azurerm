package maps

import (
	"github.com/Azure/azure-sdk-for-go/services/maps/mgmt/2018-05-01/maps"
)

type Client struct {
	AccountsClient maps.AccountsClient
}
