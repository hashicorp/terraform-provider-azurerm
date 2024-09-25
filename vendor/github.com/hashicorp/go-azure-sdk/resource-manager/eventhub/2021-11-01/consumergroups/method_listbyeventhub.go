package consumergroups

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByEventHubOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ConsumerGroup
}

type ListByEventHubCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ConsumerGroup
}

type ListByEventHubOperationOptions struct {
	Skip *int64
	Top  *int64
}

func DefaultListByEventHubOperationOptions() ListByEventHubOperationOptions {
	return ListByEventHubOperationOptions{}
}

func (o ListByEventHubOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByEventHubOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListByEventHubOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type ListByEventHubCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByEventHubCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByEventHub ...
func (c ConsumerGroupsClient) ListByEventHub(ctx context.Context, id EventhubId, options ListByEventHubOperationOptions) (result ListByEventHubOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListByEventHubCustomPager{},
		Path:          fmt.Sprintf("%s/consumerGroups", id.ID()),
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
		Values *[]ConsumerGroup `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByEventHubComplete retrieves all the results into a single object
func (c ConsumerGroupsClient) ListByEventHubComplete(ctx context.Context, id EventhubId, options ListByEventHubOperationOptions) (ListByEventHubCompleteResult, error) {
	return c.ListByEventHubCompleteMatchingPredicate(ctx, id, options, ConsumerGroupOperationPredicate{})
}

// ListByEventHubCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ConsumerGroupsClient) ListByEventHubCompleteMatchingPredicate(ctx context.Context, id EventhubId, options ListByEventHubOperationOptions, predicate ConsumerGroupOperationPredicate) (result ListByEventHubCompleteResult, err error) {
	items := make([]ConsumerGroup, 0)

	resp, err := c.ListByEventHub(ctx, id, options)
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

	result = ListByEventHubCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
