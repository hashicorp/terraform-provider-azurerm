package job

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetOutputOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *string
}

type GetOutputOperationOptions struct {
	ClientRequestId *string
}

func DefaultGetOutputOperationOptions() GetOutputOperationOptions {
	return GetOutputOperationOptions{}
}

func (o GetOutputOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.ClientRequestId != nil {
		out.Append("clientRequestId", fmt.Sprintf("%v", *o.ClientRequestId))
	}
	return &out
}

func (o GetOutputOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o GetOutputOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// GetOutput ...
func (c JobClient) GetOutput(ctx context.Context, id JobId, options GetOutputOperationOptions) (result GetOutputOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "text/plain",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/output", id.ID()),
		OptionsObject: options,
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

	if err = resp.Unmarshal(&result.Model); err != nil {
		return
	}

	return
}
