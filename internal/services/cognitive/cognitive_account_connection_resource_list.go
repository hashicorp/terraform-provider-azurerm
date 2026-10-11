// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive

import (
	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/accountconnectionresource"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type cognitiveAccountConnectionListModel struct {
	CognitiveAccountId types.String `tfsdk:"cognitive_account_id"`
}

func cognitiveAccountConnectionListResourceConfigSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"cognitive_account_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: accountconnectionresource.ValidateAccountID,
					},
				},
			},
		},
	}
}
