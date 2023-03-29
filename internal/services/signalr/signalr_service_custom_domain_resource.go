package signalr

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2023-02-01/signalr"
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

			signalRServiceId, err := signalr.ParseSignalRID(metadata.ResourceData.Get("signalr_service_id").(string))
			if err != nil {
				return fmt.Errorf("parsing signalr service id error: %+v", err)
			}

			id := signalr.NewCustomDomainID(signalRServiceId.SubscriptionId, signalRServiceId.ResourceGroupName, signalRServiceId.SignalRName, metadata.ResourceData.Get("name").(string))
			if _, err := signalr.ParseCustomCertificateIDInsensitively(customDomainSignalrServiceModel.SignalrCustomCertificateId); err != nil {
				return fmt.Errorf("parsing custom certificate for %s: %+v", id, err)
			}

			existing, err := client.CustomDomainsGet(ctx, id)
			if !response.WasNotFound(existing.HttpResponse) {
				if err != nil {
					return fmt.Errorf("retrieving %s: %v", id, err)
				}
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
			if err := client.CustomDomainsCreateOrUpdateThenPoll(ctx, id, customDomainObj); err != nil {
				return fmt.Errorf("creating signalR custom domain: %s: %+v", id, err)
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
				return fmt.Errorf("reading signalR custom domain %s: %+v", id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: got nil model", *id)
			}

			signalrServiceId := signalr.NewSignalRID(id.SubscriptionId, id.ResourceGroupName, id.SignalRName).ID()

			state := CustomDomainSignalrServiceModel{
				Name:             id.CustomDomainName,
				DomainName:       resp.Model.Properties.DomainName,
				SignalRServiceId: signalrServiceId,
			}

			if resp.Model.Properties.CustomCertificate.Id != nil {
				signalrCustomCertId, err := signalr.ParseCustomCertificateIDInsensitively(*resp.Model.Properties.CustomCertificate.Id)
				if err != nil {
					return fmt.Errorf("parsing signalr custom cert id for %s: %+v", id, err)
				}
				state.SignalrCustomCertificateId = signalrCustomCertId.ID()
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

			if err := waitForSignalrServiceStatusToBeReady(ctx, metadata, signalrId); err != nil {
				return fmt.Errorf("waiting for signalR service %s state to be ready error: %+v", signalrId, err)
			}

			if err := client.CustomDomainsDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r CustomDomainSignalrServiceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return signalr.ValidateCustomDomainID
}

func waitForSignalrServiceStatusToBeReady(ctx context.Context, metadata sdk.ResourceMetaData, id signalr.SignalRId) error {
	signalrClient := metadata.Client.SignalR.SignalRClient

	stateConf := &pluginsdk.StateChangeConf{
		Pending: []string{
			string(signalr.ProvisioningStateUpdating),
			string(signalr.ProvisioningStateCreating),
			string(signalr.ProvisioningStateMoving),
			string(signalr.ProvisioningStateRunning),
		},
		Target:                    []string{string(signalr.ProvisioningStateSucceeded)},
		Refresh:                   signalrServiceProvisioningStateRefreshFunc(ctx, signalrClient, id),
		Timeout:                   5 * time.Minute,
		PollInterval:              10 * time.Second,
		ContinuousTargetOccurence: 5,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return err
	}
	return nil
}

func signalrServiceProvisioningStateRefreshFunc(ctx context.Context, client *signalr.SignalRClient, id signalr.SignalRId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id)

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
