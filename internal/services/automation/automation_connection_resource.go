// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/connection"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/connectiontype"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceAutomationConnection() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAutomationConnectionCreateUpdate,
		Read:   resourceAutomationConnectionRead,
		Update: resourceAutomationConnectionCreateUpdate,
		Delete: resourceAutomationConnectionDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := connection.ParseConnectionID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ConnectionName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"automation_account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AutomationAccount(),
			},

			"type": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"values": {
				Type:     pluginsdk.TypeMap,
				Required: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAutomationConnectionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.Connection
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	connectionTypeClient := meta.(*clients.Client).Automation.ConnectionType
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Automation Connection creation.")

	id := connection.NewConnectionID(subscriptionId, d.Get("resource_group_name").(string), d.Get("automation_account_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
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
	values := expandStringInterfaceMap(d.Get("values").(map[string]interface{}))

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
			Description: utils.String(d.Get("description").(string)),
			ConnectionType: connection.ConnectionTypeAssociationProperty{
				Name: utils.String(connectionTypeName),
			},
			FieldDefinitionValues: &values,
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceAutomationConnectionRead(d, meta)
}

func resourceAutomationConnectionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.Connection
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := connection.ParseConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("read request on %s: %+v", *id, err)
	}

	d.Set("name", id.ConnectionName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("automation_account_name", id.AutomationAccountName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if props.ConnectionType != nil {
				d.Set("type", props.ConnectionType.Name)
			}

			if props.FieldDefinitionValues != nil {
				if err := d.Set("values", props.FieldDefinitionValues); err != nil {
					return fmt.Errorf("setting `values`: %+v", err)
				}
			}
			d.Set("description", props.Description)
		}
	}

	return nil
}

func resourceAutomationConnectionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.Connection
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := connection.ParseConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil
		}

		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
