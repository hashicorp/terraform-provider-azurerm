package streamingpoliciesandstreaminglocators

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StreamingPoliciesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]StreamingPolicy
}

type StreamingPoliciesListCompleteResult struct {
	Items []StreamingPolicy
}

type StreamingPoliciesListOperationOptions struct {
	Filter  *string
	Orderby *string
	Top     *int64
}

func DefaultStreamingPoliciesListOperationOptions() StreamingPoliciesListOperationOptions {
	return StreamingPoliciesListOperationOptions{}
}

func (o StreamingPoliciesListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o StreamingPoliciesListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o StreamingPoliciesListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Orderby != nil {
		out.Append("$orderby", fmt.Sprintf("%v", *o.Orderby))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// StreamingPoliciesList ...
func (c StreamingPoliciesAndStreamingLocatorsClient) StreamingPoliciesList(ctx context.Context, id MediaServiceId, options StreamingPoliciesListOperationOptions) (result StreamingPoliciesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/streamingPolicies", id.ID()),
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
		Values *[]StreamingPolicy `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// StreamingPoliciesListComplete retrieves all the results into a single object
func (c StreamingPoliciesAndStreamingLocatorsClient) StreamingPoliciesListComplete(ctx context.Context, id MediaServiceId, options StreamingPoliciesListOperationOptions) (StreamingPoliciesListCompleteResult, error) {
	return c.StreamingPoliciesListCompleteMatchingPredicate(ctx, id, options, StreamingPolicyOperationPredicate{})
}

// StreamingPoliciesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c StreamingPoliciesAndStreamingLocatorsClient) StreamingPoliciesListCompleteMatchingPredicate(ctx context.Context, id MediaServiceId, options StreamingPoliciesListOperationOptions, predicate StreamingPolicyOperationPredicate) (result StreamingPoliciesListCompleteResult, err error) {
	items := make([]StreamingPolicy, 0)

	resp, err := c.StreamingPoliciesList(ctx, id, options)
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

	result = StreamingPoliciesListCompleteResult{
		Items: items,
	}
	return
}
