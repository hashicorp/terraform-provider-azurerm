package projectpolicies

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByDevCenterOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ProjectPolicy
}

type ListByDevCenterCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ProjectPolicy
}

type ListByDevCenterOperationOptions struct {
	Top *int64
}

func DefaultListByDevCenterOperationOptions() ListByDevCenterOperationOptions {
	return ListByDevCenterOperationOptions{}
}

func (o ListByDevCenterOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByDevCenterOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListByDevCenterOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type ListByDevCenterCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByDevCenterCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByDevCenter ...
func (c ProjectPoliciesClient) ListByDevCenter(ctx context.Context, id DevCenterId, options ListByDevCenterOperationOptions) (result ListByDevCenterOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListByDevCenterCustomPager{},
		Path:          fmt.Sprintf("%s/projectPolicies", id.ID()),
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
		Values *[]ProjectPolicy `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByDevCenterComplete retrieves all the results into a single object
func (c ProjectPoliciesClient) ListByDevCenterComplete(ctx context.Context, id DevCenterId, options ListByDevCenterOperationOptions) (ListByDevCenterCompleteResult, error) {
	return c.ListByDevCenterCompleteMatchingPredicate(ctx, id, options, ProjectPolicyOperationPredicate{})
}

// ListByDevCenterCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ProjectPoliciesClient) ListByDevCenterCompleteMatchingPredicate(ctx context.Context, id DevCenterId, options ListByDevCenterOperationOptions, predicate ProjectPolicyOperationPredicate) (result ListByDevCenterCompleteResult, err error) {
	items := make([]ProjectPolicy, 0)

	resp, err := c.ListByDevCenter(ctx, id, options)
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

	result = ListByDevCenterCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
