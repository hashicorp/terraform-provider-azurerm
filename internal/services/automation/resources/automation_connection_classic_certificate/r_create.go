// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_connection_classic_certificate

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/connection"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAutomationConnectionClassicCertificateCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Automation.Connection
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := connection.NewConnectionID(subscriptionId, d.Get("resource_group_name").(string), d.Get("automation_account_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %s", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_automation_connection_classic_certificate", id.ID())
	}

	parameters := connection.ConnectionCreateOrUpdateParameters{
		Name: id.ConnectionName,
		Properties: connection.ConnectionCreateOrUpdateProperties{
			Description: pointer.To(d.Get("description").(string)),
			ConnectionType: connection.ConnectionTypeAssociationProperty{
				Name: pointer.To("AzureClassicCertificate"),
			},
			FieldDefinitionValues: &map[string]string{
				"SubscriptionName":     d.Get("subscription_name").(string),
				"SubscriptionId":       d.Get("subscription_id").(string),
				"CertificateAssetName": d.Get("certificate_asset_name").(string),
			},
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceAutomationConnectionClassicCertificateRead(d, meta)
}
