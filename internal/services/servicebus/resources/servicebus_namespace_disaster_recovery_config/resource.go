// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_namespace_disaster_recovery_config

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/armdisasterrecoveries"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// Default Authorization Rule/Policy created by Azure, used to populate the
// default connection strings and keys
var (
	serviceBusNamespaceDefaultAuthorizationRule = "RootManageSharedAccessKey"
	serviceBusNamespaceResourceName             = "azurerm_servicebus_namespace"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Manages a Disaster Recovery Config for a Service Bus Namespace.",
		Create:      resourceServiceBusNamespaceDisasterRecoveryConfigCreate,
		Read:        resourceServiceBusNamespaceDisasterRecoveryConfigRead,
		Update:      resourceServiceBusNamespaceDisasterRecoveryConfigUpdate,
		Delete:      resourceServiceBusNamespaceDisasterRecoveryConfigDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := armdisasterrecoveries.ParseDisasterRecoveryConfigID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: servicebusNamespaceDisasterRecoveryConfigSchema(),
	}
}
