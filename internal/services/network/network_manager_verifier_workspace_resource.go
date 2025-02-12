package network

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/ipampools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/verifierworkspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.Resource           = ManagerVerifierWorkspaceResource{}
	_ sdk.ResourceWithUpdate = ManagerVerifierWorkspaceResource{}
)

type ManagerVerifierWorkspaceResource struct{}

func (ManagerVerifierWorkspaceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return verifierworkspaces.ValidateVerifierWorkspaceID
}

func (ManagerVerifierWorkspaceResource) ResourceType() string {
	return "azurerm_network_manager_verifier_workspace"
}

func (ManagerVerifierWorkspaceResource) ModelObject() interface{} {
	return &ManagerVerifierWorkspaceResourceModel{}
}

type ManagerVerifierWorkspaceResourceModel struct {
	Description      string                 `tfschema:"description"`
	Location         string                 `tfschema:"location"`
	Name             string                 `tfschema:"name"`
	NetworkManagerId string                 `tfschema:"network_manager_id"`
	Tags             map[string]interface{} `tfschema:"tags"`
}

func (ManagerVerifierWorkspaceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9\_\.\-]{1,64}$`),
				"name must be between 1 and 64 characters long and can only contain letters, numbers, underscores(_), periods(.), and hyphens(-).",
			),
		},

		"network_manager_id": commonschema.ResourceIDReferenceRequiredForceNew(&ipampools.NetworkManagerId{}),

		"location": commonschema.Location(),

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tags": commonschema.Tags(),
	}
}

func (ManagerVerifierWorkspaceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ManagerVerifierWorkspaceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.VerifierWorkspaces

			var config ManagerVerifierWorkspaceResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			networkManagerId, err := verifierworkspaces.ParseNetworkManagerID(config.NetworkManagerId)
			if err != nil {
				return fmt.Errorf("parsing `network_manager_id`: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId
			id := verifierworkspaces.NewVerifierWorkspaceID(subscriptionId, networkManagerId.ResourceGroupName, networkManagerId.NetworkManagerName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := verifierworkspaces.VerifierWorkspace{
				Name:     pointer.To(config.Name),
				Location: location.Normalize(config.Location),
				Tags:     tags.Expand(config.Tags),
				Properties: &verifierworkspaces.VerifierWorkspaceProperties{
					Description: pointer.To(config.Description),
				},
			}

			if _, err := client.Create(ctx, id, payload); err != nil {
				return fmt.Errorf("performing create %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r ManagerVerifierWorkspaceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.VerifierWorkspaces

			id, err := verifierworkspaces.ParseVerifierWorkspaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			networkManagerId := ipampools.NewNetworkManagerID(id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName).ID()
			schema := ManagerVerifierWorkspaceResourceModel{
				Name:             id.VerifierWorkspaceName,
				NetworkManagerId: networkManagerId,
			}

			if model := resp.Model; model != nil {
				schema.Location = location.Normalize(model.Location)
				schema.Tags = tags.Flatten(model.Tags)

				if props := model.Properties; props != nil {
					schema.Description = pointer.From(props.Description)
				}
			}

			return metadata.Encode(&schema)
		},
	}
}

func (r ManagerVerifierWorkspaceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.VerifierWorkspaces

			id, err := verifierworkspaces.ParseVerifierWorkspaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ManagerVerifierWorkspaceResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parameters := verifierworkspaces.VerifierWorkspaceUpdate{
				Properties: &verifierworkspaces.VerifierWorkspaceUpdateProperties{},
			}

			if metadata.ResourceData.HasChange("tags") {
				parameters.Tags = tags.Expand(model.Tags)
			}

			if metadata.ResourceData.HasChange("description") {
				parameters.Properties.Description = pointer.To(model.Description)
			}

			if _, err := client.Update(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
			return nil
		},
	}
}

func (r ManagerVerifierWorkspaceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.VerifierWorkspaces

			id, err := verifierworkspaces.ParseVerifierWorkspaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			// https://github.com/Azure/azure-rest-api-specs/issues/31688
			if err := resourceVerifierWorkspaceWaitForDeleted(ctx, *client, *id); err != nil {
				return fmt.Errorf("waiting for %s to be deleted: %+v", id, err)
			}

			return nil
		},
	}
}

func resourceVerifierWorkspaceWaitForDeleted(ctx context.Context, client verifierworkspaces.VerifierWorkspacesClient, id verifierworkspaces.VerifierWorkspaceId) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal error: context had no deadline")
	}

	state := &pluginsdk.StateChangeConf{
		MinTimeout:                5 * time.Second,
		ContinuousTargetOccurence: 3,
		Pending:                   []string{"Exists"},
		Target:                    []string{"NotFound"},
		Refresh:                   resourceVerifierWorkspaceRefreshFunc(ctx, client, id),
		Timeout:                   time.Until(deadline),
	}

	if _, err := state.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be deleted: %+v", id, err)
	}

	return nil
}

func resourceVerifierWorkspaceRefreshFunc(ctx context.Context, client verifierworkspaces.VerifierWorkspacesClient, id verifierworkspaces.VerifierWorkspaceId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Checking %s status ..", id)

		resp, err := client.Get(ctx, id)
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return resp, "NotFound", nil
			}

			return resp, "Error", fmt.Errorf("retrieving %s: %+v", id, err)
		}

		return resp, "Exists", nil
	}
}
