package softwareupdateconfigurationrun

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetByIdOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *SoftwareUpdateConfigurationRun
}

type GetByIdOperationOptions struct {
	ClientRequestId *string
}

func DefaultGetByIdOperationOptions() GetByIdOperationOptions {
	return GetByIdOperationOptions{}
}

func (o GetByIdOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.ClientRequestId != nil {
		out.Append("clientRequestId", fmt.Sprintf("%v", *o.ClientRequestId))
	}
	return &out
}

func (o GetByIdOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o GetByIdOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// GetById ...
func (c SoftwareUpdateConfigurationRunClient) GetById(ctx context.Context, id SoftwareUpdateConfigurationRunId, options GetByIdOperationOptions) (result GetByIdOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          id.ID(),
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
