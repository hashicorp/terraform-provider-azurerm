// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sentinel

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/workspaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	securityinsight "github.com/jackofallops/kermit/sdk/securityinsights/2022-10-01-preview/securityinsights"
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

			id := parse.NewDataConnectorID(workSpaceId.SubscriptionId, workSpaceId.ResourceGroupName, workSpaceId.WorkspaceName, metaModel.Name)
			existing, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
			if err != nil {
				if !utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(s.ResourceType(), id)
			}

			tenantId := metaModel.TenantId
			if tenantId == "" {
				tenantId = metadata.Client.Account.TenantId
			}

			dataConnector := securityinsight.MSTIDataConnector{
				Name: &id.Name,
				Kind: securityinsight.KindBasicDataConnectorKindMicrosoftThreatIntelligence,
				MSTIDataConnectorProperties: &securityinsight.MSTIDataConnectorProperties{
					DataTypes: &securityinsight.MSTIDataConnectorDataTypes{
						MicrosoftEmergingThreatFeed: expandSentinelDataConnectorMicrosoftThreatIntelligenceMicrosoftEmergingThreatFeed(metaModel),
					},
					TenantID: &tenantId,
				},
			}
			_, err = client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.Name, dataConnector)
			if err != nil {
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

			dc, ok := existing.Value.(securityinsight.MSTIDataConnector)
			if !ok {
				return fmt.Errorf("%s was not an Microsoft Threat Protection Data Connector", id)
			}

			state := DataConnectorMicrosoftThreatIntelligenceModel{
				Name:        id.Name,
				WorkspaceId: workspaceId.ID(),
				TenantId:    *dc.TenantID,
			}

			if dc.TenantID != nil {
				state.TenantId = *dc.TenantID
			}

			if dt := dc.DataTypes; dt != nil {
				if dt.MicrosoftEmergingThreatFeed != nil {
					if strings.EqualFold(string(dt.MicrosoftEmergingThreatFeed.State), string(securityinsight.DataTypeStateEnabled)) {
						state.MicrosoftEmergingThreatFeedLookBackDate, err = flattenSentinelDataConnectorMicrosoftThreatIntelligenceTime(*dt.MicrosoftEmergingThreatFeed.LookbackPeriod)
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

func (s DataConnectorMicrosoftThreatIntelligenceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.DataConnectorID
}

func expandSentinelDataConnectorMicrosoftThreatIntelligenceMicrosoftEmergingThreatFeed(input DataConnectorMicrosoftThreatIntelligenceModel) *securityinsight.MSTIDataConnectorDataTypesMicrosoftEmergingThreatFeed {
	if input.MicrosoftEmergingThreatFeedLookBackDate == "" {
		return &securityinsight.MSTIDataConnectorDataTypesMicrosoftEmergingThreatFeed{
			LookbackPeriod: utils.String(""),
			State:          securityinsight.DataTypeStateDisabled,
		}
	}

	return &securityinsight.MSTIDataConnectorDataTypesMicrosoftEmergingThreatFeed{
		LookbackPeriod: utils.String(input.MicrosoftEmergingThreatFeedLookBackDate),
		State:          securityinsight.DataTypeStateEnabled,
	}
}

func flattenSentinelDataConnectorMicrosoftThreatIntelligenceTime(input string) (string, error) {
	// TODO: check if this workaround could be removed in 4.0

	t, err := time.Parse(time.RFC3339, input)
	if err != nil {
		t, err = time.Parse("01/02/2006 15:04:05", input)
		if err != nil {
			return "", err
		}
	}

	return t.Format(time.RFC3339), nil
}
