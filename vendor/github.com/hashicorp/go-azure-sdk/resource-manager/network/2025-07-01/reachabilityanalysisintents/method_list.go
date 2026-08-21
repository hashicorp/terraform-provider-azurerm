package reachabilityanalysisintents

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ReachabilityAnalysisIntent
}

type ListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ReachabilityAnalysisIntent
}

type ListOperationOptions struct {
	Skip      *int64
	SkipToken *string
	SortKey   *string
	SortValue *string
	Top       *int64
}

func DefaultListOperationOptions() ListOperationOptions {
	return ListOperationOptions{}
}

func (o ListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Skip != nil {
		out.Append("skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.SkipToken != nil {
		out.Append("skipToken", fmt.Sprintf("%v", *o.SkipToken))
	}
	if o.SortKey != nil {
		out.Append("sortKey", fmt.Sprintf("%v", *o.SortKey))
	}
	if o.SortValue != nil {
		out.Append("sortValue", fmt.Sprintf("%v", *o.SortValue))
	}
	if o.Top != nil {
		out.Append("top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type ListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// List ...
func (c ReachabilityAnalysisIntentsClient) List(ctx context.Context, id VerifierWorkspaceId, options ListOperationOptions) (result ListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListCustomPager{},
		Path:          fmt.Sprintf("%s/reachabilityAnalysisIntents", id.ID()),
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
		Values *[]ReachabilityAnalysisIntent `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListComplete retrieves all the results into a single object
func (c ReachabilityAnalysisIntentsClient) ListComplete(ctx context.Context, id VerifierWorkspaceId, options ListOperationOptions) (ListCompleteResult, error) {
	return c.ListCompleteMatchingPredicate(ctx, id, options, ReachabilityAnalysisIntentOperationPredicate{})
}

// ListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ReachabilityAnalysisIntentsClient) ListCompleteMatchingPredicate(ctx context.Context, id VerifierWorkspaceId, options ListOperationOptions, predicate ReachabilityAnalysisIntentOperationPredicate) (result ListCompleteResult, err error) {
	items := make([]ReachabilityAnalysisIntent, 0)

	resp, err := c.List(ctx, id, options)
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

	result = ListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
