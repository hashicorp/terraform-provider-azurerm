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

type EmailTemplateUpdateOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *EmailTemplateContract
}

type EmailTemplateUpdateOperationOptions struct {
	IfMatch *string
}

func DefaultEmailTemplateUpdateOperationOptions() EmailTemplateUpdateOperationOptions {
	return EmailTemplateUpdateOperationOptions{}
}

func (o EmailTemplateUpdateOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.IfMatch != nil {
		out.Append("If-Match", fmt.Sprintf("%v", *o.IfMatch))
	}
	return &out
}

func (o EmailTemplateUpdateOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o EmailTemplateUpdateOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// EmailTemplateUpdate ...
func (c EmailTemplatesClient) EmailTemplateUpdate(ctx context.Context, id TemplateId, input EmailTemplateUpdateParameters, options EmailTemplateUpdateOperationOptions) (result EmailTemplateUpdateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
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
