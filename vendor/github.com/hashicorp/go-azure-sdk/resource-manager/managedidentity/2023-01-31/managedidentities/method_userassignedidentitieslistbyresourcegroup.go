package managedidentities

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

type UserAssignedIdentitiesListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Identity
}

type UserAssignedIdentitiesListByResourceGroupCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Identity
}

type UserAssignedIdentitiesListByResourceGroupCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *UserAssignedIdentitiesListByResourceGroupCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// UserAssignedIdentitiesListByResourceGroup ...
func (c ManagedIdentitiesClient) UserAssignedIdentitiesListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (result UserAssignedIdentitiesListByResourceGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &UserAssignedIdentitiesListByResourceGroupCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.ManagedIdentity/userAssignedIdentities", id.ID()),
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
		Values *[]Identity `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// UserAssignedIdentitiesListByResourceGroupComplete retrieves all the results into a single object
func (c ManagedIdentitiesClient) UserAssignedIdentitiesListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (UserAssignedIdentitiesListByResourceGroupCompleteResult, error) {
	return c.UserAssignedIdentitiesListByResourceGroupCompleteMatchingPredicate(ctx, id, IdentityOperationPredicate{})
}

// UserAssignedIdentitiesListByResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ManagedIdentitiesClient) UserAssignedIdentitiesListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate IdentityOperationPredicate) (result UserAssignedIdentitiesListByResourceGroupCompleteResult, err error) {
	items := make([]Identity, 0)

	resp, err := c.UserAssignedIdentitiesListByResourceGroup(ctx, id)
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

	result = UserAssignedIdentitiesListByResourceGroupCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
