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

type ListDomainOwnershipIdentifiersOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Identifier
}

type ListDomainOwnershipIdentifiersCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Identifier
}

type ListDomainOwnershipIdentifiersCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListDomainOwnershipIdentifiersCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListDomainOwnershipIdentifiers ...
func (c WebAppsClient) ListDomainOwnershipIdentifiers(ctx context.Context, id commonids.AppServiceId) (result ListDomainOwnershipIdentifiersOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListDomainOwnershipIdentifiersCustomPager{},
		Path:       fmt.Sprintf("%s/domainOwnershipIdentifiers", id.ID()),
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
		Values *[]Identifier `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListDomainOwnershipIdentifiersComplete retrieves all the results into a single object
func (c WebAppsClient) ListDomainOwnershipIdentifiersComplete(ctx context.Context, id commonids.AppServiceId) (ListDomainOwnershipIdentifiersCompleteResult, error) {
	return c.ListDomainOwnershipIdentifiersCompleteMatchingPredicate(ctx, id, IdentifierOperationPredicate{})
}

// ListDomainOwnershipIdentifiersCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListDomainOwnershipIdentifiersCompleteMatchingPredicate(ctx context.Context, id commonids.AppServiceId, predicate IdentifierOperationPredicate) (result ListDomainOwnershipIdentifiersCompleteResult, err error) {
	items := make([]Identifier, 0)

	resp, err := c.ListDomainOwnershipIdentifiers(ctx, id)
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

	result = ListDomainOwnershipIdentifiersCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
