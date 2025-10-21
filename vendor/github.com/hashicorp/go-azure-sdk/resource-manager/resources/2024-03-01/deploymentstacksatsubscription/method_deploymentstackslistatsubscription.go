package deploymentstacksatsubscription

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

type DeploymentStacksListAtSubscriptionOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DeploymentStack
}

type DeploymentStacksListAtSubscriptionCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DeploymentStack
}

type DeploymentStacksListAtSubscriptionCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *DeploymentStacksListAtSubscriptionCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// DeploymentStacksListAtSubscription ...
func (c DeploymentStacksAtSubscriptionClient) DeploymentStacksListAtSubscription(ctx context.Context, id commonids.SubscriptionId) (result DeploymentStacksListAtSubscriptionOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &DeploymentStacksListAtSubscriptionCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.Resources/deploymentStacks", id.ID()),
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
		Values *[]DeploymentStack `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// DeploymentStacksListAtSubscriptionComplete retrieves all the results into a single object
func (c DeploymentStacksAtSubscriptionClient) DeploymentStacksListAtSubscriptionComplete(ctx context.Context, id commonids.SubscriptionId) (DeploymentStacksListAtSubscriptionCompleteResult, error) {
	return c.DeploymentStacksListAtSubscriptionCompleteMatchingPredicate(ctx, id, DeploymentStackOperationPredicate{})
}

// DeploymentStacksListAtSubscriptionCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c DeploymentStacksAtSubscriptionClient) DeploymentStacksListAtSubscriptionCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate DeploymentStackOperationPredicate) (result DeploymentStacksListAtSubscriptionCompleteResult, err error) {
	items := make([]DeploymentStack, 0)

	resp, err := c.DeploymentStacksListAtSubscription(ctx, id)
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

	result = DeploymentStacksListAtSubscriptionCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
