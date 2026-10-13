// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application_insights_standard_web_test

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2022-06-15/webtestsapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	tfbase64 "github.com/hashicorp/terraform-provider-azurerm/internal/tf/base64"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppInsights.StandardWebTestsClient

			subscriptionId := metadata.Client.Account.SubscriptionId

			var model ApplicationInsightsStandardWebTestResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id := webtestsapis.NewWebTestID(subscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.WebTestsGet(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			validations := expandApplicationInsightsStandardWebTestValidations(model.ValidationRules)

			appInsightsId, err := webtestsapis.ParseComponentID(model.ApplicationInsightsID)
			if err != nil {
				return err
			}

			if model.Tags == nil {
				model.Tags = make(map[string]string)
			}

			model.Tags[fmt.Sprintf("hidden-link:%s", appInsightsId.ID())] = "Resource"

			props := webtestsapis.WebTestProperties{
				Name:               id.WebTestName, // API requires this to be specified despite ARM spec guidance that it should come from the ID
				Enabled:            pointer.To(model.Enabled),
				Frequency:          pointer.To(model.Frequency),
				Kind:               webtestsapis.WebTestKindStandard,
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

			param := webtestsapis.WebTest{
				Kind:       pointer.To(webtestsapis.WebTestKindStandard),
				Location:   location.Normalize(model.Location),
				Properties: &props,
				Tags:       pointer.To(model.Tags),
			}

			if _, err := client.WebTestsCreateOrUpdate(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}
			return nil
		},
	}
}

func expandApplicationInsightsStandardWebTestRequest(input []RequestModel) (request *webtestsapis.WebTestPropertiesRequest) {
	if len(input) == 0 {
		return nil
	}
	requestInput := input[0]

	request = &webtestsapis.WebTestPropertiesRequest{
		FollowRedirects:        pointer.To(requestInput.FollowRedirects),
		HTTPVerb:               pointer.To(requestInput.HTTPVerb),
		ParseDependentRequests: pointer.To(requestInput.ParseDependentRequests),
	}

	request.Headers = expandApplicationInsightsStandardWebTestRequestHeaders(requestInput.Header)

	if v := requestInput.Body; v != "" {
		request.RequestBody = pointer.To(tfbase64.EncodeIfNot(v))
	}

	if v := requestInput.URL; v != "" {
		request.RequestURL = pointer.To(v)
	}

	return request
}

func expandApplicationInsightsStandardWebTestRequestHeaders(input []HeaderModel) *[]webtestsapis.HeaderField {
	if len(input) == 0 {
		return nil
	}

	headers := make([]webtestsapis.HeaderField, 0)

	for _, v := range input {
		h := webtestsapis.HeaderField{
			Key:   pointer.To(v.Name),
			Value: pointer.To(v.Value),
		}
		headers = append(headers, h)
	}

	return &headers
}

func expandApplicationInsightsStandardWebTestValidations(input []ValidationRuleModel) webtestsapis.WebTestPropertiesValidationRules {
	rules := webtestsapis.WebTestPropertiesValidationRules{
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

func expandApplicationInsightsStandardWebTestContentValidations(input []ContentModel) *webtestsapis.WebTestPropertiesValidationRulesContentValidation {
	if len(input) == 0 {
		return nil
	}

	contentInput := input[0]

	content := webtestsapis.WebTestPropertiesValidationRulesContentValidation{
		ContentMatch:    pointer.To(contentInput.ContentMatch),
		IgnoreCase:      pointer.To(contentInput.IgnoreCase),
		PassIfTextFound: pointer.To(contentInput.PassIfTextFound),
	}

	return &content
}

func expandApplicationInsightsStandardWebTestGeoLocations(input []string) []webtestsapis.WebTestGeolocation {
	if len(input) == 0 {
		return []webtestsapis.WebTestGeolocation{}
	}

	locations := make([]webtestsapis.WebTestGeolocation, 0)

	for _, v := range input {
		loc := webtestsapis.WebTestGeolocation{
			Id: pointer.To(v),
		}
		locations = append(locations, loc)
	}

	return locations
}
