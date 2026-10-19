package manageddatabases

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

type ListInaccessibleByInstanceOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ManagedDatabase
}

type ListInaccessibleByInstanceCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ManagedDatabase
}

type ListInaccessibleByInstanceCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListInaccessibleByInstanceCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListInaccessibleByInstance ...
func (c ManagedDatabasesClient) ListInaccessibleByInstance(ctx context.Context, id commonids.SqlManagedInstanceId) (result ListInaccessibleByInstanceOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListInaccessibleByInstanceCustomPager{},
		Path:       fmt.Sprintf("%s/inaccessibleManagedDatabases", id.ID()),
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
		Values *[]ManagedDatabase `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListInaccessibleByInstanceComplete retrieves all the results into a single object
func (c ManagedDatabasesClient) ListInaccessibleByInstanceComplete(ctx context.Context, id commonids.SqlManagedInstanceId) (ListInaccessibleByInstanceCompleteResult, error) {
	return c.ListInaccessibleByInstanceCompleteMatchingPredicate(ctx, id, ManagedDatabaseOperationPredicate{})
}

// ListInaccessibleByInstanceCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ManagedDatabasesClient) ListInaccessibleByInstanceCompleteMatchingPredicate(ctx context.Context, id commonids.SqlManagedInstanceId, predicate ManagedDatabaseOperationPredicate) (result ListInaccessibleByInstanceCompleteResult, err error) {
	items := make([]ManagedDatabase, 0)

	resp, err := c.ListInaccessibleByInstance(ctx, id)
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

	result = ListInaccessibleByInstanceCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
