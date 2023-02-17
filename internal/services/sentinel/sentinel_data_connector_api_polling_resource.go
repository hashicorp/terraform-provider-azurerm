package sentinel

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	securityinsight "github.com/tombuildsstuff/kermit/sdk/securityinsights/2022-10-01-preview/securityinsights"
)

type DataConnectorApiPollModel struct {
	Name                    string                              `tfschema:"name"`
	LogAnalyticsWorkspaceId string                              `tfschema:"log_analytics_workspace_id"`
	Active                  bool                                `tfschema:"active"`
	Auth                    []DataConnectorApiPollAuthModel     `tfschema:"auth"`
	Request                 []DataConnectorApiPollRequestModel  `tfschema:"request"`
	Paging                  []DataConnectorApiPollPagingModel   `tfschema:"paging"`
	Response                []DataConnectorApiPollResponseModel `tfschema:"response"`
	UIConfig                []DataConnectorGenericUIConfigModel `tfschema:"ui"`
}

type DataConnectorApiPollAuthModel struct {
	Type                                 string `tfschema:"type"`
	APIKeyName                           string `tfschema:"api_key_name"`
	APIKeyIdentifier                     string `tfschema:"api_key_identifier"`
	IsAPIKeyInPostPayload                bool   `tfschema:"api_key_in_header_enabled"`
	FlowName                             string `tfschema:"oauth2_flow_name"`
	TokenEndpoint                        string `tfschema:"oauth2_token_endpoint"`
	AuthorizationEndpoint                string `tfschema:"oauth2_authorization_endpoint"`
	AuthorizationEndpointQueryParameters string `tfschema:"oauth2_authorization_endpoint_query_parameters"`
	RedirectionEndpoint                  string `tfschema:"oauth2_redirection_endpoint"`
	TokenEndpointHeaders                 string `tfschema:"oauth2_token_endpoint_headers"`
	TokenEndpointQueryParameters         string `tfschema:"oauth2_token_endpoint_query_parameters"`
	IsClientSecretInHeader               bool   `tfschema:"oauth2_client_secret_in_header_enabled"`
	Scope                                string `tfschema:"oauth2_scope"`
}

type DataConnectorApiPollRequestModel struct {
	ApiEndpoint             string `tfschema:"api_endpoint"`
	RateLimitQPS            int32  `tfschema:"rate_limit_query_per_seconds"`
	QueryWindowInMin        int32  `tfschema:"query_window_in_minutes"`
	HttpMethod              string `tfschema:"http_method"`
	QueryTimeFormat         string `tfschema:"query_time_format"`
	RetryCount              int32  `tfschema:"retry_count"`
	TimeOutInSeconds        int32  `tfschema:"time_out_in_seconds"`
	Headers                 string `tfschema:"headers"`
	QueryParameters         string `tfschema:"query_parameters"`
	QueryParametersTemplate string `tfschema:"query_parameters_template"`
	StartTimeAttributeName  string `tfschema:"start_time_attribute_name"`
	EndTimeAttributeName    string `tfschema:"end_time_attribute_name"`
}

type DataConnectorApiPollPagingModel struct {
	PagingType                             string `tfschema:"type"`
	NextPageParaName                       string `tfschema:"next_page_parameter_name"`
	NextPageTokenJSONPath                  string `tfschema:"next_page_token_json_path"`
	PageCountAttributePath                 string `tfschema:"page_count_attribute_path"`
	PageTotalCountAttributePath            string `tfschema:"page_total_count_attribute_path"`
	PageTimeStampAttributePath             string `tfschema:"page_timestamp_attribute_path"`
	SearchTheLatestTimeStampFromEventsList bool   `tfschema:"search_the_latest_timestamp_from_events_list_enabled"`
	PageSizeParameterName                  string `tfschema:"page_size_parameter_name"`
	PageSize                               int32  `tfschema:"page_size"`
}

