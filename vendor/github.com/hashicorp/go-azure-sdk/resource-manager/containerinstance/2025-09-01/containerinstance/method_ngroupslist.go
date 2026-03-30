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

type NGroupsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]NGroup
}

type NGroupsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []NGroup
}

type NGroupsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *NGroupsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// NGroupsList ...
func (c ContainerInstanceClient) NGroupsList(ctx context.Context, id commonids.SubscriptionId) (result NGroupsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &NGroupsListCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.ContainerInstance/ngroups", id.ID()),
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
		Values *[]NGroup `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// NGroupsListComplete retrieves all the results into a single object
func (c ContainerInstanceClient) NGroupsListComplete(ctx context.Context, id commonids.SubscriptionId) (NGroupsListCompleteResult, error) {
	return c.NGroupsListCompleteMatchingPredicate(ctx, id, NGroupOperationPredicate{})
}

// NGroupsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ContainerInstanceClient) NGroupsListCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate NGroupOperationPredicate) (result NGroupsListCompleteResult, err error) {
	items := make([]NGroup, 0)

	resp, err := c.NGroupsList(ctx, id)
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

	result = NGroupsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
