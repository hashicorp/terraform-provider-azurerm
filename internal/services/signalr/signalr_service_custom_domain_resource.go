// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package signalr

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2023-02-01/signalr"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CustomDomainSignalrServiceModel struct {
	Name                       string `tfschema:"name"`
	SignalRServiceId           string `tfschema:"signalr_service_id"`
	DomainName                 string `tfschema:"domain_name"`
	SignalrCustomCertificateId string `tfschema:"signalr_custom_certificate_id"`
}

type CustomDomainSignalrServiceResource struct{}

var _ sdk.Resource = CustomDomainSignalrServiceResource{}

func (r CustomDomainSignalrServiceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"signalr_service_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: signalr.ValidateSignalRID,
		},

		"domain_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"signalr_custom_certificate_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: signalr.ValidateCustomCertificateID,
		},
	}
}

func (r CustomDomainSignalrServiceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r CustomDomainSignalrServiceResource) ModelObject() interface{} {
	return &CustomDomainSignalrServiceModel{}
}

func (r CustomDomainSignalrServiceResource) ResourceType() string {
	return "azurerm_signalr_service_custom_domain"
}

func (r CustomDomainSignalrServiceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var customDomainSignalrServiceModel CustomDomainSignalrServiceModel
			if err := metadata.Decode(&customDomainSignalrServiceModel); err != nil {
				return err
			}
			client := metadata.Client.SignalR.SignalRClient

			signalRServiceId, err := signalr.ParseSignalRIDInsensitively(metadata.ResourceData.Get("signalr_service_id").(string))
			if err != nil {
				return fmt.Errorf("parsing signalr service id error: %+v", err)
			}

			id := signalr.NewCustomDomainID(signalRServiceId.SubscriptionId, signalRServiceId.ResourceGroupName, signalRServiceId.SignalRName, metadata.ResourceData.Get("name").(string))

			locks.ByID(signalRServiceId.ID())
			defer locks.UnlockByID(signalRServiceId.ID())

			if _, err := signalr.ParseCustomCertificateIDInsensitively(customDomainSignalrServiceModel.SignalrCustomCertificateId); err != nil {
				return fmt.Errorf("parsing custom certificate for %s: %+v", id, err)
			}

			existing, err := client.CustomDomainsGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			customDomainObj := signalr.CustomDomain{
				Properties: signalr.CustomDomainProperties{
					DomainName: customDomainSignalrServiceModel.DomainName,
					CustomCertificate: signalr.ResourceReference{
						Id: utils.String(customDomainSignalrServiceModel.SignalrCustomCertificateId),
					},
				},
			}
			if _, err := client.CustomDomainsCreateOrUpdate(ctx, id, customDomainObj); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal-error: context had no deadline")
			}
			stateConf := &pluginsdk.StateChangeConf{
				Pending: []string{
					string(signalr.ProvisioningStateUpdating),
					string(signalr.ProvisioningStateCreating),
					string(signalr.ProvisioningStateMoving),
					string(signalr.ProvisioningStateRunning),
				},
				Target: []string{
					string(signalr.ProvisioningStateSucceeded),
					string(signalr.ProvisioningStateFailed),
				},
				Refresh:                   signalrServiceCustomDomainProvisioningStateRefreshFunc(ctx, client, id),
				Timeout:                   time.Until(deadline),
				PollInterval:              10 * time.Second,
				ContinuousTargetOccurence: 20,
			}

			if _, err := stateConf.WaitForStateContext(ctx); err != nil {
				return err
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r CustomDomainSignalrServiceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SignalR.SignalRClient
			id, err := signalr.ParseCustomDomainID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.CustomDomainsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := CustomDomainSignalrServiceModel{
				Name:             id.CustomDomainName,
				SignalRServiceId: signalr.NewSignalRID(id.SubscriptionId, id.ResourceGroupName, id.SignalRName).ID(),
			}

			if model := resp.Model; model != nil {
				props := model.Properties
				signalrCustomCertificateId := ""
				if props.CustomCertificate.Id != nil {
					signalrCustomCertificateID, err := signalr.ParseCustomCertificateIDInsensitively(*props.CustomCertificate.Id)
					if err != nil {
						return fmt.Errorf("parsing signalr custom cert id for %s: %+v", id, err)
					}
					signalrCustomCertificateId = signalrCustomCertificateID.ID()
				}

				state.SignalrCustomCertificateId = signalrCustomCertificateId
				state.DomainName = props.DomainName
			}

			return metadata.Encode(&state)
		},
	}
}

func (r CustomDomainSignalrServiceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SignalR.SignalRClient

			id, err := signalr.ParseCustomDomainID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			signalrId := signalr.NewSignalRID(id.SubscriptionId, id.ResourceGroupName, id.SignalRName)

			locks.ByID(signalrId.ID())
			defer locks.UnlockByID(signalrId.ID())

			if err := client.CustomDomainsDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal-error: context had no deadline")
			}
			stateConf := &pluginsdk.StateChangeConf{
				Pending:                   []string{"Exists"},
				Target:                    []string{"NotFound"},
				Refresh:                   signalrServiceCustomDomainDeleteRefreshFunc(ctx, client, *id),
				Timeout:                   time.Until(deadline),
				PollInterval:              10 * time.Second,
				ContinuousTargetOccurence: 20,
			}

			if _, err = stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for %s to be fully deleted: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r CustomDomainSignalrServiceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return signalr.ValidateCustomDomainID
}

func signalrServiceCustomDomainProvisioningStateRefreshFunc(ctx context.Context, client *signalr.SignalRClient, id signalr.CustomDomainId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.CustomDomainsGet(ctx, id)

		provisioningState := "Pending"
		if err != nil {
			if response.WasNotFound(res.HttpResponse) {
				return res, provisioningState, nil
			}
			return nil, "Error", fmt.Errorf("polling for the provisioning state of %s: %+v", id, err)
		}

		if res.Model != nil && res.Model.Properties.ProvisioningState != nil {
			provisioningState = string(*res.Model.Properties.ProvisioningState)
		}

		return res, provisioningState, nil
	}
}

func signalrServiceCustomDomainDeleteRefreshFunc(ctx context.Context, client *signalr.SignalRClient, id signalr.CustomDomainId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.CustomDomainsGet(ctx, id)
		if err != nil {
			if response.WasNotFound(res.HttpResponse) {
				return "NotFound", "NotFound", nil
			}

			return nil, "", fmt.Errorf("checking if %s has been deleted: %+v", id, err)
		}

		return res, "Exists", nil
	}
}
