// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_certificate

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/certificate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAutomationCertificateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Automation.Certificate
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := certificate.ParseCertificateID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving existing %s: %+v", *id, err)
	}

	if existing.Model == nil || existing.Model.Properties == nil {
		return fmt.Errorf("retrieving existing %s: model or properties were nil", *id)
	}

	parameters := certificate.CertificateCreateOrUpdateParameters{
		Name: id.CertificateName,
		Properties: certificate.CertificateCreateOrUpdateProperties{
			Description:  existing.Model.Properties.Description,
			IsExportable: existing.Model.Properties.IsExportable,
			Base64Value:  d.Get("base64").(string), // value is sensitive so we must grab it with a d.get
		},
	}

	if d.HasChange("description") {
		parameters.Properties.Description = pointer.To(d.Get("description").(string))
	}

	if d.HasChange("exportable") {
		parameters.Properties.IsExportable = pointer.To(d.Get("exportable").(bool))
	}

	if _, err := client.CreateOrUpdate(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceAutomationCertificateRead(d, meta)
}
