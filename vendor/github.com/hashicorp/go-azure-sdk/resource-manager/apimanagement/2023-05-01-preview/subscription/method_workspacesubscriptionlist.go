package subscription

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceSubscriptionListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SubscriptionContract
}

type WorkspaceSubscriptionListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SubscriptionContract
}

type WorkspaceSubscriptionListOperationOptions struct {
	Filter *string
	Skip   *int64
	Top    *int64
}

func DefaultWorkspaceSubscriptionListOperationOptions() WorkspaceSubscriptionListOperationOptions {
	return WorkspaceSubscriptionListOperationOptions{}
}

func (o WorkspaceSubscriptionListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o WorkspaceSubscriptionListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o WorkspaceSubscriptionListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type WorkspaceSubscriptionListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *WorkspaceSubscriptionListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// WorkspaceSubscriptionList ...
func (c SubscriptionClient) WorkspaceSubscriptionList(ctx context.Context, id WorkspaceId, options WorkspaceSubscriptionListOperationOptions) (result WorkspaceSubscriptionListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &WorkspaceSubscriptionListCustomPager{},
		Path:          fmt.Sprintf("%s/subscriptions", id.ID()),
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
		Values *[]SubscriptionContract `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// WorkspaceSubscriptionListComplete retrieves all the results into a single object
func (c SubscriptionClient) WorkspaceSubscriptionListComplete(ctx context.Context, id WorkspaceId, options WorkspaceSubscriptionListOperationOptions) (WorkspaceSubscriptionListCompleteResult, error) {
	return c.WorkspaceSubscriptionListCompleteMatchingPredicate(ctx, id, options, SubscriptionContractOperationPredicate{})
}

// WorkspaceSubscriptionListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c SubscriptionClient) WorkspaceSubscriptionListCompleteMatchingPredicate(ctx context.Context, id WorkspaceId, options WorkspaceSubscriptionListOperationOptions, predicate SubscriptionContractOperationPredicate) (result WorkspaceSubscriptionListCompleteResult, err error) {
	items := make([]SubscriptionContract, 0)

	resp, err := c.WorkspaceSubscriptionList(ctx, id, options)
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

	result = WorkspaceSubscriptionListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