type DataConnectorApiPollResponseModel struct {
	EventsJSONPaths       []string `tfschema:"events_json_paths"`
	SuccessStatusJSONPath string   `tfschema:"success_status_json_path"`
	SuccessStatusValue    string   `tfschema:"success_status_value"`
	GzipCompressEnabled   bool     `tfschema:"gzip_compress_enabled"`
}

type DataConnectorApiPollingResource struct{}

var _ sdk.ResourceWithUpdate = DataConnectorApiPollingResource{}

func (r DataConnectorApiPollingResource) Arguments() map[string]*pluginsdk.Schema {
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

		"auth": DataConnectorApiPollAuthSchema(),

		"request": DataConnectorApiPollRequestSchema(),

		"paging": DataConnectorApiPollPagingSchema(),

		"response": DataConnectorApiPollResponseSchema(),

		"ui": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: DataConnectorGenericUIConfigSchema(),
			},
		},
	}
}

func DataConnectorApiPollAuthSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"type": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice([]string{"Basic", "APIKey", "OAuth2"}, false),
					// the enums comes from  https://learn.microsoft.com/en-us/azure/sentinel/create-codeless-connector?tabs=deploy-via-arm-template%2Cconnect-via-the-azure-portal
				},
				"api_key_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"api_key_identifier": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"api_key_in_header_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},
				"oauth2_flow_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"oauth2_token_endpoint": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.IsURLWithHTTPorHTTPS,
				},
				"oauth2_authorization_endpoint": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.IsURLWithHTTPorHTTPS,
				},
				"oauth2_authorization_endpoint_query_parameters": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsJSON,
				},
				"oauth2_redirection_endpoint": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.IsURLWithHTTPorHTTPS,
				},
				"oauth2_token_endpoint_headers": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsJSON,
				},
				"oauth2_token_endpoint_query_parameters": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsJSON,
				},
				"oauth2_client_secret_in_header_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},
				"oauth2_scope": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func DataConnectorApiPollRequestSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"api_endpoint": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.IsURLWithHTTPorHTTPS,
				},
				"http_method": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice([]string{"GET", "POST"}, false),
				},
				"query_time_format": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"start_time_attribute_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"end_time_attribute_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"query_window_in_minutes": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntAtLeast(5),
				},
				"query_parameters": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsJSON,
				},
				"query_parameters_template": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsJSON,
				},
				"rate_limit_query_per_seconds": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntAtLeast(1),
				},
				"time_out_in_seconds": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntAtLeast(1),
				},
				"retry_count": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntAtLeast(1),
				},
				"headers": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsJSON,
				},
			},
		},
	}
}

func DataConnectorApiPollPagingSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"type": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						"None",
						"PageToken",
						"PageCount",
						"TimeStamp",
					}, false),
				},
				"next_page_parameter_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"next_page_token_json_path": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"page_count_attribute_path": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"page_total_count_attribute_path": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"page_timestamp_attribute_path": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"search_the_latest_timestamp_from_events_list_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},
				"page_size_parameter_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"page_size": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntAtLeast(1),
				},
			},
		},
	}
}

func DataConnectorApiPollResponseSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"events_json_paths": {
					Type:     pluginsdk.TypeSet,
					Required: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
				"success_status_json_path": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"success_status_value": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"gzip_compress_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},
			},
		},
	}
}

func (r DataConnectorApiPollingResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"active": {
			Type: pluginsdk.TypeBool,
		},
	}
}

func (r DataConnectorApiPollingResource) ModelObject() interface{} {
	return &DataConnectorApiPollModel{}
}

func (r DataConnectorApiPollingResource) ResourceType() string {
	return "azurerm_sentinel_data_connector_api_polling"
}

func (r DataConnectorApiPollingResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.DataConnectorID
}

