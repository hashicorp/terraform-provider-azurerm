// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventhub_namespace_disaster_recovery_config

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/disasterrecoveryconfigs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// Default Authorization Rule/Policy created by Azure, used to populate the
// default connection strings and keys
var (
	eventHubNamespaceDefaultAuthorizationRule = "RootManageSharedAccessKey"
	eventHubNamespaceResourceName             = "azurerm_eventhub_namespace"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceEventHubNamespaceDisasterRecoveryConfigCreate,
		Read:   resourceEventHubNamespaceDisasterRecoveryConfigRead,
		Update: resourceEventHubNamespaceDisasterRecoveryConfigUpdate,
		Delete: resourceEventHubNamespaceDisasterRecoveryConfigDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := disasterrecoveryconfigs.ParseDisasterRecoveryConfigID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: eventhubNamespaceDisasterRecoveryConfigSchema(),
	}
}
