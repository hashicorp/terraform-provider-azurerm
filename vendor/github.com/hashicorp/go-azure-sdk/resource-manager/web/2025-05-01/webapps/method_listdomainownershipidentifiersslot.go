package webapps

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListDomainOwnershipIdentifiersSlotOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Identifier
}

type ListDomainOwnershipIdentifiersSlotCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Identifier
}

type ListDomainOwnershipIdentifiersSlotCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListDomainOwnershipIdentifiersSlotCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListDomainOwnershipIdentifiersSlot ...
func (c WebAppsClient) ListDomainOwnershipIdentifiersSlot(ctx context.Context, id SlotId) (result ListDomainOwnershipIdentifiersSlotOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListDomainOwnershipIdentifiersSlotCustomPager{},
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

// ListDomainOwnershipIdentifiersSlotComplete retrieves all the results into a single object
func (c WebAppsClient) ListDomainOwnershipIdentifiersSlotComplete(ctx context.Context, id SlotId) (ListDomainOwnershipIdentifiersSlotCompleteResult, error) {
	return c.ListDomainOwnershipIdentifiersSlotCompleteMatchingPredicate(ctx, id, IdentifierOperationPredicate{})
}

// ListDomainOwnershipIdentifiersSlotCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListDomainOwnershipIdentifiersSlotCompleteMatchingPredicate(ctx context.Context, id SlotId, predicate IdentifierOperationPredicate) (result ListDomainOwnershipIdentifiersSlotCompleteResult, err error) {
	items := make([]Identifier, 0)

	resp, err := c.ListDomainOwnershipIdentifiersSlot(ctx, id)
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

	result = ListDomainOwnershipIdentifiersSlotCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
