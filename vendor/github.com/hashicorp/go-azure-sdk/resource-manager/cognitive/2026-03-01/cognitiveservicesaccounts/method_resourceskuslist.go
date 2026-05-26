package cognitiveservicesaccounts

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

type ResourceSkusListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ResourceSku
}

type ResourceSkusListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ResourceSku
}

type ResourceSkusListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ResourceSkusListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ResourceSkusList ...
func (c CognitiveServicesAccountsClient) ResourceSkusList(ctx context.Context, id commonids.SubscriptionId) (result ResourceSkusListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ResourceSkusListCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.CognitiveServices/skus", id.ID()),
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
		Values *[]ResourceSku `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ResourceSkusListComplete retrieves all the results into a single object
func (c CognitiveServicesAccountsClient) ResourceSkusListComplete(ctx context.Context, id commonids.SubscriptionId) (ResourceSkusListCompleteResult, error) {
	return c.ResourceSkusListCompleteMatchingPredicate(ctx, id, ResourceSkuOperationPredicate{})
}

// ResourceSkusListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c CognitiveServicesAccountsClient) ResourceSkusListCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate ResourceSkuOperationPredicate) (result ResourceSkusListCompleteResult, err error) {
	items := make([]ResourceSku, 0)

	resp, err := c.ResourceSkusList(ctx, id)
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

	result = ResourceSkusListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
