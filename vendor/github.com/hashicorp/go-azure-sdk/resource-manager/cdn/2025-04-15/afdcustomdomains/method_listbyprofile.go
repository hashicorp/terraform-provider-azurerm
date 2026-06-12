package afdcustomdomains

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByProfileOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]AFDDomain
}

type ListByProfileCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []AFDDomain
}

type ListByProfileCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByProfileCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByProfile ...
func (c AFDCustomDomainsClient) ListByProfile(ctx context.Context, id ProfileId) (result ListByProfileOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByProfileCustomPager{},
		Path:       fmt.Sprintf("%s/customDomains", id.ID()),
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
		Values *[]AFDDomain `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByProfileComplete retrieves all the results into a single object
func (c AFDCustomDomainsClient) ListByProfileComplete(ctx context.Context, id ProfileId) (ListByProfileCompleteResult, error) {
	return c.ListByProfileCompleteMatchingPredicate(ctx, id, AFDDomainOperationPredicate{})
}

// ListByProfileCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AFDCustomDomainsClient) ListByProfileCompleteMatchingPredicate(ctx context.Context, id ProfileId, predicate AFDDomainOperationPredicate) (result ListByProfileCompleteResult, err error) {
	items := make([]AFDDomain, 0)

	resp, err := c.ListByProfile(ctx, id)
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

	result = ListByProfileCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
