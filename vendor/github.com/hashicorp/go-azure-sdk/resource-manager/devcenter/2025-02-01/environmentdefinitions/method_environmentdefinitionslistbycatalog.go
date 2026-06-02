package environmentdefinitions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EnvironmentDefinitionsListByCatalogOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]EnvironmentDefinition
}

type EnvironmentDefinitionsListByCatalogCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []EnvironmentDefinition
}

type EnvironmentDefinitionsListByCatalogOperationOptions struct {
	Top *int64
}

func DefaultEnvironmentDefinitionsListByCatalogOperationOptions() EnvironmentDefinitionsListByCatalogOperationOptions {
	return EnvironmentDefinitionsListByCatalogOperationOptions{}
}

func (o EnvironmentDefinitionsListByCatalogOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o EnvironmentDefinitionsListByCatalogOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o EnvironmentDefinitionsListByCatalogOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type EnvironmentDefinitionsListByCatalogCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *EnvironmentDefinitionsListByCatalogCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// EnvironmentDefinitionsListByCatalog ...
func (c EnvironmentDefinitionsClient) EnvironmentDefinitionsListByCatalog(ctx context.Context, id DevCenterCatalogId, options EnvironmentDefinitionsListByCatalogOperationOptions) (result EnvironmentDefinitionsListByCatalogOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &EnvironmentDefinitionsListByCatalogCustomPager{},
		Path:          fmt.Sprintf("%s/environmentDefinitions", id.ID()),
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
		Values *[]EnvironmentDefinition `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// EnvironmentDefinitionsListByCatalogComplete retrieves all the results into a single object
func (c EnvironmentDefinitionsClient) EnvironmentDefinitionsListByCatalogComplete(ctx context.Context, id DevCenterCatalogId, options EnvironmentDefinitionsListByCatalogOperationOptions) (EnvironmentDefinitionsListByCatalogCompleteResult, error) {
	return c.EnvironmentDefinitionsListByCatalogCompleteMatchingPredicate(ctx, id, options, EnvironmentDefinitionOperationPredicate{})
}

// EnvironmentDefinitionsListByCatalogCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c EnvironmentDefinitionsClient) EnvironmentDefinitionsListByCatalogCompleteMatchingPredicate(ctx context.Context, id DevCenterCatalogId, options EnvironmentDefinitionsListByCatalogOperationOptions, predicate EnvironmentDefinitionOperationPredicate) (result EnvironmentDefinitionsListByCatalogCompleteResult, err error) {
	items := make([]EnvironmentDefinition, 0)

	resp, err := c.EnvironmentDefinitionsListByCatalog(ctx, id, options)
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

	result = EnvironmentDefinitionsListByCatalogCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
