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

type ListByParentOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DbServer
}

type ListByParentCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DbServer
}

type ListByParentCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByParentCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByParent ...
func (c DbServersClient) ListByParent(ctx context.Context, id CloudExadataInfrastructureId) (result ListByParentOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByParentCustomPager{},
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

// ListByParentComplete retrieves all the results into a single object
func (c DbServersClient) ListByParentComplete(ctx context.Context, id CloudExadataInfrastructureId) (ListByParentCompleteResult, error) {
	return c.ListByParentCompleteMatchingPredicate(ctx, id, DbServerOperationPredicate{})
}

// ListByParentCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c DbServersClient) ListByParentCompleteMatchingPredicate(ctx context.Context, id CloudExadataInfrastructureId, predicate DbServerOperationPredicate) (result ListByParentCompleteResult, err error) {
	items := make([]DbServer, 0)

	resp, err := c.ListByParent(ctx, id)
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

	result = ListByParentCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
