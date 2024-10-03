package componentfeaturesandpricingapis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ComponentCurrentBillingFeaturesUpdateOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *ApplicationInsightsComponentBillingFeatures
}

// ComponentCurrentBillingFeaturesUpdate ...
func (c ComponentFeaturesAndPricingAPIsClient) ComponentCurrentBillingFeaturesUpdate(ctx context.Context, id ComponentId, input ApplicationInsightsComponentBillingFeatures) (result ComponentCurrentBillingFeaturesUpdateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPut,
		Path:       fmt.Sprintf("%s/currentbillingfeatures", id.ID()),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	if err = req.Marshal(input); err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	var model ApplicationInsightsComponentBillingFeatures
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
