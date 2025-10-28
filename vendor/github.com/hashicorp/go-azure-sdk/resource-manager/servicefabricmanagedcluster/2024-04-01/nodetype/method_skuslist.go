package nodetype

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SkusListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]NodeTypeAvailableSku
}

type SkusListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []NodeTypeAvailableSku
}

type SkusListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *SkusListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// SkusList ...
func (c NodeTypeClient) SkusList(ctx context.Context, id NodeTypeId) (result SkusListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &SkusListCustomPager{},
		Path:       fmt.Sprintf("%s/skus", id.ID()),
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
		Values *[]NodeTypeAvailableSku `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// SkusListComplete retrieves all the results into a single object
func (c NodeTypeClient) SkusListComplete(ctx context.Context, id NodeTypeId) (SkusListCompleteResult, error) {
	return c.SkusListCompleteMatchingPredicate(ctx, id, NodeTypeAvailableSkuOperationPredicate{})
}

// SkusListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c NodeTypeClient) SkusListCompleteMatchingPredicate(ctx context.Context, id NodeTypeId, predicate NodeTypeAvailableSkuOperationPredicate) (result SkusListCompleteResult, err error) {
	items := make([]NodeTypeAvailableSku, 0)

	resp, err := c.SkusList(ctx, id)
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

	result = SkusListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
