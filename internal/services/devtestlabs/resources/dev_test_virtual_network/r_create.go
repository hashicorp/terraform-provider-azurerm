// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_test_virtual_network

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/virtualnetworks"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	devtestlabsClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/devtestlabs/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmDevTestVirtualNetworkCreate(d *pluginsdk.ResourceData, meta any) error {
	client := devtestlabsClient.From(meta).VirtualNetworksClient
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForCreate(account.StopContext(meta), d)
	defer cancel()

	id := virtualnetworks.NewVirtualNetworkID(subscriptionId, d.Get("resource_group_name").(string), d.Get("lab_name").(string), d.Get("name").(string))

	if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id, virtualnetworks.GetOperationOptions{})
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_dev_test_virtual_network", id.ID())
		}
	}

	description := d.Get("description").(string)
	subnetsRaw := d.Get("subnet").([]any)
	subnets := expandDevTestVirtualNetworkSubnets(subnetsRaw, subscriptionId, id.ResourceGroupName, id.VirtualNetworkName)

	parameters := virtualnetworks.VirtualNetwork{
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
		Properties: virtualnetworks.VirtualNetworkProperties{
			Description:     pointer.To(description),
			SubnetOverrides: subnets,
		},
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceArmDevTestVirtualNetworkUpdate(d, meta)
}
