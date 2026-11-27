// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_app_mysql_association

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	appplatform_rm "github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func resourceSpringCloudAppMysqlAssociationCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppPlatform.BindingsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// the ID is parsed insensitively to preserve the behaviour of the legacy parser this replaced -
	// we are ok with this remaining insensitive as Azure Spring Apps is deprecated and will be removed
	appId, err := appplatform_rm.ParseAppIDInsensitively(d.Get("spring_cloud_app_id").(string))
	if err != nil {
		return err
	}

	id := appplatform_rm.NewBindingID(appId.SubscriptionId, appId.ResourceGroupName, appId.SpringName, appId.AppName, d.Get("name").(string))
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id.ResourceGroupName, id.SpringName, id.AppName, id.BindingName)
			if err != nil {
				if !response.WasNotFound(existing.Response.Response) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.Response.Response) {
				return tf.ImportAsExistsError("azurerm_spring_cloud_app_mysql_association", id.ID())
			}
		}
	}

	bindingResource := appplatform.BindingResource{
		Properties: &appplatform.BindingResourceProperties{
			BindingParameters: map[string]*string{
				springCloudAppMysqlAssociationKeyDatabase: pointer.To(d.Get("database_name").(string)),
				springCloudAppMysqlAssociationKeyUsername: pointer.To(d.Get("username").(string)),
			},
			Key:        pointer.To(d.Get("password").(string)),
			ResourceID: pointer.To(d.Get("mysql_server_id").(string)),
		},
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroupName, id.SpringName, id.AppName, id.BindingName, bindingResource)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %q: %+v", id, err)
	}

	return resourceSpringCloudAppMysqlAssociationRead(d, meta)
}
