// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_interface_application_security_group_association

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/applicationsecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkinterfaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetworkInterfaceApplicationSecurityGroupAssociationCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.NetworkInterfaces
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	applicationSecurityGroupId, err := applicationsecuritygroups.ParseApplicationSecurityGroupID(d.Get("application_security_group_id").(string))
	if err != nil {
		return err
	}
	networkInterfaceId, err := commonids.ParseNetworkInterfaceID(d.Get("network_interface_id").(string))
	if err != nil {
		return err
	}

	locks.ByID(strings.ToLower(networkInterfaceId.ID()))
	defer locks.UnlockByID(strings.ToLower(networkInterfaceId.ID()))

	read, err := client.Get(ctx, *networkInterfaceId, networkinterfaces.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(read.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", networkInterfaceId.ID())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *networkInterfaceId, err)
	}

	if read.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", networkInterfaceId)
	}

	props := read.Model.Properties
	if props == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", networkInterfaceId)
	}
	if props.IPConfigurations == nil {
		return fmt.Errorf("retrieving %s: `properties.ipConfigurations` was nil", networkInterfaceId)
	}

	info := helpers.ParseFieldsFromNetworkInterface(*props)
	id := commonids.NewCompositeResourceID(networkInterfaceId, applicationSecurityGroupId)

	exists := false
	if slices.Contains(info.ApplicationSecurityGroupIDs, applicationSecurityGroupId.ID()) {
		exists = true
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			return tf.ImportAsExistsError("azurerm_network_interface_application_security_group_association", id.ID())
		}
	}

	if !exists {
		info.ApplicationSecurityGroupIDs = append(info.ApplicationSecurityGroupIDs, applicationSecurityGroupId.ID())
	}

	props.IPConfigurations = helpers.MapFieldsToNetworkInterface(props.IPConfigurations, info)

	if err = client.CreateOrUpdateCallbackThenPoll(ctx, *networkInterfaceId, *read.Model, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("updating Application Security Group Association for %s: %+v", *networkInterfaceId, err)
	}

	d.SetId(id.ID())

	return resourceNetworkInterfaceApplicationSecurityGroupAssociationRead(d, meta)
}
