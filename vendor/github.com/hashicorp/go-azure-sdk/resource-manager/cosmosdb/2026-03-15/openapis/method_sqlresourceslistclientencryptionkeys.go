package openapis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SqlResourcesListClientEncryptionKeysOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ClientEncryptionKeyGetResults
}

type SqlResourcesListClientEncryptionKeysCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ClientEncryptionKeyGetResults
}

type SqlResourcesListClientEncryptionKeysCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *SqlResourcesListClientEncryptionKeysCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// SqlResourcesListClientEncryptionKeys ...
func (c OpenapisClient) SqlResourcesListClientEncryptionKeys(ctx context.Context, id SqlDatabaseId) (result SqlResourcesListClientEncryptionKeysOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &SqlResourcesListClientEncryptionKeysCustomPager{},
		Path:       fmt.Sprintf("%s/clientEncryptionKeys", id.ID()),
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
		Values *[]ClientEncryptionKeyGetResults `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// SqlResourcesListClientEncryptionKeysComplete retrieves all the results into a single object
func (c OpenapisClient) SqlResourcesListClientEncryptionKeysComplete(ctx context.Context, id SqlDatabaseId) (SqlResourcesListClientEncryptionKeysCompleteResult, error) {
	return c.SqlResourcesListClientEncryptionKeysCompleteMatchingPredicate(ctx, id, ClientEncryptionKeyGetResultsOperationPredicate{})
}

// SqlResourcesListClientEncryptionKeysCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) SqlResourcesListClientEncryptionKeysCompleteMatchingPredicate(ctx context.Context, id SqlDatabaseId, predicate ClientEncryptionKeyGetResultsOperationPredicate) (result SqlResourcesListClientEncryptionKeysCompleteResult, err error) {
	items := make([]ClientEncryptionKeyGetResults, 0)

	resp, err := c.SqlResourcesListClientEncryptionKeys(ctx, id)
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

	result = SqlResourcesListClientEncryptionKeysCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
