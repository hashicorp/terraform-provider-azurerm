package webapps

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

type ListHostNameBindingsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]HostNameBinding
}

type ListHostNameBindingsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []HostNameBinding
}

type ListHostNameBindingsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListHostNameBindingsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListHostNameBindings ...
func (c WebAppsClient) ListHostNameBindings(ctx context.Context, id commonids.AppServiceId) (result ListHostNameBindingsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListHostNameBindingsCustomPager{},
		Path:       fmt.Sprintf("%s/hostNameBindings", id.ID()),
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
		Values *[]HostNameBinding `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListHostNameBindingsComplete retrieves all the results into a single object
func (c WebAppsClient) ListHostNameBindingsComplete(ctx context.Context, id commonids.AppServiceId) (ListHostNameBindingsCompleteResult, error) {
	return c.ListHostNameBindingsCompleteMatchingPredicate(ctx, id, HostNameBindingOperationPredicate{})
}

// ListHostNameBindingsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListHostNameBindingsCompleteMatchingPredicate(ctx context.Context, id commonids.AppServiceId, predicate HostNameBindingOperationPredicate) (result ListHostNameBindingsCompleteResult, err error) {
	items := make([]HostNameBinding, 0)

	resp, err := c.ListHostNameBindings(ctx, id)
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

	result = ListHostNameBindingsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
