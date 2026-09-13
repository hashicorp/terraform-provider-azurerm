package subscriptions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListLocationsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Location
}

type ListLocationsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Location
}

type ListLocationsOperationOptions struct {
	IncludeExtendedLocations *bool
}

func DefaultListLocationsOperationOptions() ListLocationsOperationOptions {
	return ListLocationsOperationOptions{}
}

func (o ListLocationsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListLocationsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListLocationsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.IncludeExtendedLocations != nil {
		out.Append("includeExtendedLocations", fmt.Sprintf("%v", *o.IncludeExtendedLocations))
	}
	return &out
}

type ListLocationsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListLocationsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListLocations ...
func (c SubscriptionsClient) ListLocations(ctx context.Context, id commonids.SubscriptionId, options ListLocationsOperationOptions) (result ListLocationsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListLocationsCustomPager{},
		Path:          fmt.Sprintf("%s/locations", id.ID()),
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
		Values *[]Location `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListLocationsComplete retrieves all the results into a single object
func (c SubscriptionsClient) ListLocationsComplete(ctx context.Context, id commonids.SubscriptionId, options ListLocationsOperationOptions) (ListLocationsCompleteResult, error) {
	return c.ListLocationsCompleteMatchingPredicate(ctx, id, options, LocationOperationPredicate{})
}

// ListLocationsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c SubscriptionsClient) ListLocationsCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, options ListLocationsOperationOptions, predicate LocationOperationPredicate) (result ListLocationsCompleteResult, err error) {
	items := make([]Location, 0)

	resp, err := c.ListLocations(ctx, id, options)
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

	result = ListLocationsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
