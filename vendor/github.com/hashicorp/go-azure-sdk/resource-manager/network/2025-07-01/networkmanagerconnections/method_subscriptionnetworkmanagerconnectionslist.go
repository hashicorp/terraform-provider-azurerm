package networkmanagerconnections

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

type SubscriptionNetworkManagerConnectionsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]NetworkManagerConnection
}

type SubscriptionNetworkManagerConnectionsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []NetworkManagerConnection
}

type SubscriptionNetworkManagerConnectionsListOperationOptions struct {
	Top *int64
}

func DefaultSubscriptionNetworkManagerConnectionsListOperationOptions() SubscriptionNetworkManagerConnectionsListOperationOptions {
	return SubscriptionNetworkManagerConnectionsListOperationOptions{}
}

func (o SubscriptionNetworkManagerConnectionsListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o SubscriptionNetworkManagerConnectionsListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o SubscriptionNetworkManagerConnectionsListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type SubscriptionNetworkManagerConnectionsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *SubscriptionNetworkManagerConnectionsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// SubscriptionNetworkManagerConnectionsList ...
func (c NetworkManagerConnectionsClient) SubscriptionNetworkManagerConnectionsList(ctx context.Context, id commonids.SubscriptionId, options SubscriptionNetworkManagerConnectionsListOperationOptions) (result SubscriptionNetworkManagerConnectionsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &SubscriptionNetworkManagerConnectionsListCustomPager{},
		Path:          fmt.Sprintf("%s/providers/Microsoft.Network/networkManagerConnections", id.ID()),
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
		Values *[]NetworkManagerConnection `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// SubscriptionNetworkManagerConnectionsListComplete retrieves all the results into a single object
func (c NetworkManagerConnectionsClient) SubscriptionNetworkManagerConnectionsListComplete(ctx context.Context, id commonids.SubscriptionId, options SubscriptionNetworkManagerConnectionsListOperationOptions) (SubscriptionNetworkManagerConnectionsListCompleteResult, error) {
	return c.SubscriptionNetworkManagerConnectionsListCompleteMatchingPredicate(ctx, id, options, NetworkManagerConnectionOperationPredicate{})
}

// SubscriptionNetworkManagerConnectionsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c NetworkManagerConnectionsClient) SubscriptionNetworkManagerConnectionsListCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, options SubscriptionNetworkManagerConnectionsListOperationOptions, predicate NetworkManagerConnectionOperationPredicate) (result SubscriptionNetworkManagerConnectionsListCompleteResult, err error) {
	items := make([]NetworkManagerConnection, 0)

	resp, err := c.SubscriptionNetworkManagerConnectionsList(ctx, id, options)
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

	result = SubscriptionNetworkManagerConnectionsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
