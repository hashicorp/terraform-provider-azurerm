// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_namespace_authorization_rule

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/namespacesauthorizationrule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/resources/servicebus_namespace_authorization_rule/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceServiceBusNamespaceAuthorizationRuleCreateUpdate,
		Read:   resourceServiceBusNamespaceAuthorizationRuleRead,
		Update: resourceServiceBusNamespaceAuthorizationRuleCreateUpdate,
		Delete: resourceServiceBusNamespaceAuthorizationRuleDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := namespacesauthorizationrule.ParseAuthorizationRuleID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ServicebusNamespaceAuthRuleV0ToV1{},
		}),

		// function takes a schema map and adds the authorization rule properties to it
		Schema: resourceServiceBusNamespaceAuthorizationRuleSchema(),

		CustomizeDiff: pluginsdk.CustomizeDiffShim(helpers.AuthorizationRuleCustomizeDiff),
	}
}
