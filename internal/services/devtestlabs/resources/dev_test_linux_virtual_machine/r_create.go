// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_test_linux_virtual_machine

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devtestlabs/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmDevTestLinuxVirtualMachineCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DevTestLabs.VirtualMachinesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := virtualmachines.NewVirtualMachineID(subscriptionId, d.Get("resource_group_name").(string), d.Get("lab_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id, virtualmachines.GetOperationOptions{})
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %s", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_dev_test_linux_virtual_machine", id.ID())
			}
		}
	}

	allowClaim := d.Get("allow_claim").(bool)
	disallowPublicIPAddress := d.Get("disallow_public_ip_address").(bool)
	labSubnetName := d.Get("lab_subnet_name").(string)
	labVirtualNetworkId := d.Get("lab_virtual_network_id").(string)
	location := location.Normalize(d.Get("location").(string))
	notes := d.Get("notes").(string)
	password := d.Get("password").(string)
	sshKey := d.Get("ssh_key").(string)
	size := d.Get("size").(string)
	storageType := d.Get("storage_type").(string)
	username := d.Get("username").(string)

	galleryImageReferenceRaw := d.Get("gallery_image_reference").([]any)
	galleryImageReference := helpers.ExpandDevTestLabVirtualMachineGalleryImageReference(galleryImageReferenceRaw, "Linux")

	natRulesRaw := d.Get("inbound_nat_rule").(*pluginsdk.Set)
	natRules := helpers.ExpandDevTestLabVirtualMachineNatRules(natRulesRaw)

	if len(natRules) > 0 && !disallowPublicIPAddress {
		return fmt.Errorf("if `inbound_nat_rule` is specified then `disallow_public_ip_address` must be set to true")
	}

	nic := virtualmachines.NetworkInterfaceProperties{}
	if disallowPublicIPAddress {
		nic.SharedPublicIPAddressConfiguration = &virtualmachines.SharedPublicIPAddressConfiguration{
			InboundNatRules: &natRules,
		}
	}

	authenticateViaSsh := sshKey != ""
	parameters := virtualmachines.LabVirtualMachine{
		Location: pointer.To(location),
		Properties: virtualmachines.LabVirtualMachineProperties{
			AllowClaim:                 pointer.To(allowClaim),
			IsAuthenticationWithSshKey: pointer.To(authenticateViaSsh),
			DisallowPublicIPAddress:    pointer.To(disallowPublicIPAddress),
			GalleryImageReference:      galleryImageReference,
			LabSubnetName:              pointer.To(labSubnetName),
			LabVirtualNetworkId:        pointer.To(labVirtualNetworkId),
			NetworkInterface:           &nic,
			OsType:                     pointer.To("Linux"),
			Notes:                      pointer.To(notes),
			Password:                   pointer.To(password),
			Size:                       pointer.To(size),
			SshKey:                     pointer.To(sshKey),
			StorageType:                pointer.To(storageType),
			UserName:                   pointer.To(username),
		},
		Tags: helpers.ExpandTags(d.Get("tags").(map[string]any)),
	}

	if d.IsNewResource() {
		if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
		d.SetId(id.ID())
	} else {
		if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	return resourceArmDevTestLinuxVirtualMachineRead(d, meta)
}
