package typefields

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FieldsListByTypeOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]TypeField
}

type FieldsListByTypeCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []TypeField
}

type FieldsListByTypeCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *FieldsListByTypeCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// FieldsListByType ...
func (c TypeFieldsClient) FieldsListByType(ctx context.Context, id TypeId) (result FieldsListByTypeOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &FieldsListByTypeCustomPager{},
		Path:       fmt.Sprintf("%s/fields", id.ID()),
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
		Values *[]TypeField `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// FieldsListByTypeComplete retrieves all the results into a single object
func (c TypeFieldsClient) FieldsListByTypeComplete(ctx context.Context, id TypeId) (FieldsListByTypeCompleteResult, error) {
	return c.FieldsListByTypeCompleteMatchingPredicate(ctx, id, TypeFieldOperationPredicate{})
}

// FieldsListByTypeCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c TypeFieldsClient) FieldsListByTypeCompleteMatchingPredicate(ctx context.Context, id TypeId, predicate TypeFieldOperationPredicate) (result FieldsListByTypeCompleteResult, err error) {
	items := make([]TypeField, 0)

	resp, err := c.FieldsListByType(ctx, id)
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

	result = FieldsListByTypeCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
