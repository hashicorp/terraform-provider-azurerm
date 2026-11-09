package managedinstancekeys

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

type ListByInstanceOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ManagedInstanceKey
}

type ListByInstanceCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ManagedInstanceKey
}

type ListByInstanceOperationOptions struct {
	Filter *string
}

func DefaultListByInstanceOperationOptions() ListByInstanceOperationOptions {
	return ListByInstanceOperationOptions{}
}

func (o ListByInstanceOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByInstanceOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListByInstanceOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

type ListByInstanceCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByInstanceCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByInstance ...
func (c ManagedInstanceKeysClient) ListByInstance(ctx context.Context, id commonids.SqlManagedInstanceId, options ListByInstanceOperationOptions) (result ListByInstanceOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListByInstanceCustomPager{},
		Path:          fmt.Sprintf("%s/keys", id.ID()),
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
		Values *[]ManagedInstanceKey `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByInstanceComplete retrieves all the results into a single object
func (c ManagedInstanceKeysClient) ListByInstanceComplete(ctx context.Context, id commonids.SqlManagedInstanceId, options ListByInstanceOperationOptions) (ListByInstanceCompleteResult, error) {
	return c.ListByInstanceCompleteMatchingPredicate(ctx, id, options, ManagedInstanceKeyOperationPredicate{})
}

// ListByInstanceCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ManagedInstanceKeysClient) ListByInstanceCompleteMatchingPredicate(ctx context.Context, id commonids.SqlManagedInstanceId, options ListByInstanceOperationOptions, predicate ManagedInstanceKeyOperationPredicate) (result ListByInstanceCompleteResult, err error) {
	items := make([]ManagedInstanceKey, 0)

	resp, err := c.ListByInstance(ctx, id, options)
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

	result = ListByInstanceCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
