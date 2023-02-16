package sentinel

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	securityinsight "github.com/tombuildsstuff/kermit/sdk/securityinsights/2022-10-01-preview/securityinsights"
)

type DataConnectorGenericUIConfigModel struct {
	Title                 string                                            `tfschema:"title"`
	Publisher             string                                            `tfschema:"publisher"`
	DescriptionMarkdown   string                                            `tfschema:"description_markdown"`
	CustomImage           string                                            `tfschema:"custom_image"`
	GraphQueriesTableName string                                            `tfschema:"graph_queries_table_name"`
	GraphQueries          []DataConnectorGenericUIGraphQueryModel           `tfschema:"graph_query"`
	SampleQueries         []DataConnectorGenericUISampleQueryModel          `tfschema:"sample_query"`
	DataTypes             []DataConnectorGenericUIDataTypeModel             `tfschema:"data_type"`
	ConnectivityCriteria  []DataConnectorGenericUIConnectivityCriteriaModel `tfschema:"connectivity_criteria"`
	Availability          []DataConnectorGenericUIAvailabilityModel         `tfschema:"availability"`
	Permissions           []DataConnectorGenericUIPermissionsModel          `tfschema:"permission"`
	InstructionSteps      []DataConnectorGenericUIInstructionStepModel      `tfschema:"instruction"`
}

func expandDataConnectorGenericUIConfigModel(input []DataConnectorGenericUIConfigModel) (*securityinsight.CodelessUIConnectorConfigProperties, error) {
	if len(input) == 0 {
		return nil, nil
	}
	item := input[0]

	instruction, err := expandDataConnectorGenericUIInstructionStepModel(item.InstructionSteps)
	if err != nil {
		return nil, err
	}
	return &securityinsight.CodelessUIConnectorConfigProperties{
		Title:                 &item.Title,
		Publisher:             &item.Publisher,
		DescriptionMarkdown:   &item.DescriptionMarkdown,
		CustomImage:           &item.CustomImage,
		GraphQueriesTableName: &item.GraphQueriesTableName,
		GraphQueries:          expandDataConnectorGenericUIGraphQueryModel(item.GraphQueries),
		SampleQueries:         expandDataConnectorGenericUISampleQueryModel(item.SampleQueries),
		DataTypes:             expandDataConnectorGenericUIDataTypeModel(item.DataTypes),
		ConnectivityCriteria:  expandDataConnectorGenericUIConnectivityCriteriaModel(item.ConnectivityCriteria),
		Availability:          expandDataConnectorGenericUIAvailabilityModel(item.Availability),
		Permissions:           expandDataConnectorGenericUIPermissionsModel(item.Permissions),
		InstructionSteps:      instruction,
	}, nil
}

func flattenDataConnectorGenericUIConfigModel(input *securityinsight.CodelessUIConnectorConfigProperties) ([]DataConnectorGenericUIConfigModel, error) {
	var output []DataConnectorGenericUIConfigModel
	if input == nil {
		return output, nil
	}

	instruction, err := flattenDataConnectorGenericUIInstructionStepModel(input.InstructionSteps)
	if err != nil {
		return output, err
	}

	return append(output, DataConnectorGenericUIConfigModel{
		Title:                 *input.Title,
		Publisher:             *input.Publisher,
		DescriptionMarkdown:   *input.DescriptionMarkdown,
		CustomImage:           *input.CustomImage,
		GraphQueriesTableName: *input.GraphQueriesTableName,
		GraphQueries:          flattenDataConnectorGenericUIGraphQueryModel(input.GraphQueries),
		SampleQueries:         flattenDataConnectorGenericUISampleQueryModel(input.SampleQueries),
		DataTypes:             flattenDataConnectorGenericUIDataTypeModel(input.DataTypes),
		ConnectivityCriteria:  flattenDataConnectorGenericUIConnectivityCriteriaModel(input.ConnectivityCriteria),
		Availability:          flattenDataConnectorGenericUIAvailabilityModel(input.Availability),
		Permissions:           flattenDataConnectorGenericUIPermissionsModel(input.Permissions),
		InstructionSteps:      instruction,
	}), nil
}

type DataConnectorGenericUIGraphQueryModel struct {
	MetricName string `tfschema:"metric_name"`
	Legend     string `tfschema:"legend"`
	BaseQuery  string `tfschema:"base_query"`
}

