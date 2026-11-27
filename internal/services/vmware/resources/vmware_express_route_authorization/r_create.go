// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vmware_express_route_authorization

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/authorizations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/privateclouds"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	vmwareClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/vmware/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVmwareExpressRouteAuthorizationCreate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := account.From(meta).SubscriptionId
	client := vmwareClient.From(meta).AuthorizationClient
	ctx, cancel := timeouts.ForCreate(account.StopContext(meta), d)
	defer cancel()

	name := d.Get("name").(string)
	privateCloudId, err := privateclouds.ParsePrivateCloudID(d.Get("private_cloud_id").(string))
	if err != nil {
		return err
	}

	id := authorizations.NewAuthorizationID(subscriptionId, privateCloudId.ResourceGroupName, privateCloudId.PrivateCloudName, name)

	if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_vmware_express_route_authorization", id.ID())
		}
	}

	props := authorizations.ExpressRouteAuthorization{}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, props, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %q: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceVmwareExpressRouteAuthorizationRead(d, meta)
}
