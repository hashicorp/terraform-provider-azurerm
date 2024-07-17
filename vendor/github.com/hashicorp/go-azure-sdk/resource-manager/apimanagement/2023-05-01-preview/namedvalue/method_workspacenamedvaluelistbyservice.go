package namedvalue

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceNamedValueListByServiceOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]NamedValueContract
}

type WorkspaceNamedValueListByServiceCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []NamedValueContract
}

type WorkspaceNamedValueListByServiceOperationOptions struct {
	Filter                  *string
	IsKeyVaultRefreshFailed *KeyVaultRefreshState
	Skip                    *int64
	Top                     *int64
}

func DefaultWorkspaceNamedValueListByServiceOperationOptions() WorkspaceNamedValueListByServiceOperationOptions {
	return WorkspaceNamedValueListByServiceOperationOptions{}
}

func (o WorkspaceNamedValueListByServiceOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o WorkspaceNamedValueListByServiceOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o WorkspaceNamedValueListByServiceOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.IsKeyVaultRefreshFailed != nil {
		out.Append("isKeyVaultRefreshFailed", fmt.Sprintf("%v", *o.IsKeyVaultRefreshFailed))
	}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type WorkspaceNamedValueListByServiceCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *WorkspaceNamedValueListByServiceCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// WorkspaceNamedValueListByService ...
func (c NamedValueClient) WorkspaceNamedValueListByService(ctx context.Context, id WorkspaceId, options WorkspaceNamedValueListByServiceOperationOptions) (result WorkspaceNamedValueListByServiceOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &WorkspaceNamedValueListByServiceCustomPager{},
		Path:          fmt.Sprintf("%s/namedValues", id.ID()),
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
		Values *[]NamedValueContract `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// WorkspaceNamedValueListByServiceComplete retrieves all the results into a single object
func (c NamedValueClient) WorkspaceNamedValueListByServiceComplete(ctx context.Context, id WorkspaceId, options WorkspaceNamedValueListByServiceOperationOptions) (WorkspaceNamedValueListByServiceCompleteResult, error) {
	return c.WorkspaceNamedValueListByServiceCompleteMatchingPredicate(ctx, id, options, NamedValueContractOperationPredicate{})
}

// WorkspaceNamedValueListByServiceCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c NamedValueClient) WorkspaceNamedValueListByServiceCompleteMatchingPredicate(ctx context.Context, id WorkspaceId, options WorkspaceNamedValueListByServiceOperationOptions, predicate NamedValueContractOperationPredicate) (result WorkspaceNamedValueListByServiceCompleteResult, err error) {
	items := make([]NamedValueContract, 0)

	resp, err := c.WorkspaceNamedValueListByService(ctx, id, options)
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

	result = WorkspaceNamedValueListByServiceCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