func expandDataConnectorGenericUIGraphQueryModel(input []DataConnectorGenericUIGraphQueryModel) *[]securityinsight.CodelessUIConnectorConfigPropertiesGraphQueriesItem {
	output := make([]securityinsight.CodelessUIConnectorConfigPropertiesGraphQueriesItem, 0)
	for _, item := range input {
		output = append(output, securityinsight.CodelessUIConnectorConfigPropertiesGraphQueriesItem{
			MetricName: &item.MetricName,
			Legend:     &item.Legend,
			BaseQuery:  &item.BaseQuery,
		})
	}

	return &output
}

func flattenDataConnectorGenericUIGraphQueryModel(input *[]securityinsight.CodelessUIConnectorConfigPropertiesGraphQueriesItem) []DataConnectorGenericUIGraphQueryModel {
	output := make([]DataConnectorGenericUIGraphQueryModel, 0)
	if input == nil {
		return output
	}

	for _, item := range *input {
		o := DataConnectorGenericUIGraphQueryModel{}
		if item.MetricName != nil {
			o.MetricName = *item.MetricName
		}
		if item.Legend != nil {
			o.Legend = *item.Legend
		}
		if item.BaseQuery != nil {
			o.BaseQuery = *item.BaseQuery
		}
		output = append(output, o)
	}

	return output
}

type DataConnectorGenericUISampleQueryModel struct {
	Description string `tfschema:"description"`
	Query       string `tfschema:"query"`
}

func expandDataConnectorGenericUISampleQueryModel(input []DataConnectorGenericUISampleQueryModel) *[]securityinsight.CodelessUIConnectorConfigPropertiesSampleQueriesItem {
	output := make([]securityinsight.CodelessUIConnectorConfigPropertiesSampleQueriesItem, 0)
	for _, item := range input {
		output = append(output, securityinsight.CodelessUIConnectorConfigPropertiesSampleQueriesItem{
			Description: &item.Description,
			Query:       &item.Query,
		})
	}

	return &output
}

func flattenDataConnectorGenericUISampleQueryModel(input *[]securityinsight.CodelessUIConnectorConfigPropertiesSampleQueriesItem) []DataConnectorGenericUISampleQueryModel {
	output := make([]DataConnectorGenericUISampleQueryModel, 0)
	if input == nil {
		return output
	}

	for _, item := range *input {
		o := DataConnectorGenericUISampleQueryModel{}
		if item.Description != nil {
			o.Description = *item.Description
		}
		if item.Query != nil {
			o.Query = *item.Query
		}
		output = append(output, o)
	}

	return output
}

type DataConnectorGenericUIDataTypeModel struct {
	Name                  string `tfschema:"name"`
	LastDataReceivedQuery string `tfschema:"last_data_received_query"`
}

func expandDataConnectorGenericUIDataTypeModel(input []DataConnectorGenericUIDataTypeModel) *[]securityinsight.CodelessUIConnectorConfigPropertiesDataTypesItem {
	output := make([]securityinsight.CodelessUIConnectorConfigPropertiesDataTypesItem, 0)
	for _, item := range input {
		output = append(output, securityinsight.CodelessUIConnectorConfigPropertiesDataTypesItem{
			Name:                  &item.Name,
			LastDataReceivedQuery: &item.LastDataReceivedQuery,
		})
	}

	return &output
}

func flattenDataConnectorGenericUIDataTypeModel(input *[]securityinsight.CodelessUIConnectorConfigPropertiesDataTypesItem) []DataConnectorGenericUIDataTypeModel {
	output := make([]DataConnectorGenericUIDataTypeModel, 0)
	if input == nil {
		return output
	}

	for _, item := range *input {
		o := DataConnectorGenericUIDataTypeModel{}
		if item.Name != nil {
			o.Name = *item.Name
		}
		if item.LastDataReceivedQuery != nil {
			o.LastDataReceivedQuery = *item.LastDataReceivedQuery
		}
		output = append(output, o)
	}

	return output
}

type DataConnectorGenericUIConnectivityCriteriaModel struct {
	Type  string   `tfschema:"type"`
	Value []string `tfschema:"value"`
}

func expandDataConnectorGenericUIConnectivityCriteriaModel(input []DataConnectorGenericUIConnectivityCriteriaModel) *[]securityinsight.CodelessUIConnectorConfigPropertiesConnectivityCriteriaItem {
	output := make([]securityinsight.CodelessUIConnectorConfigPropertiesConnectivityCriteriaItem, 0)
	for _, item := range input {
		output = append(output, securityinsight.CodelessUIConnectorConfigPropertiesConnectivityCriteriaItem{
			Type:  securityinsight.ConnectivityType(item.Type),
			Value: &item.Value,
		})
	}

	return &output
}

