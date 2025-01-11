package profiles

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

type CheckTrafficManagerNameAvailabilityV2OperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *TrafficManagerNameAvailability
}

// CheckTrafficManagerNameAvailabilityV2 ...
func (c ProfilesClient) CheckTrafficManagerNameAvailabilityV2(ctx context.Context, id commonids.SubscriptionId, input CheckTrafficManagerRelativeDnsNameAvailabilityParameters) (result CheckTrafficManagerNameAvailabilityV2OperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/providers/Microsoft.Network/checkTrafficManagerNameAvailabilityV2", id.ID()),
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

	var model TrafficManagerNameAvailability
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
