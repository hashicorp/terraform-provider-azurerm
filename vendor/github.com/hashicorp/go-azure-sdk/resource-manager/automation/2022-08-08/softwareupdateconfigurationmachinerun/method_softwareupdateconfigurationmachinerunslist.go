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

type SoftwareUpdateConfigurationMachineRunsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *SoftwareUpdateConfigurationMachineRunListResult
}

type SoftwareUpdateConfigurationMachineRunsListOperationOptions struct {
	ClientRequestId *string
	Filter          *string
	Skip            *string
	Top             *string
}

func DefaultSoftwareUpdateConfigurationMachineRunsListOperationOptions() SoftwareUpdateConfigurationMachineRunsListOperationOptions {
	return SoftwareUpdateConfigurationMachineRunsListOperationOptions{}
}

func (o SoftwareUpdateConfigurationMachineRunsListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.ClientRequestId != nil {
		out.Append("clientRequestId", fmt.Sprintf("%v", *o.ClientRequestId))
	}
	return &out
}

func (o SoftwareUpdateConfigurationMachineRunsListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o SoftwareUpdateConfigurationMachineRunsListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// SoftwareUpdateConfigurationMachineRunsList ...
func (c SoftwareUpdateConfigurationMachineRunClient) SoftwareUpdateConfigurationMachineRunsList(ctx context.Context, id AutomationAccountId, options SoftwareUpdateConfigurationMachineRunsListOperationOptions) (result SoftwareUpdateConfigurationMachineRunsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/softwareUpdateConfigurationMachineRuns", id.ID()),
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
