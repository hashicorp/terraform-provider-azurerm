package nodepools

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListBySupercomputerOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]NodePool
}

type ListBySupercomputerCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []NodePool
}

type ListBySupercomputerCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListBySupercomputerCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListBySupercomputer ...
func (c NodePoolsClient) ListBySupercomputer(ctx context.Context, id SupercomputerId) (result ListBySupercomputerOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListBySupercomputerCustomPager{},
		Path:       fmt.Sprintf("%s/nodePools", id.ID()),
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
		Values *[]NodePool `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListBySupercomputerComplete retrieves all the results into a single object
func (c NodePoolsClient) ListBySupercomputerComplete(ctx context.Context, id SupercomputerId) (ListBySupercomputerCompleteResult, error) {
	return c.ListBySupercomputerCompleteMatchingPredicate(ctx, id, NodePoolOperationPredicate{})
}

// ListBySupercomputerCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c NodePoolsClient) ListBySupercomputerCompleteMatchingPredicate(ctx context.Context, id SupercomputerId, predicate NodePoolOperationPredicate) (result ListBySupercomputerCompleteResult, err error) {
	items := make([]NodePool, 0)

	resp, err := c.ListBySupercomputer(ctx, id)
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

	result = ListBySupercomputerCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
