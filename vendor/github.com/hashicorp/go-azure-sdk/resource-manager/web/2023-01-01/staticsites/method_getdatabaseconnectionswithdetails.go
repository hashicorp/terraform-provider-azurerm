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

type GetDatabaseConnectionsWithDetailsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DatabaseConnection
}

type GetDatabaseConnectionsWithDetailsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DatabaseConnection
}

type GetDatabaseConnectionsWithDetailsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GetDatabaseConnectionsWithDetailsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GetDatabaseConnectionsWithDetails ...
func (c StaticSitesClient) GetDatabaseConnectionsWithDetails(ctx context.Context, id StaticSiteId) (result GetDatabaseConnectionsWithDetailsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Pager:      &GetDatabaseConnectionsWithDetailsCustomPager{},
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

// GetDatabaseConnectionsWithDetailsComplete retrieves all the results into a single object
func (c StaticSitesClient) GetDatabaseConnectionsWithDetailsComplete(ctx context.Context, id StaticSiteId) (GetDatabaseConnectionsWithDetailsCompleteResult, error) {
	return c.GetDatabaseConnectionsWithDetailsCompleteMatchingPredicate(ctx, id, DatabaseConnectionOperationPredicate{})
}

// GetDatabaseConnectionsWithDetailsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c StaticSitesClient) GetDatabaseConnectionsWithDetailsCompleteMatchingPredicate(ctx context.Context, id StaticSiteId, predicate DatabaseConnectionOperationPredicate) (result GetDatabaseConnectionsWithDetailsCompleteResult, err error) {
	items := make([]DatabaseConnection, 0)

	resp, err := c.GetDatabaseConnectionsWithDetails(ctx, id)
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

	result = GetDatabaseConnectionsWithDetailsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
