// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application_insights_standard_web_test

import (
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2022-06-15/webtestsapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *webtestsapis.WebTestId, model *webtestsapis.WebTest) error {
	state := ApplicationInsightsStandardWebTestResourceModel{
		Name:              id.WebTestName,
		ResourceGroupName: id.ResourceGroupName,
	}

	if model != nil {
		tags := pointer.From(model.Tags)
		for i := range tags {
			if strings.HasPrefix(i, "hidden-link") {
				appInsightsId := strings.Split(i, ":")[1]

				parsedAppInsightsId, err := webtestsapis.ParseComponentIDInsensitively(appInsightsId)
				if err != nil {
					// there might be more than one hidden-link https://github.com/hashicorp/terraform-provider-azurerm/issues/27994
					log.Printf("[DEBUG] Error parsing hidden-link id: %+v", err)
					delete(tags, i)
					continue
				}
				state.ApplicationInsightsID = parsedAppInsightsId.ID()
				delete(tags, i)
			}
		}

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

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func flattenApplicationInsightsStandardWebTestRequest(input *webtestsapis.WebTestPropertiesRequest) ([]RequestModel, error) {
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

func flattenApplicationInsightsStandardWebTestRequestHeaders(input *[]webtestsapis.HeaderField) []HeaderModel {
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

func flattenApplicationInsightsStandardWebTestValidations(input *webtestsapis.WebTestPropertiesValidationRules) []ValidationRuleModel {
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

func flattenApplicationInsightsStandardWebTestContentValidations(input *webtestsapis.WebTestPropertiesValidationRulesContentValidation) []ContentModel {
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

func flattenApplicationInsightsStandardWebTestGeoLocations(input []webtestsapis.WebTestGeolocation) []string {
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
