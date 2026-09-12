// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package sreagent

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2026-01-01/agents"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity -id "/subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.App/agents/{name}"

type SreAgentResource struct{}

var (
	_ sdk.ResourceWithIdentity      = SreAgentResource{}
	_ sdk.ResourceWithUpdate        = SreAgentResource{}
	_ sdk.ResourceWithCustomizeDiff = SreAgentResource{}
)

type SreAgentModel struct {
	Name                   string                                     `tfschema:"name"`
	ResourceGroupName      string                                     `tfschema:"resource_group_name"`
	Location               string                                     `tfschema:"location"`
	ActionConfiguration    []SreAgentActionConfiguration              `tfschema:"action_configuration"`
	Identity               []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Networking             []SreAgentNetworking                       `tfschema:"networking"`
	ResourcesConfiguration []SreAgentResourcesConfiguration           `tfschema:"resources_configuration"`
	Tags                   map[string]string                          `tfschema:"tags"`
	DefaultModel           []SreAgentDefaultModel                     `tfschema:"default_model"`
	Endpoint               string                                     `tfschema:"endpoint"`
	PowerState             string                                     `tfschema:"power_state"`
	ProvisioningState      string                                     `tfschema:"provisioning_state"`
	RunningState           string                                     `tfschema:"running_state"`
}

type SreAgentActionConfiguration struct {
	AccessLevel string `tfschema:"access_level"`
	IdentityID  string `tfschema:"identity_id"`
	Mode        string `tfschema:"mode"`
}

type SreAgentResourcesConfiguration struct {
	IdentityID  string   `tfschema:"identity_id"`
	ResourceIDs []string `tfschema:"resource_ids"`
}

type SreAgentDefaultModel struct {
	Name     string `tfschema:"name"`
	Provider string `tfschema:"provider"`
}

func (SreAgentResource) ModelObject() interface{} {
	return &SreAgentModel{}
}

func (SreAgentResource) ResourceType() string {
	return "azurerm_sre_agent"
}

func (SreAgentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return agents.ValidateAgentID
}

func (SreAgentResource) Identity() resourceids.ResourceId {
	return &agents.AgentId{}
}

