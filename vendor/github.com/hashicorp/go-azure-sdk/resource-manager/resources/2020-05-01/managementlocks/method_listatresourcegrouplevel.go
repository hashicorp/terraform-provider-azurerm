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

type ListAtResourceGroupLevelOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ManagementLockObject
}

type ListAtResourceGroupLevelCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ManagementLockObject
}

type ListAtResourceGroupLevelOperationOptions struct {
	Filter *string
}

func DefaultListAtResourceGroupLevelOperationOptions() ListAtResourceGroupLevelOperationOptions {
	return ListAtResourceGroupLevelOperationOptions{}
}

func (o ListAtResourceGroupLevelOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListAtResourceGroupLevelOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListAtResourceGroupLevelOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

type ListAtResourceGroupLevelCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListAtResourceGroupLevelCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListAtResourceGroupLevel ...
func (c ManagementLocksClient) ListAtResourceGroupLevel(ctx context.Context, id commonids.ResourceGroupId, options ListAtResourceGroupLevelOperationOptions) (result ListAtResourceGroupLevelOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListAtResourceGroupLevelCustomPager{},
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

// ListAtResourceGroupLevelComplete retrieves all the results into a single object
func (c ManagementLocksClient) ListAtResourceGroupLevelComplete(ctx context.Context, id commonids.ResourceGroupId, options ListAtResourceGroupLevelOperationOptions) (ListAtResourceGroupLevelCompleteResult, error) {
	return c.ListAtResourceGroupLevelCompleteMatchingPredicate(ctx, id, options, ManagementLockObjectOperationPredicate{})
}

// ListAtResourceGroupLevelCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ManagementLocksClient) ListAtResourceGroupLevelCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, options ListAtResourceGroupLevelOperationOptions, predicate ManagementLockObjectOperationPredicate) (result ListAtResourceGroupLevelCompleteResult, err error) {
	items := make([]ManagementLockObject, 0)

	resp, err := c.ListAtResourceGroupLevel(ctx, id, options)
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

	result = ListAtResourceGroupLevelCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
