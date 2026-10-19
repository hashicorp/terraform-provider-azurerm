// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventhub_authorization_rule

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/eventhubs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/resources/eventhub_authorization_rule/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var eventHubResourceName = "azurerm_eventhub"

var eventHubNamespaceResourceName = "azurerm_eventhub_namespace"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceEventHubAuthorizationRuleCreateUpdate,
		Read:   resourceEventHubAuthorizationRuleRead,
		Update: resourceEventHubAuthorizationRuleCreateUpdate,
		Delete: resourceEventHubAuthorizationRuleDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := eventhubs.ParseEventhubAuthorizationRuleID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.EventHubAuthorizationRuleV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: helpers.EventHubAuthorizationRuleSchemaFrom(map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ValidateEventHubAuthorizationRuleName(),
			},

			"namespace_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ValidateEventHubNamespaceName(),
			},

			"eventhub_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ValidateEventHubName(),
			},

			"resource_group_name": commonschema.ResourceGroupName(),
		}),

		CustomizeDiff: pluginsdk.CustomizeDiffShim(helpers.EventHubAuthorizationRuleCustomizeDiff),
	}
}
