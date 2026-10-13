// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_api_release

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apirelease"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementApiReleaseDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.ApiReleasesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := apirelease.ParseReleaseID(d.Id())
	if err != nil {
		return err
	}

	name := helpers.GetApiName(id.ApiId)

	newId := apirelease.NewReleaseID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, name, id.ReleaseId)
	if _, err := client.Delete(ctx, newId, apirelease.DeleteOperationOptions{IfMatch: pointer.To("*")}); err != nil {
		return fmt.Errorf("deleting %s: %+v", newId, err)
	}
	return nil
}