func flattenDataConnectorGenericUIConnectivityCriteriaModel(input *[]securityinsight.CodelessUIConnectorConfigPropertiesConnectivityCriteriaItem) []DataConnectorGenericUIConnectivityCriteriaModel {
	output := make([]DataConnectorGenericUIConnectivityCriteriaModel, 0)
	if input == nil {
		return output
	}

	for _, item := range *input {
		o := DataConnectorGenericUIConnectivityCriteriaModel{
			Type: string(item.Type),
		}
		if item.Value != nil && len(*item.Value) > 0 {
			o.Value = *item.Value
		}
		output = append(output, o)
	}

	return output
}

type DataConnectorGenericUIAvailabilityModel struct {
	Enabled bool `tfschema:"enabled"`
	Preview bool `tfschema:"preview"`
}

func expandDataConnectorGenericUIAvailabilityModel(input []DataConnectorGenericUIAvailabilityModel) *securityinsight.Availability {
	if len(input) == 0 {
		return nil
	}
	item := input[0]
	var status int32 = 0
	if item.Enabled {
		status = 1
	}
	output := &securityinsight.Availability{
		Status:    &status,
		IsPreview: &item.Preview,
	}

	return output
}

func flattenDataConnectorGenericUIAvailabilityModel(input *securityinsight.Availability) []DataConnectorGenericUIAvailabilityModel {
	output := DataConnectorGenericUIAvailabilityModel{}
	if input.Status != nil {
		output.Enabled = *input.Status == 1
	}
	if input.IsPreview != nil {
		output.Preview = *input.IsPreview
	}

	return append([]DataConnectorGenericUIAvailabilityModel{}, output)
}

type DataConnectorGenericUIPermissionsModel struct {
	ResourceProviders []DataConnectorGenericUIPermissionsResourceProviderModel `tfschema:"resource_provider"`
	Customs           []DataConnectorGenericUIPermissionsCustomModel           `tfschema:"custom"`
}

func expandDataConnectorGenericUIPermissionsModel(input []DataConnectorGenericUIPermissionsModel) *securityinsight.Permissions {
	if len(input) == 0 {
		return nil
	}
	item := input[0]

	output := &securityinsight.Permissions{
		ResourceProvider: expandDataConnectorGenericUIPermissionsResourceProviderModel(item.ResourceProviders),
		Customs:          expandDataConnectorGenericUIPermissionsCustomModel(item.Customs),
	}

	return output
}

func flattenDataConnectorGenericUIPermissionsModel(input *securityinsight.Permissions) []DataConnectorGenericUIPermissionsModel {
	output := DataConnectorGenericUIPermissionsModel{}
	if input.ResourceProvider != nil {
		output.ResourceProviders = flattenDataConnectorGenericUIPermissionsResourceProviderModel(input.ResourceProvider)
	}
	if input.Customs != nil {
		output.Customs = flattenDataConnectorGenericUIPermissionsCustomModel(input.Customs)
	}

	return append([]DataConnectorGenericUIPermissionsModel{}, output)
}

type DataConnectorGenericUIPermissionsResourceProviderModel struct {
	Name                  string                                                     `tfschema:"name"`
	ProviderDisplayName   string                                                     `tfschema:"display_name"`
	PermissionDisplayText string                                                     `tfschema:"display_text"`
	Scope                 string                                                     `tfschema:"scope"`
	RequiredPermissions   []DataConnectorGenericUIPermissionRequiredPermissionsModel `tfschema:"required_permissions"`
}

func expandDataConnectorGenericUIPermissionsResourceProviderModel(input []DataConnectorGenericUIPermissionsResourceProviderModel) *[]securityinsight.PermissionsResourceProviderItem {
	output := make([]securityinsight.PermissionsResourceProviderItem, 0)
	for _, item := range input {
		output = append(output, securityinsight.PermissionsResourceProviderItem{
			Provider:               securityinsight.ProviderName(item.Name),
			ProviderDisplayName:    &item.ProviderDisplayName,
			PermissionsDisplayText: &item.PermissionDisplayText,
			Scope:                  securityinsight.PermissionProviderScope(item.Scope),
			RequiredPermissions:    expandDataConnectorGenericUIPermissionRequiredPermissionsModel(item.RequiredPermissions),
		})
	}

	return &output
}

