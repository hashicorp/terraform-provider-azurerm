package azuretrafficcollectors

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

type ByResourceGroupListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]AzureTrafficCollector
}

type ByResourceGroupListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []AzureTrafficCollector
}

type ByResourceGroupListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ByResourceGroupListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ByResourceGroupList ...
func (c AzureTrafficCollectorsClient) ByResourceGroupList(ctx context.Context, id commonids.ResourceGroupId) (result ByResourceGroupListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ByResourceGroupListCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.NetworkFunction/azureTrafficCollectors", id.ID()),
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
		Values *[]AzureTrafficCollector `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ByResourceGroupListComplete retrieves all the results into a single object
func (c AzureTrafficCollectorsClient) ByResourceGroupListComplete(ctx context.Context, id commonids.ResourceGroupId) (ByResourceGroupListCompleteResult, error) {
	return c.ByResourceGroupListCompleteMatchingPredicate(ctx, id, AzureTrafficCollectorOperationPredicate{})
}

// ByResourceGroupListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AzureTrafficCollectorsClient) ByResourceGroupListCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate AzureTrafficCollectorOperationPredicate) (result ByResourceGroupListCompleteResult, err error) {
	items := make([]AzureTrafficCollector, 0)

	resp, err := c.ByResourceGroupList(ctx, id)
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

	result = ByResourceGroupListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
