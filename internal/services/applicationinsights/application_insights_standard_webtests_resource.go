// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package applicationinsights

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	webtests "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2022-06-15/webtestsapis"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var _ sdk.Resource = ApplicationInsightsStandardWebTestResource{}

type ApplicationInsightsStandardWebTestResource struct{}

func (ApplicationInsightsStandardWebTestResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.NoZeroValues,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"application_insights_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ComponentID,
		},

		"location": commonschema.Location(),

		"frequency": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Default:  300,
			ValidateFunc: validation.IntInSlice([]int{
				300,
				600,
				900,
			}),
		},

		"timeout": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Default:  30,
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"retry_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"request": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"follow_redirects_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"http_verb": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "GET",
						ValidateFunc: validation.StringInSlice([]string{
							"GET", "POST", "PUT", "PATCH", "DELETE",
						}, false),
					},

					"parse_dependent_requests_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"header": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"value": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
							},
						},
					},

					"body": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"url": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"validation_rules": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"expected_status_code": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  200,
					},

					// Typo in API spec, issue: https://github.com/Azure/azure-rest-api-specs/issues/22136
					// "ignore_status_code": {
					// 	Type:     pluginsdk.TypeBool,
					// 	Optional: true,
					// 	Default:  false,
					// },

					"ssl_cert_remaining_lifetime": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(1, 365),
					},

					"ssl_check_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"content": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"content_match": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"ignore_case": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},
								"pass_if_text_found": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},
							},
						},
					},
				},
			},
		},

		"geo_locations": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type:             pluginsdk.TypeString,
				ValidateFunc:     validation.StringIsNotEmpty,
				StateFunc:        location.StateFunc,
				DiffSuppressFunc: location.DiffSuppressFunc,
			},
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"tags": commonschema.Tags(),
	}
}

func (ApplicationInsightsStandardWebTestResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"synthetic_monitor_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (ApplicationInsightsStandardWebTestResource) ModelObject() interface{} {
	return nil
}

func (ApplicationInsightsStandardWebTestResource) ResourceType() string {
	return "azurerm_application_insights_standard_web_test"
}

func (r ApplicationInsightsStandardWebTestResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppInsights.StandardWebTestsClient

			subscriptionId := metadata.Client.Account.SubscriptionId
			name := metadata.ResourceData.Get("name").(string)
			resourceGroupName := metadata.ResourceData.Get("resource_group_name").(string)
			id := webtests.NewWebTestID(subscriptionId, resourceGroupName, name)

			existing, err := client.WebTestsGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			location := location.Normalize(metadata.ResourceData.Get("location").(string))
			description := metadata.ResourceData.Get("description").(string)
			frequency := int64(metadata.ResourceData.Get("frequency").(int))
			timeout := int64(metadata.ResourceData.Get("timeout").(int))
			isEnabled := metadata.ResourceData.Get("enabled").(bool)
			retryEnabled := metadata.ResourceData.Get("retry_enabled").(bool)
			geoLocationsRaw := metadata.ResourceData.Get("geo_locations").([]interface{})
			geoLocations := expandApplicationInsightsStandardWebTestGeoLocations(geoLocationsRaw)

			requestRaw := metadata.ResourceData.Get("request").([]interface{})
			request, isHttps := expandApplicationInsightsStandardWebTestRequest(requestRaw)

			validationsRaw := metadata.ResourceData.Get("validation_rules").([]interface{})
			validations := expandApplicationInsightsStandardWebTestValidations(validationsRaw, isHttps)

			appInsightsId, err := webtests.ParseComponentID(metadata.ResourceData.Get("application_insights_id").(string))
			if err != nil {
				return err
			}
			t := metadata.ResourceData.Get("tags").(map[string]interface{})
			tagKey := fmt.Sprintf("hidden-link:%s", appInsightsId.ID())
			t[tagKey] = "Resource"

			param := webtests.WebTest{
				Kind:     utils.ToPtr(webtests.WebTestKindStandard),
				Location: location,
				Properties: &webtests.WebTestProperties{
					Description:        utils.String(description),
					Enabled:            utils.Bool(isEnabled),
					Frequency:          utils.Int64(frequency),
					Kind:               webtests.WebTestKindStandard,
					Locations:          geoLocations,
					Name:               id.WebTestName,
					RetryEnabled:       utils.Bool(retryEnabled),
					SyntheticMonitorId: id.WebTestName,
					Timeout:            utils.Int64(timeout),
					Request:            &request,
					ValidationRules:    &validations,
				},
				Tags: tags.Expand(t),
			}
			if _, err := client.WebTestsCreateOrUpdate(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ApplicationInsightsStandardWebTestResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppInsights.StandardWebTestsClient
			id, err := webtests.ParseWebTestID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			location := location.Normalize(metadata.ResourceData.Get("location").(string))
			description := metadata.ResourceData.Get("description").(string)
			frequency := int64(metadata.ResourceData.Get("frequency").(int))
			timeout := int64(metadata.ResourceData.Get("timeout").(int))
			isEnabled := metadata.ResourceData.Get("enabled").(bool)
			retryEnabled := metadata.ResourceData.Get("retry_enabled").(bool)

			geoLocationsRaw := metadata.ResourceData.Get("geo_locations").([]interface{})
			geoLocations := expandApplicationInsightsStandardWebTestGeoLocations(geoLocationsRaw)

			requestRaw := metadata.ResourceData.Get("request").([]interface{})
			request, isHttps := expandApplicationInsightsStandardWebTestRequest(requestRaw)

			validationsRaw := metadata.ResourceData.Get("validation_rules").([]interface{})
			validations := expandApplicationInsightsStandardWebTestValidations(validationsRaw, isHttps)

			appInsightsId, err := webtests.ParseComponentID(metadata.ResourceData.Get("application_insights_id").(string))
			if err != nil {
				return err
			}
			t := metadata.ResourceData.Get("tags").(map[string]interface{})
			tagKey := fmt.Sprintf("hidden-link:%s", appInsightsId.ID())
			t[tagKey] = "Resource"

			param := webtests.WebTest{
				Kind:     utils.ToPtr(webtests.WebTestKindStandard),
				Location: location,
				Properties: &webtests.WebTestProperties{
					Description:        utils.String(description),
					Enabled:            utils.Bool(isEnabled),
					Frequency:          utils.Int64(frequency),
					Kind:               webtests.WebTestKindStandard,
					Locations:          geoLocations,
					Name:               id.WebTestName,
					RetryEnabled:       utils.Bool(retryEnabled),
					SyntheticMonitorId: id.WebTestName,
					Timeout:            utils.Int64(timeout),
					Request:            &request,
					ValidationRules:    &validations,
				},
				Tags: tags.Expand(t),
			}
			if _, err := client.WebTestsCreateOrUpdate(ctx, *id, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (ApplicationInsightsStandardWebTestResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppInsights.StandardWebTestsClient

			id, err := webtests.ParseWebTestID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.WebTestsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if model := resp.Model; model != nil {
				appInsightsId := ""
				if model.Tags != nil {
					for i := range *model.Tags {
						if strings.HasPrefix(i, "hidden-link") {
							appInsightsId = strings.Split(i, ":")[1]
						}
					}
				}

				parsedAppInsightsId, err := webtests.ParseComponentIDInsensitively(appInsightsId)
				if err != nil {
					return fmt.Errorf("parsing `application_insights_id`: %+v", err)
				}

				metadata.ResourceData.Set("application_insights_id", parsedAppInsightsId.ID())
				metadata.ResourceData.Set("name", id.WebTestName)
				metadata.ResourceData.Set("resource_group_name", id.ResourceGroupName)
				metadata.ResourceData.Set("location", location.NormalizeNilable(&model.Location))
				if props := model.Properties; props != nil {
					metadata.ResourceData.Set("synthetic_monitor_id", props.SyntheticMonitorId)
					metadata.ResourceData.Set("description", props.Description)
					metadata.ResourceData.Set("enabled", props.Enabled)
					metadata.ResourceData.Set("frequency", props.Frequency)
					metadata.ResourceData.Set("timeout", props.Timeout)
					metadata.ResourceData.Set("retry_enabled", props.RetryEnabled)
					if props.Request != nil {
						request, err := flattenApplicationInsightsStandardWebTestRequest(*props.Request)
						if err != nil {
							return fmt.Errorf("setting `request`: %+v", err)
						}
						metadata.ResourceData.Set("request", request)
					}
					if props.ValidationRules != nil {
						rules := flattenApplicationInsightsStandardWebTestValidations(*props.ValidationRules)
						metadata.ResourceData.Set("validation_rules", rules)
					}

					if err := metadata.ResourceData.Set("geo_locations", flattenApplicationInsightsStandardWebTestGeoLocations(props.Locations)); err != nil {
						return fmt.Errorf("setting `geo_locations`: %+v", err)
					}
				}
				return tags.FlattenAndSet(metadata.ResourceData, model.Tags)
			}

			return nil
		},
	}
}

func (ApplicationInsightsStandardWebTestResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppInsights.StandardWebTestsClient

			id, err := webtests.ParseWebTestID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			_, err = client.WebTestsDelete(ctx, *id)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (ApplicationInsightsStandardWebTestResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return webtests.ValidateWebTestID
}

func expandApplicationInsightsStandardWebTestRequest(input []interface{}) (webtests.WebTestPropertiesRequest, bool) {
	requestInput := input[0].(map[string]interface{})

	followRedirects := true
	if v, ok := requestInput["follow_redirects_enabled"].(bool); ok {
		followRedirects = v
	}
	httpVerb := "GET"
	if v, ok := requestInput["http_verb"].(string); ok {
		httpVerb = v
	}
	parseDependentRequests := true
	if v, ok := requestInput["parse_dependent_requests_enabled"].(bool); ok {
		parseDependentRequests = v
	}

	request := webtests.WebTestPropertiesRequest{
		FollowRedirects:        utils.Bool(followRedirects),
		HTTPVerb:               utils.String(httpVerb),
		ParseDependentRequests: utils.Bool(parseDependentRequests),
	}

	request.Headers = expandApplicationInsightsStandardWebTestRequestHeaders(requestInput["header"].([]interface{}))

	if v, ok := requestInput["body"].(string); ok && v != "" {
		request.RequestBody = utils.String(utils.Base64EncodeIfNot(v))
	}
	isHttps := true
	if v, ok := requestInput["url"].(string); ok {
		request.RequestUrl = utils.String(v)
		isHttps = strings.HasPrefix(v, "https://")
	}

	return request, isHttps
}

func expandApplicationInsightsStandardWebTestRequestHeaders(input []interface{}) *[]webtests.HeaderField {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	headers := make([]webtests.HeaderField, len(input))

	for i, v := range input {
		header := v.(map[string]interface{})
		headers[i] = webtests.HeaderField{
			Key:   utils.String(header["name"].(string)),
			Value: utils.String(header["value"].(string)),
		}
	}

	return &headers
}

func flattenApplicationInsightsStandardWebTestRequest(req webtests.WebTestPropertiesRequest) ([]interface{}, error) {
	result := make(map[string]interface{})

	followRedirects := true
	if req.FollowRedirects != nil {
		followRedirects = *req.FollowRedirects
	}
	result["follow_redirects_enabled"] = followRedirects

	httpVerb := "GET"
	if req.HTTPVerb != nil {
		httpVerb = *req.HTTPVerb
	}
	result["http_verb"] = httpVerb

	parseDependentRequests := true
	if req.ParseDependentRequests != nil {
		parseDependentRequests = *req.ParseDependentRequests
	}
	result["parse_dependent_requests_enabled"] = parseDependentRequests

	if req.RequestUrl != nil {
		result["url"] = *req.RequestUrl
	}
	if req.RequestBody != nil {
		body, err := base64.StdEncoding.DecodeString(*req.RequestBody)
		if err != nil {
			return nil, fmt.Errorf("decoding `body`: %+v", err)
		}
		result["body"] = string(body)
	}
	if req.Headers != nil {
		result["header"] = flattenApplicationInsightsStandardWebTestRequestHeaders(req.Headers)
	}

	return []interface{}{result}, nil
}

func flattenApplicationInsightsStandardWebTestRequestHeaders(input *[]webtests.HeaderField) []interface{} {
	result := make([]interface{}, 0)
	if input == nil {
		return result
	}

	headers := *input
	if len(headers) == 0 {
		return result
	}

	for _, v := range headers {
		header := make(map[string]string, 2)
		header["name"] = *v.Key
		header["value"] = *v.Value
		result = append(result, header)
	}

	return result
}

func flattenApplicationInsightsStandardWebTestValidations(rules webtests.WebTestPropertiesValidationRules) []interface{} {
	result := make(map[string]interface{})

	if rules.ExpectedHTTPStatusCode != nil {
		result["expected_status_code"] = *rules.ExpectedHTTPStatusCode
	}
	// if rules.IgnoreHTTPSStatusCode != nil {
	// 	result["ignore_status_code"] = *rules.IgnoreHTTPSStatusCode
	// }
	if rules.SSLCertRemainingLifetimeCheck != nil {
		result["ssl_cert_remaining_lifetime"] = *rules.SSLCertRemainingLifetimeCheck
	}
	if rules.SSLCheck != nil {
		result["ssl_check_enabled"] = *rules.SSLCheck
	}

	if rules.ContentValidation != nil {
		result["content"] = flattenApplicationInsightsStandardWebTestContentValidations(rules.ContentValidation)
	}

	return []interface{}{result}
}

func flattenApplicationInsightsStandardWebTestContentValidations(input *webtests.WebTestPropertiesValidationRulesContentValidation) []interface{} {
	result := make(map[string]interface{})

	if input.ContentMatch != nil {
		result["content_match"] = *input.ContentMatch
	}
	if input.IgnoreCase != nil {
		result["ignore_case"] = *input.IgnoreCase
	}
	if input.PassIfTextFound != nil {
		result["pass_if_text_found"] = *input.PassIfTextFound
	}

	return []interface{}{result}
}

func expandApplicationInsightsStandardWebTestValidations(input []interface{}, isHttps bool) webtests.WebTestPropertiesValidationRules {
	rules := webtests.WebTestPropertiesValidationRules{
		ExpectedHTTPStatusCode: utils.Int64(200),
		// IgnoreHTTPSStatusCode:  utils.Bool(false),
		SSLCheck: utils.Bool(false),
	}
	if len(input) == 0 {
		return rules
	}

	validationsInput := input[0].(map[string]interface{})
	if v, ok := validationsInput["expected_status_code"].(int); ok {
		rules.ExpectedHTTPStatusCode = utils.Int64(int64(v))
	}
	// if v, ok := validationsInput["ignore_status_code"].(bool); ok {
	// 	rules.IgnoreHTTPSStatusCode = utils.Bool(v)
	// }

	// if URL http, sslCheck cannot be enabled
	sslCheckEnabled := false
	if v, ok := validationsInput["ssl_check_enabled"].(bool); ok && isHttps {
		rules.SSLCheck = utils.Bool(v)
		sslCheckEnabled = true
	}
	// if sslCheck not enabled, SSLCertRemainingLifetimeCheck cannot be enabled
	if v, ok := validationsInput["ssl_cert_remaining_lifetime"].(int); ok && v != 0 && sslCheckEnabled {
		rules.SSLCertRemainingLifetimeCheck = utils.Int64(int64(v))
	}
	if contentValidation, ok := validationsInput["content"].([]interface{}); ok {
		rules.ContentValidation = expandApplicationInsightsStandardWebTestContentValidations(contentValidation)
	}

	return rules
}

func expandApplicationInsightsStandardWebTestContentValidations(input []interface{}) *webtests.WebTestPropertiesValidationRulesContentValidation {
	content := webtests.WebTestPropertiesValidationRulesContentValidation{}
	if len(input) == 0 {
		return nil
	}

	contentInput := input[0].(map[string]interface{})
	content.ContentMatch = utils.String(contentInput["content_match"].(string))
	if v, ok := contentInput["ignore_case"].(bool); ok {
		content.IgnoreCase = utils.Bool(v)
	}
	if v, ok := contentInput["pass_if_text_found"].(bool); ok {
		content.PassIfTextFound = utils.Bool(v)
	}

	return &content
}

func expandApplicationInsightsStandardWebTestGeoLocations(input []interface{}) []webtests.WebTestGeolocation {
	locations := make([]webtests.WebTestGeolocation, 0)

	for _, v := range input {
		lc := v.(string)
		loc := webtests.WebTestGeolocation{
			Id: &lc,
		}
		locations = append(locations, loc)
	}

	return locations
}

func flattenApplicationInsightsStandardWebTestGeoLocations(input []webtests.WebTestGeolocation) []string {
	results := make([]string, 0)
	if input == nil {
		return results
	}

	for _, prop := range input {
		if prop.Id != nil {
			results = append(results, azure.NormalizeLocation(*prop.Id))
		}
	}

	return results
}
