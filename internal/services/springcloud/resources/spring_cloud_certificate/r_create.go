// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_certificate

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	appplatform_rm "github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func resourceSpringCloudCertificateCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppPlatform.CertificatesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("service_name").(string)

	resourceId := appplatform_rm.NewCertificateID(subscriptionId, resourceGroup, serviceName, name).ID()

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, resourceGroup, serviceName, name)
		if err != nil {
			if !response.WasNotFound(existing.Response.Response) {
				return fmt.Errorf("checking for presence of existing Spring Cloud Service Certificate %q (Spring Cloud Service %q / Resource Group %q): %+v", name, serviceName, resourceGroup, err)
			}
		}
		if !response.WasNotFound(existing.Response.Response) {
			return tf.ImportAsExistsError("azurerm_spring_cloud_certificate", resourceId)
		}
	}

	cert := appplatform.CertificateResource{}
	if value, ok := d.GetOk("key_vault_certificate_id"); ok {
		keyVaultCertificateId, err := keyvault.ParseNestedItemID(value.(string), keyvault.VersionTypeVersioned, keyvault.NestedItemTypeAny)
		if err != nil {
			return err
		}
		cert.Properties = &appplatform.KeyVaultCertificateProperties{
			VaultURI:          pointer.To(strings.TrimSuffix(keyVaultCertificateId.KeyVaultBaseURL, "/")),
			KeyVaultCertName:  &keyVaultCertificateId.Name,
			ExcludePrivateKey: pointer.To(d.Get("exclude_private_key").(bool)),
		}
	}
	if value, ok := d.GetOk("certificate_content"); ok {
		cert.Properties = &appplatform.ContentCertificateProperties{
			Content: pointer.To(value.(string)),
		}
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, name, cert)
	if err != nil {
		return fmt.Errorf("creating Spring Cloud Certificate %q (Spring Cloud Service %q / Resource Group %q): %+v", name, serviceName, resourceGroup, err)
	}

	d.SetId(resourceId)

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %q(Spring Cloud Service %q / Resource Group %q): %+v", name, serviceName, resourceGroup, err)
	}

	return resourceSpringCloudCertificateRead(d, meta)
}
