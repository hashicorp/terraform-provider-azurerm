package deletedconfigurationstores

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

type ConfigurationStoresListDeletedOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DeletedConfigurationStore
}

type ConfigurationStoresListDeletedCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DeletedConfigurationStore
}

type ConfigurationStoresListDeletedCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ConfigurationStoresListDeletedCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ConfigurationStoresListDeleted ...
func (c DeletedConfigurationStoresClient) ConfigurationStoresListDeleted(ctx context.Context, id commonids.SubscriptionId) (result ConfigurationStoresListDeletedOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ConfigurationStoresListDeletedCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.AppConfiguration/deletedConfigurationStores", id.ID()),
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
		Values *[]DeletedConfigurationStore `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ConfigurationStoresListDeletedComplete retrieves all the results into a single object
func (c DeletedConfigurationStoresClient) ConfigurationStoresListDeletedComplete(ctx context.Context, id commonids.SubscriptionId) (ConfigurationStoresListDeletedCompleteResult, error) {
	return c.ConfigurationStoresListDeletedCompleteMatchingPredicate(ctx, id, DeletedConfigurationStoreOperationPredicate{})
}

// ConfigurationStoresListDeletedCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c DeletedConfigurationStoresClient) ConfigurationStoresListDeletedCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate DeletedConfigurationStoreOperationPredicate) (result ConfigurationStoresListDeletedCompleteResult, err error) {
	items := make([]DeletedConfigurationStore, 0)

	resp, err := c.ConfigurationStoresListDeleted(ctx, id)
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

	result = ConfigurationStoresListDeletedCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
