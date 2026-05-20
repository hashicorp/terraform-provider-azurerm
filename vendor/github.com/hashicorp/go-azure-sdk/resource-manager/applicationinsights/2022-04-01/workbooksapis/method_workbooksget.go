package workbooksapis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkbooksGetOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *Workbook
}

type WorkbooksGetOperationOptions struct {
	CanFetchContent *bool
}

func DefaultWorkbooksGetOperationOptions() WorkbooksGetOperationOptions {
	return WorkbooksGetOperationOptions{}
}

func (o WorkbooksGetOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o WorkbooksGetOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o WorkbooksGetOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.CanFetchContent != nil {
		out.Append("canFetchContent", fmt.Sprintf("%v", *o.CanFetchContent))
	}
	return &out
}

// WorkbooksGet ...
func (c WorkbooksAPIsClient) WorkbooksGet(ctx context.Context, id WorkbookId, options WorkbooksGetOperationOptions) (result WorkbooksGetOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Path:          id.ID(),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	var model Workbook
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
