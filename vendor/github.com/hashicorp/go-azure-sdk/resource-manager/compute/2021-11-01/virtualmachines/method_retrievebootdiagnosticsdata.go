package virtualmachines

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RetrieveBootDiagnosticsDataOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *RetrieveBootDiagnosticsDataResult
}

type RetrieveBootDiagnosticsDataOperationOptions struct {
	SasUriExpirationTimeInMinutes *int64
}

func DefaultRetrieveBootDiagnosticsDataOperationOptions() RetrieveBootDiagnosticsDataOperationOptions {
	return RetrieveBootDiagnosticsDataOperationOptions{}
}

func (o RetrieveBootDiagnosticsDataOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o RetrieveBootDiagnosticsDataOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o RetrieveBootDiagnosticsDataOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.SasUriExpirationTimeInMinutes != nil {
		out.Append("sasUriExpirationTimeInMinutes", fmt.Sprintf("%v", *o.SasUriExpirationTimeInMinutes))
	}
	return &out
}

// RetrieveBootDiagnosticsData ...
func (c VirtualMachinesClient) RetrieveBootDiagnosticsData(ctx context.Context, id VirtualMachineId, options RetrieveBootDiagnosticsDataOperationOptions) (result RetrieveBootDiagnosticsDataOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		Path:          fmt.Sprintf("%s/retrieveBootDiagnosticsData", id.ID()),
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

	if err = resp.Unmarshal(&result.Model); err != nil {
		return
	}

	return
}
