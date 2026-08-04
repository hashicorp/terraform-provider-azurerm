package suppressionlists

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SuppressionListAddressesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SuppressionListAddressResource
}

type SuppressionListAddressesListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SuppressionListAddressResource
}

type SuppressionListAddressesListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *SuppressionListAddressesListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// SuppressionListAddressesList ...
func (c SuppressionListsClient) SuppressionListAddressesList(ctx context.Context, id SuppressionListId) (result SuppressionListAddressesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &SuppressionListAddressesListCustomPager{},
		Path:       fmt.Sprintf("%s/suppressionListAddresses", id.ID()),
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
		Values *[]SuppressionListAddressResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// SuppressionListAddressesListComplete retrieves all the results into a single object
func (c SuppressionListsClient) SuppressionListAddressesListComplete(ctx context.Context, id SuppressionListId) (SuppressionListAddressesListCompleteResult, error) {
	return c.SuppressionListAddressesListCompleteMatchingPredicate(ctx, id, SuppressionListAddressResourceOperationPredicate{})
}

// SuppressionListAddressesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c SuppressionListsClient) SuppressionListAddressesListCompleteMatchingPredicate(ctx context.Context, id SuppressionListId, predicate SuppressionListAddressResourceOperationPredicate) (result SuppressionListAddressesListCompleteResult, err error) {
	items := make([]SuppressionListAddressResource, 0)

	resp, err := c.SuppressionListAddressesList(ctx, id)
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

	result = SuppressionListAddressesListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
