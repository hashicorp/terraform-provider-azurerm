// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_interface_security_group_association

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkinterfaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networksecuritygroups"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetworkInterfaceSecurityGroupAssociationCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.NetworkInterfaces
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	nicId, err := commonids.ParseNetworkInterfaceID(d.Get("network_interface_id").(string))
	if err != nil {
		return err
	}

	locks.ByID(strings.ToLower(nicId.ID()))
	defer locks.UnlockByID(strings.ToLower(nicId.ID()))

	nsgId, err := networksecuritygroups.ParseNetworkSecurityGroupID(d.Get("network_security_group_id").(string))
	if err != nil {
		return err
	}

	locks.ByID(nsgId.ID())
	defer locks.UnlockByID(nsgId.ID())

	read, err := client.Get(ctx, *nicId, networkinterfaces.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(read.HttpResponse) {
			return fmt.Errorf("%s was not found", *nicId)
		}
		return fmt.Errorf("retrieving %s: %+v", *nicId, err)
	}

	if read.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", nicId)
	}
	if read.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", nicId)
	}

	id := commonids.NewCompositeResourceID(nicId, nsgId)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		if read.Model.Properties.NetworkSecurityGroup != nil {
			return tf.ImportAsExistsError("azurerm_network_interface_security_group_association", id.ID())
		}
	}

	read.Model.Properties.NetworkSecurityGroup = &networkinterfaces.NetworkSecurityGroup{
		Id: pointer.To(nsgId.ID()),
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, *nicId, *read.Model, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("updating Security Group Association for %s: %+v", *nicId, err)
	}

	d.SetId(id.ID())

	return resourceNetworkInterfaceSecurityGroupAssociationRead(d, meta)
}
