package origins

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByEndpointOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Origin
}

type ListByEndpointCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Origin
}

type ListByEndpointCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByEndpointCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByEndpoint ...
func (c OriginsClient) ListByEndpoint(ctx context.Context, id EndpointId) (result ListByEndpointOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByEndpointCustomPager{},
		Path:       fmt.Sprintf("%s/origins", id.ID()),
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
		Values *[]Origin `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByEndpointComplete retrieves all the results into a single object
func (c OriginsClient) ListByEndpointComplete(ctx context.Context, id EndpointId) (ListByEndpointCompleteResult, error) {
	return c.ListByEndpointCompleteMatchingPredicate(ctx, id, OriginOperationPredicate{})
}

// ListByEndpointCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OriginsClient) ListByEndpointCompleteMatchingPredicate(ctx context.Context, id EndpointId, predicate OriginOperationPredicate) (result ListByEndpointCompleteResult, err error) {
	items := make([]Origin, 0)

	resp, err := c.ListByEndpoint(ctx, id)
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

	result = ListByEndpointCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
