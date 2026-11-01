package scheduledactionoperationgroup

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

type ScheduledActionsListByScopeOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ScheduledAction
}

type ScheduledActionsListByScopeCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ScheduledAction
}

type ScheduledActionsListByScopeOperationOptions struct {
	Filter *string
}

func DefaultScheduledActionsListByScopeOperationOptions() ScheduledActionsListByScopeOperationOptions {
	return ScheduledActionsListByScopeOperationOptions{}
}

func (o ScheduledActionsListByScopeOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ScheduledActionsListByScopeOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ScheduledActionsListByScopeOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

type ScheduledActionsListByScopeCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ScheduledActionsListByScopeCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ScheduledActionsListByScope ...
func (c ScheduledActionOperationGroupClient) ScheduledActionsListByScope(ctx context.Context, id commonids.ScopeId, options ScheduledActionsListByScopeOperationOptions) (result ScheduledActionsListByScopeOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ScheduledActionsListByScopeCustomPager{},
		Path:          fmt.Sprintf("%s/providers/Microsoft.CostManagement/scheduledActions", id.ID()),
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
		Values *[]ScheduledAction `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ScheduledActionsListByScopeComplete retrieves all the results into a single object
func (c ScheduledActionOperationGroupClient) ScheduledActionsListByScopeComplete(ctx context.Context, id commonids.ScopeId, options ScheduledActionsListByScopeOperationOptions) (ScheduledActionsListByScopeCompleteResult, error) {
	return c.ScheduledActionsListByScopeCompleteMatchingPredicate(ctx, id, options, ScheduledActionOperationPredicate{})
}

// ScheduledActionsListByScopeCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ScheduledActionOperationGroupClient) ScheduledActionsListByScopeCompleteMatchingPredicate(ctx context.Context, id commonids.ScopeId, options ScheduledActionsListByScopeOperationOptions, predicate ScheduledActionOperationPredicate) (result ScheduledActionsListByScopeCompleteResult, err error) {
	items := make([]ScheduledAction, 0)

	resp, err := c.ScheduledActionsListByScope(ctx, id, options)
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

	result = ScheduledActionsListByScopeCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
