package helper

import (
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func ExpandEntityStatus(input interface{}) servicebus.EntityStatus {
	entityStatus := servicebus.Unknown

	if status, ok := input.(string); ok {
		entityStatus = servicebus.EntityStatus(strings.Title(status))
	}

	return entityStatus
}

func FlattenEntityStatus(status servicebus.EntityStatus) *string {
	return utils.String(string(status))
}
