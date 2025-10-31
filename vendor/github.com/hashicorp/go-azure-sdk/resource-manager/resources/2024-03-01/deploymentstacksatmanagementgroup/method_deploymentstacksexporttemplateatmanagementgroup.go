package deploymentstacksatmanagementgroup

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentStacksExportTemplateAtManagementGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *DeploymentStackTemplateDefinition
}

// DeploymentStacksExportTemplateAtManagementGroup ...
func (c DeploymentStacksAtManagementGroupClient) DeploymentStacksExportTemplateAtManagementGroup(ctx context.Context, id Providers2DeploymentStackId) (result DeploymentStacksExportTemplateAtManagementGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/exportTemplate", id.ID()),
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

	var model DeploymentStackTemplateDefinition
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
