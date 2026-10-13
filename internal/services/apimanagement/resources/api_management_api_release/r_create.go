// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_api_release

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/api"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apirelease"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementApiReleaseCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).ApiManagement.ApiReleasesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	apiId, err := api.ParseApiID(d.Get("api_id").(string))
	if err != nil {
		return err
	}

	id := apirelease.NewReleaseID(subscriptionId, apiId.ResourceGroupName, apiId.ServiceName, apiId.ApiId, name)
	ifMatch := "*"

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_api_management_api_release", id.ID())
			}
		}
		ifMatch = ""
	}

	parameters := apirelease.ApiReleaseContract{
		Properties: &apirelease.ApiReleaseContractProperties{
			ApiId: pointer.To(d.Get("api_id").(string)),
			Notes: pointer.To(d.Get("notes").(string)),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters, apirelease.CreateOrUpdateOperationOptions{IfMatch: pointer.To(ifMatch)}); err != nil {
		return fmt.Errorf("creating/ updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceApiManagementApiReleaseRead(d, meta)
}
