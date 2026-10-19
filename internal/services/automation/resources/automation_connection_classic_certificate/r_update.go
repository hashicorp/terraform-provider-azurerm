// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_connection_classic_certificate

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/connection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAutomationConnectionClassicCertificateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Automation.Connection
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := connection.ParseConnectionID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("reading existing %s: %+v", *id, err)
	}

	if existing.Model == nil || existing.Model.Properties == nil {
		return fmt.Errorf("reading existing %s: model or properties was nil", *id)
	}

	props := existing.Model.Properties

	fieldDefinitionValues := map[string]string{}
	if props.FieldDefinitionValues != nil {
		fieldDefinitionValues = *props.FieldDefinitionValues
	}

	if d.HasChange("subscription_name") {
		fieldDefinitionValues["SubscriptionName"] = d.Get("subscription_name").(string)
	}

	if d.HasChange("subscription_id") {
		fieldDefinitionValues["SubscriptionId"] = d.Get("subscription_id").(string)
	}

	if d.HasChange("certificate_asset_name") {
		fieldDefinitionValues["CertificateAssetName"] = d.Get("certificate_asset_name").(string)
	}

	parameters := connection.ConnectionCreateOrUpdateParameters{
		Name: id.ConnectionName,
		Properties: connection.ConnectionCreateOrUpdateProperties{
			Description: props.Description,
			ConnectionType: connection.ConnectionTypeAssociationProperty{
				Name: pointer.To("AzureClassicCertificate"),
			},
			FieldDefinitionValues: &fieldDefinitionValues,
		},
	}

	if d.HasChange("description") {
		parameters.Properties.Description = pointer.To(d.Get("description").(string))
	}

	if _, err := client.CreateOrUpdate(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceAutomationConnectionClassicCertificateRead(d, meta)
}
