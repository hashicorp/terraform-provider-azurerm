package automation

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/automation/mgmt/2020-01-13-preview/automation"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2021-06-22/automationaccount"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/parse"
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
	ResourceGrup          string  `json:"resource_grup" tfschema:"resource_group_name"`
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
			client := meta.Client.Automation.ConnectionTypeClient
			connClient := meta.Client.Automation.AccountClient

			var model AutomationConnectionTypeModel
			if err := meta.Decode(&model); err != nil {
				return err
			}
			subscriptionID := meta.Client.Account.SubscriptionId

			accountID := automationaccount.NewAutomationAccountID(subscriptionID, model.ResourceGrup, model.AutomationAccountName)
			account, err := connClient.Get(ctx, accountID)
			if err != nil {
				return fmt.Errorf("retrieving automation account %q: %+v", accountID, err)
			}
			if response.WasNotFound(account.HttpResponse) {
				return fmt.Errorf("automation account %q was not found", accountID)
			}

			id := parse.NewConnectionTypeID(accountID.SubscriptionId, model.ResourceGrup, model.AutomationAccountName, model.Name)
			existing, err := client.Get(ctx, id.ResourceGroup, model.AutomationAccountName, model.Name)
			if !utils.ResponseWasNotFound(existing.Response) {
				if err != nil {
					return fmt.Errorf("retreiving %s: %v", id, err)
				}
				return meta.ResourceRequiresImport(m.ResourceType(), id)
			}
			param := automation.ConnectionTypeCreateOrUpdateParameters{
				Name: utils.String(model.Name),
				ConnectionTypeCreateOrUpdateProperties: &automation.ConnectionTypeCreateOrUpdateProperties{
					IsGlobal:         utils.Bool(model.IsGlobal),
					FieldDefinitions: map[string]*automation.FieldDefinition{},
				},
			}
			for _, field := range model.Field {
				param.ConnectionTypeCreateOrUpdateProperties.FieldDefinitions[field.Name] = &automation.FieldDefinition{
					IsEncrypted: utils.Bool(field.IsEncrypted),
					IsOptional:  utils.Bool(field.IsOptional),
					Type:        utils.String(field.Type),
				}
			}
			_, err = client.CreateOrUpdate(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name, param)
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
			id, err := parse.ConnectionTypeID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			client := meta.Client.Automation.ConnectionTypeClient
			result, err := client.Get(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name)
			if err != nil {
				return err
			}

			var output AutomationConnectionTypeModel
			output.IsGlobal = utils.NormaliseNilableBool(result.IsGlobal)
			output.Name = utils.NormalizeNilableString(result.Name)
			output.AutomationAccountName = id.AutomationAccountName
			output.ResourceGrup = id.ResourceGroup
			for name, prop := range result.FieldDefinitions {
				output.Field = append(output.Field, Field{
					Name:        name,
					Type:        utils.NormalizeNilableString(prop.Type),
					IsEncrypted: utils.NormaliseNilableBool(prop.IsEncrypted),
					IsOptional:  utils.NormaliseNilableBool(prop.IsOptional),
				})
			}

			return meta.Encode(&output)
		},
	}
}

func (m AutomationConnectionTypeResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := parse.ConnectionTypeID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			client := meta.Client.Automation.ConnectionTypeClient
			if _, err = client.Delete(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name); err != nil {
				return fmt.Errorf("deleting %s: %v", id, err)
			}
			return nil
		},
	}
}

func (m AutomationConnectionTypeResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ConnectionTypeID
}
