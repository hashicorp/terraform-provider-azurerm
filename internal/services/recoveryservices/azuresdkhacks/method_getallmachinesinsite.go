// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azuresdkhacks

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	vmwaremachines "github.com/hashicorp/go-azure-sdk/resource-manager/migrate/2020-01-01/machines"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
)

// workaround for https://github.com/hashicorp/go-azure-sdk/issues/492
// TODO4.0: check if this could be removed.
// the method has been re-written to read `nextLink`

type MachinesClient struct {
	Client *resourcemanager.Client
}

type Values struct {
	Values   *[]vmwaremachines.VMwareMachine `json:"value"`
	NextLink *string                         `json:"nextLink"`
}

func (c MachinesClient) GetAllVMWareMachinesInSite(ctx context.Context, id vmwaremachines.VMwareSiteId, options vmwaremachines.GetAllMachinesInSiteOperationOptions) (result vmwaremachines.GetAllMachinesInSiteOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/machines", id.ID()),
		OptionsObject: options,
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	return wrapExecutePaged(ctx, req)
}

func wrapExecutePaged(ctx context.Context, req *client.Request) (result vmwaremachines.GetAllMachinesInSiteOperationResponse, err error) {
	resp, err := req.ExecutePaged(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	var values Values

	if err = resp.Unmarshal(&values); err != nil {
		return
	}
	result.Model = values.Values

	if values.NextLink != nil {
		nextReq := req
		u, err := url.Parse(*values.NextLink)
		if err != nil {
			return result, err
		}
		nextReq.URL = u
		nextResp, err := wrapExecutePaged(ctx, nextReq)
		if err != nil {
			return result, err
		}
		if nextResp.Model != nil {
			result.Model = pointer.To(append(*result.Model, *nextResp.Model...))
		}
	}

	return
}
