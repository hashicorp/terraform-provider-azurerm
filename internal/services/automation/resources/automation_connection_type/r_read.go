// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_connection_type

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/connectiontype"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (m Resource) Read() sdk.ResourceFunc {
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
					output.IsGlobal = pointer.From(props.IsGlobal)
					if props.FieldDefinitions != nil {
						for name, field := range *props.FieldDefinitions {
							output.Field = append(output.Field, Field{
								Name:        name,
								Type:        field.Type,
								IsEncrypted: pointer.From(field.IsEncrypted),
								IsOptional:  pointer.From(field.IsOptional),
							})
						}
					}
				}
			}
			return meta.Encode(&output)
		},
	}
}

type Field struct {
	Name        string `tfschema:"name"`
	IsOptional  bool   `tfschema:"is_optional"`
	IsEncrypted bool   `tfschema:"is_encrypted"`
	Type        string `tfschema:"type"`
}
