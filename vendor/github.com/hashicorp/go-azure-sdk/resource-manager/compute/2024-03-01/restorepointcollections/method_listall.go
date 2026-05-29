package restorepointcollections

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

type ListAllOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RestorePointCollection
}

type ListAllCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []RestorePointCollection
}

type ListAllCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListAllCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListAll ...
func (c RestorePointCollectionsClient) ListAll(ctx context.Context, id commonids.SubscriptionId) (result ListAllOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListAllCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.Compute/restorePointCollections", id.ID()),
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
		Values *[]RestorePointCollection `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListAllComplete retrieves all the results into a single object
func (c RestorePointCollectionsClient) ListAllComplete(ctx context.Context, id commonids.SubscriptionId) (ListAllCompleteResult, error) {
	return c.ListAllCompleteMatchingPredicate(ctx, id, RestorePointCollectionOperationPredicate{})
}

// ListAllCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c RestorePointCollectionsClient) ListAllCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate RestorePointCollectionOperationPredicate) (result ListAllCompleteResult, err error) {
	items := make([]RestorePointCollection, 0)

	resp, err := c.ListAll(ctx, id)
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

	result = ListAllCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
