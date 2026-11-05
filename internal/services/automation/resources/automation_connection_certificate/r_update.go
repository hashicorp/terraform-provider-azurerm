// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_connection_certificate

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/connection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAutomationConnectionCertificateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Automation.Connection
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := connection.ParseConnectionID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving existing %s: %+v", *id, err)
	}

	// Start from the existing field definition values so that fields managed
	// via ignore_changes retain their server-side value.
	fieldDefinitionValues := make(map[string]string)
	if existing.Model != nil && existing.Model.Properties != nil && existing.Model.Properties.FieldDefinitionValues != nil {
		fieldDefinitionValues = *existing.Model.Properties.FieldDefinitionValues
	}

	if d.HasChange("automation_certificate_name") {
		fieldDefinitionValues["AutomationCertificateName"] = d.Get("automation_certificate_name").(string)
	}

	if d.HasChange("subscription_id") {
		fieldDefinitionValues["SubscriptionID"] = d.Get("subscription_id").(string)
	}

	parameters := connection.ConnectionUpdateParameters{
		Name: &id.ConnectionName,
		Properties: &connection.ConnectionUpdateProperties{
			FieldDefinitionValues: &fieldDefinitionValues,
		},
	}

	if d.HasChange("description") {
		parameters.Properties.Description = pointer.To(d.Get("description").(string))
	}

	if _, err := client.Update(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceAutomationConnectionCertificateRead(d, meta)
}
