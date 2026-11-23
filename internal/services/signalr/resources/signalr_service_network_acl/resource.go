// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package signalr_service_network_acl

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2024-03-01/signalr"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/resources/signalr_service_network_acl/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var defaultRequestTypes = []signalr.SignalRRequestType{
	signalr.SignalRRequestTypeRESTAPI,
	signalr.SignalRRequestTypeClientConnection,
	signalr.SignalRRequestTypeServerConnection,
	signalr.SignalRRequestTypeTrace,
}

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSignalRServiceNetworkACLCreateUpdate,
		Read:   resourceSignalRServiceNetworkACLRead,
		Update: resourceSignalRServiceNetworkACLCreateUpdate,
		Delete: resourceSignalRServiceNetworkACLDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.NetworkAclV0ToV1{},
		}),
		SchemaVersion: 1,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := signalr.ParseSignalRID(id)
			return err
		}),

		Schema: signalrServiceNetworkAclSchema(),
	}
}
