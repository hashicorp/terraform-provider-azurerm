package workbooksapis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkbooksListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Workbook
}

type WorkbooksListByResourceGroupCompleteResult struct {
	Items []Workbook
}

type WorkbooksListByResourceGroupOperationOptions struct {
	CanFetchContent *bool
	Category        *CategoryType
	SourceId        *string
	Tags            *string
}

func DefaultWorkbooksListByResourceGroupOperationOptions() WorkbooksListByResourceGroupOperationOptions {
	return WorkbooksListByResourceGroupOperationOptions{}
}

func (o WorkbooksListByResourceGroupOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o WorkbooksListByResourceGroupOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o WorkbooksListByResourceGroupOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.CanFetchContent != nil {
		out.Append("canFetchContent", fmt.Sprintf("%v", *o.CanFetchContent))
	}
	if o.Category != nil {
		out.Append("category", fmt.Sprintf("%v", *o.Category))
	}
	if o.SourceId != nil {
		out.Append("sourceId", fmt.Sprintf("%v", *o.SourceId))
	}
	if o.Tags != nil {
		out.Append("tags", fmt.Sprintf("%v", *o.Tags))
	}
	return &out
}

// WorkbooksListByResourceGroup ...
func (c WorkbooksAPIsClient) WorkbooksListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId, options WorkbooksListByResourceGroupOperationOptions) (result WorkbooksListByResourceGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/providers/Microsoft.Insights/workbooks", id.ID()),
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
		Values *[]Workbook `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// WorkbooksListByResourceGroupComplete retrieves all the results into a single object
func (c WorkbooksAPIsClient) WorkbooksListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId, options WorkbooksListByResourceGroupOperationOptions) (WorkbooksListByResourceGroupCompleteResult, error) {
	return c.WorkbooksListByResourceGroupCompleteMatchingPredicate(ctx, id, options, WorkbookOperationPredicate{})
}

// WorkbooksListByResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WorkbooksAPIsClient) WorkbooksListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, options WorkbooksListByResourceGroupOperationOptions, predicate WorkbookOperationPredicate) (result WorkbooksListByResourceGroupCompleteResult, err error) {
	items := make([]Workbook, 0)

	resp, err := c.WorkbooksListByResourceGroup(ctx, id, options)
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

	result = WorkbooksListByResourceGroupCompleteResult{
		Items: items,
	}
	return
}
