package domains

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByEmailServiceResourceOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DomainResource
}

type ListByEmailServiceResourceCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DomainResource
}

type ListByEmailServiceResourceCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByEmailServiceResourceCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByEmailServiceResource ...
func (c DomainsClient) ListByEmailServiceResource(ctx context.Context, id EmailServiceId) (result ListByEmailServiceResourceOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByEmailServiceResourceCustomPager{},
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
		Values *[]DomainResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByEmailServiceResourceComplete retrieves all the results into a single object
func (c DomainsClient) ListByEmailServiceResourceComplete(ctx context.Context, id EmailServiceId) (ListByEmailServiceResourceCompleteResult, error) {
	return c.ListByEmailServiceResourceCompleteMatchingPredicate(ctx, id, DomainResourceOperationPredicate{})
}

// ListByEmailServiceResourceCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c DomainsClient) ListByEmailServiceResourceCompleteMatchingPredicate(ctx context.Context, id EmailServiceId, predicate DomainResourceOperationPredicate) (result ListByEmailServiceResourceCompleteResult, err error) {
	items := make([]DomainResource, 0)

	resp, err := c.ListByEmailServiceResource(ctx, id)
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

	result = ListByEmailServiceResourceCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