func (SreAgentResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type: pluginsdk.TypeString, Required: true, ForceNew: true,
			ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[A-Za-z]([-A-Za-z0-9]{0,30}[A-Za-z0-9])$`),
				"must be 2 to 32 letters, digits or hyphens, start with a letter and end with a letter or digit"),
		},
		"resource_group_name": commonschema.ResourceGroupName(),
		"location":            commonschema.Location(),
		"identity":            commonschema.SystemAssignedUserAssignedIdentityRequired(),
		"networking":          sreAgentNetworkingSchema(),
		"action_configuration": {
			Type: pluginsdk.TypeList, Optional: true, Computed: true, MaxItems: 1,
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"access_level": {
					Type: pluginsdk.TypeString, Required: true,
					ValidateFunc: validation.StringInSlice(agents.PossibleValuesForAgentAccessLevel(), false),
				},
				"identity_id": {
					Type: pluginsdk.TypeString, Required: true, ValidateFunc: azure.ValidateResourceID,
				},
				"mode": {
					Type: pluginsdk.TypeString, Required: true,
					ValidateFunc: validation.StringInSlice([]string{string(agents.AgentModeReview), string(agents.AgentModeAutonomous)}, false),
				},
			}},
		},
		"resources_configuration": {
			Type: pluginsdk.TypeList, Optional: true, Computed: true, MaxItems: 1,
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"identity_id": {
					Type: pluginsdk.TypeString, Required: true, ValidateFunc: azure.ValidateResourceID,
				},
				"resource_ids": {
					Type: pluginsdk.TypeSet, Required: true, MinItems: 1,
					Elem: &pluginsdk.Schema{Type: pluginsdk.TypeString, ValidateFunc: azure.ValidateResourceID},
				},
			}},
		},
		"tags": commonschema.Tags(),
	}
}

func (SreAgentResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"default_model": {
			Type: pluginsdk.TypeList, Computed: true,
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"name":     {Type: pluginsdk.TypeString, Computed: true},
				"provider": {Type: pluginsdk.TypeString, Computed: true},
			}},
		},
		"endpoint":           {Type: pluginsdk.TypeString, Computed: true},
		"power_state":        {Type: pluginsdk.TypeString, Computed: true},
		"provisioning_state": {Type: pluginsdk.TypeString, Computed: true},
		"running_state":      {Type: pluginsdk.TypeString, Computed: true},
	}
}

func (r SreAgentResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SreAgent.Agents
			var config SreAgentModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			id := agents.NewAgentID(metadata.Client.Account.SubscriptionId, config.ResourceGroupName, config.Name)
			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}
			payload, err := r.expandCreate(config)
			if err != nil {
				return err
			}
			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, payload, metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}

func (r SreAgentResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := agents.ParseAgentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			resp, err := metadata.Client.SreAgent.Agents.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			return r.flatten(metadata, id, resp.Model)
		},
	}
}

func (r SreAgentResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := agents.ParseAgentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			var config SreAgentModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			payload, err := r.expandPatch(config, metadata.ResourceData.HasChange)
			if err != nil {
				return err
			}
			var previousIdentity *identity.LegacySystemAndUserAssignedMap
			if metadata.ResourceData.HasChange("identity") {
				oldIdentity, _ := metadata.ResourceData.GetChange("identity")
				previousIdentity, err = identity.ExpandLegacySystemAndUserAssignedMap(oldIdentity.([]interface{}))
				if err != nil {
					return fmt.Errorf("expanding previous identity: %+v", err)
				}
			}
			if err := updateSreAgent(ctx, metadata.Client.SreAgent.Agents, *id, payload, previousIdentity); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
			return nil
		},
	}
}

func (SreAgentResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := agents.ParseAgentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			result, err := metadata.Client.SreAgent.Agents.Delete(ctx, *id)
			if response.WasNotFound(result.HttpResponse) {
				return nil
			}
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}
			if err := result.Poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("polling deletion of %s: %+v", id, err)
			}
			return nil
		},
	}
}

func (SreAgentResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(_ context.Context, metadata sdk.ResourceMetaData) error {
			diff := metadata.ResourceDiff
			// Legacy diff readers mark the set count computed when any identity ID is unknown.
			identityKnown := diff.NewValueKnown("identity") &&
				diff.NewValueKnown("identity.0.type") &&
				diff.NewValueKnown("identity.0.identity_ids.#")
			rawConfig := diff.GetRawConfig()
			if !rawConfig.IsNull() && rawConfig.Type().IsObjectType() && rawConfig.Type().HasAttribute("identity") {
				// A known block or set count does not imply that each identity reference is known.
				identityKnown = identityKnown && rawConfig.GetAttr("identity").IsWhollyKnown()
			}
			if identityKnown {
				desired, err := identity.ExpandLegacySystemAndUserAssignedMap(diff.Get("identity").([]interface{}))
				if err != nil {
					return err
				}
				if err := validateSreAgentIdentity(desired); err != nil {
					return err
				}
			}
			for _, block := range []string{"action_configuration", "resources_configuration"} {
				if diff.Id() == "" || !diff.HasChange(block) || !diff.NewValueKnown(block) {
					continue
				}
				oldValue, newValue := diff.GetChange(block)
				if len(oldValue.([]interface{})) > 0 && len(newValue.([]interface{})) == 0 {
					return sreAgentBlockRemovalError(block)
				}
			}
			if diff.NewValueKnown("networking") && (diff.Id() == "" || diff.HasChange("networking")) {
				values := diff.Get("networking").([]interface{})
				if len(values) == 0 {
					oldValue, _ := diff.GetChange("networking")
					if len(oldValue.([]interface{})) > 0 {
						return sreAgentNetworkingRemovalError()
					}
				} else if diff.NewValueKnown("networking.0.egress_mode") && diff.NewValueKnown("networking.0.subnet_id") {
					network := values[0].(map[string]interface{})
					if err := validateSreAgentNetworking(SreAgentNetworking{
						EgressMode: network["egress_mode"].(string), SubnetID: network["subnet_id"].(string),
					}); err != nil {
						return err
					}
				}
			}
			if diff.HasChange("networking.0.private_dns") && diff.NewValueKnown("networking.0.private_dns") {
				oldValue, newValue := diff.GetChange("networking.0.private_dns")
				if len(oldValue.([]interface{})) > 0 && len(newValue.([]interface{})) == 0 {
					return sreAgentPrivateDNSRemovalError()
				}
			}
			return nil
		},
	}
}

func sreAgentBlockRemovalError(block string) error {
	return fmt.Errorf("removing `%s` is not supported by this local draft: the stable API reset contract has not been established; retain the block", block)
}

func (r SreAgentResource) expandCreate(config SreAgentModel) (agents.Agent, error) {
	payload := agents.Agent{Location: location.Normalize(config.Location)}
	if config.Tags != nil {
		payload.Tags = pointer.To(config.Tags)
	}
	patch, err := r.expandPatch(config, func(key string) bool {
		switch key {
		case "identity":
			return true
		case "action_configuration":
			return len(config.ActionConfiguration) > 0
		case "resources_configuration":
			return len(config.ResourcesConfiguration) > 0
		case "networking", "networking.0.egress_mode", "networking.0.subnet_id":
			return len(config.Networking) > 0
		case "networking.0.private_dns":
			return len(config.Networking) > 0 && len(config.Networking[0].PrivateDNS) > 0
		}
		return false
	})
	if err != nil {
		return payload, err
	}
	payload.Identity = patch.Identity
	if patch.Properties != nil {
		payload.Properties = &agents.AgentProperties{
			ActionConfiguration:         patch.Properties.ActionConfiguration,
			KnowledgeGraphConfiguration: patch.Properties.KnowledgeGraphConfiguration,
			VnetConfiguration:           patch.Properties.VnetConfiguration,
			SandboxConfiguration:        patch.Properties.SandboxConfiguration,
		}
	}
	return payload, nil
}

func (SreAgentResource) expandPatch(config SreAgentModel, changed func(string) bool) (agents.AgentPatch, error) {
	payload := agents.AgentPatch{}
	if changed("tags") {
		desiredTags := config.Tags
		if desiredTags == nil {
			desiredTags = map[string]string{}
		}
		payload.Tags = pointer.To(desiredTags)
	}
	if changed("identity") {
		expanded, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(config.Identity)
		if err != nil {
			return payload, fmt.Errorf("expanding identity: %+v", err)
		}
		if err := validateSreAgentIdentity(expanded); err != nil {
			return payload, err
		}
		payload.Identity = expanded
	}
	if changed("action_configuration") {
		if len(config.ActionConfiguration) == 0 {
			return payload, sreAgentBlockRemovalError("action_configuration")
		}
		action := config.ActionConfiguration[0]
		if action.IdentityID == "" || (action.Mode != string(agents.AgentModeReview) && action.Mode != string(agents.AgentModeAutonomous)) ||
			(action.AccessLevel != string(agents.AgentAccessLevelLow) && action.AccessLevel != string(agents.AgentAccessLevelHigh)) {
			return payload, fmt.Errorf("`action_configuration` requires an identity, mode Review or Autonomous, and access_level Low or High")
		}
		payload.Properties = &agents.AgentPatchProperties{
			ActionConfiguration: &agents.ActionConfiguration{
				Identity: pointer.To(action.IdentityID), Mode: pointer.To(agents.AgentMode(action.Mode)), AccessLevel: pointer.To(agents.AgentAccessLevel(action.AccessLevel)),
			},
		}
	}
	if changed("resources_configuration") {
		if len(config.ResourcesConfiguration) == 0 {
			return payload, sreAgentBlockRemovalError("resources_configuration")
		}
		resources := config.ResourcesConfiguration[0]
		if resources.IdentityID == "" || len(resources.ResourceIDs) == 0 {
			return payload, fmt.Errorf("`resources_configuration` requires an identity and at least one explicit resource ID; an empty scope is not a removal request")
		}
		if payload.Properties == nil {
			payload.Properties = &agents.AgentPatchProperties{}
		}
		payload.Properties.KnowledgeGraphConfiguration = &agents.KnowledgeGraphConfiguration{
			Identity: pointer.To(resources.IdentityID), ManagedResources: pointer.To(resources.ResourceIDs),
		}
	}
	if changed("networking") {
		if len(config.Networking) != 1 {
			return payload, sreAgentNetworkingRemovalError()
		}
		if payload.Properties == nil {
			payload.Properties = &agents.AgentPatchProperties{}
		}
		attachmentChanged := changed("networking.0.egress_mode") || changed("networking.0.subnet_id")
		if err := expandSreAgentNetworking(config.Networking[0], payload.Properties, attachmentChanged, changed("networking.0.private_dns")); err != nil {
			return payload, err
		}
	}
	return payload, nil
}

func (SreAgentResource) flatten(metadata sdk.ResourceMetaData, id *agents.AgentId, model *agents.Agent) error {
	if model == nil {
		return fmt.Errorf("retrieving %s: model was nil", id)
	}
	state := SreAgentModel{
		Name: id.AgentName, ResourceGroupName: id.ResourceGroupName,
		Location: location.Normalize(model.Location), Tags: pointer.From(model.Tags),
	}
	flattenedIdentity, err := identity.FlattenLegacySystemAndUserAssignedMapToModel(model.Identity)
	if err != nil {
		return fmt.Errorf("flattening identity: %+v", err)
	}
	state.Identity = flattenedIdentity
	if props := model.Properties; props != nil {
		state.Networking = flattenSreAgentNetworking(props)
		if len(state.Networking) > 0 {
			if err := validateSreAgentNetworking(state.Networking[0]); err != nil {
				metadata.Logger.Warnf("SRE Agent returned networking configuration this experiment cannot configure; preserving returned state: %s", err)
			}
		}
		state.Endpoint = pointer.From(props.AgentEndpoint)
		state.PowerState = string(pointer.From(props.PowerState))
		state.ProvisioningState = string(pointer.From(props.ProvisioningState))
		state.RunningState = pointer.From(props.RunningState)
		if action := props.ActionConfiguration; action != nil {
			state.ActionConfiguration = []SreAgentActionConfiguration{{
				IdentityID: pointer.From(action.Identity), Mode: string(pointer.From(action.Mode)), AccessLevel: string(pointer.From(action.AccessLevel)),
			}}
			if mode := pointer.From(action.Mode); mode != agents.AgentModeReview && mode != agents.AgentModeAutonomous {
				metadata.Logger.Warnf("SRE Agent returned action mode %q, which this local draft cannot configure. The returned value is preserved in state.", mode)
			}
		}
		if resources := props.KnowledgeGraphConfiguration; resources != nil {
			state.ResourcesConfiguration = []SreAgentResourcesConfiguration{{
				IdentityID: pointer.From(resources.Identity), ResourceIDs: pointer.From(resources.ManagedResources),
			}}
			if len(pointer.From(resources.ManagedResources)) == 0 {
				metadata.Logger.Warn("SRE Agent returned an empty managed-resource scope, which this local draft cannot configure. The returned value is preserved in state.")
			}
		}
		if model := props.DefaultModel; model != nil {
			state.DefaultModel = []SreAgentDefaultModel{{Name: pointer.From(model.Name), Provider: pointer.From(model.Provider)}}
		}
	}
	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}
	return metadata.Encode(&state)
}
