package guestconfigurationassignments

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

type RGListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *GuestConfigurationAssignmentList
}

// RGList ...
func (c GuestConfigurationAssignmentsClient) RGList(ctx context.Context, id commonids.ResourceGroupId) (result RGListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusNoContent,
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/providers/Microsoft.GuestConfiguration/guestConfigurationAssignments", id.ID()),
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

	var model GuestConfigurationAssignmentList
	result.Model = &model

	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
