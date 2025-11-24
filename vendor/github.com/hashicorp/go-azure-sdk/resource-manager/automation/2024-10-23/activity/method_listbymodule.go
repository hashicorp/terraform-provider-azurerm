package activity

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByModuleOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Activity
}

type ListByModuleCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Activity
}

type ListByModuleCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByModuleCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByModule ...
func (c ActivityClient) ListByModule(ctx context.Context, id ModuleId) (result ListByModuleOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByModuleCustomPager{},
		Path:       fmt.Sprintf("%s/activities", id.ID()),
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
		Values *[]Activity `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByModuleComplete retrieves all the results into a single object
func (c ActivityClient) ListByModuleComplete(ctx context.Context, id ModuleId) (ListByModuleCompleteResult, error) {
	return c.ListByModuleCompleteMatchingPredicate(ctx, id, ActivityOperationPredicate{})
}

// ListByModuleCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ActivityClient) ListByModuleCompleteMatchingPredicate(ctx context.Context, id ModuleId, predicate ActivityOperationPredicate) (result ListByModuleCompleteResult, err error) {
	items := make([]Activity, 0)

	resp, err := c.ListByModule(ctx, id)
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

	result = ListByModuleCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
