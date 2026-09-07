package jobs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobPatchOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
}

type JobPatchOperationOptions struct {
	ClientRequestId       *string
	IfMatch               *string
	IfModifiedSince       *string
	IfNoneMatch           *string
	IfUnmodifiedSince     *string
	OcpDate               *string
	ReturnClientRequestId *bool
	Timeout               *int64
}

func DefaultJobPatchOperationOptions() JobPatchOperationOptions {
	return JobPatchOperationOptions{}
}

func (o JobPatchOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.ClientRequestId != nil {
		out.Append("client-request-id", fmt.Sprintf("%v", *o.ClientRequestId))
	}
	if o.IfMatch != nil {
		out.Append("If-Match", fmt.Sprintf("%v", *o.IfMatch))
	}
	if o.IfModifiedSince != nil {
		out.Append("If-Modified-Since", fmt.Sprintf("%v", *o.IfModifiedSince))
	}
	if o.IfNoneMatch != nil {
		out.Append("If-None-Match", fmt.Sprintf("%v", *o.IfNoneMatch))
	}
	if o.IfUnmodifiedSince != nil {
		out.Append("If-Unmodified-Since", fmt.Sprintf("%v", *o.IfUnmodifiedSince))
	}
	if o.OcpDate != nil {
		out.Append("ocp-date", fmt.Sprintf("%v", *o.OcpDate))
	}
	if o.ReturnClientRequestId != nil {
		out.Append("return-client-request-id", fmt.Sprintf("%v", *o.ReturnClientRequestId))
	}
	return &out
}

func (o JobPatchOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o JobPatchOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Timeout != nil {
		out.Append("timeout", fmt.Sprintf("%v", *o.Timeout))
	}
	return &out
}

// JobPatch ...
func (c JobsClient) JobPatch(ctx context.Context, id JobId, input JobPatchParameter, options JobPatchOperationOptions) (result JobPatchOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; odata=minimalmetadata; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPatch,
		OptionsObject: options,
		Path:          id.Path(),
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

	return
}
