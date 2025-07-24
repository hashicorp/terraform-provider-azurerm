// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azuresdkhacks

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-01-01/managedenvironments"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
)

// This `azuresdkhack` only exists because `go-azure-sdk` currently does not support sending explicit `null` values.
// If that changes, this can be removed.

type ManagedEnvironmentsClient struct {
	client *managedenvironments.ManagedEnvironmentsClient
}

func NewManagedEnvironmentWorkaroundClient(client *managedenvironments.ManagedEnvironmentsClient) ManagedEnvironmentsClient {
	return ManagedEnvironmentsClient{
		client: client,
	}
}

func (c ManagedEnvironmentsClient) Update(ctx context.Context, id managedenvironments.ManagedEnvironmentId, input ManagedEnvironment) (result managedenvironments.UpdateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPatch,
		Path:       id.ID(),
	}

	req, err := c.client.Client.NewRequest(ctx, opts)
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

	result.Poller, err = resourcemanager.PollerFromResponse(resp, c.client.Client)
	if err != nil {
		return
	}

	return
}

func (c ManagedEnvironmentsClient) UpdateThenPoll(ctx context.Context, id managedenvironments.ManagedEnvironmentId, input ManagedEnvironment) error {
	result, err := c.Update(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing Update: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after Update: %+v", err)
	}

	return nil
}

type ManagedEnvironment struct {
	Id         *string                       `json:"id,omitempty"`
	Kind       *string                       `json:"kind,omitempty"`
	Location   string                        `json:"location"`
	Name       *string                       `json:"name,omitempty"`
	Properties *ManagedEnvironmentProperties `json:"properties,omitempty"`
	SystemData *systemdata.SystemData        `json:"systemData,omitempty"`
	Tags       *map[string]string            `json:"tags,omitempty"`
	Type       *string                       `json:"type,omitempty"`
}

type ManagedEnvironmentProperties struct {
	AppLogsConfiguration        *AppLogsConfiguration                                                     `json:"appLogsConfiguration,omitempty"`
	CustomDomainConfiguration   *managedenvironments.CustomDomainConfiguration                            `json:"customDomainConfiguration,omitempty"`
	DaprAIConnectionString      *string                                                                   `json:"daprAIConnectionString,omitempty"`
	DaprAIInstrumentationKey    *string                                                                   `json:"daprAIInstrumentationKey,omitempty"`
	DaprConfiguration           *managedenvironments.DaprConfiguration                                    `json:"daprConfiguration,omitempty"`
	DefaultDomain               *string                                                                   `json:"defaultDomain,omitempty"`
	DeploymentErrors            *string                                                                   `json:"deploymentErrors,omitempty"`
	EventStreamEndpoint         *string                                                                   `json:"eventStreamEndpoint,omitempty"`
	InfrastructureResourceGroup *string                                                                   `json:"infrastructureResourceGroup,omitempty"`
	KedaConfiguration           *managedenvironments.KedaConfiguration                                    `json:"kedaConfiguration,omitempty"`
	PeerAuthentication          *managedenvironments.ManagedEnvironmentPropertiesPeerAuthentication       `json:"peerAuthentication,omitempty"`
	PeerTrafficConfiguration    *managedenvironments.ManagedEnvironmentPropertiesPeerTrafficConfiguration `json:"peerTrafficConfiguration,omitempty"`
	ProvisioningState           *managedenvironments.EnvironmentProvisioningState                         `json:"provisioningState,omitempty"`
	StaticIP                    *string                                                                   `json:"staticIp,omitempty"`
	VnetConfiguration           *managedenvironments.VnetConfiguration                                    `json:"vnetConfiguration,omitempty"`
	WorkloadProfiles            *[]managedenvironments.WorkloadProfile                                    `json:"workloadProfiles,omitempty"`
	ZoneRedundant               *bool                                                                     `json:"zoneRedundant,omitempty"`
}

type AppLogsConfiguration struct {
	Destination *string `json:"destination,omitempty"`
	// The only difference is the removed `omitempty` tag so that we can explicitly set `logAnalyticsConfiguration` to `null` in the API request.
	LogAnalyticsConfiguration *managedenvironments.LogAnalyticsConfiguration `json:"logAnalyticsConfiguration"`
}
