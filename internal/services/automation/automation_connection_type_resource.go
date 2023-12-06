// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/automationaccount"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/connectiontype"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type Field struct {
	Name        string `tfschema:"name"`
	IsOptional  bool   `tfschema:"is_optional"`
	IsEncrypted bool   `tfschema:"is_encrypted"`
	Type        string `tfschema:"type"`
}

type AutomationConnectionTypeModel struct {
	ResourceGroup         string  `json:"resource_group" tfschema:"resource_group_name"`
	AutomationAccountName string  `json:"automation_account_name" tfschema:"automation_account_name"`
	Name                  string  `json:"name" tfschema:"name"`
	IsGlobal              bool    `json:"is_global" tfschema:"is_global"`
	Field                 []Field `json:"field" tfschema:"field"`
}

type AutomationConnectionTypeResource struct{}

var _ sdk.Resource = (*AutomationConnectionTypeResource)(nil)

func (m AutomationConnectionTypeResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupName(),

		"automation_account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AutomationAccount(),
		},

		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ConnectionTypeName,
		},

		"is_global": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
		},

		"field": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"is_encrypted": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
					"is_optional": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
					"type": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},
	}
}

func (m AutomationConnectionTypeResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (m AutomationConnectionTypeResource) ModelObject() interface{} {
	return &AutomationConnectionTypeModel{}
}

func (m AutomationConnectionTypeResource) ResourceType() string {
	return "azurerm_automation_connection_type"
}

func (m AutomationConnectionTypeResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.Automation.ConnectionType
			connClient := meta.Client.Automation.AutomationAccount

			var model AutomationConnectionTypeModel
			if err := meta.Decode(&model); err != nil {
				return err
			}
			subscriptionID := meta.Client.Account.SubscriptionId

			accountID := automationaccount.NewAutomationAccountID(subscriptionID, model.ResourceGroup, model.AutomationAccountName)
			account, err := connClient.Get(ctx, accountID)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", accountID, err)
			}
			if response.WasNotFound(account.HttpResponse) {
				return fmt.Errorf("%s was not found", accountID)
			}

			id := connectiontype.NewConnectionTypeID(accountID.SubscriptionId, model.ResourceGroup, model.AutomationAccountName, model.Name)
			existing, err := client.Get(ctx, id)
			if !response.WasNotFound(existing.HttpResponse) {
				if err != nil {
					return fmt.Errorf("retrieving %s: %v", id, err)
				}
				return meta.ResourceRequiresImport(m.ResourceType(), id)
			}
			param := connectiontype.ConnectionTypeCreateOrUpdateParameters{
				Name: model.Name,
				Properties: connectiontype.ConnectionTypeCreateOrUpdateProperties{
					IsGlobal:         utils.Bool(model.IsGlobal),
					FieldDefinitions: map[string]connectiontype.FieldDefinition{},
				},
			}
			for _, field := range model.Field {
				param.Properties.FieldDefinitions[field.Name] = connectiontype.FieldDefinition{
					IsEncrypted: utils.Bool(field.IsEncrypted),
					IsOptional:  utils.Bool(field.IsOptional),
					Type:        field.Type,
				}
			}
			_, err = client.CreateOrUpdate(ctx, id, param)
			if err != nil {
				return fmt.Errorf("creating %s: %v", id, err)
			}

			meta.SetID(id)
			return nil
		},
	}
}

func (m AutomationConnectionTypeResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := connectiontype.ParseConnectionTypeID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			client := meta.Client.Automation.ConnectionType
			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return meta.MarkAsGone(id)
				}
				return err
			}

			var output AutomationConnectionTypeModel
			output.Name = id.ConnectionTypeName
			output.AutomationAccountName = id.AutomationAccountName
			output.ResourceGroup = id.ResourceGroupName

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					output.IsGlobal = utils.NormaliseNilableBool(props.IsGlobal)
					if props.FieldDefinitions != nil {
						for name, field := range *props.FieldDefinitions {
							output.Field = append(output.Field, Field{
								Name:        name,
								Type:        field.Type,
								IsEncrypted: utils.NormaliseNilableBool(field.IsEncrypted),
								IsOptional:  utils.NormaliseNilableBool(field.IsOptional),
							})
						}
					}
				}
			}
			return meta.Encode(&output)
		},
	}
}

func (m AutomationConnectionTypeResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := connectiontype.ParseConnectionTypeID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			client := meta.Client.Automation.ConnectionType
			if _, err = client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %v", *id, err)
			}
			return nil
		},
	}
}

func (m AutomationConnectionTypeResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return connectiontype.ValidateConnectionTypeID
}
