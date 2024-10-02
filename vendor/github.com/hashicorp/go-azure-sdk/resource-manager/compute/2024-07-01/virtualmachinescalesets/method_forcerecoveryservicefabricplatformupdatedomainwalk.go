package virtualmachinescalesets

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ForceRecoveryServiceFabricPlatformUpdateDomainWalkOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *RecoveryWalkResponse
}

type ForceRecoveryServiceFabricPlatformUpdateDomainWalkOperationOptions struct {
	PlacementGroupId     *string
	PlatformUpdateDomain *int64
	Zone                 *string
}

func DefaultForceRecoveryServiceFabricPlatformUpdateDomainWalkOperationOptions() ForceRecoveryServiceFabricPlatformUpdateDomainWalkOperationOptions {
	return ForceRecoveryServiceFabricPlatformUpdateDomainWalkOperationOptions{}
}

func (o ForceRecoveryServiceFabricPlatformUpdateDomainWalkOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ForceRecoveryServiceFabricPlatformUpdateDomainWalkOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ForceRecoveryServiceFabricPlatformUpdateDomainWalkOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.PlacementGroupId != nil {
		out.Append("placementGroupId", fmt.Sprintf("%v", *o.PlacementGroupId))
	}
	if o.PlatformUpdateDomain != nil {
		out.Append("platformUpdateDomain", fmt.Sprintf("%v", *o.PlatformUpdateDomain))
	}
	if o.Zone != nil {
		out.Append("zone", fmt.Sprintf("%v", *o.Zone))
	}
	return &out
}

// ForceRecoveryServiceFabricPlatformUpdateDomainWalk ...
func (c VirtualMachineScaleSetsClient) ForceRecoveryServiceFabricPlatformUpdateDomainWalk(ctx context.Context, id VirtualMachineScaleSetId, options ForceRecoveryServiceFabricPlatformUpdateDomainWalkOperationOptions) (result ForceRecoveryServiceFabricPlatformUpdateDomainWalkOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/forceRecoveryServiceFabricPlatformUpdateDomainWalk", id.ID()),
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

	var model RecoveryWalkResponse
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
