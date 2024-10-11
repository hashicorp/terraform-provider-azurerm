package roleeligibilityschedulerequests

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

type ListForScopeOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RoleEligibilityScheduleRequest
}

type ListForScopeCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []RoleEligibilityScheduleRequest
}

type ListForScopeOperationOptions struct {
	Filter *string
}

func DefaultListForScopeOperationOptions() ListForScopeOperationOptions {
	return ListForScopeOperationOptions{}
}

func (o ListForScopeOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListForScopeOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListForScopeOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

type ListForScopeCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListForScopeCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListForScope ...
func (c RoleEligibilityScheduleRequestsClient) ListForScope(ctx context.Context, id commonids.ScopeId, options ListForScopeOperationOptions) (result ListForScopeOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListForScopeCustomPager{},
		Path:          fmt.Sprintf("%s/providers/Microsoft.Authorization/roleEligibilityScheduleRequests", id.ID()),
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
		Values *[]RoleEligibilityScheduleRequest `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListForScopeComplete retrieves all the results into a single object
func (c RoleEligibilityScheduleRequestsClient) ListForScopeComplete(ctx context.Context, id commonids.ScopeId, options ListForScopeOperationOptions) (ListForScopeCompleteResult, error) {
	return c.ListForScopeCompleteMatchingPredicate(ctx, id, options, RoleEligibilityScheduleRequestOperationPredicate{})
}

// ListForScopeCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c RoleEligibilityScheduleRequestsClient) ListForScopeCompleteMatchingPredicate(ctx context.Context, id commonids.ScopeId, options ListForScopeOperationOptions, predicate RoleEligibilityScheduleRequestOperationPredicate) (result ListForScopeCompleteResult, err error) {
	items := make([]RoleEligibilityScheduleRequest, 0)

	resp, err := c.ListForScope(ctx, id, options)
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

	result = ListForScopeCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
