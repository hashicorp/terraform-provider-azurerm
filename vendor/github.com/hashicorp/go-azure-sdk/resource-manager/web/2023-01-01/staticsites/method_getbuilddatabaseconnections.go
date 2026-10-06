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

type GetBuildDatabaseConnectionsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DatabaseConnection
}

type GetBuildDatabaseConnectionsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DatabaseConnection
}

type GetBuildDatabaseConnectionsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GetBuildDatabaseConnectionsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GetBuildDatabaseConnections ...
func (c StaticSitesClient) GetBuildDatabaseConnections(ctx context.Context, id BuildId) (result GetBuildDatabaseConnectionsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &GetBuildDatabaseConnectionsCustomPager{},
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

// GetBuildDatabaseConnectionsComplete retrieves all the results into a single object
func (c StaticSitesClient) GetBuildDatabaseConnectionsComplete(ctx context.Context, id BuildId) (GetBuildDatabaseConnectionsCompleteResult, error) {
	return c.GetBuildDatabaseConnectionsCompleteMatchingPredicate(ctx, id, DatabaseConnectionOperationPredicate{})
}

// GetBuildDatabaseConnectionsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c StaticSitesClient) GetBuildDatabaseConnectionsCompleteMatchingPredicate(ctx context.Context, id BuildId, predicate DatabaseConnectionOperationPredicate) (result GetBuildDatabaseConnectionsCompleteResult, err error) {
	items := make([]DatabaseConnection, 0)

	resp, err := c.GetBuildDatabaseConnections(ctx, id)
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

	result = GetBuildDatabaseConnectionsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
