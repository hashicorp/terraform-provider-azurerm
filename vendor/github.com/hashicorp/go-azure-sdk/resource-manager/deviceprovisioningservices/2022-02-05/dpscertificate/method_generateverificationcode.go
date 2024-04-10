package dpscertificate

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GenerateVerificationCodeOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *VerificationCodeResponse
}

type GenerateVerificationCodeOperationOptions struct {
	CertificateCreated       *string
	CertificateHasPrivateKey *bool
	CertificateIsVerified    *bool
	CertificateLastUpdated   *string
	CertificateName          *string
	CertificateNonce         *string
	CertificatePurpose       *CertificatePurpose
	CertificateRawBytes      *string
	IfMatch                  *string
}

func DefaultGenerateVerificationCodeOperationOptions() GenerateVerificationCodeOperationOptions {
	return GenerateVerificationCodeOperationOptions{}
}

func (o GenerateVerificationCodeOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.IfMatch != nil {
		out.Append("If-Match", fmt.Sprintf("%v", *o.IfMatch))
	}
	return &out
}

func (o GenerateVerificationCodeOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o GenerateVerificationCodeOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.CertificateCreated != nil {
		out.Append("certificate.created", fmt.Sprintf("%v", *o.CertificateCreated))
	}
	if o.CertificateHasPrivateKey != nil {
		out.Append("certificate.hasPrivateKey", fmt.Sprintf("%v", *o.CertificateHasPrivateKey))
	}
	if o.CertificateIsVerified != nil {
		out.Append("certificate.isVerified", fmt.Sprintf("%v", *o.CertificateIsVerified))
	}
	if o.CertificateLastUpdated != nil {
		out.Append("certificate.lastUpdated", fmt.Sprintf("%v", *o.CertificateLastUpdated))
	}
	if o.CertificateName != nil {
		out.Append("certificate.name", fmt.Sprintf("%v", *o.CertificateName))
	}
	if o.CertificateNonce != nil {
		out.Append("certificate.nonce", fmt.Sprintf("%v", *o.CertificateNonce))
	}
	if o.CertificatePurpose != nil {
		out.Append("certificate.purpose", fmt.Sprintf("%v", *o.CertificatePurpose))
	}
	if o.CertificateRawBytes != nil {
		out.Append("certificate.rawBytes", fmt.Sprintf("%v", *o.CertificateRawBytes))
	}
	return &out
}

// GenerateVerificationCode ...
func (c DpsCertificateClient) GenerateVerificationCode(ctx context.Context, id CertificateId, options GenerateVerificationCodeOperationOptions) (result GenerateVerificationCodeOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		Path:          fmt.Sprintf("%s/generateVerificationCode", id.ID()),
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

	var model VerificationCodeResponse
	result.Model = &model

	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
