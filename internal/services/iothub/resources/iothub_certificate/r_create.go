// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_certificate

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	devices "github.com/jackofallops/kermit/sdk/iothub/2022-04-30-preview/iothub" // azignore:AZG010 - package name does not match its path
)

func resourceIotHubCertificateCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).IoTHub.IotHubCertificateClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewIotHubCertificateID(subscriptionId, d.Get("resource_group_name").(string), d.Get("iothub_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id.ResourceGroup, id.IotHubName, id.CertificateName)
		if err != nil {
			if !response.WasNotFound(existing.Response.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.Response.Response) {
			return tf.ImportAsExistsError("azurerm_iothub_certificate", id.ID())
		}
	}

	certificate := devices.CertificateDescription{
		Properties: &devices.CertificateProperties{
			IsVerified:  pointer.To(d.Get("is_verified").(bool)),
			Certificate: pointer.To(d.Get("certificate_content").(string)),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.IotHubName, id.CertificateName, certificate, ""); err != nil {
		return fmt.Errorf("creating  %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceIotHubCertificateRead(d, meta)
}
