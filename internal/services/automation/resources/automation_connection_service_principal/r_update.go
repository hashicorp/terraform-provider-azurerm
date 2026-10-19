// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_connection_service_principal

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/connection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAutomationConnectionServicePrincipalUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Automation.Connection
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := connection.ParseConnectionID(d.Id())
	if err != nil {
		return err
	}

	// The API represents application_id, certificate_thumbprint, subscription_id, and tenant_id
	// as a single FieldDefinitionValues map — the map is replaced as a whole on update. We GET the
	// existing values first, then overlay only the changed fields so that ignore_changes on
	// individual properties is respected.
	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving existing %s: %+v", *id, err)
	}

	fieldDefinitionValues := make(map[string]string)
	if existing.Model != nil && existing.Model.Properties != nil && existing.Model.Properties.FieldDefinitionValues != nil {
		fieldDefinitionValues = *existing.Model.Properties.FieldDefinitionValues
	}

	if d.HasChange("application_id") {
		fieldDefinitionValues["ApplicationId"] = d.Get("application_id").(string)
	}

	if d.HasChange("certificate_thumbprint") {
		fieldDefinitionValues["CertificateThumbprint"] = d.Get("certificate_thumbprint").(string)
	}

	if d.HasChange("subscription_id") {
		fieldDefinitionValues["SubscriptionId"] = d.Get("subscription_id").(string)
	}

	if d.HasChange("tenant_id") {
		fieldDefinitionValues["TenantId"] = d.Get("tenant_id").(string)
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

	return resourceAutomationConnectionServicePrincipalRead(d, meta)
}
