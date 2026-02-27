package storage

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetSasDefinitionsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SasDefinitionItem
}

type GetSasDefinitionsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SasDefinitionItem
}

type GetSasDefinitionsOperationOptions struct {
	Maxresults *int64
}

func DefaultGetSasDefinitionsOperationOptions() GetSasDefinitionsOperationOptions {
	return GetSasDefinitionsOperationOptions{}
}

func (o GetSasDefinitionsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o GetSasDefinitionsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o GetSasDefinitionsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Maxresults != nil {
		out.Append("maxresults", fmt.Sprintf("%v", *o.Maxresults))
	}
	return &out
}

type GetSasDefinitionsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GetSasDefinitionsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GetSasDefinitions ...
func (c StorageClient) GetSasDefinitions(ctx context.Context, id StorageId, options GetSasDefinitionsOperationOptions) (result GetSasDefinitionsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &GetSasDefinitionsCustomPager{},
		Path:          fmt.Sprintf("%s/sas", id.Path()),
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
		Values *[]SasDefinitionItem `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// GetSasDefinitionsComplete retrieves all the results into a single object
func (c StorageClient) GetSasDefinitionsComplete(ctx context.Context, id StorageId, options GetSasDefinitionsOperationOptions) (GetSasDefinitionsCompleteResult, error) {
	return c.GetSasDefinitionsCompleteMatchingPredicate(ctx, id, options, SasDefinitionItemOperationPredicate{})
}

// GetSasDefinitionsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c StorageClient) GetSasDefinitionsCompleteMatchingPredicate(ctx context.Context, id StorageId, options GetSasDefinitionsOperationOptions, predicate SasDefinitionItemOperationPredicate) (result GetSasDefinitionsCompleteResult, err error) {
	items := make([]SasDefinitionItem, 0)

	resp, err := c.GetSasDefinitions(ctx, id, options)
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

	result = GetSasDefinitionsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
