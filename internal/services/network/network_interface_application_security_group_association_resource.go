// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/applicationsecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/networkinterfaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceNetworkInterfaceApplicationSecurityGroupAssociation() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNetworkInterfaceApplicationSecurityGroupAssociationCreate,
		Read:   resourceNetworkInterfaceApplicationSecurityGroupAssociationRead,
		Delete: resourceNetworkInterfaceApplicationSecurityGroupAssociationDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			splitId := strings.Split(id, "|")
			if _, err := commonids.ParseNetworkInterfaceID(splitId[0]); err != nil {
				return err
			}
			if _, err := applicationsecuritygroups.ParseApplicationSecurityGroupID(splitId[1]); err != nil {
				return err
			}
			return nil
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.NetworkInterfaceApplicationSecurityGroupAssociationV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"network_interface_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NetworkInterfaceID,
			},

			"application_security_group_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: applicationsecuritygroups.ValidateApplicationSecurityGroupID,
			},
		},
	}
}

func resourceNetworkInterfaceApplicationSecurityGroupAssociationCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NetworkInterfaces
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Network Interface <-> Application Security Group Association creation.")

	networkInterfaceId := d.Get("network_interface_id").(string)
	applicationSecurityGroupId := d.Get("application_security_group_id").(string)

	id, err := commonids.ParseNetworkInterfaceID(networkInterfaceId)
	if err != nil {
		return err
	}

	locks.ByName(id.NetworkInterfaceName, networkInterfaceResourceName)
	defer locks.UnlockByName(id.NetworkInterfaceName, networkInterfaceResourceName)

	read, err := client.Get(ctx, *id, networkinterfaces.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(read.HttpResponse) {
			log.Printf("[INFO] Network Interface %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if read.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *id)
	}

	props := read.Model.Properties
	if props == nil {
		return fmt.Errorf("Error: `properties` was nil for %s", *id)
	}
	if props.IPConfigurations == nil {
		return fmt.Errorf("Error: `properties.ipConfigurations` was nil for %s", *id)
	}

	info := parseFieldsFromNetworkInterface(*props)
	resourceId := fmt.Sprintf("%s|%s", networkInterfaceId, applicationSecurityGroupId)
	if utils.SliceContainsValue(info.applicationSecurityGroupIDs, applicationSecurityGroupId) {
		return tf.ImportAsExistsError("azurerm_network_interface_application_security_group_association", resourceId)
	}

	info.applicationSecurityGroupIDs = append(info.applicationSecurityGroupIDs, applicationSecurityGroupId)

	props.IPConfigurations = mapFieldsToNetworkInterface(props.IPConfigurations, info)

	err = client.CreateOrUpdateThenPoll(ctx, *id, *read.Model)
	if err != nil {
		return fmt.Errorf("updating Application Security Group Association for %s: %+v", *id, err)
	}

	d.SetId(resourceId)

	return resourceNetworkInterfaceApplicationSecurityGroupAssociationRead(d, meta)
}

func resourceNetworkInterfaceApplicationSecurityGroupAssociationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NetworkInterfaces
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	splitId := strings.Split(d.Id(), "|")
	if len(splitId) != 2 {
		return fmt.Errorf("Expected ID to be in the format {networkInterfaceId}|{applicationSecurityGroupId} but got %q", d.Id())
	}

	nicID, err := commonids.ParseNetworkInterfaceID(splitId[0])
	if err != nil {
		return err
	}

	applicationSecurityGroupId := splitId[1]

	read, err := client.Get(ctx, *nicID, networkinterfaces.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(read.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *nicID)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *nicID, err)
	}

	if model := read.Model; model != nil {
		nicProps := read.Model.Properties
		if nicProps == nil {
			return fmt.Errorf("Error: `properties` was nil for %s", *nicID)
		}

		info := parseFieldsFromNetworkInterface(*nicProps)
		exists := false
		for _, groupId := range info.applicationSecurityGroupIDs {
			if groupId == applicationSecurityGroupId {
				exists = true
			}
		}

		if !exists {
			log.Printf("[DEBUG] Association between %s and Application Security Group %q was not found - removing from state!", *nicID, applicationSecurityGroupId)
			d.SetId("")
			return nil
		}

		d.Set("application_security_group_id", applicationSecurityGroupId)
		d.Set("network_interface_id", model.Id)
	}

	return nil
}

func resourceNetworkInterfaceApplicationSecurityGroupAssociationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NetworkInterfaces
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	splitId := strings.Split(d.Id(), "|")
	if len(splitId) != 2 {
		return fmt.Errorf("Expected ID to be in the format {networkInterfaceId}|{applicationSecurityGroupId} but got %q", d.Id())
	}

	nicID, err := commonids.ParseNetworkInterfaceID(splitId[0])
	if err != nil {
		return err
	}

	applicationSecurityGroupId := splitId[1]

	locks.ByName(nicID.NetworkInterfaceName, networkInterfaceResourceName)
	defer locks.UnlockByName(nicID.NetworkInterfaceName, networkInterfaceResourceName)

	read, err := client.Get(ctx, *nicID, networkinterfaces.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(read.HttpResponse) {
			return fmt.Errorf("%s was not found!", *nicID)
		}

		return fmt.Errorf("retrieving  %s: %+v", *nicID, err)
	}

	if read.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *nicID)
	}

	props := read.Model.Properties
	if props == nil {
		return fmt.Errorf("Error: `properties` was nil for %s", *nicID)
	}

	if props.IPConfigurations == nil {
		return fmt.Errorf("Error: `properties.ipConfigurations` was nil for %s)", *nicID)
	}

	info := parseFieldsFromNetworkInterface(*props)

	applicationSecurityGroupIds := make([]string, 0)
	for _, v := range info.applicationSecurityGroupIDs {
		if v != applicationSecurityGroupId {
			applicationSecurityGroupIds = append(applicationSecurityGroupIds, v)
		}
	}
	info.applicationSecurityGroupIDs = applicationSecurityGroupIds
	props.IPConfigurations = mapFieldsToNetworkInterface(props.IPConfigurations, info)

	err = client.CreateOrUpdateThenPoll(ctx, *nicID, *read.Model)
	if err != nil {
		return fmt.Errorf("removing Application Security Group for %s: %+v", *nicID, err)
	}

	return nil
}
