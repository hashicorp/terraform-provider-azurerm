package softwareupdateconfigurationmachinerun

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SoftwareUpdateConfigurationMachineRunsGetByIdOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *SoftwareUpdateConfigurationMachineRun
}

type SoftwareUpdateConfigurationMachineRunsGetByIdOperationOptions struct {
	ClientRequestId *string
}

func DefaultSoftwareUpdateConfigurationMachineRunsGetByIdOperationOptions() SoftwareUpdateConfigurationMachineRunsGetByIdOperationOptions {
	return SoftwareUpdateConfigurationMachineRunsGetByIdOperationOptions{}
}

func (o SoftwareUpdateConfigurationMachineRunsGetByIdOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.ClientRequestId != nil {
		out.Append("clientRequestId", fmt.Sprintf("%v", *o.ClientRequestId))
	}
	return &out
}

func (o SoftwareUpdateConfigurationMachineRunsGetByIdOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o SoftwareUpdateConfigurationMachineRunsGetByIdOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// SoftwareUpdateConfigurationMachineRunsGetById ...
func (c SoftwareUpdateConfigurationMachineRunClient) SoftwareUpdateConfigurationMachineRunsGetById(ctx context.Context, id SoftwareUpdateConfigurationMachineRunId, options SoftwareUpdateConfigurationMachineRunsGetByIdOperationOptions) (result SoftwareUpdateConfigurationMachineRunsGetByIdOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          id.ID(),
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
