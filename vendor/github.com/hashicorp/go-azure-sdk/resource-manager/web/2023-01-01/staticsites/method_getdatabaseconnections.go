package staticsites

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetDatabaseConnectionsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DatabaseConnection
}

type GetDatabaseConnectionsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DatabaseConnection
}

type GetDatabaseConnectionsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GetDatabaseConnectionsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GetDatabaseConnections ...
func (c StaticSitesClient) GetDatabaseConnections(ctx context.Context, id StaticSiteId) (result GetDatabaseConnectionsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &GetDatabaseConnectionsCustomPager{},
		Path:       fmt.Sprintf("%s/databaseConnections", id.ID()),
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
		Values *[]DatabaseConnection `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// GetDatabaseConnectionsComplete retrieves all the results into a single object
func (c StaticSitesClient) GetDatabaseConnectionsComplete(ctx context.Context, id StaticSiteId) (GetDatabaseConnectionsCompleteResult, error) {
	return c.GetDatabaseConnectionsCompleteMatchingPredicate(ctx, id, DatabaseConnectionOperationPredicate{})
}

// GetDatabaseConnectionsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c StaticSitesClient) GetDatabaseConnectionsCompleteMatchingPredicate(ctx context.Context, id StaticSiteId, predicate DatabaseConnectionOperationPredicate) (result GetDatabaseConnectionsCompleteResult, err error) {
	items := make([]DatabaseConnection, 0)

	resp, err := c.GetDatabaseConnections(ctx, id)
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

	result = GetDatabaseConnectionsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
