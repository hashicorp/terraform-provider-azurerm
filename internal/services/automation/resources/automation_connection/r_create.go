// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_connection

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/connection"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/connectiontype"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAutomationConnectionCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Automation.Connection
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	connectionTypeClient := meta.(*clients.Client).Automation.ConnectionType
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := connection.NewConnectionID(subscriptionId, d.Get("resource_group_name").(string), d.Get("automation_account_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_automation_connection", id.ID())
		}
	}

	connectionTypeName := d.Get("type").(string)
	values := helpers.ExpandStringInterfaceMap(d.Get("values").(map[string]any))

	connectionTypeId := connectiontype.NewConnectionTypeID(subscriptionId, id.ResourceGroupName, id.AutomationAccountName, connectionTypeName)
	// check `type` exists and required fields are passed by users
	connectionType, err := connectionTypeClient.Get(ctx, connectionTypeId)
	if err != nil {
		return fmt.Errorf("retrieving %s: %s", id, err)
	}
	if connectionType.Model != nil && connectionType.Model.Properties != nil && connectionType.Model.Properties.FieldDefinitions != nil {
		var missingFields []string
		for key := range *connectionType.Model.Properties.FieldDefinitions {
			if _, ok := values[key]; !ok {
				missingFields = append(missingFields, key)
			}
		}
		if len(missingFields) > 0 {
			return fmt.Errorf("%q should be specified in `values` when type is %q for `azurerm_automation_connection`", strings.Join(missingFields, ", "), connectionTypeName)
		}
	}

	parameters := connection.ConnectionCreateOrUpdateParameters{
		Name: id.ConnectionName,
		Properties: connection.ConnectionCreateOrUpdateProperties{
			Description: pointer.To(d.Get("description").(string)),
			ConnectionType: connection.ConnectionTypeAssociationProperty{
				Name: pointer.To(connectionTypeName),
			},
			FieldDefinitionValues: &values,
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceAutomationConnectionRead(d, meta)
}
