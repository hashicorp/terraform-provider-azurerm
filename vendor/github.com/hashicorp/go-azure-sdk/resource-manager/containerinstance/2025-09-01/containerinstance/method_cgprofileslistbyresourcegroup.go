package containerinstance

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

type CGProfilesListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ContainerGroupProfile
}

type CGProfilesListByResourceGroupCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ContainerGroupProfile
}

type CGProfilesListByResourceGroupCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *CGProfilesListByResourceGroupCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// CGProfilesListByResourceGroup ...
func (c ContainerInstanceClient) CGProfilesListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (result CGProfilesListByResourceGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &CGProfilesListByResourceGroupCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.ContainerInstance/containerGroupProfiles", id.ID()),
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
		Values *[]ContainerGroupProfile `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// CGProfilesListByResourceGroupComplete retrieves all the results into a single object
func (c ContainerInstanceClient) CGProfilesListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (CGProfilesListByResourceGroupCompleteResult, error) {
	return c.CGProfilesListByResourceGroupCompleteMatchingPredicate(ctx, id, ContainerGroupProfileOperationPredicate{})
}

// CGProfilesListByResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ContainerInstanceClient) CGProfilesListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate ContainerGroupProfileOperationPredicate) (result CGProfilesListByResourceGroupCompleteResult, err error) {
	items := make([]ContainerGroupProfile, 0)

	resp, err := c.CGProfilesListByResourceGroup(ctx, id)
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

	result = CGProfilesListByResourceGroupCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
