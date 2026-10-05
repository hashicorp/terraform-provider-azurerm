package hostpool

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

type ListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]HostPool
}

type ListByResourceGroupCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []HostPool
}

type ListByResourceGroupOperationOptions struct {
	InitialSkip  *int64
	IsDescending *bool
	PageSize     *int64
}

func DefaultListByResourceGroupOperationOptions() ListByResourceGroupOperationOptions {
	return ListByResourceGroupOperationOptions{}
}

func (o ListByResourceGroupOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByResourceGroupOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListByResourceGroupOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.InitialSkip != nil {
		out.Append("initialSkip", fmt.Sprintf("%v", *o.InitialSkip))
	}
	if o.IsDescending != nil {
		out.Append("isDescending", fmt.Sprintf("%v", *o.IsDescending))
	}
	if o.PageSize != nil {
		out.Append("pageSize", fmt.Sprintf("%v", *o.PageSize))
	}
	return &out
}

type ListByResourceGroupCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByResourceGroupCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByResourceGroup ...
func (c HostPoolClient) ListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId, options ListByResourceGroupOperationOptions) (result ListByResourceGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListByResourceGroupCustomPager{},
		Path:          fmt.Sprintf("%s/providers/Microsoft.DesktopVirtualization/hostPools", id.ID()),
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
		Values *[]HostPool `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByResourceGroupComplete retrieves all the results into a single object
func (c HostPoolClient) ListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId, options ListByResourceGroupOperationOptions) (ListByResourceGroupCompleteResult, error) {
	return c.ListByResourceGroupCompleteMatchingPredicate(ctx, id, options, HostPoolOperationPredicate{})
}

// ListByResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c HostPoolClient) ListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, options ListByResourceGroupOperationOptions, predicate HostPoolOperationPredicate) (result ListByResourceGroupCompleteResult, err error) {
	items := make([]HostPool, 0)

	resp, err := c.ListByResourceGroup(ctx, id, options)
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

	result = ListByResourceGroupCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
