package endpoints

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DigitalTwinsEndpointListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DigitalTwinsEndpointResource
}

type DigitalTwinsEndpointListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DigitalTwinsEndpointResource
}

type DigitalTwinsEndpointListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *DigitalTwinsEndpointListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// DigitalTwinsEndpointList ...
func (c EndpointsClient) DigitalTwinsEndpointList(ctx context.Context, id DigitalTwinsInstanceId) (result DigitalTwinsEndpointListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &DigitalTwinsEndpointListCustomPager{},
		Path:       fmt.Sprintf("%s/endpoints", id.ID()),
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
		Values *[]DigitalTwinsEndpointResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// DigitalTwinsEndpointListComplete retrieves all the results into a single object
func (c EndpointsClient) DigitalTwinsEndpointListComplete(ctx context.Context, id DigitalTwinsInstanceId) (DigitalTwinsEndpointListCompleteResult, error) {
	return c.DigitalTwinsEndpointListCompleteMatchingPredicate(ctx, id, DigitalTwinsEndpointResourceOperationPredicate{})
}

// DigitalTwinsEndpointListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c EndpointsClient) DigitalTwinsEndpointListCompleteMatchingPredicate(ctx context.Context, id DigitalTwinsInstanceId, predicate DigitalTwinsEndpointResourceOperationPredicate) (result DigitalTwinsEndpointListCompleteResult, err error) {
	items := make([]DigitalTwinsEndpointResource, 0)

	resp, err := c.DigitalTwinsEndpointList(ctx, id)
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

	result = DigitalTwinsEndpointListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
