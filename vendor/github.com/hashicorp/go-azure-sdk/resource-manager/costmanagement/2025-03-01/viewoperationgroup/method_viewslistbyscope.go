package viewoperationgroup

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

type ViewsListByScopeOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]View
}

type ViewsListByScopeCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []View
}

type ViewsListByScopeCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ViewsListByScopeCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ViewsListByScope ...
func (c ViewOperationGroupClient) ViewsListByScope(ctx context.Context, id commonids.ScopeId) (result ViewsListByScopeOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ViewsListByScopeCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.CostManagement/views", id.ID()),
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
		Values *[]View `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ViewsListByScopeComplete retrieves all the results into a single object
func (c ViewOperationGroupClient) ViewsListByScopeComplete(ctx context.Context, id commonids.ScopeId) (ViewsListByScopeCompleteResult, error) {
	return c.ViewsListByScopeCompleteMatchingPredicate(ctx, id, ViewOperationPredicate{})
}

// ViewsListByScopeCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ViewOperationGroupClient) ViewsListByScopeCompleteMatchingPredicate(ctx context.Context, id commonids.ScopeId, predicate ViewOperationPredicate) (result ViewsListByScopeCompleteResult, err error) {
	items := make([]View, 0)

	resp, err := c.ViewsListByScope(ctx, id)
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

	result = ViewsListByScopeCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
