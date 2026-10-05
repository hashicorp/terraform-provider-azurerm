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

type WorkbooksCreateOrUpdateOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *Workbook
}

type WorkbooksCreateOrUpdateOperationOptions struct {
	SourceId *string
}

func DefaultWorkbooksCreateOrUpdateOperationOptions() WorkbooksCreateOrUpdateOperationOptions {
	return WorkbooksCreateOrUpdateOperationOptions{}
}

func (o WorkbooksCreateOrUpdateOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o WorkbooksCreateOrUpdateOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o WorkbooksCreateOrUpdateOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.SourceId != nil {
		out.Append("sourceId", fmt.Sprintf("%v", *o.SourceId))
	}
	return &out
}

// WorkbooksCreateOrUpdate ...
func (c WorkbooksAPIsClient) WorkbooksCreateOrUpdate(ctx context.Context, id WorkbookId, input Workbook, options WorkbooksCreateOrUpdateOperationOptions) (result WorkbooksCreateOrUpdateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusCreated,
			http.StatusOK,
		},
		HttpMethod:    http.MethodPut,
		OptionsObject: options,
		Path:          id.ID(),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	if err = req.Marshal(input); err != nil {
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
