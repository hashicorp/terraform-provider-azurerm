package containerinstance

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LocationListCapabilitiesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Capabilities
}

type LocationListCapabilitiesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Capabilities
}

type LocationListCapabilitiesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *LocationListCapabilitiesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// LocationListCapabilities ...
func (c ContainerInstanceClient) LocationListCapabilities(ctx context.Context, id LocationId) (result LocationListCapabilitiesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &LocationListCapabilitiesCustomPager{},
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
		Values *[]Capabilities `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// LocationListCapabilitiesComplete retrieves all the results into a single object
func (c ContainerInstanceClient) LocationListCapabilitiesComplete(ctx context.Context, id LocationId) (LocationListCapabilitiesCompleteResult, error) {
	return c.LocationListCapabilitiesCompleteMatchingPredicate(ctx, id, CapabilitiesOperationPredicate{})
}

// LocationListCapabilitiesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ContainerInstanceClient) LocationListCapabilitiesCompleteMatchingPredicate(ctx context.Context, id LocationId, predicate CapabilitiesOperationPredicate) (result LocationListCapabilitiesCompleteResult, err error) {
	items := make([]Capabilities, 0)

	resp, err := c.LocationListCapabilities(ctx, id)
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

	result = LocationListCapabilitiesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
