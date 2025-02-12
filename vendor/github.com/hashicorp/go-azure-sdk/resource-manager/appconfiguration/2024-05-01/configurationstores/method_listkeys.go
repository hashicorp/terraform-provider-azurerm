package configurationstores

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListKeysOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ApiKey
}

type ListKeysCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ApiKey
}

type ListKeysCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListKeysCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListKeys ...
func (c ConfigurationStoresClient) ListKeys(ctx context.Context, id ConfigurationStoreId) (result ListKeysOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Pager:      &ListKeysCustomPager{},
		Path:       fmt.Sprintf("%s/listKeys", id.ID()),
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
		Values *[]ApiKey `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListKeysComplete retrieves all the results into a single object
func (c ConfigurationStoresClient) ListKeysComplete(ctx context.Context, id ConfigurationStoreId) (ListKeysCompleteResult, error) {
	return c.ListKeysCompleteMatchingPredicate(ctx, id, ApiKeyOperationPredicate{})
}

// ListKeysCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ConfigurationStoresClient) ListKeysCompleteMatchingPredicate(ctx context.Context, id ConfigurationStoreId, predicate ApiKeyOperationPredicate) (result ListKeysCompleteResult, err error) {
	items := make([]ApiKey, 0)

	resp, err := c.ListKeys(ctx, id)
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

	result = ListKeysCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
