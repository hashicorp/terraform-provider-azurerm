package managementlocks

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

type ListAtSubscriptionLevelOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ManagementLockObject
}

type ListAtSubscriptionLevelCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ManagementLockObject
}

type ListAtSubscriptionLevelOperationOptions struct {
	Filter *string
}

func DefaultListAtSubscriptionLevelOperationOptions() ListAtSubscriptionLevelOperationOptions {
	return ListAtSubscriptionLevelOperationOptions{}
}

func (o ListAtSubscriptionLevelOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListAtSubscriptionLevelOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListAtSubscriptionLevelOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

type ListAtSubscriptionLevelCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListAtSubscriptionLevelCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListAtSubscriptionLevel ...
func (c ManagementLocksClient) ListAtSubscriptionLevel(ctx context.Context, id commonids.SubscriptionId, options ListAtSubscriptionLevelOperationOptions) (result ListAtSubscriptionLevelOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListAtSubscriptionLevelCustomPager{},
		Path:          fmt.Sprintf("%s/providers/Microsoft.Authorization/locks", id.ID()),
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
		Values *[]ManagementLockObject `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListAtSubscriptionLevelComplete retrieves all the results into a single object
func (c ManagementLocksClient) ListAtSubscriptionLevelComplete(ctx context.Context, id commonids.SubscriptionId, options ListAtSubscriptionLevelOperationOptions) (ListAtSubscriptionLevelCompleteResult, error) {
	return c.ListAtSubscriptionLevelCompleteMatchingPredicate(ctx, id, options, ManagementLockObjectOperationPredicate{})
}

// ListAtSubscriptionLevelCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ManagementLocksClient) ListAtSubscriptionLevelCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, options ListAtSubscriptionLevelOperationOptions, predicate ManagementLockObjectOperationPredicate) (result ListAtSubscriptionLevelCompleteResult, err error) {
	items := make([]ManagementLockObject, 0)

	resp, err := c.ListAtSubscriptionLevel(ctx, id, options)
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

	result = ListAtSubscriptionLevelCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
