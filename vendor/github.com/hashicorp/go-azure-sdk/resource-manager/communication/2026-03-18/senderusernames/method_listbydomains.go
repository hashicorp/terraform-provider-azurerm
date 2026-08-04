package senderusernames

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByDomainsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SenderUsernameResource
}

type ListByDomainsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SenderUsernameResource
}

type ListByDomainsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByDomainsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByDomains ...
func (c SenderUsernamesClient) ListByDomains(ctx context.Context, id DomainId) (result ListByDomainsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByDomainsCustomPager{},
		Path:       fmt.Sprintf("%s/senderUsernames", id.ID()),
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
		Values *[]SenderUsernameResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByDomainsComplete retrieves all the results into a single object
func (c SenderUsernamesClient) ListByDomainsComplete(ctx context.Context, id DomainId) (ListByDomainsCompleteResult, error) {
	return c.ListByDomainsCompleteMatchingPredicate(ctx, id, SenderUsernameResourceOperationPredicate{})
}

// ListByDomainsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c SenderUsernamesClient) ListByDomainsCompleteMatchingPredicate(ctx context.Context, id DomainId, predicate SenderUsernameResourceOperationPredicate) (result ListByDomainsCompleteResult, err error) {
	items := make([]SenderUsernameResource, 0)

	resp, err := c.ListByDomains(ctx, id)
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

	result = ListByDomainsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
