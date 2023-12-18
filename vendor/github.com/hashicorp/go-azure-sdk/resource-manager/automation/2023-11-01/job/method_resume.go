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

type ResumeOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
}

type ResumeOperationOptions struct {
	ClientRequestId *string
}

func DefaultResumeOperationOptions() ResumeOperationOptions {
	return ResumeOperationOptions{}
}

func (o ResumeOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.ClientRequestId != nil {
		out.Append("clientRequestId", fmt.Sprintf("%v", *o.ClientRequestId))
	}
	return &out
}

func (o ResumeOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ResumeOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// Resume ...
func (c JobClient) Resume(ctx context.Context, id JobId, options ResumeOperationOptions) (result ResumeOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		Path:          fmt.Sprintf("%s/resume", id.ID()),
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

	return
}
