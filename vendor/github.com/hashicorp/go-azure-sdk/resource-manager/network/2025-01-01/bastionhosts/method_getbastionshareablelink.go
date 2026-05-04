package bastionhosts

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetBastionShareableLinkOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]BastionShareableLink
}

type GetBastionShareableLinkCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []BastionShareableLink
}

type GetBastionShareableLinkCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GetBastionShareableLinkCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GetBastionShareableLink ...
func (c BastionHostsClient) GetBastionShareableLink(ctx context.Context, id BastionHostId, input BastionShareableLinkListRequest) (result GetBastionShareableLinkOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Pager:      &GetBastionShareableLinkCustomPager{},
		Path:       fmt.Sprintf("%s/getShareableLinks", id.ID()),
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
		Values *[]BastionShareableLink `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// GetBastionShareableLinkComplete retrieves all the results into a single object
func (c BastionHostsClient) GetBastionShareableLinkComplete(ctx context.Context, id BastionHostId, input BastionShareableLinkListRequest) (GetBastionShareableLinkCompleteResult, error) {
	return c.GetBastionShareableLinkCompleteMatchingPredicate(ctx, id, input, BastionShareableLinkOperationPredicate{})
}

// GetBastionShareableLinkCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c BastionHostsClient) GetBastionShareableLinkCompleteMatchingPredicate(ctx context.Context, id BastionHostId, input BastionShareableLinkListRequest, predicate BastionShareableLinkOperationPredicate) (result GetBastionShareableLinkCompleteResult, err error) {
	items := make([]BastionShareableLink, 0)

	resp, err := c.GetBastionShareableLink(ctx, id, input)
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

	result = GetBastionShareableLinkCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
