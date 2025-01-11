package dbservers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByCloudExadataInfrastructureOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DbServer
}

type ListByCloudExadataInfrastructureCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DbServer
}

type ListByCloudExadataInfrastructureCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByCloudExadataInfrastructureCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByCloudExadataInfrastructure ...
func (c DbServersClient) ListByCloudExadataInfrastructure(ctx context.Context, id CloudExadataInfrastructureId) (result ListByCloudExadataInfrastructureOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByCloudExadataInfrastructureCustomPager{},
		Path:       fmt.Sprintf("%s/dbServers", id.ID()),
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
		Values *[]DbServer `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByCloudExadataInfrastructureComplete retrieves all the results into a single object
func (c DbServersClient) ListByCloudExadataInfrastructureComplete(ctx context.Context, id CloudExadataInfrastructureId) (ListByCloudExadataInfrastructureCompleteResult, error) {
	return c.ListByCloudExadataInfrastructureCompleteMatchingPredicate(ctx, id, DbServerOperationPredicate{})
}

// ListByCloudExadataInfrastructureCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c DbServersClient) ListByCloudExadataInfrastructureCompleteMatchingPredicate(ctx context.Context, id CloudExadataInfrastructureId, predicate DbServerOperationPredicate) (result ListByCloudExadataInfrastructureCompleteResult, err error) {
	items := make([]DbServer, 0)

	resp, err := c.ListByCloudExadataInfrastructure(ctx, id)
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

	result = ListByCloudExadataInfrastructureCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
