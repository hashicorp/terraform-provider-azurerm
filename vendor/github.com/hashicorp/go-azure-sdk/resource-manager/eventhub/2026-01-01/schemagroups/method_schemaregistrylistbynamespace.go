package schemagroups

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SchemaRegistryListByNamespaceOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SchemaGroup
}

type SchemaRegistryListByNamespaceCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SchemaGroup
}

type SchemaRegistryListByNamespaceOperationOptions struct {
	Skip *int64
	Top  *int64
}

func DefaultSchemaRegistryListByNamespaceOperationOptions() SchemaRegistryListByNamespaceOperationOptions {
	return SchemaRegistryListByNamespaceOperationOptions{}
}

func (o SchemaRegistryListByNamespaceOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o SchemaRegistryListByNamespaceOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o SchemaRegistryListByNamespaceOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type SchemaRegistryListByNamespaceCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *SchemaRegistryListByNamespaceCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// SchemaRegistryListByNamespace ...
func (c SchemaGroupsClient) SchemaRegistryListByNamespace(ctx context.Context, id NamespaceId, options SchemaRegistryListByNamespaceOperationOptions) (result SchemaRegistryListByNamespaceOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &SchemaRegistryListByNamespaceCustomPager{},
		Path:          fmt.Sprintf("%s/schemaGroups", id.ID()),
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
		Values *[]SchemaGroup `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// SchemaRegistryListByNamespaceComplete retrieves all the results into a single object
func (c SchemaGroupsClient) SchemaRegistryListByNamespaceComplete(ctx context.Context, id NamespaceId, options SchemaRegistryListByNamespaceOperationOptions) (SchemaRegistryListByNamespaceCompleteResult, error) {
	return c.SchemaRegistryListByNamespaceCompleteMatchingPredicate(ctx, id, options, SchemaGroupOperationPredicate{})
}

// SchemaRegistryListByNamespaceCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c SchemaGroupsClient) SchemaRegistryListByNamespaceCompleteMatchingPredicate(ctx context.Context, id NamespaceId, options SchemaRegistryListByNamespaceOperationOptions, predicate SchemaGroupOperationPredicate) (result SchemaRegistryListByNamespaceCompleteResult, err error) {
	items := make([]SchemaGroup, 0)

	resp, err := c.SchemaRegistryListByNamespace(ctx, id, options)
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

	result = SchemaRegistryListByNamespaceCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
