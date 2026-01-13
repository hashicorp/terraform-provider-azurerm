package certificates

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificateCreateOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *Certificate
}

type CertificateCreateOperationOptions struct {
	IfMatch     *string
	IfNoneMatch *string
}

func DefaultCertificateCreateOperationOptions() CertificateCreateOperationOptions {
	return CertificateCreateOperationOptions{}
}

func (o CertificateCreateOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.IfMatch != nil {
		out.Append("If-Match", fmt.Sprintf("%v", *o.IfMatch))
	}
	if o.IfNoneMatch != nil {
		out.Append("If-None-Match", fmt.Sprintf("%v", *o.IfNoneMatch))
	}
	return &out
}

func (o CertificateCreateOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o CertificateCreateOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// CertificateCreate ...
func (c CertificatesClient) CertificateCreate(ctx context.Context, id CertificateId, input CertificateCreateOrUpdateParameters, options CertificateCreateOperationOptions) (result CertificateCreateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
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

	var model Certificate
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
