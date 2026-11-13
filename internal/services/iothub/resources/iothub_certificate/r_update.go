// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_certificate

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceIotHubCertificateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).IoTHub.IotHubCertificateClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewIotHubCertificateID(subscriptionId, d.Get("resource_group_name").(string), d.Get("iothub_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id.ResourceGroup, id.IotHubName, id.CertificateName)
	if err != nil {
		return fmt.Errorf("reading %s: %v", id, err)
	}

	etag := pointer.From(existing.Etag)

	if d.HasChange("is_verified") {
		existing.Properties.IsVerified = pointer.To(d.Get("is_verified").(bool))
	}

	if d.HasChange("certificate_content") {
		existing.Properties.Certificate = pointer.To(d.Get("certificate_content").(string))
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.IotHubName, id.CertificateName, existing, etag); err != nil {
		return fmt.Errorf("updating  %s: %+v", id, err)
	}

	return resourceIotHubCertificateRead(d, meta)
}
