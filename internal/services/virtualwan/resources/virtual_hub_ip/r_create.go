// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_hub_ip

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	virtualwanClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualHubIPCreate(d *pluginsdk.ResourceData, meta any) error {
	client := virtualwanClient.From(meta).VirtualWANs
	ctx, cancel := timeouts.ForCreate(account.StopContext(meta), d)
	defer cancel()

	virtualHubId, err := virtualwans.ParseVirtualHubID(d.Get("virtual_hub_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(virtualHubId.VirtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(virtualHubId.VirtualHubName, virtualHubResourceName)

	id := commonids.NewVirtualHubIPConfigurationID(virtualHubId.SubscriptionId, virtualHubId.ResourceGroupName, virtualHubId.VirtualHubName, d.Get("name").(string))

	if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.VirtualHubIPConfigurationGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_virtual_hub_ip", id.ID())
		}
	}

	parameters := virtualwans.HubIPConfiguration{
		Name: pointer.To(id.IpConfigName),
		Properties: &virtualwans.HubIPConfigurationPropertiesFormat{
			Subnet: &virtualwans.Subnet{
				Id: pointer.To(d.Get("subnet_id").(string)),
			},
		},
	}

	if v, ok := d.GetOk("private_ip_address"); ok {
		parameters.Properties.PrivateIPAddress = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("private_ip_allocation_method"); ok {
		parameters.Properties.PrivateIPAllocationMethod = pointer.ToEnum[virtualwans.IPAllocationMethod](v.(string))
	}

	if v, ok := d.GetOk("public_ip_address_id"); ok {
		parameters.Properties.PublicIPAddress = &virtualwans.PublicIPAddress{
			Id: pointer.To(v.(string)),
		}
	}

	if err := client.VirtualHubIPConfigurationCreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceVirtualHubIPRead(d, meta)
}
