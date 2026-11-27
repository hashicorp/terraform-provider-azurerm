// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_service_managed_certificate

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/certificates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	webClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/web/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAppServiceManagedCertificateRead(d *pluginsdk.ResourceData, meta any) error {
	client := webClient.From(meta).CertificatesClient

	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id, err := certificates.ParseCertificateID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if model := resp.Model; model != nil {
		d.Set("tags", tags.Flatten(model.Tags))
		if props := model.Properties; props != nil {
			d.Set("canonical_name", props.CanonicalName)
			d.Set("friendly_name", props.FriendlyName)
			d.Set("subject_name", props.SubjectName)
			d.Set("host_names", props.HostNames)
			d.Set("issuer", props.Issuer)
			d.Set("issue_date", props.IssueDate)
			d.Set("expiration_date", props.ExpirationDate)
			d.Set("thumbprint", props.Thumbprint)
		}
	}

	return nil
}
