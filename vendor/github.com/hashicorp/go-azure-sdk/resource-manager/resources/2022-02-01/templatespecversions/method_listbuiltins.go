package templatespecversions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListBuiltInsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]TemplateSpecVersion
}

type ListBuiltInsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []TemplateSpecVersion
}

type ListBuiltInsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListBuiltInsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListBuiltIns ...
func (c TemplateSpecVersionsClient) ListBuiltIns(ctx context.Context, id BuiltInTemplateSpecId) (result ListBuiltInsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListBuiltInsCustomPager{},
		Path:       fmt.Sprintf("%s/versions", id.ID()),
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
		Values *[]TemplateSpecVersion `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListBuiltInsComplete retrieves all the results into a single object
func (c TemplateSpecVersionsClient) ListBuiltInsComplete(ctx context.Context, id BuiltInTemplateSpecId) (ListBuiltInsCompleteResult, error) {
	return c.ListBuiltInsCompleteMatchingPredicate(ctx, id, TemplateSpecVersionOperationPredicate{})
}

// ListBuiltInsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c TemplateSpecVersionsClient) ListBuiltInsCompleteMatchingPredicate(ctx context.Context, id BuiltInTemplateSpecId, predicate TemplateSpecVersionOperationPredicate) (result ListBuiltInsCompleteResult, err error) {
	items := make([]TemplateSpecVersion, 0)

	resp, err := c.ListBuiltIns(ctx, id)
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

	result = ListBuiltInsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
