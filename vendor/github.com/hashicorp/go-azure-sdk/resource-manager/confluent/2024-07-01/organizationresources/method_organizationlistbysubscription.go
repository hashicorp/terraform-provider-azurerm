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

type OrganizationListBySubscriptionOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]OrganizationResource
}

type OrganizationListBySubscriptionCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []OrganizationResource
}

type OrganizationListBySubscriptionCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *OrganizationListBySubscriptionCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// OrganizationListBySubscription ...
func (c OrganizationResourcesClient) OrganizationListBySubscription(ctx context.Context, id commonids.SubscriptionId) (result OrganizationListBySubscriptionOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &OrganizationListBySubscriptionCustomPager{},
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

// OrganizationListBySubscriptionComplete retrieves all the results into a single object
func (c OrganizationResourcesClient) OrganizationListBySubscriptionComplete(ctx context.Context, id commonids.SubscriptionId) (OrganizationListBySubscriptionCompleteResult, error) {
	return c.OrganizationListBySubscriptionCompleteMatchingPredicate(ctx, id, OrganizationResourceOperationPredicate{})
}

// OrganizationListBySubscriptionCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OrganizationResourcesClient) OrganizationListBySubscriptionCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate OrganizationResourceOperationPredicate) (result OrganizationListBySubscriptionCompleteResult, err error) {
	items := make([]OrganizationResource, 0)

	resp, err := c.OrganizationListBySubscription(ctx, id)
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

	result = OrganizationListBySubscriptionCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
