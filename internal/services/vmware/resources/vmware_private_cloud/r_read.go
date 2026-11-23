// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vmware_private_cloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/privateclouds"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVmwarePrivateCloudRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Vmware.PrivateCloudClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := privateclouds.ParsePrivateCloudID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	d.Set("name", id.PrivateCloudName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		props := model.Properties

		if err := d.Set("management_cluster", flattenPrivateCloudManagementCluster(props.ManagementCluster)); err != nil {
			return fmt.Errorf("setting `management_cluster`: %+v", err)
		}
		d.Set("network_subnet_cidr", props.NetworkBlock)
		if err := d.Set("circuit", flattenPrivateCloudCircuit(props.Circuit)); err != nil {
			return fmt.Errorf("setting `circuit`: %+v", err)
		}

		internetConnectionEnabled := false
		if props.Internet != nil {
			internetConnectionEnabled = *props.Internet == privateclouds.InternetEnumEnabled
		}
		d.Set("internet_connection_enabled", internetConnectionEnabled)

		d.Set("hcx_cloud_manager_endpoint", props.Endpoints.HcxCloudManager)
		d.Set("nsxt_manager_endpoint", props.Endpoints.NsxtManager)
		d.Set("vcsa_endpoint", props.Endpoints.Vcsa)
		d.Set("management_subnet_cidr", props.ManagementNetwork)
		d.Set("nsxt_certificate_thumbprint", props.NsxtCertificateThumbprint)
		d.Set("provisioning_subnet_cidr", props.ProvisioningNetwork)
		d.Set("vcenter_certificate_thumbprint", props.VcenterCertificateThumbprint)
		d.Set("vmotion_subnet_cidr", props.VMotionNetwork)

		d.Set("sku_name", model.Sku.Name)

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func privateCloudStateRefreshFunc(ctx context.Context, client *privateclouds.PrivateCloudsClient, id privateclouds.PrivateCloudId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		res, err := client.Get(ctx, id)
		if err != nil {
			return nil, "", fmt.Errorf("polling for status of vmware private cloud %s error: %+v", id, err)
		}

		if model := res.Model; model != nil {
			if model.Properties.ProvisioningState != nil {
				return res, string(*res.Model.Properties.ProvisioningState), nil
			}
		}
		return nil, "", fmt.Errorf("unable to read the provisioning state")
	}
}
