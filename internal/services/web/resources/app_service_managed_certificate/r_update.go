// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_service_managed_certificate

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/certificates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	webClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/web/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAppServiceManagedCertificateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := webClient.From(meta).CertificatesClient

	ctx, cancel := timeouts.ForCreate(account.StopContext(meta), d)
	defer cancel()

	id, err := certificates.ParseCertificateID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Get(ctx, *id); err != nil {
		return fmt.Errorf("retrieving %s: %w", id, err)
	}

	payload := certificates.CertificatePatchResource{
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if _, err := client.Update(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating %s: %w", id, err)
	}

	return resourceAppServiceManagedCertificateRead(d, meta)
}
