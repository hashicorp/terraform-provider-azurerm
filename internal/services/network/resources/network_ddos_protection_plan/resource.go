// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_ddos_protection_plan

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/ddosprotectionplans"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -test-name basicConfigIdentity

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNetworkDDoSProtectionPlanCreate,
		Read:   resourceNetworkDDoSProtectionPlanRead,
		Update: resourceNetworkDDoSProtectionPlanUpdate,
		Delete: resourceNetworkDDoSProtectionPlanDelete,

		Importer: pluginsdk.ImporterValidatingIdentity(&ddosprotectionplans.DdosProtectionPlanId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&ddosprotectionplans.DdosProtectionPlanId{}),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: networkDdosProtectionPlanSchema(),
	}
}
