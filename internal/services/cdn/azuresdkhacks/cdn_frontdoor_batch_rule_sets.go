// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package azuresdkhacks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
	rulesets "github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/rulesets"
	batchRules "github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rules"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

const batchRuleSetsDefaultApiVersion = "2025-12-01"

type BatchRuleSetsClient struct {
	Client *resourcemanager.Client
}

func NewBatchRuleSetsClientWithBaseURI(sdkApi sdkEnv.Api) (*BatchRuleSetsClient, error) {
	resourceManagerClient, err := resourcemanager.NewClient(sdkApi, "rulesets", batchRuleSetsDefaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating BatchRuleSetsClient: %+v", err)
	}

	return &BatchRuleSetsClient{Client: resourceManagerClient}, nil
}

type BatchRuleSetResource struct {
	Id         *string                 `json:"id,omitempty"`
	Name       *string                 `json:"name,omitempty"`
	Properties *BatchRuleSetProperties `json:"properties,omitempty"`
	Type       *string                 `json:"type,omitempty"`
}

type BatchRuleSetProperties struct {
	BatchMode         *bool                  `json:"batchMode,omitempty"`
	DeploymentStatus  *string                `json:"deploymentStatus,omitempty"`
	ProfileName       *string                `json:"profileName,omitempty"`
	ProvisioningState *string                `json:"provisioningState,omitempty"`
	Rules             *[]BatchRuleProperties `json:"rules,omitempty"`
	ForceSendRules    bool                   `json:"-"`
}

var _ json.Marshaler = &BatchRuleSetProperties{}

func (b BatchRuleSetProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if b.BatchMode != nil {
		objectMap["batchMode"] = b.BatchMode
	}
	if b.DeploymentStatus != nil {
		objectMap["deploymentStatus"] = b.DeploymentStatus
	}
	if b.ProfileName != nil {
		objectMap["profileName"] = b.ProfileName
	}
	if b.ProvisioningState != nil {
		objectMap["provisioningState"] = b.ProvisioningState
	}
	if b.ForceSendRules || b.Rules != nil {
		objectMap["rules"] = b.Rules
	}

	return json.Marshal(objectMap)
}

type BatchRuleProperties struct {
	Actions                 *[]batchRules.DeliveryRuleAction    `json:"actions,omitempty"`
	Conditions              *[]batchRules.DeliveryRuleCondition `json:"conditions,omitempty"`
	DeploymentStatus        *batchRules.DeploymentStatus        `json:"deploymentStatus,omitempty"`
	Id                      *string                             `json:"id,omitempty"`
	MatchProcessingBehavior *batchRules.MatchProcessingBehavior `json:"matchProcessingBehavior,omitempty"`
	Name                    *string                             `json:"-"`
	Order                   *int64                              `json:"order,omitempty"`
	ProvisioningState       *batchRules.AfdProvisioningState    `json:"provisioningState,omitempty"`
	RuleName                *string                             `json:"ruleName,omitempty"`
	RuleSetName             *string                             `json:"-"`
	SystemData              *systemdata.SystemData              `json:"systemData,omitempty"`
	Type                    *string                             `json:"type,omitempty"`
}

var _ json.Unmarshaler = &BatchRuleProperties{}

func (b *BatchRuleProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		DeploymentStatus        *batchRules.DeploymentStatus        `json:"deploymentStatus,omitempty"`
		Id                      *string                             `json:"id,omitempty"`
		MatchProcessingBehavior *batchRules.MatchProcessingBehavior `json:"matchProcessingBehavior,omitempty"`
		Name                    *string                             `json:"name,omitempty"`
		Order                   *int64                              `json:"order,omitempty"`
		ProvisioningState       *batchRules.AfdProvisioningState    `json:"provisioningState,omitempty"`
		RuleName                *string                             `json:"ruleName,omitempty"`
		RuleSetName             *string                             `json:"ruleSetName,omitempty"`
		SystemData              *systemdata.SystemData              `json:"systemData,omitempty"`
		Type                    *string                             `json:"type,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	b.DeploymentStatus = decoded.DeploymentStatus
	b.Id = decoded.Id
	b.MatchProcessingBehavior = decoded.MatchProcessingBehavior
	b.Name = decoded.Name
	b.Order = decoded.Order
	b.ProvisioningState = decoded.ProvisioningState
	b.RuleName = decoded.RuleName
	b.RuleSetName = decoded.RuleSetName
	b.SystemData = decoded.SystemData
	b.Type = decoded.Type

	if b.Name == nil && b.RuleName != nil {
		b.Name = b.RuleName
	}
	if b.RuleName == nil && b.Name != nil {
		b.RuleName = b.Name
	}

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling BatchRuleProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["actions"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling Actions into list []json.RawMessage: %+v", err)
		}

		output := make([]batchRules.DeliveryRuleAction, 0)
		for i, val := range listTemp {
			impl, err := batchRules.UnmarshalDeliveryRuleActionImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'Actions' for 'BatchRuleProperties': %+v", i, err)
			}
			output = append(output, impl)
		}
		b.Actions = &output
	}

	if v, ok := temp["conditions"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling Conditions into list []json.RawMessage: %+v", err)
		}

		output := make([]batchRules.DeliveryRuleCondition, 0)
		for i, val := range listTemp {
			impl, err := batchRules.UnmarshalDeliveryRuleConditionImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'Conditions' for 'BatchRuleProperties': %+v", i, err)
			}
			output = append(output, impl)
		}
		b.Conditions = &output
	}

	return nil
}

type BatchRuleSetsCreateOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *BatchRuleSetResource
}

func (c BatchRuleSetsClient) Create(ctx context.Context, id rulesets.RuleSetId, input BatchRuleSetResource) (result BatchRuleSetsCreateOperationResponse, err error) {
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

	var model BatchRuleSetResource
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}

func (c BatchRuleSetsClient) CreateThenPoll(ctx context.Context, id rulesets.RuleSetId, input BatchRuleSetResource) error {
	result, err := c.Create(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing Create: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after Create: %+v", err)
	}

	return nil
}

type BatchRuleSetsGetOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *BatchRuleSetResource
}

func (c BatchRuleSetsClient) Get(ctx context.Context, id rulesets.RuleSetId) (result BatchRuleSetsGetOperationResponse, err error) {
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

	var model BatchRuleSetResource
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}

type BatchRuleSetsDeleteOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

func (c BatchRuleSetsClient) Delete(ctx context.Context, id rulesets.RuleSetId) (result BatchRuleSetsDeleteOperationResponse, err error) {
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

func (c BatchRuleSetsClient) DeleteThenPoll(ctx context.Context, id rulesets.RuleSetId) error {
	result, err := c.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing Delete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after Delete: %+v", err)
	}

	return nil
}