func flattenDataConnectorGenericUIPermissionsResourceProviderModel(input *[]securityinsight.PermissionsResourceProviderItem) []DataConnectorGenericUIPermissionsResourceProviderModel {
	output := make([]DataConnectorGenericUIPermissionsResourceProviderModel, 0)
	if input == nil {
		return output
	}

	for _, item := range *input {
		o := DataConnectorGenericUIPermissionsResourceProviderModel{}
		if item.Provider != "" {
			o.Name = string(item.Provider)
		}
		if item.ProviderDisplayName != nil {
			o.ProviderDisplayName = *item.ProviderDisplayName
		}
		if item.PermissionsDisplayText != nil {
			o.PermissionDisplayText = *item.PermissionsDisplayText
		}
		if item.Scope != "" {
			o.Scope = string(item.Scope)
		}
		if item.RequiredPermissions != nil {
			o.RequiredPermissions = flattenDataConnectorGenericUIPermissionRequiredPermissionsModel(item.RequiredPermissions)
		}
		output = append(output, o)
	}

	return output
}

type DataConnectorGenericUIPermissionRequiredPermissionsModel struct {
	Action bool `tfschema:"action"`
	Write  bool `tfschema:"write"`
	Read   bool `tfschema:"read"`
	Delete bool `tfschema:"delete"`
}

func expandDataConnectorGenericUIPermissionRequiredPermissionsModel(input []DataConnectorGenericUIPermissionRequiredPermissionsModel) *securityinsight.RequiredPermissions {
	if len(input) == 0 {
		return nil
	}
	item := input[0]

	output := &securityinsight.RequiredPermissions{
		Action: &item.Action,
		Write:  &item.Write,
		Read:   &item.Read,
		Delete: &item.Delete,
	}

	return output
}

func flattenDataConnectorGenericUIPermissionRequiredPermissionsModel(input *securityinsight.RequiredPermissions) []DataConnectorGenericUIPermissionRequiredPermissionsModel {
	output := DataConnectorGenericUIPermissionRequiredPermissionsModel{}
	if input.Action != nil {
		output.Action = *input.Action
	}
	if input.Write != nil {
		output.Write = *input.Write
	}
	if input.Read != nil {
		output.Read = *input.Read
	}
	if input.Delete != nil {
		output.Delete = *input.Delete
	}

	return append([]DataConnectorGenericUIPermissionRequiredPermissionsModel{}, output)
}

type DataConnectorGenericUIPermissionsCustomModel struct {
	Name        string `tfschema:"name"`
	Description string `tfschema:"description"`
}

func expandDataConnectorGenericUIPermissionsCustomModel(input []DataConnectorGenericUIPermissionsCustomModel) *[]securityinsight.PermissionsCustomsItem {
	output := make([]securityinsight.PermissionsCustomsItem, 0)
	for _, v := range input {
		output = append(output, securityinsight.PermissionsCustomsItem{
			Name:        &v.Name,
			Description: &v.Description,
		})
	}
	return &output
}

func flattenDataConnectorGenericUIPermissionsCustomModel(input *[]securityinsight.PermissionsCustomsItem) []DataConnectorGenericUIPermissionsCustomModel {
	output := make([]DataConnectorGenericUIPermissionsCustomModel, 0)
	if input != nil {
		for _, v := range *input {
			o := DataConnectorGenericUIPermissionsCustomModel{}
			if v.Name != nil {
				o.Name = *v.Name
			}
			if v.Description != nil {
				o.Description = *v.Description
			}
			output = append(output, o)
		}
	}
	return output
}

type DataConnectorGenericUIInstructionStepModel struct {
	Title        string                                                  `tfschema:"title"`
	Description  string                                                  `tfschema:"description"`
	Instructions []DataConnectorGenericUIInstructionStepInstructionModel `tfschema:"step"`
}

func expandDataConnectorGenericUIInstructionStepModel(input []DataConnectorGenericUIInstructionStepModel) (*[]securityinsight.CodelessUIConnectorConfigPropertiesInstructionStepsItem, error) {
	output := make([]securityinsight.CodelessUIConnectorConfigPropertiesInstructionStepsItem, 0)
	for _, v := range input {
		instruction, err := expandDataConnectorGenericUIInstructionStepInstructionModel(v.Instructions)
		if err != nil {
			return nil, err
		}
		output = append(output, securityinsight.CodelessUIConnectorConfigPropertiesInstructionStepsItem{
			Title:        &v.Title,
			Description:  &v.Description,
			Instructions: instruction,
		})
	}
	return &output, nil
}

