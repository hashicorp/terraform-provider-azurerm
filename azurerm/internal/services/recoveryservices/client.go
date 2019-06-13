package recoveryservices

import (
	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2016-06-01/recoveryservices"
	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2017-07-01/backup"
)

type Client struct {
	ProtectedItemsClient     backup.ProtectedItemsGroupClient
	ProtectionPoliciesClient backup.ProtectionPoliciesClient
	VaultsClient             recoveryservices.VaultsClient
}
