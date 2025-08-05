package groupuser

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceGroupUserListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]UserContract
}

type WorkspaceGroupUserListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []UserContract
}

type WorkspaceGroupUserListOperationOptions struct {
	Filter *string
	Skip   *int64
	Top    *int64
}

func DefaultWorkspaceGroupUserListOperationOptions() WorkspaceGroupUserListOperationOptions {
	return WorkspaceGroupUserListOperationOptions{}
}

func (o WorkspaceGroupUserListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o WorkspaceGroupUserListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o WorkspaceGroupUserListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type WorkspaceGroupUserListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *WorkspaceGroupUserListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// WorkspaceGroupUserList ...
func (c GroupUserClient) WorkspaceGroupUserList(ctx context.Context, id WorkspaceGroupId, options WorkspaceGroupUserListOperationOptions) (result WorkspaceGroupUserListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &WorkspaceGroupUserListCustomPager{},
		Path:          fmt.Sprintf("%s/users", id.ID()),
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
		Values *[]UserContract `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// WorkspaceGroupUserListComplete retrieves all the results into a single object
func (c GroupUserClient) WorkspaceGroupUserListComplete(ctx context.Context, id WorkspaceGroupId, options WorkspaceGroupUserListOperationOptions) (WorkspaceGroupUserListCompleteResult, error) {
	return c.WorkspaceGroupUserListCompleteMatchingPredicate(ctx, id, options, UserContractOperationPredicate{})
}

// WorkspaceGroupUserListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c GroupUserClient) WorkspaceGroupUserListCompleteMatchingPredicate(ctx context.Context, id WorkspaceGroupId, options WorkspaceGroupUserListOperationOptions, predicate UserContractOperationPredicate) (result WorkspaceGroupUserListCompleteResult, err error) {
	items := make([]UserContract, 0)

	resp, err := c.WorkspaceGroupUserList(ctx, id, options)
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

	result = WorkspaceGroupUserListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
