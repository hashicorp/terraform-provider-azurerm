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

type ListByDataCollectionEndpointOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DataCollectionRuleAssociationProxyOnlyResource
}

type ListByDataCollectionEndpointCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DataCollectionRuleAssociationProxyOnlyResource
}

type ListByDataCollectionEndpointCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByDataCollectionEndpointCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByDataCollectionEndpoint ...
func (c DataCollectionRuleAssociationsClient) ListByDataCollectionEndpoint(ctx context.Context, id DataCollectionEndpointId) (result ListByDataCollectionEndpointOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByDataCollectionEndpointCustomPager{},
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

// ListByDataCollectionEndpointComplete retrieves all the results into a single object
func (c DataCollectionRuleAssociationsClient) ListByDataCollectionEndpointComplete(ctx context.Context, id DataCollectionEndpointId) (ListByDataCollectionEndpointCompleteResult, error) {
	return c.ListByDataCollectionEndpointCompleteMatchingPredicate(ctx, id, DataCollectionRuleAssociationProxyOnlyResourceOperationPredicate{})
}

// ListByDataCollectionEndpointCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c DataCollectionRuleAssociationsClient) ListByDataCollectionEndpointCompleteMatchingPredicate(ctx context.Context, id DataCollectionEndpointId, predicate DataCollectionRuleAssociationProxyOnlyResourceOperationPredicate) (result ListByDataCollectionEndpointCompleteResult, err error) {
	items := make([]DataCollectionRuleAssociationProxyOnlyResource, 0)

	resp, err := c.ListByDataCollectionEndpoint(ctx, id)
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

	result = ListByDataCollectionEndpointCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
