package networkmanagereffectiveconnectivityconfiguration

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

type ListNetworkManagerEffectiveConnectivityConfigurationsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *NetworkManagerEffectiveConnectivityConfigurationListResult
}

type ListNetworkManagerEffectiveConnectivityConfigurationsOperationOptions struct {
	Top *int64
}

func DefaultListNetworkManagerEffectiveConnectivityConfigurationsOperationOptions() ListNetworkManagerEffectiveConnectivityConfigurationsOperationOptions {
	return ListNetworkManagerEffectiveConnectivityConfigurationsOperationOptions{}
}

func (o ListNetworkManagerEffectiveConnectivityConfigurationsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListNetworkManagerEffectiveConnectivityConfigurationsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListNetworkManagerEffectiveConnectivityConfigurationsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// ListNetworkManagerEffectiveConnectivityConfigurations ...
func (c NetworkManagerEffectiveConnectivityConfigurationClient) ListNetworkManagerEffectiveConnectivityConfigurations(ctx context.Context, id commonids.VirtualNetworkId, input QueryRequestOptions, options ListNetworkManagerEffectiveConnectivityConfigurationsOperationOptions) (result ListNetworkManagerEffectiveConnectivityConfigurationsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/listNetworkManagerEffectiveConnectivityConfigurations", id.ID()),
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

	var model NetworkManagerEffectiveConnectivityConfigurationListResult
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
