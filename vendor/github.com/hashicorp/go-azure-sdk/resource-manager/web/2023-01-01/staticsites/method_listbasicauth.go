package staticsites

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListBasicAuthOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]StaticSiteBasicAuthPropertiesARMResource
}

type ListBasicAuthCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []StaticSiteBasicAuthPropertiesARMResource
}

type ListBasicAuthCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListBasicAuthCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListBasicAuth ...
func (c StaticSitesClient) ListBasicAuth(ctx context.Context, id StaticSiteId) (result ListBasicAuthOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListBasicAuthCustomPager{},
		Path:       fmt.Sprintf("%s/basicAuth", id.ID()),
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
		Values *[]StaticSiteBasicAuthPropertiesARMResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListBasicAuthComplete retrieves all the results into a single object
func (c StaticSitesClient) ListBasicAuthComplete(ctx context.Context, id StaticSiteId) (ListBasicAuthCompleteResult, error) {
	return c.ListBasicAuthCompleteMatchingPredicate(ctx, id, StaticSiteBasicAuthPropertiesARMResourceOperationPredicate{})
}

// ListBasicAuthCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c StaticSitesClient) ListBasicAuthCompleteMatchingPredicate(ctx context.Context, id StaticSiteId, predicate StaticSiteBasicAuthPropertiesARMResourceOperationPredicate) (result ListBasicAuthCompleteResult, err error) {
	items := make([]StaticSiteBasicAuthPropertiesARMResource, 0)

	resp, err := c.ListBasicAuth(ctx, id)
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

	result = ListBasicAuthCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