func (r DataConnectorApiPollingResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.DataConnectorsClient

			var plan DataConnectorApiPollModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			workspaceId, err := workspaces.ParseWorkspaceID(plan.LogAnalyticsWorkspaceId)
			if err != nil {
				return err
			}

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

			pollAuth, err := expandDataConnectorApiPollAuth(plan.Auth)
			if err != nil {
				return fmt.Errorf("expanding `auth`: %+v", err)
			}

			pollRequest, err := expandDataConnectorApiPollRequest(plan.Request)
			if err != nil {
				return fmt.Errorf("expanding `request`: %+v", err)
			}

			uiConfig, err := expandDataConnectorGenericUIConfigModel(plan.UIConfig)
			if err != nil {
				return fmt.Errorf("expanding `ui`: %+v", err)
			}

			params := securityinsight.CodelessAPIPollingDataConnector{
				Name: &plan.Name,
				APIPollingParameters: &securityinsight.APIPollingParameters{
					PollingConfig: &securityinsight.CodelessConnectorPollingConfigProperties{
						Auth:     pollAuth,
						Request:  pollRequest,
						Paging:   expandDataConnectorApiPollPaging(plan.Paging),
						Response: expandDataConnectorApiPollResponse(plan.Response),
					},
					ConnectorUIConfig: uiConfig,
				},
				Kind: securityinsight.KindBasicDataConnectorKindAPIPolling,
			}

			if _, err = client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.Name, params); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r DataConnectorApiPollingResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.DataConnectorsClient
			id, err := parse.DataConnectorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
			if err != nil {
				if utils.ResponseWasNotFound(existing.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			dc, ok := existing.Value.(securityinsight.CodelessAPIPollingDataConnector)
			if !ok {
				return fmt.Errorf("%s was not an IoT Data Connector", id)
			}

			if dc.PollingConfig == nil {
				return fmt.Errorf("retrieving %s: `polling_config` was nil", id)
			}

			model := DataConnectorApiPollModel{}

			if dc.Name != nil {
				model.Name = *dc.Name
			}

			if dc.PollingConfig.IsActive != nil {
				model.Active = *dc.PollingConfig.IsActive
			}

			model.Auth, err = flattenDataConnectorApiPollAuth(dc.PollingConfig.Auth)
			if err != nil {
				return fmt.Errorf("flattening `auth`: %+v", err)
			}

			model.Request, err = flattenDataConnectorApiPollRequest(dc.PollingConfig.Request)
			if err != nil {
				return fmt.Errorf("flattening `request`: %+v", err)
			}

			model.Response = flattenDataConnectorApiPollResponse(dc.PollingConfig.Response)

			model.Paging = flattenDataConnectorAPiPollPaging(dc.PollingConfig.Paging)

			model.UIConfig, err = flattenDataConnectorGenericUIConfigModel(dc.ConnectorUIConfig)
			if err != nil {
				return fmt.Errorf("flattening `ui`: %+v", err)
			}

			workspaceId := workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName)
			model.LogAnalyticsWorkspaceId = workspaceId.ID()

			return metadata.Encode(&model)
		},
	}
}

