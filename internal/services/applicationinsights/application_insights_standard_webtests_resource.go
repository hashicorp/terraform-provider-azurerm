// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package applicationinsights

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	components "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	webtests "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2022-06-15/webtestsapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var (
	_ sdk.ResourceWithUpdate        = ApplicationInsightsStandardWebTestResource{}
	_ sdk.ResourceWithCustomizeDiff = ApplicationInsightsStandardWebTestResource{}
)

type ApplicationInsightsStandardWebTestResource struct{}

type ApplicationInsightsStandardWebTestResourceModel struct {
	Name                  string                `tfschema:"name"`
	ResourceGroupName     string                `tfschema:"resource_group_name"`
	ApplicationInsightsID string                `tfschema:"application_insights_id"`
	Location              string                `tfschema:"location"`
	Frequency             int64                 `tfschema:"frequency"`
	Timeout               int64                 `tfschema:"timeout"`
	Enabled               bool                  `tfschema:"enabled"`
	Retry                 bool                  `tfschema:"retry_enabled"`
	Request               []RequestModel        `tfschema:"request"`
	ValidationRules       []ValidationRuleModel `tfschema:"validation_rules"`
	GeoLocations          []string              `tfschema:"geo_locations"`
	Description           string                `tfschema:"description"`
	Tags                  map[string]string     `tfschema:"tags"`

	// ComputedOnly
	SyntheticMonitorID string `tfschema:"synthetic_monitor_id"`
}

type RequestModel struct {
	FollowRedirects        bool          `tfschema:"follow_redirects_enabled"`
	HTTPVerb               string        `tfschema:"http_verb"`
	ParseDependentRequests bool          `tfschema:"parse_dependent_requests_enabled"`
	Header                 []HeaderModel `tfschema:"header"`
	Body                   string        `tfschema:"body"`
	URL                    string        `tfschema:"url"`
}

type ValidationRuleModel struct {
	ExpectedStatusCode           int64          `tfschema:"expected_status_code"`
	CertificateRemainingLifetime int64          `tfschema:"ssl_cert_remaining_lifetime"`
	SSLCheck                     bool           `tfschema:"ssl_check_enabled"`
	Content                      []ContentModel `tfschema:"content"`
}

type HeaderModel struct {
	Name  string `tfschema:"name"`
	Value string `tfschema:"value"`
}

type ContentModel struct {
	ContentMatch    string `tfschema:"content_match"`
	IgnoreCase      bool   `tfschema:"ignore_case"`
	PassIfTextFound bool   `tfschema:"pass_if_text_found"`
}

func (r ApplicationInsightsStandardWebTestResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff

			// SSLCheck conditions
			url, ok := rd.GetOk("request.0.url")
			if ok {
				if !strings.HasPrefix(strings.ToLower(url.(string)), "https://") {
					if v, ok := rd.GetOkExists("validation_rules.0.ssl_check_enabled"); ok && v.(bool) {
						return fmt.Errorf("cannot set ssl_check_enabled to true if request.0.url is not https")
					}
					if v, ok := rd.GetOkExists("validation_rules.0.ssl_cert_remaining_lifetime"); ok && v.(int) != 0 {
						return fmt.Errorf("cannot set ssl_cert_remaining_lifetime if request.0.url is not https")
					}
				}
			}

			return nil
		},
	}
}

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
			ValidateFunc: components.ValidateComponentID,
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
							"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS",
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
	return &ApplicationInsightsStandardWebTestResourceModel{}
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

			var model ApplicationInsightsStandardWebTestResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id := webtests.NewWebTestID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.WebTestsGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			validations := expandApplicationInsightsStandardWebTestValidations(model.ValidationRules)

			appInsightsId, err := webtests.ParseComponentID(model.ApplicationInsightsID)
			if err != nil {
				return err
			}

			if model.Tags == nil {
				model.Tags = make(map[string]string)
			}

			model.Tags[fmt.Sprintf("hidden-link:%s", appInsightsId.ID())] = "Resource"

			props := webtests.WebTestProperties{
				Name:               id.WebTestName, // API requires this to be specified despite ARM spec guidance that it should come from the ID
				Enabled:            pointer.To(model.Enabled),
				Frequency:          pointer.To(model.Frequency),
				Kind:               webtests.WebTestKindStandard,
				SyntheticMonitorId: id.WebTestName,
				RetryEnabled:       pointer.To(model.Retry),
				Timeout:            pointer.To(model.Timeout),
				Locations:          expandApplicationInsightsStandardWebTestGeoLocations(model.GeoLocations),
				ValidationRules:    pointer.To(validations),
				Request:            expandApplicationInsightsStandardWebTestRequest(model.Request),
			}

			if model.Description != "" {
				props.Description = pointer.To(model.Description)
			}

			param := webtests.WebTest{
				Kind:       pointer.To(webtests.WebTestKindStandard),
				Location:   location.Normalize(model.Location),
				Properties: &props,
				Tags:       pointer.To(model.Tags),
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

			model := ApplicationInsightsStandardWebTestResourceModel{
				Name:              id.WebTestName,
				ResourceGroupName: id.ResourceGroupName,
			}
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			existing, err := client.WebTestsGet(ctx, *id)
			if err != nil || existing.Model == nil {
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			props := pointer.From(existing.Model.Properties)

			if metadata.ResourceData.HasChange("description") {
				props.Description = pointer.To(model.Description)
			}

			if metadata.ResourceData.HasChange("frequency") {
				props.Frequency = pointer.To(model.Frequency)
			}

			if metadata.ResourceData.HasChange("timeout") {
				props.Timeout = pointer.To(model.Timeout)
			}

			props.Enabled = pointer.To(model.Enabled)
			props.RetryEnabled = pointer.To(model.Retry)

			// API requires that ths `Locations` property is always set, even if it is an empty list
			props.Locations = expandApplicationInsightsStandardWebTestGeoLocations(model.GeoLocations)

			if metadata.ResourceData.HasChange("request") {
				props.Request = expandApplicationInsightsStandardWebTestRequest(model.Request)
			}

			if metadata.ResourceData.HasChange("validation_rules") {
				props.ValidationRules = pointer.To(expandApplicationInsightsStandardWebTestValidations(model.ValidationRules))
			}

			existing.Model.Properties = &props

			appInsightsId, err := webtests.ParseComponentID(metadata.ResourceData.Get("application_insights_id").(string))
			if err != nil {
				return err
			}
			// Since we set the hidden tag, we always update them
			if model.Tags == nil {
				model.Tags = make(map[string]string)
			}
			t := model.Tags
			t[fmt.Sprintf("hidden-link:%s", appInsightsId.ID())] = "Resource"
			existing.Model.Tags = pointer.To(t)

			if _, err := client.WebTestsCreateOrUpdate(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

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

			state := ApplicationInsightsStandardWebTestResourceModel{
				Name:              id.WebTestName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				tags := pointer.From(model.Tags)
				appInsightsId := ""
				for i := range tags {
					if strings.HasPrefix(i, "hidden-link") {
						appInsightsId = strings.Split(i, ":")[1]
						delete(tags, i)
					}
				}

				parsedAppInsightsId, err := webtests.ParseComponentIDInsensitively(appInsightsId)
				if err != nil {
					return fmt.Errorf("parsing `application_insights_id` for %s: %+v", *id, err)
				}
				state.ApplicationInsightsID = parsedAppInsightsId.ID()
				state.Tags = tags
				state.Location = location.Normalize(model.Location)

				if props := model.Properties; props != nil {
					state.SyntheticMonitorID = props.SyntheticMonitorId
					state.Description = pointer.From(props.Description)
					state.Enabled = pointer.From(props.Enabled)
					state.Frequency = pointer.From(props.Frequency)
					state.Timeout = pointer.From(props.Timeout)
					state.Retry = pointer.From(props.RetryEnabled)
					req, err := flattenApplicationInsightsStandardWebTestRequest(props.Request)
					if err != nil {
						return fmt.Errorf("flattening request for %s: %+v", *id, err)
					}
					state.Request = req
					state.ValidationRules = flattenApplicationInsightsStandardWebTestValidations(props.ValidationRules)
					state.GeoLocations = flattenApplicationInsightsStandardWebTestGeoLocations(props.Locations)
				}
			}

			return metadata.Encode(&state)
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

func expandApplicationInsightsStandardWebTestRequest(input []RequestModel) (request *webtests.WebTestPropertiesRequest) {
	if len(input) == 0 {
		return nil
	}
	requestInput := input[0]

	request = &webtests.WebTestPropertiesRequest{
		FollowRedirects:        pointer.To(requestInput.FollowRedirects),
		HTTPVerb:               pointer.To(requestInput.HTTPVerb),
		ParseDependentRequests: pointer.To(requestInput.ParseDependentRequests),
	}

	request.Headers = expandApplicationInsightsStandardWebTestRequestHeaders(requestInput.Header)

	if v := requestInput.Body; v != "" {
		request.RequestBody = pointer.To(utils.Base64EncodeIfNot(v))
	}

	if v := requestInput.URL; v != "" {
		request.RequestURL = pointer.To(v)
	}

	return request
}

func expandApplicationInsightsStandardWebTestRequestHeaders(input []HeaderModel) *[]webtests.HeaderField {
	if len(input) == 0 {
		return nil
	}

	headers := make([]webtests.HeaderField, 0)

	for _, v := range input {
		h := webtests.HeaderField{
			Key:   utils.String(v.Name),
			Value: utils.String(v.Value),
		}
		headers = append(headers, h)
	}

	return &headers
}

func flattenApplicationInsightsStandardWebTestRequest(input *webtests.WebTestPropertiesRequest) ([]RequestModel, error) {
	if input == nil {
		return []RequestModel{}, nil
	}

	req := pointer.From(input)

	result := RequestModel{
		FollowRedirects:        pointer.From(req.FollowRedirects),
		HTTPVerb:               pointer.From(req.HTTPVerb),
		ParseDependentRequests: pointer.From(req.ParseDependentRequests),
		URL:                    pointer.From(req.RequestURL),
		Header:                 flattenApplicationInsightsStandardWebTestRequestHeaders(req.Headers),
	}

	if req.RequestBody != nil {
		body, err := base64.StdEncoding.DecodeString(pointer.From(req.RequestBody))
		if err != nil {
			return nil, err
		}
		result.Body = string(body)
	}

	return []RequestModel{result}, nil
}

func flattenApplicationInsightsStandardWebTestRequestHeaders(input *[]webtests.HeaderField) []HeaderModel {
	if input == nil || len(*input) == 0 {
		return []HeaderModel{}
	}

	result := make([]HeaderModel, 0)

	headers := *input

	for _, v := range headers {
		header := HeaderModel{
			Name:  pointer.From(v.Key),
			Value: pointer.From(v.Value),
		}
		result = append(result, header)
	}

	return result
}

func flattenApplicationInsightsStandardWebTestValidations(input *webtests.WebTestPropertiesValidationRules) []ValidationRuleModel {
	if input == nil {
		return []ValidationRuleModel{}
	}

	rules := pointer.From(input)

	// API Always returns this block, despite being a pointer as the `SSLCheck` property is always set. It is required, despite being marked as optional in the swagger
	if rules.ContentValidation == nil && rules.ExpectedHTTPStatusCode == nil && rules.IgnoreHTTPStatusCode == nil && rules.SSLCertRemainingLifetimeCheck == nil && (rules.SSLCheck == nil || !*rules.SSLCheck) {
		return []ValidationRuleModel{}
	}

	result := ValidationRuleModel{
		ExpectedStatusCode:           pointer.From(rules.ExpectedHTTPStatusCode),
		CertificateRemainingLifetime: pointer.From(rules.SSLCertRemainingLifetimeCheck),
		SSLCheck:                     pointer.From(rules.SSLCheck),
		Content:                      flattenApplicationInsightsStandardWebTestContentValidations(rules.ContentValidation),
	}

	return []ValidationRuleModel{result}
}

func flattenApplicationInsightsStandardWebTestContentValidations(input *webtests.WebTestPropertiesValidationRulesContentValidation) []ContentModel {
	if input == nil {
		return []ContentModel{}
	}

	result := ContentModel{
		ContentMatch:    pointer.From(input.ContentMatch),
		IgnoreCase:      pointer.From(input.IgnoreCase),
		PassIfTextFound: pointer.From(input.PassIfTextFound),
	}

	return []ContentModel{result}
}

func expandApplicationInsightsStandardWebTestValidations(input []ValidationRuleModel) webtests.WebTestPropertiesValidationRules {
	rules := webtests.WebTestPropertiesValidationRules{
		SSLCheck: pointer.To(false),
	}

	if len(input) == 0 {
		return rules
	}

	validationsInput := input[0]
	rules.ExpectedHTTPStatusCode = pointer.To(validationsInput.ExpectedStatusCode)

	// if URL http, sslCheck cannot be enabled - Catch in CustomiseDiff
	rules.SSLCheck = pointer.To(validationsInput.SSLCheck)
	// if sslCheck not enabled, SSLCertRemainingLifetimeCheck cannot be enabled
	if validationsInput.CertificateRemainingLifetime != 0 && validationsInput.SSLCheck {
		rules.SSLCertRemainingLifetimeCheck = pointer.To(validationsInput.CertificateRemainingLifetime)
	}
	rules.ContentValidation = expandApplicationInsightsStandardWebTestContentValidations(validationsInput.Content)

	return rules
}

func expandApplicationInsightsStandardWebTestContentValidations(input []ContentModel) *webtests.WebTestPropertiesValidationRulesContentValidation {
	if len(input) == 0 {
		return nil
	}

	contentInput := input[0]

	content := webtests.WebTestPropertiesValidationRulesContentValidation{
		ContentMatch:    pointer.To(contentInput.ContentMatch),
		IgnoreCase:      pointer.To(contentInput.IgnoreCase),
		PassIfTextFound: pointer.To(contentInput.PassIfTextFound),
	}

	return &content
}

func expandApplicationInsightsStandardWebTestGeoLocations(input []string) []webtests.WebTestGeolocation {
	if len(input) == 0 {
		return []webtests.WebTestGeolocation{}
	}

	locations := make([]webtests.WebTestGeolocation, 0)

	for _, v := range input {
		loc := webtests.WebTestGeolocation{
			Id: pointer.To(v),
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
			results = append(results, location.NormalizeNilable(prop.Id))
		}
	}

	return results
}