func flattenDataConnectorGenericUIInstructionStepModel(input *[]securityinsight.CodelessUIConnectorConfigPropertiesInstructionStepsItem) ([]DataConnectorGenericUIInstructionStepModel, error) {
	output := make([]DataConnectorGenericUIInstructionStepModel, 0)
	for _, v := range *input {
		o := DataConnectorGenericUIInstructionStepModel{}
		instruction, err := flattenDataConnectorGenericUIInstructionStepInstructionModel(v.Instructions)
		if err != nil {
			return nil, err
		}
		o.Instructions = instruction
		if v.Title != nil {
			o.Title = *v.Title
		}
		if v.Description != nil {
			o.Description = *v.Description
		}
		output = append(output, o)
	}
	return output, nil
}

type DataConnectorGenericUIInstructionStepInstructionModel struct {
	Type       string `tfschema:"type"`
	Parameters string `tfschema:"parameters"`
}

func expandDataConnectorGenericUIInstructionStepInstructionModel(input []DataConnectorGenericUIInstructionStepInstructionModel) (*[]securityinsight.InstructionStepsInstructionsItem, error) {
	output := make([]securityinsight.InstructionStepsInstructionsItem, 0)

	for _, v := range input {
		param, err := pluginsdk.ExpandJsonFromString(v.Parameters)
		if err != nil {
			return nil, err
		}
		output = append(output, securityinsight.InstructionStepsInstructionsItem{
			Type:       securityinsight.SettingType(v.Type),
			Parameters: &param,
		})
	}
	return &output, nil
}

func flattenDataConnectorGenericUIInstructionStepInstructionModel(input *[]securityinsight.InstructionStepsInstructionsItem) ([]DataConnectorGenericUIInstructionStepInstructionModel, error) {
	output := make([]DataConnectorGenericUIInstructionStepInstructionModel, 0)
	if input != nil {
		for _, v := range *input {
			o := DataConnectorGenericUIInstructionStepInstructionModel{
				Type: string(v.Type),
			}
			p, err := pluginsdk.FlattenJsonToString(v.Parameters.(map[string]interface{}))
			if err != nil {
				return nil, err
			}
			o.Parameters = p
			output = append(output, o)
		}
	}
	return output, nil
}

type DataConnectorGenericUIResource struct{}

func DataConnectorGenericUIConfigSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"title": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"publisher": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"description_markdown": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"custom_image": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},
		"graph_queries_table_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringMatch(regexp.MustCompile(`.+_CL$`), "must end with _CL"),
		},
		"graph_query": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"metric_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"legend": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"base_query": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},
		"sample_query": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"description": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"query": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},
		"data_type": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"last_data_received_query": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},
		"connectivity_criteria": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							"IsConnectedQuery",
						}, false),
					},
					"value": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
		},
		"availability": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
					"preview": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
				},
			},
		},
		"permission": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"resource_provider": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(securityinsight.ProviderNameMicrosoftAuthorizationpolicyAssignments),
										string(securityinsight.ProviderNameMicrosoftOperationalInsightssolutions),
										string(securityinsight.ProviderNameMicrosoftOperationalInsightsworkspaces),
										string(securityinsight.ProviderNameMicrosoftOperationalInsightsworkspacesdatasources),
										string(securityinsight.ProviderNameMicrosoftOperationalInsightsworkspacessharedKeys),
										string(securityinsight.ProviderNameMicrosoftaadiamdiagnosticSettings),
									}, false),
								},
								"display_name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"display_text": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"scope": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice(
										[]string{
											string(securityinsight.PermissionProviderScopeResourceGroup),
											string(securityinsight.PermissionProviderScopeSubscription),
											string(securityinsight.PermissionProviderScopeWorkspace),
										}, false),
								},
								"required_permissions": {
									Type:     pluginsdk.TypeList,
									Required: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*schema.Schema{
											"action": {
												Type:     pluginsdk.TypeBool,
												Optional: true,
											},
											"write": {
												Type:     pluginsdk.TypeBool,
												Optional: true,
											},
											"read": {
												Type:     pluginsdk.TypeBool,
												Optional: true,
											},
											"delete": {
												Type:     pluginsdk.TypeBool,
												Optional: true,
											},
										},
									},
								},
							},
						},
					},
					"custom": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"description": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},
				},
			},
		},
		"instruction": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"title": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"description": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"step": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"type": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(securityinsight.SettingTypeCopyableLabel),
										string(securityinsight.SettingTypeInfoMessage),
										string(securityinsight.SettingTypeInstructionStepsGroup),
									}, false),
								},
								"parameters": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsJSON,
								},
							},
						},
					},
				},
			},
		},
	}
}
