// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourcegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/cognitiveservicesprojects"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type cognitiveAccountProjectConnectionListModel struct {
	CognitiveAccountName types.String `tfsdk:"cognitive_account_name"`
	ProjectName          types.String `tfsdk:"project_name"`
	ResourceGroupName    types.String `tfsdk:"resource_group_name"`
	SubscriptionId       types.String `tfsdk:"subscription_id"`
}

func cognitiveAccountProjectConnectionListResourceConfigSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"cognitive_account_name": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: validate.AccountName(),
					},
				},
			},

			"project_name": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: validate.AccountProjectName(),
					},
				},
			},

			"resource_group_name": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: resourcegroups.ValidateName,
					},
				},
			},

			"subscription_id": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: validation.IsUUID,
					},
				},
			},
		},
	}
}

func cognitiveAccountProjectConnectionListProjects(ctx context.Context, metadata sdk.ResourceMetadata, data cognitiveAccountProjectConnectionListModel) ([]cognitiveservicesprojects.Project, error) {
	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	switch {
	case !data.ProjectName.IsNull():
		if data.CognitiveAccountName.IsNull() {
			return nil, errors.New("`cognitive_account_name` is required when `project_name` is specified")
		}
		if data.ResourceGroupName.IsNull() {
			return nil, errors.New("`resource_group_name` is required when `project_name` is specified")
		}

		id := cognitiveservicesprojects.NewProjectID(subscriptionID, data.ResourceGroupName.ValueString(), data.CognitiveAccountName.ValueString(), data.ProjectName.ValueString())
		resp, err := metadata.Client.Cognitive.ProjectsClient.ProjectsGet(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("retrieving %s: %+v", id, err)
		}

		if resp.Model == nil {
			return []cognitiveservicesprojects.Project{}, nil
		}

		return []cognitiveservicesprojects.Project{*resp.Model}, nil

	case !data.CognitiveAccountName.IsNull():
		if data.ResourceGroupName.IsNull() {
			return nil, errors.New("`resource_group_name` is required when `cognitive_account_name` is specified")
		}

		accountId := cognitiveservicesprojects.NewAccountID(subscriptionID, data.ResourceGroupName.ValueString(), data.CognitiveAccountName.ValueString())
		resp, err := metadata.Client.Cognitive.ProjectsClient.ProjectsListComplete(ctx, accountId)
		if err != nil {
			return nil, fmt.Errorf("listing projects for %s: %+v", accountId, err)
		}

		return resp.Items, nil

	default:
		return nil, errors.New("`cognitive_account_name` and `resource_group_name` are required to list project connections")
	}
}
