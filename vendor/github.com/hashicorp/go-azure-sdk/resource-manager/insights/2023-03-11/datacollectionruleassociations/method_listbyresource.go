package datacollectionruleassociations

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

type ListByResourceOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DataCollectionRuleAssociationProxyOnlyResource
}

type ListByResourceCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DataCollectionRuleAssociationProxyOnlyResource
}

type ListByResourceCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByResourceCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByResource ...
func (c DataCollectionRuleAssociationsClient) ListByResource(ctx context.Context, id commonids.ScopeId) (result ListByResourceOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByResourceCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.Insights/dataCollectionRuleAssociations", id.ID()),
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

// ListByResourceComplete retrieves all the results into a single object
func (c DataCollectionRuleAssociationsClient) ListByResourceComplete(ctx context.Context, id commonids.ScopeId) (ListByResourceCompleteResult, error) {
	return c.ListByResourceCompleteMatchingPredicate(ctx, id, DataCollectionRuleAssociationProxyOnlyResourceOperationPredicate{})
}

// ListByResourceCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c DataCollectionRuleAssociationsClient) ListByResourceCompleteMatchingPredicate(ctx context.Context, id commonids.ScopeId, predicate DataCollectionRuleAssociationProxyOnlyResourceOperationPredicate) (result ListByResourceCompleteResult, err error) {
	items := make([]DataCollectionRuleAssociationProxyOnlyResource, 0)

	resp, err := c.ListByResource(ctx, id)
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

	result = ListByResourceCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
