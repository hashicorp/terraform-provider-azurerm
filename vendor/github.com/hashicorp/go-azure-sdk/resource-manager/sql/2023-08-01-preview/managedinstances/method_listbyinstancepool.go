package managedinstances

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByInstancePoolOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ManagedInstance
}

type ListByInstancePoolCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ManagedInstance
}

type ListByInstancePoolOperationOptions struct {
	Expand *string
}

func DefaultListByInstancePoolOperationOptions() ListByInstancePoolOperationOptions {
	return ListByInstancePoolOperationOptions{}
}

func (o ListByInstancePoolOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByInstancePoolOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListByInstancePoolOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Expand != nil {
		out.Append("$expand", fmt.Sprintf("%v", *o.Expand))
	}
	return &out
}

type ListByInstancePoolCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByInstancePoolCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByInstancePool ...
func (c ManagedInstancesClient) ListByInstancePool(ctx context.Context, id InstancePoolId, options ListByInstancePoolOperationOptions) (result ListByInstancePoolOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListByInstancePoolCustomPager{},
		Path:          fmt.Sprintf("%s/managedInstances", id.ID()),
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
		Values *[]ManagedInstance `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByInstancePoolComplete retrieves all the results into a single object
func (c ManagedInstancesClient) ListByInstancePoolComplete(ctx context.Context, id InstancePoolId, options ListByInstancePoolOperationOptions) (ListByInstancePoolCompleteResult, error) {
	return c.ListByInstancePoolCompleteMatchingPredicate(ctx, id, options, ManagedInstanceOperationPredicate{})
}

// ListByInstancePoolCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ManagedInstancesClient) ListByInstancePoolCompleteMatchingPredicate(ctx context.Context, id InstancePoolId, options ListByInstancePoolOperationOptions, predicate ManagedInstanceOperationPredicate) (result ListByInstancePoolCompleteResult, err error) {
	items := make([]ManagedInstance, 0)

	resp, err := c.ListByInstancePool(ctx, id, options)
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

	result = ListByInstancePoolCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
