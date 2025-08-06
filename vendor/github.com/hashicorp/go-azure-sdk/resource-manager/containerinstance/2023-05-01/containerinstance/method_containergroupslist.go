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

type ContainerGroupsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ListResultContainerGroup
}

type ContainerGroupsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ListResultContainerGroup
}

type ContainerGroupsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ContainerGroupsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ContainerGroupsList ...
func (c ContainerInstanceClient) ContainerGroupsList(ctx context.Context, id commonids.SubscriptionId) (result ContainerGroupsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ContainerGroupsListCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.ContainerInstance/containerGroups", id.ID()),
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
		Values *[]ListResultContainerGroup `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ContainerGroupsListComplete retrieves all the results into a single object
func (c ContainerInstanceClient) ContainerGroupsListComplete(ctx context.Context, id commonids.SubscriptionId) (ContainerGroupsListCompleteResult, error) {
	return c.ContainerGroupsListCompleteMatchingPredicate(ctx, id, ListResultContainerGroupOperationPredicate{})
}

// ContainerGroupsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ContainerInstanceClient) ContainerGroupsListCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate ListResultContainerGroupOperationPredicate) (result ContainerGroupsListCompleteResult, err error) {
	items := make([]ListResultContainerGroup, 0)

	resp, err := c.ContainerGroupsList(ctx, id)
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

	result = ContainerGroupsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
