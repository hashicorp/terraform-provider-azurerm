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

type GetBuildDatabaseConnectionsWithDetailsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DatabaseConnection
}

type GetBuildDatabaseConnectionsWithDetailsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DatabaseConnection
}

type GetBuildDatabaseConnectionsWithDetailsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GetBuildDatabaseConnectionsWithDetailsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GetBuildDatabaseConnectionsWithDetails ...
func (c StaticSitesClient) GetBuildDatabaseConnectionsWithDetails(ctx context.Context, id BuildId) (result GetBuildDatabaseConnectionsWithDetailsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Pager:      &GetBuildDatabaseConnectionsWithDetailsCustomPager{},
		Path:       fmt.Sprintf("%s/showDatabaseConnections", id.ID()),
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

// GetBuildDatabaseConnectionsWithDetailsComplete retrieves all the results into a single object
func (c StaticSitesClient) GetBuildDatabaseConnectionsWithDetailsComplete(ctx context.Context, id BuildId) (GetBuildDatabaseConnectionsWithDetailsCompleteResult, error) {
	return c.GetBuildDatabaseConnectionsWithDetailsCompleteMatchingPredicate(ctx, id, DatabaseConnectionOperationPredicate{})
}

// GetBuildDatabaseConnectionsWithDetailsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c StaticSitesClient) GetBuildDatabaseConnectionsWithDetailsCompleteMatchingPredicate(ctx context.Context, id BuildId, predicate DatabaseConnectionOperationPredicate) (result GetBuildDatabaseConnectionsWithDetailsCompleteResult, err error) {
	items := make([]DatabaseConnection, 0)

	resp, err := c.GetBuildDatabaseConnectionsWithDetails(ctx, id)
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

	result = GetBuildDatabaseConnectionsWithDetailsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
