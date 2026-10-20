// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_pool

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/capacitypools"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetAppPoolCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).NetApp.PoolClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := capacitypools.NewCapacityPoolID(subscriptionId, d.Get("resource_group_name").(string), d.Get("account_name").(string), d.Get("name").(string))
	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.PoolsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_netapp_pool", id.ID())
		}
	}

	sizeInTB := int64(d.Get("size_in_tb").(int))
	sizeInMB := sizeInTB * 1024 * 1024
	sizeInBytes := sizeInMB * 1024 * 1024

	encryptionType := capacitypools.EncryptionType(d.Get("encryption_type").(string))

	capacityPoolParameters := capacitypools.CapacityPool{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: capacitypools.PoolProperties{
			ServiceLevel:   capacitypools.ServiceLevel(d.Get("service_level").(string)),
			Size:           sizeInBytes,
			EncryptionType: &encryptionType,
			CoolAccess:     pointer.To(d.Get("cool_access_enabled").(bool)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if qosType, ok := d.GetOk("qos_type"); ok {
		capacityPoolParameters.Properties.QosType = pointer.ToEnum[capacitypools.QosType](qosType.(string))
	}

	if customThroughputMibps, ok := d.GetOk("custom_throughput_mibps"); ok {
		capacityPoolParameters.Properties.CustomThroughputMibps = pointer.To(int64(customThroughputMibps.(int)))
	}

	if err := client.PoolsCreateOrUpdateCallbackThenPoll(ctx, id, capacityPoolParameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceNetAppPoolRead(d, meta)
}
