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

type AnalyzeCustomHostnameOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *CustomHostnameAnalysisResult
}

type AnalyzeCustomHostnameOperationOptions struct {
	HostName *string
}

func DefaultAnalyzeCustomHostnameOperationOptions() AnalyzeCustomHostnameOperationOptions {
	return AnalyzeCustomHostnameOperationOptions{}
}

func (o AnalyzeCustomHostnameOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o AnalyzeCustomHostnameOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o AnalyzeCustomHostnameOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.HostName != nil {
		out.Append("hostName", fmt.Sprintf("%v", *o.HostName))
	}
	return &out
}

// AnalyzeCustomHostname ...
func (c WebAppsClient) AnalyzeCustomHostname(ctx context.Context, id commonids.AppServiceId, options AnalyzeCustomHostnameOperationOptions) (result AnalyzeCustomHostnameOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/analyzeCustomHostname", id.ID()),
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

	var model CustomHostnameAnalysisResult
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
