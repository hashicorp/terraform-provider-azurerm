package certificate

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceCertificateDeleteOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
}

type WorkspaceCertificateDeleteOperationOptions struct {
	IfMatch *string
}

func DefaultWorkspaceCertificateDeleteOperationOptions() WorkspaceCertificateDeleteOperationOptions {
	return WorkspaceCertificateDeleteOperationOptions{}
}

func (o WorkspaceCertificateDeleteOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.IfMatch != nil {
		out.Append("If-Match", fmt.Sprintf("%v", *o.IfMatch))
	}
	return &out
}

func (o WorkspaceCertificateDeleteOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o WorkspaceCertificateDeleteOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// WorkspaceCertificateDelete ...
func (c CertificateClient) WorkspaceCertificateDelete(ctx context.Context, id WorkspaceCertificateId, options WorkspaceCertificateDeleteOperationOptions) (result WorkspaceCertificateDeleteOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusNoContent,
			http.StatusOK,
		},
		HttpMethod:    http.MethodDelete,
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

	return
}
