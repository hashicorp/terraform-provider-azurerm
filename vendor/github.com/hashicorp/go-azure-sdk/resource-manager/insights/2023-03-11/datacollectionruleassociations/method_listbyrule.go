package datacollectionruleassociations

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByRuleOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DataCollectionRuleAssociationProxyOnlyResource
}

type ListByRuleCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DataCollectionRuleAssociationProxyOnlyResource
}

type ListByRuleCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByRuleCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByRule ...
func (c DataCollectionRuleAssociationsClient) ListByRule(ctx context.Context, id DataCollectionRuleId) (result ListByRuleOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByRuleCustomPager{},
		Path:       fmt.Sprintf("%s/associations", id.ID()),
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
		Values *[]DataCollectionRuleAssociationProxyOnlyResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByRuleComplete retrieves all the results into a single object
func (c DataCollectionRuleAssociationsClient) ListByRuleComplete(ctx context.Context, id DataCollectionRuleId) (ListByRuleCompleteResult, error) {
	return c.ListByRuleCompleteMatchingPredicate(ctx, id, DataCollectionRuleAssociationProxyOnlyResourceOperationPredicate{})
}

// ListByRuleCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c DataCollectionRuleAssociationsClient) ListByRuleCompleteMatchingPredicate(ctx context.Context, id DataCollectionRuleId, predicate DataCollectionRuleAssociationProxyOnlyResourceOperationPredicate) (result ListByRuleCompleteResult, err error) {
	items := make([]DataCollectionRuleAssociationProxyOnlyResource, 0)

	resp, err := c.ListByRule(ctx, id)
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

	result = ListByRuleCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
