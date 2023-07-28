package checkdnsavailabilities

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CheckDnsNameAvailabilityOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *DnsNameAvailabilityResult
}

type CheckDnsNameAvailabilityOperationOptions struct {
	DomainNameLabel *string
}

func DefaultCheckDnsNameAvailabilityOperationOptions() CheckDnsNameAvailabilityOperationOptions {
	return CheckDnsNameAvailabilityOperationOptions{}
}

func (o CheckDnsNameAvailabilityOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o CheckDnsNameAvailabilityOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o CheckDnsNameAvailabilityOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.DomainNameLabel != nil {
		out.Append("domainNameLabel", fmt.Sprintf("%v", *o.DomainNameLabel))
	}
	return &out
}

// CheckDnsNameAvailability ...
func (c CheckDnsAvailabilitiesClient) CheckDnsNameAvailability(ctx context.Context, id LocationId, options CheckDnsNameAvailabilityOperationOptions) (result CheckDnsNameAvailabilityOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/checkDnsNameAvailability", id.ID()),
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
