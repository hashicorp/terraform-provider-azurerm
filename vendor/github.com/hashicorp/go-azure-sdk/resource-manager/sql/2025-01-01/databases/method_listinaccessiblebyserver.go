package databases

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

type ListInaccessibleByServerOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Database
}

type ListInaccessibleByServerCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Database
}

type ListInaccessibleByServerCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListInaccessibleByServerCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListInaccessibleByServer ...
func (c DatabasesClient) ListInaccessibleByServer(ctx context.Context, id commonids.SqlServerId) (result ListInaccessibleByServerOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListInaccessibleByServerCustomPager{},
		Path:       fmt.Sprintf("%s/inaccessibleDatabases", id.ID()),
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
		Values *[]Database `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListInaccessibleByServerComplete retrieves all the results into a single object
func (c DatabasesClient) ListInaccessibleByServerComplete(ctx context.Context, id commonids.SqlServerId) (ListInaccessibleByServerCompleteResult, error) {
	return c.ListInaccessibleByServerCompleteMatchingPredicate(ctx, id, DatabaseOperationPredicate{})
}

// ListInaccessibleByServerCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c DatabasesClient) ListInaccessibleByServerCompleteMatchingPredicate(ctx context.Context, id commonids.SqlServerId, predicate DatabaseOperationPredicate) (result ListInaccessibleByServerCompleteResult, err error) {
	items := make([]Database, 0)

	resp, err := c.ListInaccessibleByServer(ctx, id)
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

	result = ListInaccessibleByServerCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
