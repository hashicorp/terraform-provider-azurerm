package v2workspaceconnectionresource

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceConnectionsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]WorkspaceConnectionPropertiesV2BasicResource
}

type WorkspaceConnectionsListCompleteResult struct {
	Items []WorkspaceConnectionPropertiesV2BasicResource
}

type WorkspaceConnectionsListOperationOptions struct {
	Category *string
	Target   *string
}

func DefaultWorkspaceConnectionsListOperationOptions() WorkspaceConnectionsListOperationOptions {
	return WorkspaceConnectionsListOperationOptions{}
}

func (o WorkspaceConnectionsListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o WorkspaceConnectionsListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o WorkspaceConnectionsListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Category != nil {
		out.Append("category", fmt.Sprintf("%v", *o.Category))
	}
	if o.Target != nil {
		out.Append("target", fmt.Sprintf("%v", *o.Target))
	}
	return &out
}

// WorkspaceConnectionsList ...
func (c V2WorkspaceConnectionResourceClient) WorkspaceConnectionsList(ctx context.Context, id WorkspaceId, options WorkspaceConnectionsListOperationOptions) (result WorkspaceConnectionsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/connections", id.ID()),
		OptionsObject: options,
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
		Values *[]WorkspaceConnectionPropertiesV2BasicResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// WorkspaceConnectionsListComplete retrieves all the results into a single object
func (c V2WorkspaceConnectionResourceClient) WorkspaceConnectionsListComplete(ctx context.Context, id WorkspaceId, options WorkspaceConnectionsListOperationOptions) (WorkspaceConnectionsListCompleteResult, error) {
	return c.WorkspaceConnectionsListCompleteMatchingPredicate(ctx, id, options, WorkspaceConnectionPropertiesV2BasicResourceOperationPredicate{})
}

// WorkspaceConnectionsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c V2WorkspaceConnectionResourceClient) WorkspaceConnectionsListCompleteMatchingPredicate(ctx context.Context, id WorkspaceId, options WorkspaceConnectionsListOperationOptions, predicate WorkspaceConnectionPropertiesV2BasicResourceOperationPredicate) (result WorkspaceConnectionsListCompleteResult, err error) {
	items := make([]WorkspaceConnectionPropertiesV2BasicResource, 0)

	resp, err := c.WorkspaceConnectionsList(ctx, id, options)
	if err != nil {
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

	result = WorkspaceConnectionsListCompleteResult{
		Items: items,
	}
	return
}
