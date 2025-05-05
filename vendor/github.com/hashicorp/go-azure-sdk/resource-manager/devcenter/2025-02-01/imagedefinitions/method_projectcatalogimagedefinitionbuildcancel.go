package imagedefinitions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProjectCatalogImageDefinitionBuildCancelOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// ProjectCatalogImageDefinitionBuildCancel ...
func (c ImageDefinitionsClient) ProjectCatalogImageDefinitionBuildCancel(ctx context.Context, id BuildId) (result ProjectCatalogImageDefinitionBuildCancelOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/cancel", id.ID()),
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

	result.Poller, err = resourcemanager.PollerFromResponse(resp, c.Client)
	if err != nil {
		return
	}

	return
}

// ProjectCatalogImageDefinitionBuildCancelThenPoll performs ProjectCatalogImageDefinitionBuildCancel then polls until it's completed
func (c ImageDefinitionsClient) ProjectCatalogImageDefinitionBuildCancelThenPoll(ctx context.Context, id BuildId) error {
	result, err := c.ProjectCatalogImageDefinitionBuildCancel(ctx, id)
	if err != nil {
		return fmt.Errorf("performing ProjectCatalogImageDefinitionBuildCancel: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after ProjectCatalogImageDefinitionBuildCancel: %+v", err)
	}

	return nil
}
