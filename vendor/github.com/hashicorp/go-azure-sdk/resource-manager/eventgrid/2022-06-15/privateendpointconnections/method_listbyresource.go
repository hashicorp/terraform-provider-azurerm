package privateendpointconnections

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

type ListByResourceOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]PrivateEndpointConnection
}

type ListByResourceCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []PrivateEndpointConnection
}

type ListByResourceOperationOptions struct {
	Filter *string
	Top    *int64
}

func DefaultListByResourceOperationOptions() ListByResourceOperationOptions {
	return ListByResourceOperationOptions{}
}

func (o ListByResourceOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByResourceOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ListByResourceOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// ListByResource ...
func (c PrivateEndpointConnectionsClient) ListByResource(ctx context.Context, id commonids.ScopeId, options ListByResourceOperationOptions) (result ListByResourceOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/privateEndpointConnections", id.ID()),
		OptionsObject: options,
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
		Values *[]PrivateEndpointConnection `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByResourceComplete retrieves all the results into a single object
func (c PrivateEndpointConnectionsClient) ListByResourceComplete(ctx context.Context, id commonids.ScopeId, options ListByResourceOperationOptions) (ListByResourceCompleteResult, error) {
	return c.ListByResourceCompleteMatchingPredicate(ctx, id, options, PrivateEndpointConnectionOperationPredicate{})
}

// ListByResourceCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c PrivateEndpointConnectionsClient) ListByResourceCompleteMatchingPredicate(ctx context.Context, id commonids.ScopeId, options ListByResourceOperationOptions, predicate PrivateEndpointConnectionOperationPredicate) (result ListByResourceCompleteResult, err error) {
	items := make([]PrivateEndpointConnection, 0)

	resp, err := c.ListByResource(ctx, id, options)
	if err != nil {
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

	result = ListByResourceCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
