// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package relay_hybrid_connection_authorization_rule

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/hybridconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/relay/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceRelayHybridConnectionAuthorizationRuleCreateUpdate,
		Read:   resourceRelayHybridConnectionAuthorizationRuleRead,
		Update: resourceRelayHybridConnectionAuthorizationRuleCreateUpdate,
		Delete: resourceRelayHybridConnectionAuthorizationRuleDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := hybridconnections.ParseHybridConnectionAuthorizationRuleID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		// function takes a schema map and adds the authorization rule properties to it
		Schema: helpers.AuthorizationRuleSchemaFrom(map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"namespace_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"hybrid_connection_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": commonschema.ResourceGroupName(),
		}),

		CustomizeDiff: pluginsdk.CustomizeDiffShim(helpers.AuthorizationRuleCustomizeDiff),
	}
}
