package appplatform

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomDomainsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]CustomDomainResource
}

type CustomDomainsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []CustomDomainResource
}

type CustomDomainsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *CustomDomainsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// CustomDomainsList ...
func (c AppPlatformClient) CustomDomainsList(ctx context.Context, id AppId) (result CustomDomainsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &CustomDomainsListCustomPager{},
		Path:       fmt.Sprintf("%s/domains", id.ID()),
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
		Values *[]CustomDomainResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// CustomDomainsListComplete retrieves all the results into a single object
func (c AppPlatformClient) CustomDomainsListComplete(ctx context.Context, id AppId) (CustomDomainsListCompleteResult, error) {
	return c.CustomDomainsListCompleteMatchingPredicate(ctx, id, CustomDomainResourceOperationPredicate{})
}

// CustomDomainsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppPlatformClient) CustomDomainsListCompleteMatchingPredicate(ctx context.Context, id AppId, predicate CustomDomainResourceOperationPredicate) (result CustomDomainsListCompleteResult, err error) {
	items := make([]CustomDomainResource, 0)

	resp, err := c.CustomDomainsList(ctx, id)
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

	result = CustomDomainsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
