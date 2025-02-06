// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/runbook"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AutomationRunbookDataSource struct{}

type AutomationRunbookDataSourceModel struct {
	RunbookName           string            `tfschema:"name"`
	AutomationAccountName string            `tfschema:"automation_account_name"`
	ResourceGroupName     string            `tfschema:"resource_group_name"`
	Location              string            `tfschema:"location"`
	Description           string            `tfschema:"description"`
	LogProgress           bool              `tfschema:"log_progress"`
	LogVerbose            bool              `tfschema:"log_verbose"`
	RunbookType           string            `tfschema:"runbook_type"`
	LogActivityTrace      int64             `tfschema:"log_activity_trace_level"`
	Content               string            `tfschema:"content"`
	Tags                  map[string]string `tfschema:"tags "`
}

var _ sdk.DataSource = AutomationRunbookDataSource{}

func (d AutomationRunbookDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.RunbookName(),
		},

		"automation_account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.AutomationAccount(),
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (d AutomationRunbookDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"description": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"log_progress": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"log_verbose": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"runbook_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"log_activity_trace_level": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"content": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (d AutomationRunbookDataSource) ModelObject() interface{} {
	return &AutomationRunbookDataSourceModel{}
}

func (d AutomationRunbookDataSource) ResourceType() string {
	return "azurerm_automation_runbook"
}

func (d AutomationRunbookDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automation.Runbook
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state AutomationRunbookDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			id := runbook.NewRunbookID(subscriptionId, state.ResourceGroupName, state.AutomationAccountName, state.RunbookName)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			contentResp, err := client.GetContent(ctx, id)
			if err != nil {
				if !response.WasNotFound(contentResp.HttpResponse) {
					return fmt.Errorf("retrieving runbook content %s: %+v", id, err)
				}
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state.Location = location.NormalizeNilable(model.Location)

			if model.Properties.Description != nil {
				state.Description = pointer.From(model.Properties.Description)
			}

			if model.Properties.LogProgress != nil {
				state.LogProgress = pointer.From(model.Properties.LogProgress)
			}

			if model.Properties.RunbookType != nil {
				state.RunbookType = string(pointer.From(model.Properties.RunbookType))
			}

			if model.Properties.LogVerbose != nil {
				state.LogVerbose = pointer.From(model.Properties.LogVerbose)
			}

			if model.Properties.LogActivityTrace != nil {
				state.LogActivityTrace = pointer.From(model.Properties.LogActivityTrace)
			}

			if contentResp.Model != nil {
				state.Content = string(pointer.From(contentResp.Model))
			}

			state.Tags = pointer.From(model.Tags)

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
