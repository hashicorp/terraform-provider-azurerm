// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package sentinel

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2023-09-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/dataconnectors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type DataConnectorMicrosoftThreatIntelligenceResource struct{}

type DataConnectorMicrosoftThreatIntelligenceModel struct {
	Name                                    string `tfschema:"name"`
	WorkspaceId                             string `tfschema:"log_analytics_workspace_id"`
	TenantId                                string `tfschema:"tenant_id"`
	MicrosoftEmergingThreatFeedLookBackDate string `tfschema:"microsoft_emerging_threat_feed_lookback_date"`
}

type DataConnectorMicrosoftThreatIntelligenceDataType struct {
	Enabled      bool   `tfschema:"enabled"`
	LookbackDate string `tfschema:"lookback_date"`
}

func (s DataConnectorMicrosoftThreatIntelligenceResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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

		"tenant_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		// lintignore: S013
		"microsoft_emerging_threat_feed_lookback_date": {
			Type:         pluginsdk.TypeString,
			ForceNew:     true,
			Required:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},
	}
}

func (s DataConnectorMicrosoftThreatIntelligenceResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (s DataConnectorMicrosoftThreatIntelligenceResource) ModelObject() interface{} {
	return &DataConnectorMicrosoftThreatIntelligenceModel{}
}

func (s DataConnectorMicrosoftThreatIntelligenceResource) ResourceType() string {
	return "azurerm_sentinel_data_connector_microsoft_threat_intelligence"
}

func (s DataConnectorMicrosoftThreatIntelligenceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.DataConnectorsClient
			var metaModel DataConnectorMicrosoftThreatIntelligenceModel
			if err := metadata.Decode(&metaModel); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			workSpaceId, err := workspaces.ParseWorkspaceID(metaModel.WorkspaceId)
			if err != nil {
				return fmt.Errorf("parsing workspace id %+v", err)
			}

			id := dataconnectors.NewDataConnectorID(workSpaceId.SubscriptionId, workSpaceId.ResourceGroupName, workSpaceId.WorkspaceName, metaModel.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(s.ResourceType(), id)
				}
			}

			tenantId := metaModel.TenantId
			if tenantId == "" {
				tenantId = metadata.Client.Account.TenantId
			}

			dataConnector := dataconnectors.MSTIDataConnector{
				Name: &id.DataConnectorId,
				Kind: dataconnectors.DataConnectorKindMicrosoftThreatIntelligence,
				Properties: &dataconnectors.MSTIDataConnectorProperties{
					DataTypes: dataconnectors.MSTIDataConnectorDataTypes{
						MicrosoftEmergingThreatFeed: expandSentinelDataConnectorMicrosoftThreatIntelligenceMicrosoftEmergingThreatFeed(metaModel),
					},
					TenantId: tenantId,
				},
			}
			if _, err = client.CreateOrUpdate(ctx, id, dataConnector); err != nil {
				return fmt.Errorf("creating %+v", err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (s DataConnectorMicrosoftThreatIntelligenceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.DataConnectorsClient
			id, err := dataconnectors.ParseDataConnectorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			dc, ok := existing.Model.(dataconnectors.MSTIDataConnector)
			if !ok {
				return fmt.Errorf("%s was not an Microsoft Threat Protection Data Connector", id)
			}

			state := DataConnectorMicrosoftThreatIntelligenceModel{
				Name:        id.DataConnectorId,
				WorkspaceId: workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName).ID(),
			}

			if props := dc.Properties; props != nil {
				state.TenantId = props.TenantId

				dt := props.DataTypes
				if metf := dt.MicrosoftEmergingThreatFeed; metf.State != nil {
					if strings.EqualFold(string(*metf.State), string(dataconnectors.DataTypeStateEnabled)) {
						state.MicrosoftEmergingThreatFeedLookBackDate, err = flattenSentinelDataConnectorMicrosoftThreatIntelligenceTime(metf.LookbackPeriod)
						if err != nil {
							return fmt.Errorf("flattening `microsoft_emerging_threat_feed`: %+v", err)
						}
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (s DataConnectorMicrosoftThreatIntelligenceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.DataConnectorsClient

			id, err := dataconnectors.ParseDataConnectorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (s DataConnectorMicrosoftThreatIntelligenceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return dataconnectors.ValidateDataConnectorID
}

func expandSentinelDataConnectorMicrosoftThreatIntelligenceMicrosoftEmergingThreatFeed(input DataConnectorMicrosoftThreatIntelligenceModel) dataconnectors.MSTIDataConnectorDataTypesMicrosoftEmergingThreatFeed {
	if input.MicrosoftEmergingThreatFeedLookBackDate == "" {
		return dataconnectors.MSTIDataConnectorDataTypesMicrosoftEmergingThreatFeed{
			LookbackPeriod: "",
			State:          pointer.To(dataconnectors.DataTypeStateDisabled),
		}
	}

	return dataconnectors.MSTIDataConnectorDataTypesMicrosoftEmergingThreatFeed{
		LookbackPeriod: input.MicrosoftEmergingThreatFeedLookBackDate,
		State:          pointer.To(dataconnectors.DataTypeStateEnabled),
	}
}

func flattenSentinelDataConnectorMicrosoftThreatIntelligenceTime(input string) (string, error) {
	// TODO: check if this workaround could be removed in the future
	t, err := time.Parse(time.RFC3339, input)
	if err != nil {
		t, err = time.Parse("01/02/2006 15:04:05", input)
		if err != nil {
			return "", err
		}
	}

	return t.Format(time.RFC3339), nil
}
