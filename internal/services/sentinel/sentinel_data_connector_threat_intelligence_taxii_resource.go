// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sentinel

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	securityinsight "github.com/jackofallops/kermit/sdk/securityinsights/2022-10-01-preview/securityinsights"
)

type DataConnectorThreatIntelligenceTAXIIResource struct{}

var _ sdk.ResourceWithUpdate = DataConnectorThreatIntelligenceTAXIIResource{}

type DataConnectorThreatIntelligenceTAXIIModel struct {
	Name                    string `tfschema:"name"`
	LogAnalyticsWorkspaceId string `tfschema:"log_analytics_workspace_id"`
	DisplayName             string `tfschema:"display_name"`
	APIRootURL              string `tfschema:"api_root_url"`
	CollectionID            string `tfschema:"collection_id"`
	UserName                string `tfschema:"user_name"`
	Password                string `tfschema:"password"`
	PollingFrequency        string `tfschema:"polling_frequency"`
	LookbackDate            string `tfschema:"lookback_date"`
	TenantId                string `tfschema:"tenant_id"`
}

func (r DataConnectorThreatIntelligenceTAXIIResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"log_analytics_workspace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
		},
		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"api_root_url": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"collection_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"user_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"password": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"polling_frequency": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(securityinsight.PollingFrequencyOnceAnHour),
			ValidateFunc: validation.StringInSlice([]string{
				string(securityinsight.PollingFrequencyOnceAMinute),
				string(securityinsight.PollingFrequencyOnceAnHour),
				string(securityinsight.PollingFrequencyOnceADay),
			},
				false),
		},
		"lookback_date": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "1970-01-01T00:00:00Z",
			ValidateFunc: validation.IsRFC3339Time,
		},
		"tenant_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},
	}
}

func (r DataConnectorThreatIntelligenceTAXIIResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r DataConnectorThreatIntelligenceTAXIIResource) ResourceType() string {
	return "azurerm_sentinel_data_connector_threat_intelligence_taxii"
}

func (r DataConnectorThreatIntelligenceTAXIIResource) ModelObject() interface{} {
	return &DataConnectorThreatIntelligenceTAXIIModel{}
}

func (r DataConnectorThreatIntelligenceTAXIIResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.DataConnectorID
}

func (r DataConnectorThreatIntelligenceTAXIIResource) CustomImporter() sdk.ResourceRunFunc {
	return importDataConnectorTyped(securityinsight.DataConnectorKindThreatIntelligenceTaxii)
}

