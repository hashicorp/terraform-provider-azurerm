// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package azuresdkhacks

import (
	"context"
	"fmt"
	"net/http"

	v2024 "github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/rulesets"
	rules2025 "github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rules"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

const rulesets2025DefaultApiVersion = "2025-12-01"

type RuleSets2025Client struct {
	Client *resourcemanager.Client
}

func NewRuleSets2025ClientWithBaseURI(sdkApi sdkEnv.Api) (*RuleSets2025Client, error) {
	client, err := resourcemanager.NewClient(sdkApi, "rulesets", rulesets2025DefaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating RuleSets2025Client: %+v", err)
	}

	return &RuleSets2025Client{Client: client}, nil
}

type RuleSet2025 struct {
	Id         *string                `json:"id,omitempty"`
	Name       *string                `json:"name,omitempty"`
	Properties *RuleSetProperties2025 `json:"properties,omitempty"`
	Type       *string                `json:"type,omitempty"`
}

type RuleSetProperties2025 struct {
	BatchMode         *bool             `json:"batchMode,omitempty"`
	DeploymentStatus  *string           `json:"deploymentStatus,omitempty"`
	ProfileName       *string           `json:"profileName,omitempty"`
	ProvisioningState *string           `json:"provisioningState,omitempty"`
	Rules             *[]rules2025.Rule `json:"rules,omitempty"`
}

type RuleSets2025CreateOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *RuleSet2025
}

func (c RuleSets2025Client) Create(ctx context.Context, id v2024.RuleSetId, input RuleSet2025) (result RuleSets2025CreateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType:         "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{http.StatusAccepted, http.StatusCreated, http.StatusOK},
		HttpMethod:          http.MethodPut,
		Path:                id.ID(),
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

	result.Poller, err = resourcemanager.PollerFromResponse(resp, c.Client)
	if err != nil {
		return
	}

	var model RuleSet2025
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}

func (c RuleSets2025Client) CreateThenPoll(ctx context.Context, id v2024.RuleSetId, input RuleSet2025) error {
	result, err := c.Create(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing Create: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after Create: %+v", err)
	}

	return nil
}

type RuleSets2025GetOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *RuleSet2025
}

func (c RuleSets2025Client) Get(ctx context.Context, id v2024.RuleSetId) (result RuleSets2025GetOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType:         "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{http.StatusOK},
		HttpMethod:          http.MethodGet,
		Path:                id.ID(),
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

	var model RuleSet2025
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}

type RuleSets2025DeleteOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

func (c RuleSets2025Client) Delete(ctx context.Context, id v2024.RuleSetId) (result RuleSets2025DeleteOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType:         "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{http.StatusAccepted, http.StatusNoContent, http.StatusOK},
		HttpMethod:          http.MethodDelete,
		Path:                id.ID(),
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

func (c RuleSets2025Client) DeleteThenPoll(ctx context.Context, id v2024.RuleSetId) error {
	result, err := c.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing Delete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after Delete: %+v", err)
	}

	return nil
}
