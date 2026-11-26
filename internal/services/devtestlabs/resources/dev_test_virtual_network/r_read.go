// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_test_virtual_network

import (
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/virtualnetworks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmDevTestVirtualNetworkRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DevTestLabs.VirtualNetworksClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualnetworks.ParseVirtualNetworkID(d.Id())
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, *id, virtualnetworks.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(read.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on %s: %+v", *id, err)
	}

	d.Set("name", id.VirtualNetworkName)
	d.Set("lab_name", id.LabName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := read.Model; model != nil {
		props := model.Properties
		d.Set("description", props.Description)

		if err := d.Set("subnet", flattenDevTestVirtualNetworkSubnets(props.SubnetOverrides)); err != nil {
			return fmt.Errorf("setting `subnet`: %+v", err)
		}

		// Computed fields
		d.Set("unique_identifier", props.UniqueIdentifier)

		if err = tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}
	return nil
}

func ValidateDevTestVirtualNetworkName() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^[A-Za-z0-9_-]+$"),
		"Virtual Network Name can only include alphanumeric characters, underscores, hyphens.",
	)
}