func (r DataConnectorThreatIntelligenceTAXIIResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.DataConnectorsClient
			wspClient := metadata.Client.LogAnalytics.WorkspaceClient

			var plan DataConnectorThreatIntelligenceTAXIIModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			workspaceId, err := workspaces.ParseWorkspaceID(plan.LogAnalyticsWorkspaceId)
			if err != nil {
				return err
			}

			wsp, err := wspClient.Get(ctx, *workspaceId)
			if err != nil {
				return fmt.Errorf("retrieving the workspace %s: %+v", workspaceId, err)
			}
			if wsp.Model == nil {
				return fmt.Errorf("nil model of the workspace %s", workspaceId)
			}
			if wsp.Model.Properties == nil {
				return fmt.Errorf("nil properties of the workspace %s", workspaceId)
			}
			if wsp.Model.Properties.CustomerId == nil {
				return fmt.Errorf("nil workspace id of the workspace %s", workspaceId)
			}
			wspId := *wsp.Model.Properties.CustomerId

			id := parse.NewDataConnectorID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, plan.Name)
			existing, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
			if err != nil {
				if !utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			tenantId := plan.TenantId
			if tenantId == "" {
				tenantId = metadata.Client.Account.TenantId
			}

			// Format is guaranteed by schema validation
			lookbackDate, _ := time.Parse(time.RFC3339, plan.LookbackDate)

			params := securityinsight.TiTaxiiDataConnector{
				Name: &plan.Name,
				TiTaxiiDataConnectorProperties: &securityinsight.TiTaxiiDataConnectorProperties{
					WorkspaceID:      &wspId,
					FriendlyName:     &plan.DisplayName,
					TaxiiServer:      &plan.APIRootURL,
					CollectionID:     &plan.CollectionID,
					PollingFrequency: securityinsight.PollingFrequency(plan.PollingFrequency),
					TaxiiLookbackPeriod: &date.Time{
						Time: lookbackDate,
					},
					DataTypes: &securityinsight.TiTaxiiDataConnectorDataTypes{
						TaxiiClient: &securityinsight.TiTaxiiDataConnectorDataTypesTaxiiClient{
							State: securityinsight.DataTypeStateEnabled,
						},
					},
					TenantID: &tenantId,
				},
				Kind: securityinsight.KindBasicDataConnectorKindThreatIntelligenceTaxii,
			}

			if plan.UserName != "" {
				params.TiTaxiiDataConnectorProperties.UserName = &plan.UserName
			}

			if plan.Password != "" {
				params.TiTaxiiDataConnectorProperties.Password = &plan.Password
			}

			if _, err = client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.Name, params); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r DataConnectorThreatIntelligenceTAXIIResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.DataConnectorsClient

			var state DataConnectorThreatIntelligenceTAXIIModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			id, err := parse.DataConnectorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			workspaceId := workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName)

			existing, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
			if err != nil {
				if utils.ResponseWasNotFound(existing.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			dc, ok := existing.Value.(securityinsight.TiTaxiiDataConnector)
			if !ok {
				return fmt.Errorf("%s was not an Threat Intelligence TAXII Data Connector", id)
			}

			model := DataConnectorThreatIntelligenceTAXIIModel{
				Name:                    id.Name,
				LogAnalyticsWorkspaceId: workspaceId.ID(),
				UserName:                state.UserName, // setting the user name from state, as it is not returned from API
				Password:                state.Password, // setting the password from state, as it is not returned from API
			}

			if props := dc.TiTaxiiDataConnectorProperties; props != nil {
				if props.FriendlyName != nil {
					model.DisplayName = *props.FriendlyName
				}

				if props.TaxiiServer != nil {
					model.APIRootURL = *props.TaxiiServer
				}

				if props.CollectionID != nil {
					model.CollectionID = *props.CollectionID
				}

				model.PollingFrequency = string(props.PollingFrequency)

				if props.TaxiiLookbackPeriod != nil {
					model.LookbackDate = props.TaxiiLookbackPeriod.Format(time.RFC3339)
				}

				if props.TenantID != nil {
					model.TenantId = *props.TenantID
				}
			}

			return metadata.Encode(&model)
		},
	}
}

func (r DataConnectorThreatIntelligenceTAXIIResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.DataConnectorsClient

			id, err := parse.DataConnectorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var plan DataConnectorThreatIntelligenceTAXIIModel
			if err := metadata.Decode(&plan); err != nil {
				return err
			}

			existing, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			dc, ok := existing.Value.(securityinsight.TiTaxiiDataConnector)
			if !ok {
				return fmt.Errorf("%s was not an Threat Intelligence TAXII Data Connector", id)
			}

			if props := dc.TiTaxiiDataConnectorProperties; props != nil {
				if metadata.ResourceData.HasChange("display_name") {
					props.FriendlyName = &plan.DisplayName
				}
				if metadata.ResourceData.HasChange("api_root_url") {
					props.TaxiiServer = &plan.APIRootURL
				}
				if metadata.ResourceData.HasChange("collection_id") {
					props.CollectionID = &plan.CollectionID
				}
				if metadata.ResourceData.HasChange("user_name") {
					props.UserName = &plan.UserName
				}
				if metadata.ResourceData.HasChange("password") {
					props.Password = &plan.Password
				}
				if metadata.ResourceData.HasChange("polling_frequency") {
					props.PollingFrequency = securityinsight.PollingFrequency(plan.PollingFrequency)
				}
				if metadata.ResourceData.HasChange("lookback_date") {
					// Format is guaranteed by schema validation
					lookbackDate, _ := time.Parse(time.RFC3339, plan.LookbackDate)
					props.TaxiiLookbackPeriod = &date.Time{
						Time: lookbackDate,
					}
				}

				// Setting the user name and password if non empty in plan, which are required by the API.
				if plan.UserName != "" {
					props.UserName = &plan.UserName
				}
				if plan.Password != "" {
					props.Password = &plan.Password
				}
			}

			if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.Name, dc); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r DataConnectorThreatIntelligenceTAXIIResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.DataConnectorsClient

			id, err := parse.DataConnectorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, id.ResourceGroup, id.WorkspaceName, id.Name); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
