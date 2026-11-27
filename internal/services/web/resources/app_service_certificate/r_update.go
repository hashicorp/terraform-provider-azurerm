// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_service_certificate

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/certificates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	webClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/web/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAppServiceCertificateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := webClient.From(meta).CertificatesClient

	ctx, cancel := timeouts.ForUpdate(account.StopContext(meta), d)
	defer cancel()

	id := certificates.NewCertificateID(account.From(meta).SubscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %w", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	existing.Model.Tags = tags.Expand(d.Get("tags").(map[string]any))

	if _, err := client.CreateOrUpdate(ctx, id, *existing.Model); err != nil {
		return fmt.Errorf("updating %s: %s", id, err)
	}

	return resourceAppServiceCertificateRead(d, meta)
}
