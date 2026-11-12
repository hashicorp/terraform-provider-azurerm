// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_test_linux_virtual_machine

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devtestlabs/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmDevTestLinuxVirtualMachineRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DevTestLabs.VirtualMachinesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualmachines.ParseVirtualMachineID(d.Id())
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, *id, virtualmachines.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(read.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on %s: %+v", *id, err)
	}

	d.Set("name", id.VirtualMachineName)
	d.Set("lab_name", id.LabName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := read.Model; model != nil {
		if loc := model.Location; loc != nil {
			d.Set("location", location.Normalize(*loc))
		}

		props := model.Properties
		d.Set("allow_claim", props.AllowClaim)
		d.Set("disallow_public_ip_address", props.DisallowPublicIPAddress)
		d.Set("notes", props.Notes)
		d.Set("size", props.Size)
		d.Set("storage_type", props.StorageType)
		d.Set("username", props.UserName)

		if err := d.Set("gallery_image_reference", helpers.FlattenDevTestVirtualMachineGalleryImage(props.GalleryImageReference)); err != nil {
			return fmt.Errorf("setting `gallery_image_reference`: %+v", err)
		}

		// Computed fields
		d.Set("fqdn", props.Fqdn)
		d.Set("unique_identifier", props.UniqueIdentifier)

		if err = tags.FlattenAndSet(d, helpers.FlattenTags(model.Tags)); err != nil {
			return err
		}
	}
	return nil
}