func (r DataConnectorApiPollingResource) Delete() sdk.ResourceFunc {
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

func (r DataConnectorApiPollingResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.DataConnectorsClient
			id, err := parse.DataConnectorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var plan DataConnectorApiPollModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			_, err = client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
			if err != nil {
				return fmt.Errorf("reading %s: %+v", id, err)
			}

			pollAuth, err := expandDataConnectorApiPollAuth(plan.Auth)
			if err != nil {
				return fmt.Errorf("expanding `auth`: %+v", err)
			}

			pollRequest, err := expandDataConnectorApiPollRequest(plan.Request)
			if err != nil {
				return fmt.Errorf("expanding `request`: %+v", err)
			}

			uiConfig, err := expandDataConnectorGenericUIConfigModel(plan.UIConfig)
			if err != nil {
				return fmt.Errorf("expanding `ui`: %+v", err)
			}

			params := securityinsight.CodelessAPIPollingDataConnector{
				Name: &plan.Name,
				APIPollingParameters: &securityinsight.APIPollingParameters{
					PollingConfig: &securityinsight.CodelessConnectorPollingConfigProperties{
						Auth:     pollAuth,
						Request:  pollRequest,
						Paging:   expandDataConnectorApiPollPaging(plan.Paging),
						Response: expandDataConnectorApiPollResponse(plan.Response),
					},
					ConnectorUIConfig: uiConfig,
				},
				Kind: securityinsight.KindBasicDataConnectorKindAPIPolling,
			}

			if _, err = client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.Name, params); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandDataConnectorApiPollAuth(input []DataConnectorApiPollAuthModel) (*securityinsight.CodelessConnectorPollingAuthProperties, error) {
	if len(input) == 0 {
		return nil, nil
	}
	item := input[0]
	apiKeyInPostPayload := strconv.FormatBool(item.IsAPIKeyInPostPayload)

	output := securityinsight.CodelessConnectorPollingAuthProperties{
		AuthType:               &item.Type,
		APIKeyName:             &item.APIKeyName,
		APIKeyIdentifier:       &item.APIKeyIdentifier,
		IsAPIKeyInPostPayload:  &apiKeyInPostPayload,
		FlowName:               &item.FlowName,
		TokenEndpoint:          &item.TokenEndpoint,
		AuthorizationEndpoint:  &item.AuthorizationEndpoint,
		RedirectionEndpoint:    &item.RedirectionEndpoint,
		IsClientSecretInHeader: utils.Bool(item.IsClientSecretInHeader),
		Scope:                  &item.Scope,
	}

	if item.AuthorizationEndpointQueryParameters != "" {
		authorizationEndpointQueryParameters, err := pluginsdk.ExpandJsonFromString(item.AuthorizationEndpointQueryParameters)
		if err != nil {
			return nil, err
		}
		output.AuthorizationEndpointQueryParameters = authorizationEndpointQueryParameters
	}

	if item.TokenEndpointHeaders != "" {
		tokenEndpointHeaders, err := pluginsdk.ExpandJsonFromString(item.TokenEndpointHeaders)
		if err != nil {
			return nil, err
		}
		output.TokenEndpointHeaders = tokenEndpointHeaders
	}
	if item.TokenEndpointQueryParameters != "" {
		tokenEndpointQueryParameters, err := pluginsdk.ExpandJsonFromString(item.TokenEndpointQueryParameters)
		if err != nil {
			return nil, err
		}
		output.TokenEndpointQueryParameters = tokenEndpointQueryParameters
	}

	return &output, nil
}

func flattenDataConnectorApiPollAuth(input *securityinsight.CodelessConnectorPollingAuthProperties) ([]DataConnectorApiPollAuthModel, error) {
	var output []DataConnectorApiPollAuthModel
	if input == nil {
		return output, nil
	}

	item := DataConnectorApiPollAuthModel{}
	if input.AuthType != nil {
		item.Type = *input.AuthType
	}
	if input.APIKeyName != nil {
		item.APIKeyName = *input.APIKeyName
	}
	if input.APIKeyIdentifier != nil {
		item.APIKeyIdentifier = *input.APIKeyIdentifier
	}
	if input.IsAPIKeyInPostPayload != nil {
		item.IsAPIKeyInPostPayload = *input.IsAPIKeyInPostPayload == "1"
	}
	if input.FlowName != nil {
		item.FlowName = *input.FlowName
	}
	if input.TokenEndpoint != nil {
		item.TokenEndpoint = *input.TokenEndpoint
	}
	if input.AuthorizationEndpoint != nil {
		item.AuthorizationEndpoint = *input.AuthorizationEndpoint
	}
	if input.AuthorizationEndpointQueryParameters != nil {
		json, err := pluginsdk.FlattenJsonToString(input.AuthorizationEndpointQueryParameters.(map[string]interface{}))
		if err != nil {
			return output, err
		}
		item.AuthorizationEndpointQueryParameters = json
	}
	if input.RedirectionEndpoint != nil {
		item.RedirectionEndpoint = *input.RedirectionEndpoint
	}
	if input.TokenEndpointHeaders != nil {
		json, err := pluginsdk.FlattenJsonToString(input.TokenEndpointHeaders.(map[string]interface{}))
		if err != nil {
			return output, err
		}
		item.TokenEndpointHeaders = json
	}
	if input.TokenEndpointQueryParameters != nil {
		json, err := pluginsdk.FlattenJsonToString(input.TokenEndpointQueryParameters.(map[string]interface{}))
		if err != nil {
			return output, err
		}
		item.TokenEndpointQueryParameters = json
	}
	if input.IsClientSecretInHeader != nil {
		item.IsClientSecretInHeader = *input.IsClientSecretInHeader
	}
	if input.Scope != nil {
		item.Scope = *input.Scope
	}

	return append(output, item), nil
}

func expandDataConnectorApiPollRequest(input []DataConnectorApiPollRequestModel) (*securityinsight.CodelessConnectorPollingRequestProperties, error) {
	if len(input) == 0 {
		return nil, nil
	}
	item := input[0]

	output := securityinsight.CodelessConnectorPollingRequestProperties{
		APIEndpoint:             &item.ApiEndpoint,
		RateLimitQPS:            &item.RateLimitQPS,
		QueryWindowInMin:        &item.QueryWindowInMin,
		HTTPMethod:              &item.HttpMethod,
		QueryTimeFormat:         &item.QueryTimeFormat,
		RetryCount:              &item.RetryCount,
		TimeoutInSeconds:        &item.TimeOutInSeconds,
		StartTimeAttributeName:  &item.StartTimeAttributeName,
		EndTimeAttributeName:    &item.EndTimeAttributeName,
		QueryParametersTemplate: &item.QueryParametersTemplate,
	}

	if item.QueryParameters != "" {
		queryParameters, err := pluginsdk.ExpandJsonFromString(item.QueryParameters)
		if err != nil {
			return nil, err
		}
		output.QueryParameters = queryParameters
	}

	if item.Headers != "" {
		headers, err := pluginsdk.ExpandJsonFromString(item.Headers)
		if err != nil {
			return nil, err
		}
		output.Headers = headers
	}

	return &output, nil
}

func flattenDataConnectorApiPollRequest(input *securityinsight.CodelessConnectorPollingRequestProperties) ([]DataConnectorApiPollRequestModel, error) {
	var output []DataConnectorApiPollRequestModel
	if input == nil {
		return output, nil
	}

	item := DataConnectorApiPollRequestModel{}
	if input.APIEndpoint != nil {
		item.ApiEndpoint = *input.APIEndpoint
	}
	if input.RateLimitQPS != nil {
		item.RateLimitQPS = *input.RateLimitQPS
	}
	if input.QueryWindowInMin != nil {
		item.QueryWindowInMin = *input.QueryWindowInMin
	}
	if input.HTTPMethod != nil {
		item.HttpMethod = *input.HTTPMethod
	}
	if input.QueryTimeFormat != nil {
		item.QueryTimeFormat = *input.QueryTimeFormat
	}
	if input.RetryCount != nil {
		item.RetryCount = *input.RetryCount
	}
	if input.TimeoutInSeconds != nil {
		item.TimeOutInSeconds = *input.TimeoutInSeconds
	}
	if input.Headers != nil {
		json, err := pluginsdk.FlattenJsonToString(input.Headers.(map[string]interface{}))
		if err != nil {
			return output, err
		}
		item.Headers = json
	}
	if input.QueryParameters != nil {
		json, err := pluginsdk.FlattenJsonToString(input.QueryParameters.(map[string]interface{}))
		if err != nil {
			return output, err
		}
		item.QueryParameters = json
	}
	if input.StartTimeAttributeName != nil {
		item.StartTimeAttributeName = *input.StartTimeAttributeName
	}
	if input.EndTimeAttributeName != nil {
		item.EndTimeAttributeName = *input.EndTimeAttributeName
	}
	if input.QueryParametersTemplate != nil {
		item.QueryParametersTemplate = *input.QueryParametersTemplate
	}
	return append(output, item), nil
}

func expandDataConnectorApiPollPaging(input []DataConnectorApiPollPagingModel) *securityinsight.CodelessConnectorPollingPagingProperties {
	if len(input) == 0 {
		return nil
	}
	item := input[0]
	searchTheLatestTimestampFromEventsListEnabledSpecified := strconv.FormatBool(item.SearchTheLatestTimeStampFromEventsList)
	return &securityinsight.CodelessConnectorPollingPagingProperties{
		PagingType:                             &item.PagingType,
		NextPageParaName:                       &item.NextPageParaName,
		NextPageTokenJSONPath:                  &item.NextPageTokenJSONPath,
		PageCountAttributePath:                 &item.PageCountAttributePath,
		PageTotalCountAttributePath:            &item.PageTotalCountAttributePath,
		PageTimeStampAttributePath:             &item.PageTimeStampAttributePath,
		SearchTheLatestTimeStampFromEventsList: &searchTheLatestTimestampFromEventsListEnabledSpecified,
		PageSizeParaName:                       &item.PageSizeParameterName,
		PageSize:                               &item.PageSize,
	}
}

func flattenDataConnectorAPiPollPaging(input *securityinsight.CodelessConnectorPollingPagingProperties) []DataConnectorApiPollPagingModel {
	var output []DataConnectorApiPollPagingModel
	item := DataConnectorApiPollPagingModel{}
	if input == nil {
		return output
	}
	if input.PagingType != nil {
		item.PagingType = *input.PagingType
	}
	if input.NextPageParaName != nil {
		item.NextPageParaName = *input.NextPageParaName
	}
	if input.NextPageTokenJSONPath != nil {
		item.NextPageTokenJSONPath = *input.NextPageTokenJSONPath
	}
	if input.PageCountAttributePath != nil {
		item.PageCountAttributePath = *input.PageCountAttributePath
	}
	if input.PageTotalCountAttributePath != nil {
		item.PageTotalCountAttributePath = *input.PageTotalCountAttributePath
	}
	if input.PageTimeStampAttributePath != nil {
		item.PageTimeStampAttributePath = *input.PageTimeStampAttributePath
	}
	if input.SearchTheLatestTimeStampFromEventsList != nil {
		item.SearchTheLatestTimeStampFromEventsList, _ = strconv.ParseBool(*input.SearchTheLatestTimeStampFromEventsList)
	}
	if input.PageSizeParaName != nil {
		item.PageSizeParameterName = *input.PageSizeParaName
	}
	if input.PageSize != nil {
		item.PageSize = *input.PageSize
	}
	return append(output, item)
}

func expandDataConnectorApiPollResponse(input []DataConnectorApiPollResponseModel) *securityinsight.CodelessConnectorPollingResponseProperties {
	if len(input) == 0 {
		return nil
	}
	item := input[0]
	return &securityinsight.CodelessConnectorPollingResponseProperties{
		EventsJSONPaths:       &item.EventsJSONPaths,
		SuccessStatusJSONPath: &item.SuccessStatusJSONPath,
		SuccessStatusValue:    &item.SuccessStatusValue,
		IsGzipCompressed:      utils.Bool(item.GzipCompressEnabled),
	}
}

func flattenDataConnectorApiPollResponse(input *securityinsight.CodelessConnectorPollingResponseProperties) []DataConnectorApiPollResponseModel {
	var output []DataConnectorApiPollResponseModel
	if input == nil {
		return output
	}
	item := DataConnectorApiPollResponseModel{}
	if input.EventsJSONPaths != nil {
		item.EventsJSONPaths = *input.EventsJSONPaths
	}
	if input.SuccessStatusJSONPath != nil {
		item.SuccessStatusJSONPath = *input.SuccessStatusJSONPath
	}
	if input.SuccessStatusValue != nil {
		item.SuccessStatusValue = *input.SuccessStatusValue
	}
	if input.IsGzipCompressed != nil {
		item.GzipCompressEnabled = *input.IsGzipCompressed
	}
	return append(output, item)
}
