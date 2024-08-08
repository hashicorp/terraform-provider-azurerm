package provisionedclusterinstances

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

type VMSkusListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]VMSkuProfile
}

type VMSkusListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []VMSkuProfile
}

type VMSkusListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *VMSkusListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// VMSkusList ...
func (c ProvisionedClusterInstancesClient) VMSkusList(ctx context.Context, id commonids.ScopeId) (result VMSkusListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &VMSkusListCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.HybridContainerService/skus", id.ID()),
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
		Values *[]VMSkuProfile `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// VMSkusListComplete retrieves all the results into a single object
func (c ProvisionedClusterInstancesClient) VMSkusListComplete(ctx context.Context, id commonids.ScopeId) (VMSkusListCompleteResult, error) {
	return c.VMSkusListCompleteMatchingPredicate(ctx, id, VMSkuProfileOperationPredicate{})
}

// VMSkusListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ProvisionedClusterInstancesClient) VMSkusListCompleteMatchingPredicate(ctx context.Context, id commonids.ScopeId, predicate VMSkuProfileOperationPredicate) (result VMSkusListCompleteResult, err error) {
	items := make([]VMSkuProfile, 0)

	resp, err := c.VMSkusList(ctx, id)
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

	result = VMSkusListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
