package localrulestackresources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LocalRulestacksgetSupportInfoOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *SupportInfo
}

type LocalRulestacksgetSupportInfoOperationOptions struct {
	Email *string
}

func DefaultLocalRulestacksgetSupportInfoOperationOptions() LocalRulestacksgetSupportInfoOperationOptions {
	return LocalRulestacksgetSupportInfoOperationOptions{}
}

func (o LocalRulestacksgetSupportInfoOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o LocalRulestacksgetSupportInfoOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o LocalRulestacksgetSupportInfoOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Email != nil {
		out.Append("email", fmt.Sprintf("%v", *o.Email))
	}
	return &out
}

// LocalRulestacksgetSupportInfo ...
func (c LocalRulestackResourcesClient) LocalRulestacksgetSupportInfo(ctx context.Context, id LocalRulestackId, options LocalRulestacksgetSupportInfoOperationOptions) (result LocalRulestacksgetSupportInfoOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/getSupportInfo", id.ID()),
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

	var model SupportInfo
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
