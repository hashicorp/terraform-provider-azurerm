package emailtemplates

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EmailTemplateDeleteOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
}

type EmailTemplateDeleteOperationOptions struct {
	IfMatch *string
}

func DefaultEmailTemplateDeleteOperationOptions() EmailTemplateDeleteOperationOptions {
	return EmailTemplateDeleteOperationOptions{}
}

func (o EmailTemplateDeleteOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.IfMatch != nil {
		out.Append("If-Match", fmt.Sprintf("%v", *o.IfMatch))
	}
	return &out
}

func (o EmailTemplateDeleteOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o EmailTemplateDeleteOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// EmailTemplateDelete ...
func (c EmailTemplatesClient) EmailTemplateDelete(ctx context.Context, id TemplateId, options EmailTemplateDeleteOperationOptions) (result EmailTemplateDeleteOperationResponse, err error) {
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
