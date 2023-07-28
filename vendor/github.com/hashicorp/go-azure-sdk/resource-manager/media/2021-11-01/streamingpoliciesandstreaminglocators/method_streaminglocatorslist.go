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

type StreamingLocatorsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]StreamingLocator
}

type StreamingLocatorsListCompleteResult struct {
	Items []StreamingLocator
}

type StreamingLocatorsListOperationOptions struct {
	Filter  *string
	Orderby *string
	Top     *int64
}

func DefaultStreamingLocatorsListOperationOptions() StreamingLocatorsListOperationOptions {
	return StreamingLocatorsListOperationOptions{}
}

func (o StreamingLocatorsListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o StreamingLocatorsListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o StreamingLocatorsListOperationOptions) ToQuery() *client.QueryParams {
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

// StreamingLocatorsList ...
func (c StreamingPoliciesAndStreamingLocatorsClient) StreamingLocatorsList(ctx context.Context, id MediaServiceId, options StreamingLocatorsListOperationOptions) (result StreamingLocatorsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/streamingLocators", id.ID()),
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
		Values *[]StreamingLocator `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// StreamingLocatorsListComplete retrieves all the results into a single object
func (c StreamingPoliciesAndStreamingLocatorsClient) StreamingLocatorsListComplete(ctx context.Context, id MediaServiceId, options StreamingLocatorsListOperationOptions) (StreamingLocatorsListCompleteResult, error) {
	return c.StreamingLocatorsListCompleteMatchingPredicate(ctx, id, options, StreamingLocatorOperationPredicate{})
}

// StreamingLocatorsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c StreamingPoliciesAndStreamingLocatorsClient) StreamingLocatorsListCompleteMatchingPredicate(ctx context.Context, id MediaServiceId, options StreamingLocatorsListOperationOptions, predicate StreamingLocatorOperationPredicate) (result StreamingLocatorsListCompleteResult, err error) {
	items := make([]StreamingLocator, 0)

	resp, err := c.StreamingLocatorsList(ctx, id, options)
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

	result = StreamingLocatorsListCompleteResult{
		Items: items,
	}
	return
}
