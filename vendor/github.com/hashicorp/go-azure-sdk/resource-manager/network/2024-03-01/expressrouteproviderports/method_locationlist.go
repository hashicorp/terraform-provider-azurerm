package expressrouteproviderports

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

type LocationListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ExpressRouteProviderPort
}

type LocationListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ExpressRouteProviderPort
}

type LocationListOperationOptions struct {
	Filter *string
}

func DefaultLocationListOperationOptions() LocationListOperationOptions {
	return LocationListOperationOptions{}
}

func (o LocationListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o LocationListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o LocationListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

type LocationListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *LocationListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// LocationList ...
func (c ExpressRouteProviderPortsClient) LocationList(ctx context.Context, id commonids.SubscriptionId, options LocationListOperationOptions) (result LocationListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &LocationListCustomPager{},
		Path:          fmt.Sprintf("%s/providers/Microsoft.Network/expressRouteProviderPorts", id.ID()),
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
		Values *[]ExpressRouteProviderPort `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// LocationListComplete retrieves all the results into a single object
func (c ExpressRouteProviderPortsClient) LocationListComplete(ctx context.Context, id commonids.SubscriptionId, options LocationListOperationOptions) (LocationListCompleteResult, error) {
	return c.LocationListCompleteMatchingPredicate(ctx, id, options, ExpressRouteProviderPortOperationPredicate{})
}

// LocationListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ExpressRouteProviderPortsClient) LocationListCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, options LocationListOperationOptions, predicate ExpressRouteProviderPortOperationPredicate) (result LocationListCompleteResult, err error) {
	items := make([]ExpressRouteProviderPort, 0)

	resp, err := c.LocationList(ctx, id, options)
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

	result = LocationListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
