// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_next_generation_firewall_virtual_network_strata_cloud_manager

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Palo Alto Next Generation Firewall Virtual Network Strata Cloud Manager.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.NextGenerationFirewallName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"network_profile": schema.VnetNetworkProfileSchema(),

		"strata_cloud_manager_tenant_name": {
			Description:  "Strata Cloud Manager name which is intended to manage the policy for this firewall.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"dns_settings": schema.DNSSettingsSchema(),

		"destination_nat": schema.DestinationNATSchema(),

		"marketplace_offer_id": {
			Description:  "The marketplace offer ID.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      "pan_swfw_cloud_ngfw",
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"plan_id": {
			Description:  "The billing plan ID as published by Liftr.PAN.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "panw-cngfw-payg",
			ValidateFunc: validation.StringLenBetween(1, 50),
		},

		"identity": commonschema.UserAssignedIdentityOptional(),

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff

			if rd.Id() != "" && !rd.HasChange("strata_cloud_manager_tenant_name") {
				return nil
			}

			tenantName := rd.Get("strata_cloud_manager_tenant_name").(string)

			if err := validate.ValidateStrataCloudManagerTenantNameExists(ctx, metadata.Client, metadata.Client.Account.SubscriptionId, tenantName); err != nil {
				return fmt.Errorf("validating strata_cloud_manager_tenant_name: %+v", err)
			}

			return nil
		},
	}
}
