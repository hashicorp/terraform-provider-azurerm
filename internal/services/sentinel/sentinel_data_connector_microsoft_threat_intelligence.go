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
	securityinsight "github.com/tombuildsstuff/kermit/sdk/securityinsights/2022-10-01-preview/securityinsights"
)

type DataConnectorMicrosoftThreatIntelligence struct{}

type DataConnectorMicrosoftThreatIntelligenceModel struct {
	Name                        string                                             `tfschema:"name"`
	WorkspaceId                 string                                             `tfschema:"log_analytics_workspace_id"`
	TenantId                    string                                             `tfschema:"tenant_id"`
	BingSafetyPhishingUrl       []DataConnectorMicrosoftThreatIntelligenceDataType `tfschema:"bing_safety_phishing_url"`
	MicrosoftEmergingThreatFeed []DataConnectorMicrosoftThreatIntelligenceDataType `tfschema:"microsoft_emerging_threat_feed"`
}

type DataConnectorMicrosoftThreatIntelligenceDataType struct {
	Enabled      bool   `tfschema:"enabled"`
	LookbackDate string `tfschema:"lookback_date"`
}

func (s DataConnectorMicrosoftThreatIntelligence) Arguments() map[string]*schema.Schema {
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

		"bing_safety_phishing_url": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
					"lookback_date": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "1970-01-01T00:00:00Z",
						ValidateFunc: validation.IsRFC3339Time,
					},
				},
			},
		},

		"microsoft_emerging_threat_feed": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
					"lookback_date": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "1970-01-01T00:00:00Z",
						ValidateFunc: validation.IsRFC3339Time,
					},
				},
			},
		},
	}
}

func (s DataConnectorMicrosoftThreatIntelligence) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (s DataConnectorMicrosoftThreatIntelligence) ModelObject() interface{} {
	return &DataConnectorMicrosoftThreatIntelligenceModel{}
}

func (s DataConnectorMicrosoftThreatIntelligence) ResourceType() string {
	return "azurerm_sentinel_data_connector_microsoft_threat_intelligence"
}

func (s DataConnectorMicrosoftThreatIntelligence) Create() sdk.ResourceFunc {
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

			tenantId := metaModel.TenantId
			if tenantId == "" {
				tenantId = metadata.Client.Account.TenantId
			}

			dataConnector := securityinsight.MSTIDataConnector{
				Name: &id.Name,
				Kind: securityinsight.KindBasicDataConnectorKindMicrosoftThreatIntelligence,
				MSTIDataConnectorProperties: &securityinsight.MSTIDataConnectorProperties{
					DataTypes: &securityinsight.MSTIDataConnectorDataTypes{
						BingSafetyPhishingURL:       expandSentinelDataConnectorMicrosoftThreatIntelligenceBingSafetyPhishingUrl(metaModel.BingSafetyPhishingUrl),
						MicrosoftEmergingThreatFeed: expandSentinelDataConnectorMicrosoftThreatIntelligenceMicrosoftEmergingThreatFeed(metaModel.MicrosoftEmergingThreatFeed),
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

func (s DataConnectorMicrosoftThreatIntelligence) Read() sdk.ResourceFunc {
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
				if dt.BingSafetyPhishingURL != nil {
					bingSafetyPhishingUrl, err := flattenSentinelDataConnectorMicrosoftThreatIntelligenceBingSafetyPhishingUrl(dt.BingSafetyPhishingURL)
					if err != nil {
						return fmt.Errorf("flattening `bing_safety_phishing_url`: %+v", err)
					}
					state.BingSafetyPhishingUrl = bingSafetyPhishingUrl
				}
				if dt.MicrosoftEmergingThreatFeed != nil {
					microsoftEmergingThreatFeed, err := flattenSentinelDataConnectorMicrosoftThreatIntelligenceMicrosoftEmergingThreatFeed(dt.MicrosoftEmergingThreatFeed)
					if err != nil {
						return fmt.Errorf("flattening `microsoft_emerging_threat_feed`: %+v", err)
					}
					state.MicrosoftEmergingThreatFeed = microsoftEmergingThreatFeed
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (s DataConnectorMicrosoftThreatIntelligence) Delete() sdk.ResourceFunc {
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

func (s DataConnectorMicrosoftThreatIntelligence) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.DataConnectorID
}

func expandSentinelDataConnectorMicrosoftThreatIntelligenceBingSafetyPhishingUrl(input []DataConnectorMicrosoftThreatIntelligenceDataType) *securityinsight.MSTIDataConnectorDataTypesBingSafetyPhishingURL {
	if len(input) == 0 {
		return nil
	}

	output := securityinsight.MSTIDataConnectorDataTypesBingSafetyPhishingURL{
		LookbackPeriod: utils.String(input[0].LookbackDate),
	}
	if input[0].Enabled {
		output.State = securityinsight.DataTypeStateEnabled
	} else {
		output.State = securityinsight.DataTypeStateDisabled
	}

	return &output
}

func flattenSentinelDataConnectorMicrosoftThreatIntelligenceBingSafetyPhishingUrl(input *securityinsight.MSTIDataConnectorDataTypesBingSafetyPhishingURL) ([]DataConnectorMicrosoftThreatIntelligenceDataType, error) {
	if input == nil {
		return []DataConnectorMicrosoftThreatIntelligenceDataType{}, nil
	}

	t, err := flattenSentinelDataConnectorMicrosoftThreatIntelligenceTime(*input.LookbackPeriod)
	if err != nil {
		return []DataConnectorMicrosoftThreatIntelligenceDataType{}, err
	}

	output := DataConnectorMicrosoftThreatIntelligenceDataType{
		Enabled:      strings.EqualFold(string(input.State), string(securityinsight.DataTypeStateEnabled)),
		LookbackDate: t,
	}

	return []DataConnectorMicrosoftThreatIntelligenceDataType{output}, nil
}

func expandSentinelDataConnectorMicrosoftThreatIntelligenceMicrosoftEmergingThreatFeed(input []DataConnectorMicrosoftThreatIntelligenceDataType) *securityinsight.MSTIDataConnectorDataTypesMicrosoftEmergingThreatFeed {
	if len(input) == 0 {
		return nil
	}

	output := securityinsight.MSTIDataConnectorDataTypesMicrosoftEmergingThreatFeed{
		LookbackPeriod: utils.String(input[0].LookbackDate),
	}
	if input[0].Enabled {
		output.State = securityinsight.DataTypeStateEnabled
	} else {
		output.State = securityinsight.DataTypeStateDisabled
	}

	return &output
}

func flattenSentinelDataConnectorMicrosoftThreatIntelligenceMicrosoftEmergingThreatFeed(input *securityinsight.MSTIDataConnectorDataTypesMicrosoftEmergingThreatFeed) ([]DataConnectorMicrosoftThreatIntelligenceDataType, error) {
	if input == nil {
		return []DataConnectorMicrosoftThreatIntelligenceDataType{}, nil
	}

	t, err := flattenSentinelDataConnectorMicrosoftThreatIntelligenceTime(*input.LookbackPeriod)
	if err != nil {
		return []DataConnectorMicrosoftThreatIntelligenceDataType{}, err
	}

	output := DataConnectorMicrosoftThreatIntelligenceDataType{
		Enabled:      strings.EqualFold(string(input.State), string(securityinsight.DataTypeStateEnabled)),
		LookbackDate: t,
	}

	return []DataConnectorMicrosoftThreatIntelligenceDataType{output}, nil
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
