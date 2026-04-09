package appserviceenvironments

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

type ListMultiRoleUsagesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Usage
}

type ListMultiRoleUsagesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Usage
}

type ListMultiRoleUsagesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListMultiRoleUsagesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListMultiRoleUsages ...
func (c AppServiceEnvironmentsClient) ListMultiRoleUsages(ctx context.Context, id commonids.AppServiceEnvironmentId) (result ListMultiRoleUsagesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListMultiRoleUsagesCustomPager{},
		Path:       fmt.Sprintf("%s/multiRolePools/default/usages", id.ID()),
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
		Values *[]Usage `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListMultiRoleUsagesComplete retrieves all the results into a single object
func (c AppServiceEnvironmentsClient) ListMultiRoleUsagesComplete(ctx context.Context, id commonids.AppServiceEnvironmentId) (ListMultiRoleUsagesCompleteResult, error) {
	return c.ListMultiRoleUsagesCompleteMatchingPredicate(ctx, id, UsageOperationPredicate{})
}

// ListMultiRoleUsagesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppServiceEnvironmentsClient) ListMultiRoleUsagesCompleteMatchingPredicate(ctx context.Context, id commonids.AppServiceEnvironmentId, predicate UsageOperationPredicate) (result ListMultiRoleUsagesCompleteResult, err error) {
	items := make([]Usage, 0)

	resp, err := c.ListMultiRoleUsages(ctx, id)
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

	result = ListMultiRoleUsagesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
