// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package automation_connection_type

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/connectiontype"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
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
					output.Global = pointer.From(props.IsGlobal)
					if !features.SixPointOh() {
						output.IsGlobal = output.Global
					}

					if props.FieldDefinitions != nil {
						for name, field := range *props.FieldDefinitions {
							f := Field{
								Name:      name,
								Type:      field.Type,
								Encrypted: pointer.From(field.IsEncrypted),
								Optional:  pointer.From(field.IsOptional),
							}

							if !features.SixPointOh() {
								f.IsEncrypted = f.Encrypted
								f.IsOptional = f.Optional
							}

							output.Field = append(output.Field, f)
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
	Optional    bool   `tfschema:"optional"`
	IsOptional  bool   `tfschema:"is_optional,removedInNextMajorVersion"`
	Encrypted   bool   `tfschema:"encrypted"`
	IsEncrypted bool   `tfschema:"is_encrypted,removedInNextMajorVersion"`
	Type        string `tfschema:"type"`
}
