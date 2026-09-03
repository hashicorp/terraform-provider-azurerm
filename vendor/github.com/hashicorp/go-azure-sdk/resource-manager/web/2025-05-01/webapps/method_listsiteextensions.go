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

type ListSiteExtensionsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SiteExtensionInfo
}

type ListSiteExtensionsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SiteExtensionInfo
}

type ListSiteExtensionsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListSiteExtensionsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListSiteExtensions ...
func (c WebAppsClient) ListSiteExtensions(ctx context.Context, id commonids.AppServiceId) (result ListSiteExtensionsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListSiteExtensionsCustomPager{},
		Path:       fmt.Sprintf("%s/siteExtensions", id.ID()),
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
		Values *[]SiteExtensionInfo `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListSiteExtensionsComplete retrieves all the results into a single object
func (c WebAppsClient) ListSiteExtensionsComplete(ctx context.Context, id commonids.AppServiceId) (ListSiteExtensionsCompleteResult, error) {
	return c.ListSiteExtensionsCompleteMatchingPredicate(ctx, id, SiteExtensionInfoOperationPredicate{})
}

// ListSiteExtensionsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListSiteExtensionsCompleteMatchingPredicate(ctx context.Context, id commonids.AppServiceId, predicate SiteExtensionInfoOperationPredicate) (result ListSiteExtensionsCompleteResult, err error) {
	items := make([]SiteExtensionInfo, 0)

	resp, err := c.ListSiteExtensions(ctx, id)
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

	result = ListSiteExtensionsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
