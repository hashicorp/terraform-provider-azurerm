package organizationresources

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

type OrganizationListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]OrganizationResource
}

type OrganizationListByResourceGroupCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []OrganizationResource
}

type OrganizationListByResourceGroupCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *OrganizationListByResourceGroupCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// OrganizationListByResourceGroup ...
func (c OrganizationResourcesClient) OrganizationListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (result OrganizationListByResourceGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &OrganizationListByResourceGroupCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.Confluent/organizations", id.ID()),
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
		Values *[]OrganizationResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// OrganizationListByResourceGroupComplete retrieves all the results into a single object
func (c OrganizationResourcesClient) OrganizationListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (OrganizationListByResourceGroupCompleteResult, error) {
	return c.OrganizationListByResourceGroupCompleteMatchingPredicate(ctx, id, OrganizationResourceOperationPredicate{})
}

// OrganizationListByResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OrganizationResourcesClient) OrganizationListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate OrganizationResourceOperationPredicate) (result OrganizationListByResourceGroupCompleteResult, err error) {
	items := make([]OrganizationResource, 0)

	resp, err := c.OrganizationListByResourceGroup(ctx, id)
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

	result = OrganizationListByResourceGroupCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
