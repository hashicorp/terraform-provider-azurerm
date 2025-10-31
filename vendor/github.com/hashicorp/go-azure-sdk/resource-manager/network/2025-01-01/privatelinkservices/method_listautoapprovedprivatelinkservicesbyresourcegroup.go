package privatelinkservices

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListAutoApprovedPrivateLinkServicesByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]AutoApprovedPrivateLinkService
}

type ListAutoApprovedPrivateLinkServicesByResourceGroupCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []AutoApprovedPrivateLinkService
}

type ListAutoApprovedPrivateLinkServicesByResourceGroupCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListAutoApprovedPrivateLinkServicesByResourceGroupCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListAutoApprovedPrivateLinkServicesByResourceGroup ...
func (c PrivateLinkServicesClient) ListAutoApprovedPrivateLinkServicesByResourceGroup(ctx context.Context, id ProviderLocationId) (result ListAutoApprovedPrivateLinkServicesByResourceGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListAutoApprovedPrivateLinkServicesByResourceGroupCustomPager{},
		Path:       fmt.Sprintf("%s/autoApprovedPrivateLinkServices", id.ID()),
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
		Values *[]AutoApprovedPrivateLinkService `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListAutoApprovedPrivateLinkServicesByResourceGroupComplete retrieves all the results into a single object
func (c PrivateLinkServicesClient) ListAutoApprovedPrivateLinkServicesByResourceGroupComplete(ctx context.Context, id ProviderLocationId) (ListAutoApprovedPrivateLinkServicesByResourceGroupCompleteResult, error) {
	return c.ListAutoApprovedPrivateLinkServicesByResourceGroupCompleteMatchingPredicate(ctx, id, AutoApprovedPrivateLinkServiceOperationPredicate{})
}

// ListAutoApprovedPrivateLinkServicesByResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c PrivateLinkServicesClient) ListAutoApprovedPrivateLinkServicesByResourceGroupCompleteMatchingPredicate(ctx context.Context, id ProviderLocationId, predicate AutoApprovedPrivateLinkServiceOperationPredicate) (result ListAutoApprovedPrivateLinkServicesByResourceGroupCompleteResult, err error) {
	items := make([]AutoApprovedPrivateLinkService, 0)

	resp, err := c.ListAutoApprovedPrivateLinkServicesByResourceGroup(ctx, id)
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

	result = ListAutoApprovedPrivateLinkServicesByResourceGroupCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
