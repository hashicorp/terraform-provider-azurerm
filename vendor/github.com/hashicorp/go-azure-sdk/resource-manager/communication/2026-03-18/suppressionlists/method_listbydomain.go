package suppressionlists

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByDomainOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SuppressionListResource
}

type ListByDomainCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SuppressionListResource
}

type ListByDomainCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByDomainCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByDomain ...
func (c SuppressionListsClient) ListByDomain(ctx context.Context, id DomainId) (result ListByDomainOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByDomainCustomPager{},
		Path:       fmt.Sprintf("%s/suppressionLists", id.ID()),
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
		Values *[]SuppressionListResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByDomainComplete retrieves all the results into a single object
func (c SuppressionListsClient) ListByDomainComplete(ctx context.Context, id DomainId) (ListByDomainCompleteResult, error) {
	return c.ListByDomainCompleteMatchingPredicate(ctx, id, SuppressionListResourceOperationPredicate{})
}

// ListByDomainCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c SuppressionListsClient) ListByDomainCompleteMatchingPredicate(ctx context.Context, id DomainId, predicate SuppressionListResourceOperationPredicate) (result ListByDomainCompleteResult, err error) {
	items := make([]SuppressionListResource, 0)

	resp, err := c.ListByDomain(ctx, id)
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

	result = ListByDomainCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
