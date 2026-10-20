// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_pool

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/capacitypools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetAppPoolRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).NetApp.PoolClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := capacitypools.ParseCapacityPoolID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.PoolsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.CapacityPoolName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("account_name", id.NetAppAccountName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		poolProperties := model.Properties
		d.Set("service_level", poolProperties.ServiceLevel)

		sizeInBytes := poolProperties.Size
		sizeInMB := sizeInBytes / 1024 / 1024
		sizeInTB := sizeInMB / 1024 / 1024
		d.Set("size_in_tb", int(sizeInTB))
		qosType := ""
		if poolProperties.QosType != nil {
			qosType = string(*poolProperties.QosType)
		}
		d.Set("qos_type", qosType)
		d.Set("encryption_type", pointer.FromEnum(poolProperties.EncryptionType))
		d.Set("cool_access_enabled", pointer.From(poolProperties.CoolAccess))
		d.Set("custom_throughput_mibps", int(pointer.From(poolProperties.CustomThroughputMibps)))

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func netappPoolDeleteStateRefreshFunc(ctx context.Context, client *capacitypools.CapacityPoolsClient, id capacitypools.CapacityPoolId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		res, err := client.PoolsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(res.HttpResponse) {
				return nil, "", fmt.Errorf("retrieving %s: %+v", id, err)
			}
		}

		statusCode := "dropped connection"
		if res.HttpResponse != nil {
			statusCode = strconv.Itoa(res.HttpResponse.StatusCode)
		}
		return res, statusCode, nil
	}
}
