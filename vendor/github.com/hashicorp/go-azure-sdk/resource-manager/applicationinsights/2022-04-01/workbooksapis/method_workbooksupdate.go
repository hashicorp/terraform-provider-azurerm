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

type WorkbooksUpdateOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *Workbook
}

type WorkbooksUpdateOperationOptions struct {
	SourceId *string
}

func DefaultWorkbooksUpdateOperationOptions() WorkbooksUpdateOperationOptions {
	return WorkbooksUpdateOperationOptions{}
}

func (o WorkbooksUpdateOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o WorkbooksUpdateOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o WorkbooksUpdateOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.SourceId != nil {
		out.Append("sourceId", fmt.Sprintf("%v", *o.SourceId))
	}
	return &out
}

// WorkbooksUpdate ...
func (c WorkbooksAPIsClient) WorkbooksUpdate(ctx context.Context, id WorkbookId, input WorkbookUpdateParameters, options WorkbooksUpdateOperationOptions) (result WorkbooksUpdateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusCreated,
			http.StatusOK,
		},
		HttpMethod:    http.MethodPatch,
		Path:          id.ID(),
		OptionsObject: options,
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

	if err = resp.Unmarshal(&result.Model); err != nil {
		return
	}

	return
}
