// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application_insights_standard_web_test

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2022-06-15/webtestsapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppInsights.StandardWebTestsClient
			id, err := webtestsapis.ParseWebTestID(metadata.ResourceData.Id())
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

			// API requires that the `Locations` property is always set, even if it is an empty list
			props.Locations = expandApplicationInsightsStandardWebTestGeoLocations(model.GeoLocations)

			if metadata.ResourceData.HasChange("request") {
				props.Request = expandApplicationInsightsStandardWebTestRequest(model.Request)
			}

			if metadata.ResourceData.HasChange("validation_rules") {
				props.ValidationRules = pointer.To(expandApplicationInsightsStandardWebTestValidations(model.ValidationRules))
			}

			existing.Model.Properties = &props

			appInsightsId, err := webtestsapis.ParseComponentID(metadata.ResourceData.Get("application_insights_id").(string))
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
