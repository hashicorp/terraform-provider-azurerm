package webapps

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeleteSourceControlOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
}

type DeleteSourceControlOperationOptions struct {
	AdditionalFlags *string
}

func DefaultDeleteSourceControlOperationOptions() DeleteSourceControlOperationOptions {
	return DeleteSourceControlOperationOptions{}
}

func (o DeleteSourceControlOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o DeleteSourceControlOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o DeleteSourceControlOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.AdditionalFlags != nil {
		out.Append("additionalFlags", fmt.Sprintf("%v", *o.AdditionalFlags))
	}
	return &out
}

// DeleteSourceControl ...
func (c WebAppsClient) DeleteSourceControl(ctx context.Context, id commonids.AppServiceId, options DeleteSourceControlOperationOptions) (result DeleteSourceControlOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod:    http.MethodDelete,
		Path:          fmt.Sprintf("%s/sourceControls/web", id.ID()),
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
