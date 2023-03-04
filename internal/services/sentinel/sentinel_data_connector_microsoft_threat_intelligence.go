package sentinel

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
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
	BingSafetyPhishingUrlLookBackDate       string `tfschema:"bing_safety_phishing_url_lookback_date"`
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
			ValidateFunc: dataconnectors.ValidateWorkspaceID,
		},

		"tenant_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"bing_safety_phishing_url_lookback_date": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsRFC3339Time,
			AtLeastOneOf: []string{"bing_safety_phishing_url_lookback_date", "microsoft_emerging_threat_feed_lookback_date"},
		},

		"microsoft_emerging_threat_feed_lookback_date": {
			Type:         pluginsdk.TypeString,
			ForceNew:     true,
			Optional:     true,
			ValidateFunc: validation.IsRFC3339Time,
			AtLeastOneOf: []string{"bing_safety_phishing_url_lookback_date", "microsoft_emerging_threat_feed_lookback_date"},
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

			workSpaceId, err := dataconnectors.ParseWorkspaceID(metaModel.WorkspaceId)
			if err != nil {
				return fmt.Errorf("parsing workspace id %+v", err)
			}

			id := dataconnectors.NewDataConnectorID(workSpaceId.SubscriptionId, workSpaceId.ResourceGroupName, workSpaceId.WorkspaceName, metaModel.Name)
			existing, err := client.DataConnectorsGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(s.ResourceType(), id)
			}

			tenantId := metaModel.TenantId
			if tenantId == "" {
				tenantId = metadata.Client.Account.TenantId
			}

			dataConnector := dataconnectors.MSTIDataConnector{
				Name: &id.DataConnectorId,
				Properties: &dataconnectors.MSTIDataConnectorProperties{
					DataTypes: dataconnectors.MSTIDataConnectorDataTypes{
						BingSafetyPhishingURL:       expandSentinelDataConnectorMicrosoftThreatIntelligenceBingSafetyPhishingUrl(metaModel),
						MicrosoftEmergingThreatFeed: expandSentinelDataConnectorMicrosoftThreatIntelligenceMicrosoftEmergingThreatFeed(metaModel),
					},
					TenantId: tenantId,
				},
			}
			_, err = client.DataConnectorsCreateOrUpdate(ctx, id, dataConnector)
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
			id, err := dataconnectors.ParseDataConnectorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			workspaceId := dataconnectors.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName)

			existing, err := client.DataConnectorsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("model was nil for %s", id)
			}

			modelPtr := *existing.Model
			dc, ok := modelPtr.(dataconnectors.MSTIDataConnector)
			if !ok {
				return fmt.Errorf("%s was not an Microsoft Threat Protection Data Connector", id)
			}

			state := DataConnectorMicrosoftThreatIntelligenceModel{
				Name:        id.DataConnectorId,
				WorkspaceId: workspaceId.ID(),
			}

			if props := dc.Properties; props != nil {
				state.TenantId = props.TenantId

				dt := props.DataTypes
				if dt.BingSafetyPhishingURL.State != nil {
					if strings.EqualFold(string(*dt.BingSafetyPhishingURL.State), string(dataconnectors.DataTypeStateEnabled)) {
						state.BingSafetyPhishingUrlLookBackDate, err = flattenSentinelDataConnectorMicrosoftThreatIntelligenceTime(dt.BingSafetyPhishingURL.LookbackPeriod)
						if err != nil {
							return fmt.Errorf("flattening `bing_safety_phishing_url`: %+v", err)
						}
					}
				}
				if dt.MicrosoftEmergingThreatFeed.State != nil {
					if strings.EqualFold(string(*dt.MicrosoftEmergingThreatFeed.State), string(dataconnectors.DataTypeStateEnabled)) {
						state.MicrosoftEmergingThreatFeedLookBackDate, err = flattenSentinelDataConnectorMicrosoftThreatIntelligenceTime(dt.MicrosoftEmergingThreatFeed.LookbackPeriod)
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

			if _, err := client.DataConnectorsDelete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (s DataConnectorMicrosoftThreatIntelligenceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return dataconnectors.ValidateDataConnectorID
}

func expandSentinelDataConnectorMicrosoftThreatIntelligenceBingSafetyPhishingUrl(input DataConnectorMicrosoftThreatIntelligenceModel) dataconnectors.MSTIDataConnectorDataTypesBingSafetyPhishingURL {
	if input.BingSafetyPhishingUrlLookBackDate == "" {
		disabled := dataconnectors.DataTypeStateDisabled
		return dataconnectors.MSTIDataConnectorDataTypesBingSafetyPhishingURL{
			LookbackPeriod: "",
			State:          &disabled,
		}
	}

	enabled := dataconnectors.DataTypeStateEnabled
	return dataconnectors.MSTIDataConnectorDataTypesBingSafetyPhishingURL{
		LookbackPeriod: input.BingSafetyPhishingUrlLookBackDate,
		State:          &enabled,
	}
}

func expandSentinelDataConnectorMicrosoftThreatIntelligenceMicrosoftEmergingThreatFeed(input DataConnectorMicrosoftThreatIntelligenceModel) dataconnectors.MSTIDataConnectorDataTypesMicrosoftEmergingThreatFeed {
	if input.MicrosoftEmergingThreatFeedLookBackDate == "" {
		disabled := dataconnectors.DataTypeStateDisabled
		return dataconnectors.MSTIDataConnectorDataTypesMicrosoftEmergingThreatFeed{
			LookbackPeriod: "",
			State:          &disabled,
		}
	}

	enabled := dataconnectors.DataTypeStateEnabled
	return dataconnectors.MSTIDataConnectorDataTypesMicrosoftEmergingThreatFeed{
		LookbackPeriod: input.MicrosoftEmergingThreatFeedLookBackDate,
		State:          &enabled,
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
