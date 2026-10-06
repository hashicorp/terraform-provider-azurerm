package webapps

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AnalyzeCustomHostnameSlotOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *CustomHostnameAnalysisResult
}

type AnalyzeCustomHostnameSlotOperationOptions struct {
	HostName *string
}

func DefaultAnalyzeCustomHostnameSlotOperationOptions() AnalyzeCustomHostnameSlotOperationOptions {
	return AnalyzeCustomHostnameSlotOperationOptions{}
}

func (o AnalyzeCustomHostnameSlotOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o AnalyzeCustomHostnameSlotOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o AnalyzeCustomHostnameSlotOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.HostName != nil {
		out.Append("hostName", fmt.Sprintf("%v", *o.HostName))
	}
	return &out
}

// AnalyzeCustomHostnameSlot ...
func (c WebAppsClient) AnalyzeCustomHostnameSlot(ctx context.Context, id SlotId, options AnalyzeCustomHostnameSlotOperationOptions) (result AnalyzeCustomHostnameSlotOperationResponse, err error) {
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
