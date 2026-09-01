package containerapps

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListCustomHostNameAnalysisOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *CustomHostnameAnalysisResult
}

type ListCustomHostNameAnalysisOperationOptions struct {
	CustomHostname *string
}

func DefaultListCustomHostNameAnalysisOperationOptions() ListCustomHostNameAnalysisOperationOptions {
	return ListCustomHostNameAnalysisOperationOptions{}
}

func (o ListCustomHostNameAnalysisOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListCustomHostNameAnalysisOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListCustomHostNameAnalysisOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.CustomHostname != nil {
		out.Append("customHostname", fmt.Sprintf("%v", *o.CustomHostname))
	}
	return &out
}

// ListCustomHostNameAnalysis ...
func (c ContainerAppsClient) ListCustomHostNameAnalysis(ctx context.Context, id ContainerAppId, options ListCustomHostNameAnalysisOperationOptions) (result ListCustomHostNameAnalysisOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/listCustomHostNameAnalysis", id.ID()),
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
