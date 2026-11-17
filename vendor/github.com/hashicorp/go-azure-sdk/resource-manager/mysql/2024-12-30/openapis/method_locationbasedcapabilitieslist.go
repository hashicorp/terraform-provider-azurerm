package openapis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LocationBasedCapabilitiesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]CapabilityProperties
}

type LocationBasedCapabilitiesListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []CapabilityProperties
}

type LocationBasedCapabilitiesListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *LocationBasedCapabilitiesListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// LocationBasedCapabilitiesList ...
func (c OpenapisClient) LocationBasedCapabilitiesList(ctx context.Context, id LocationId) (result LocationBasedCapabilitiesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &LocationBasedCapabilitiesListCustomPager{},
		Path:       fmt.Sprintf("%s/capabilities", id.ID()),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.ExecutePaged(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	var values struct {
		Values *[]CapabilityProperties `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// LocationBasedCapabilitiesListComplete retrieves all the results into a single object
func (c OpenapisClient) LocationBasedCapabilitiesListComplete(ctx context.Context, id LocationId) (LocationBasedCapabilitiesListCompleteResult, error) {
	return c.LocationBasedCapabilitiesListCompleteMatchingPredicate(ctx, id, CapabilityPropertiesOperationPredicate{})
}

// LocationBasedCapabilitiesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) LocationBasedCapabilitiesListCompleteMatchingPredicate(ctx context.Context, id LocationId, predicate CapabilityPropertiesOperationPredicate) (result LocationBasedCapabilitiesListCompleteResult, err error) {
	items := make([]CapabilityProperties, 0)

	resp, err := c.LocationBasedCapabilitiesList(ctx, id)
	if err != nil {
		result.LatestHttpResponse = resp.HttpResponse
		err = fmt.Errorf("loading results: %+v", err)
		return
	}
	if resp.Model != nil {
		for _, v := range *resp.Model {
			if predicate.Matches(v) {
				items = append(items, v)
			}
		}
	}

	result = LocationBasedCapabilitiesListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
