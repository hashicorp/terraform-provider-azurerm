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

type ListConfigurationsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SiteConfigResource
}

type ListConfigurationsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SiteConfigResource
}

type ListConfigurationsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListConfigurationsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListConfigurations ...
func (c WebAppsClient) ListConfigurations(ctx context.Context, id commonids.AppServiceId) (result ListConfigurationsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListConfigurationsCustomPager{},
		Path:       fmt.Sprintf("%s/config", id.ID()),
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
		Values *[]SiteConfigResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListConfigurationsComplete retrieves all the results into a single object
func (c WebAppsClient) ListConfigurationsComplete(ctx context.Context, id commonids.AppServiceId) (ListConfigurationsCompleteResult, error) {
	return c.ListConfigurationsCompleteMatchingPredicate(ctx, id, SiteConfigResourceOperationPredicate{})
}

// ListConfigurationsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListConfigurationsCompleteMatchingPredicate(ctx context.Context, id commonids.AppServiceId, predicate SiteConfigResourceOperationPredicate) (result ListConfigurationsCompleteResult, err error) {
	items := make([]SiteConfigResource, 0)

	resp, err := c.ListConfigurations(ctx, id)
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

	result = ListConfigurationsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
