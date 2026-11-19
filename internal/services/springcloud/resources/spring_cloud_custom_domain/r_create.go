// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_custom_domain

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

func resourceSpringCloudCustomDomainCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppPlatform.CustomDomainsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	// the ID is parsed insensitively to preserve the behaviour of the legacy parser this replaced -
	// we are ok with this remaining insensitive as Azure Spring Apps is deprecated and will be removed
	appId, err := appplatform_rm.ParseAppIDInsensitively(d.Get("spring_cloud_app_id").(string))
	if err != nil {
		return err
	}

	resourceId := appplatform_rm.NewDomainID(appId.SubscriptionId, appId.ResourceGroupName, appId.SpringName, appId.AppName, name).ID()
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, appId.ResourceGroupName, appId.SpringName, appId.AppName, name)
			if err != nil {
				if !response.WasNotFound(existing.Response.Response) {
					return fmt.Errorf("making Read request on AzureRM Spring Cloud Custom Domain %q (Spring Cloud service %q / App %q / rcsource group %q): %+v", name, appId.SpringName, appId.AppName, appId.ResourceGroupName, err)
				}
			}
			if !response.WasNotFound(existing.Response.Response) {
				return tf.ImportAsExistsError("azurerm_spring_cloud_custom_domain", resourceId)
			}
		}
	}

	domain := appplatform.CustomDomainResource{
		Properties: &appplatform.CustomDomainProperties{
			Thumbprint: pointer.To(d.Get("thumbprint").(string)),
			CertName:   pointer.To(d.Get("certificate_name").(string)),
		},
	}

	future, err := client.CreateOrUpdate(ctx, appId.ResourceGroupName, appId.SpringName, appId.AppName, name, domain)
	if err != nil {
		return fmt.Errorf("creating/update Spring Cloud Custom Domain %q (Spring Cloud service %q / App %q / rcsource group %q): %+v", name, appId.SpringName, appId.AppName, appId.ResourceGroupName, err)
	}

	if d.IsNewResource() {
		d.SetId(resourceId)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %q(Spring Cloud service %q / App %q / rcsource group %q): %+v", name, appId.SpringName, appId.AppName, appId.ResourceGroupName, err)
	}
	return resourceSpringCloudCustomDomainRead(d, meta)
}
