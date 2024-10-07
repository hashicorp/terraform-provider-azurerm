// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/customipprefixes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type CustomIpPrefixModel struct {
	CIDR                        string                 `tfschema:"cidr"`
	CommissioningEnabled        bool                   `tfschema:"commissioning_enabled"`
	InternetAdvertisingDisabled bool                   `tfschema:"internet_advertising_disabled"`
	Location                    string                 `tfschema:"location"`
	Name                        string                 `tfschema:"name"`
	ParentCustomIPPrefixID      string                 `tfschema:"parent_custom_ip_prefix_id"`
	ROAValidityEndDate          string                 `tfschema:"roa_validity_end_date"`
	ResourceGroupName           string                 `tfschema:"resource_group_name"`
	Tags                        map[string]interface{} `tfschema:"tags"`
	WANValidationSignedMessage  string                 `tfschema:"wan_validation_signed_message"`
	Zones                       []string               `tfschema:"zones"`
}

var (
	_ sdk.ResourceWithUpdate = CustomIpPrefixResource{}
)

type CustomIpPrefixResource struct {
	client *customipprefixes.CustomIPPrefixesClient
}

func (CustomIpPrefixResource) ResourceType() string {
	return "azurerm_custom_ip_prefix"
}

func (CustomIpPrefixResource) ModelObject() interface{} {
	return &CustomIpPrefixModel{}
}

func (CustomIpPrefixResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return customipprefixes.ValidateCustomIPPrefixID
}
func (r CustomIpPrefixResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CustomIpPrefixName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"cidr": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: func(i interface{}, k string) (warnings []string, errors []error) {
				v, ok := i.(string)
				if !ok {
					errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
					return
				}

				if _, _, err := net.ParseCIDR(v); err != nil {
					errors = append(errors, fmt.Errorf("expected %q to be a valid IPv4 or IPv6 network, got %v: %v", k, i, err))
				}

				return
			},
		},

		"parent_custom_ip_prefix_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: customipprefixes.ValidateCustomIPPrefixID,
		},

		"roa_validity_end_date": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			ValidateFunc: func(i interface{}, k string) (warnings []string, errors []error) {
				v, ok := i.(string)
				if !ok {
					errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
					return warnings, errors
				}

				if _, err := time.Parse("2006-01-02", v); err != nil {
					errors = append(errors, fmt.Errorf("expected %q to be a valid date in the format YYYY-MM-DD, got %q: %+v", k, i, err))
				}

				return warnings, errors
			},
		},

		"wan_validation_signed_message": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"commissioning_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"internet_advertising_disabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"tags": commonschema.Tags(),

		"zones": commonschema.ZonesMultipleOptionalForceNew(),
	}
}

func (r CustomIpPrefixResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r CustomIpPrefixResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 9 * time.Hour,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			r.client = metadata.Client.Network.Client.CustomIPPrefixes
			subscriptionId := metadata.Client.Account.SubscriptionId

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal-error: context has no deadline")
			}

			var model CustomIpPrefixModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := customipprefixes.NewCustomIPPrefixID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := r.client.Get(ctx, id, customipprefixes.DefaultGetOperationOptions())
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			_, cidr, err := net.ParseCIDR(model.CIDR)
			if err != nil {
				return fmt.Errorf("parsing `cidr`: %+v", err)
			}

			results, err := r.client.ListAll(ctx, commonids.NewSubscriptionID(subscriptionId))
			if err != nil {
				return fmt.Errorf("listing existing %s: %+v", id, err)
			}

			if prefixes := results.Model; prefixes != nil {
				for _, prefix := range *prefixes {
					if prefix.Properties != nil && prefix.Properties.Cidr != nil {
						_, netw, err := net.ParseCIDR(*prefix.Properties.Cidr)
						if err != nil {
							// couldn't parse the existing custom prefix, so skip it
							continue
						}
						if cidr == netw {
							return metadata.ResourceRequiresImport(r.ResourceType(), id)
						}
					}
				}
			}

			payload := customipprefixes.CustomIPPrefix{
				Name:             &model.Name,
				Location:         pointer.To(location.Normalize(model.Location)),
				Tags:             tags.Expand(model.Tags),
				ExtendedLocation: nil,
				Properties: &customipprefixes.CustomIPPrefixPropertiesFormat{
					Cidr:              &model.CIDR,
					CommissionedState: pointer.To(customipprefixes.CommissionedStateProvisioning),
				},
			}

			if model.ParentCustomIPPrefixID != "" {
				payload.Properties.CustomIPPrefixParent = &customipprefixes.SubResource{
					Id: &model.ParentCustomIPPrefixID,
				}
			}

			if model.WANValidationSignedMessage != "" {
				payload.Properties.SignedMessage = &model.WANValidationSignedMessage
			}

			if model.ROAValidityEndDate != "" {
				roaValidityEndDate, err := time.Parse("2006-01-02", model.ROAValidityEndDate)
				if err != nil {
					return err
				}
				authorizationMessage := fmt.Sprintf("%s|%s|%s", subscriptionId, model.CIDR, roaValidityEndDate.Format("20060102"))
				payload.Properties.AuthorizationMessage = &authorizationMessage
			}

			if len(model.Zones) > 0 {
				payload.Zones = &model.Zones
			}

			if err := r.client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			stateConf := &pluginsdk.StateChangeConf{
				Pending:    []string{string(customipprefixes.ProvisioningStateUpdating)},
				Target:     []string{string(customipprefixes.ProvisioningStateSucceeded)},
				Refresh:    r.provisioningStateRefreshFunc(ctx, id),
				MinTimeout: 2 * time.Minute,
				Timeout:    time.Until(deadline),
			}
			if _, err = stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for ProvisioningState of %s: %+v", id, err)
			}

			desiredState := customipprefixes.CommissionedStateProvisioned
			if model.CommissioningEnabled {
				if model.InternetAdvertisingDisabled {
					desiredState = customipprefixes.CommissionedStateCommissionedNoInternetAdvertise
				} else {
					desiredState = customipprefixes.CommissionedStateCommissioned
				}
			}

			commissionedState, err := r.updateCommissionedState(ctx, id, desiredState)
			if err != nil {
				return err
			}
			if commissionedState == nil {
				return fmt.Errorf("waiting for CommissionedState: final commissionedState was nil")
			}

			log.Printf("[DEBUG] Final CommissionedState is %q for %s..", *commissionedState, id)
			metadata.SetID(id)
			return nil
		},
	}
}

func (r CustomIpPrefixResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 17 * time.Hour,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			r.client = metadata.Client.Network.Client.CustomIPPrefixes

			id, err := customipprefixes.ParseCustomIPPrefixID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Info("Decoding state...")
			var state CustomIpPrefixModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			desiredState := customipprefixes.CommissionedStateProvisioned
			if state.CommissioningEnabled {
				if state.InternetAdvertisingDisabled {
					desiredState = customipprefixes.CommissionedStateCommissionedNoInternetAdvertise
				} else {
					desiredState = customipprefixes.CommissionedStateCommissioned
				}
			}

			commissionedState, err := r.updateCommissionedState(ctx, *id, desiredState)
			if err != nil {
				return err
			}
			if commissionedState == nil {
				return fmt.Errorf("waiting for CommissionedState: final commissionedState was nil")
			}

			log.Printf("[DEBUG] Final CommissionedState is %q for %s..", *commissionedState, id)
			return nil
		},
	}
}

func (r CustomIpPrefixResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			r.client = metadata.Client.Network.Client.CustomIPPrefixes

			id, err := customipprefixes.ParseCustomIPPrefixID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := r.client.Get(ctx, *id, customipprefixes.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := CustomIpPrefixModel{
				Name:              id.CustomIPPrefixName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := existing.Model; model != nil {
				state.Location = location.NormalizeNilable(model.Location)
				state.Tags = tags.Flatten(model.Tags)
				state.Zones = pointer.From(model.Zones)

				if props := model.Properties; props != nil {
					state.CIDR = pointer.From(props.Cidr)
					state.InternetAdvertisingDisabled = pointer.From(props.NoInternetAdvertise)
					state.WANValidationSignedMessage = pointer.From(props.SignedMessage)

					if parent := props.CustomIPPrefixParent; parent != nil {
						state.ParentCustomIPPrefixID = pointer.From(parent.Id)
					}

					if props.AuthorizationMessage != nil {
						authMessage := strings.Split(*props.AuthorizationMessage, "|")
						if len(authMessage) == 3 {
							if roaValidityEndDate, err := time.Parse("20060102", authMessage[2]); err == nil {
								state.ROAValidityEndDate = roaValidityEndDate.Format("2006-01-02")
							}
						}
					}

					switch pointer.From(props.CommissionedState) {
					case customipprefixes.CommissionedStateCommissioning, customipprefixes.CommissionedStateCommissioned, customipprefixes.CommissionedStateCommissionedNoInternetAdvertise:
						state.CommissioningEnabled = true
					default:
						state.CommissioningEnabled = false
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r CustomIpPrefixResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 17 * time.Hour,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			r.client = metadata.Client.Network.Client.CustomIPPrefixes

			id, err := customipprefixes.ParseCustomIPPrefixID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			// Must be de-provisioned before deleting
			if _, err = r.updateCommissionedState(ctx, *id, customipprefixes.CommissionedStateDeprovisioned); err != nil {
				return err
			}

			if err := r.client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

type commissionedStates []customipprefixes.CommissionedState

func (t commissionedStates) contains(i customipprefixes.CommissionedState) bool {
	for _, s := range t {
		if i == s {
			return true
		}
	}
	return false
}

func (t commissionedStates) strings() (out []string) {
	for _, s := range t {
		out = append(out, string(s))
	}
	return
}

// updateCommissionedState implements a state machine to coordinate transitions between different values of CommissionedState for both v4 and v6 prefixes.
// The provided desiredState should be the sought after end state, and the method will work out a path to achieving that state and walk the resource to get there.
func (r CustomIpPrefixResource) updateCommissionedState(ctx context.Context, id customipprefixes.CustomIPPrefixId, desiredState customipprefixes.CommissionedState) (*customipprefixes.CommissionedState, error) {
	existing, err := r.client.Get(ctx, id, customipprefixes.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving existing %s: %+v", id, err)
	}
	if existing.Model == nil {
		return nil, fmt.Errorf("retrieving existing %s: `model` was nil", id)
	}
	if existing.Model.Properties == nil {
		return nil, fmt.Errorf("retrieving existing %s: `properties` was nil", id)
	}

	initialState := existing.Model.Properties.CommissionedState

	log.Printf("[DEBUG] Updating CommissionedState for %s from current value %q to desired value %q..", id, *initialState, desiredState)

	// stateTree is a map of desired state, to a map of current state, to the list of transition states needed to get there
	stateTree := map[customipprefixes.CommissionedState]map[customipprefixes.CommissionedState][]customipprefixes.CommissionedState{
		customipprefixes.CommissionedStateDeprovisioned: {
			customipprefixes.CommissionedStateProvisioned:                     {customipprefixes.CommissionedStateDeprovisioning},
			customipprefixes.CommissionedStateCommissioned:                    {customipprefixes.CommissionedStateDecommissioning, customipprefixes.CommissionedStateDeprovisioning},
			customipprefixes.CommissionedStateCommissionedNoInternetAdvertise: {customipprefixes.CommissionedStateDecommissioning, customipprefixes.CommissionedStateDeprovisioning},
		},
		customipprefixes.CommissionedStateProvisioned: {
			customipprefixes.CommissionedStateDeprovisioned:                   {customipprefixes.CommissionedStateProvisioning},
			customipprefixes.CommissionedStateCommissioned:                    {customipprefixes.CommissionedStateDecommissioning},
			customipprefixes.CommissionedStateCommissionedNoInternetAdvertise: {customipprefixes.CommissionedStateDecommissioning},
		},
		customipprefixes.CommissionedStateCommissioned: {
			customipprefixes.CommissionedStateDeprovisioned:                   {customipprefixes.CommissionedStateProvisioning, customipprefixes.CommissionedStateCommissioning},
			customipprefixes.CommissionedStateProvisioned:                     {customipprefixes.CommissionedStateCommissioning},
			customipprefixes.CommissionedStateCommissionedNoInternetAdvertise: {customipprefixes.CommissionedStateCommissioning},
		},
		customipprefixes.CommissionedStateCommissionedNoInternetAdvertise: {
			customipprefixes.CommissionedStateDeprovisioned: {customipprefixes.CommissionedStateProvisioning, customipprefixes.CommissionedStateCommissioning},
			customipprefixes.CommissionedStateProvisioned:   {customipprefixes.CommissionedStateCommissioning},
			customipprefixes.CommissionedStateCommissioned:  {customipprefixes.CommissionedStateDecommissioning, customipprefixes.CommissionedStateCommissioning},
		},
	}

	// transitioningStatesFor returns the known transitioning states for the desired goal state
	transitioningStatesFor := func(finalState customipprefixes.CommissionedState) (out commissionedStates) {
		switch finalState {
		case customipprefixes.CommissionedStateProvisioned:
			out = commissionedStates{customipprefixes.CommissionedStateProvisioning, customipprefixes.CommissionedStateDecommissioning}
		case customipprefixes.CommissionedStateDeprovisioned:
			out = commissionedStates{customipprefixes.CommissionedStateDeprovisioning}
		case customipprefixes.CommissionedStateCommissioned:
			out = commissionedStates{customipprefixes.CommissionedStateCommissioning}
		}
		return
	}

	// finalStatesFor returns the known final states for the current transitioning state
	finalStatesFor := func(transitioningState customipprefixes.CommissionedState) (out commissionedStates) {
		switch transitioningState {
		case customipprefixes.CommissionedStateProvisioning:
			out = commissionedStates{customipprefixes.CommissionedStateProvisioned}
		case customipprefixes.CommissionedStateDeprovisioning:
			out = commissionedStates{customipprefixes.CommissionedStateDeprovisioned}
		case customipprefixes.CommissionedStateCommissioning:
			out = commissionedStates{customipprefixes.CommissionedStateCommissioned, customipprefixes.CommissionedStateCommissionedNoInternetAdvertise}
		case customipprefixes.CommissionedStateDecommissioning:
			out = commissionedStates{customipprefixes.CommissionedStateProvisioned}
		}
		return
	}

	// shouldNotAdvertise determines whether to set the noInternetAdvertise flag, which can only be set at the point of transitioning to `Commissioning`
	shouldNotAdvertise := func(steppingState customipprefixes.CommissionedState) *bool {
		if steppingState == customipprefixes.CommissionedStateCommissioning {
			switch desiredState {
			case customipprefixes.CommissionedStateCommissioned:
				return pointer.To(false)
			case customipprefixes.CommissionedStateCommissionedNoInternetAdvertise:
				return pointer.To(true)
			}
		}
		return nil
	}

	if plan, ok := stateTree[desiredState]; ok {
		lastKnownState := initialState

		// If we're already transitioning to the desiredState, wait for this to complete
		if transitioningStatesFor(desiredState).contains(pointer.From(initialState)) {
			if lastKnownState, err = r.waitForCommissionedState(ctx, id, transitioningStatesFor(desiredState), commissionedStates{desiredState}); err != nil {
				return lastKnownState, err
			}
		}

		// Return early if the desiredState was already reached
		if *lastKnownState == desiredState {
			return lastKnownState, nil
		}

		for startingState, path := range plan {
			// Look for a plan that works from our lastKnownState
			if *lastKnownState == startingState || transitioningStatesFor(startingState).contains(*lastKnownState) {

				// If we're currently transitioning to the startingState for this plan, wait for this to complete before proceeding
				if lastKnownState, err = r.waitForCommissionedState(ctx, id, transitioningStatesFor(startingState), commissionedStates{startingState}); err != nil {
					return lastKnownState, err
				}

				retries := 0
				const maxRetries = 2

				// Iterate the plan
				for i := 0; i < len(path); i++ {
					steppingState := path[i]

					// Instruct the resource to transition to the next CommissionedState according to the plan
					if err = r.setCommissionedState(ctx, id, steppingState, shouldNotAdvertise(steppingState)); err != nil {
						return lastKnownState, err
					}

					// Wait for the CommissionedState to be reached
					latestState, err := r.waitForCommissionedState(ctx, id, commissionedStates{steppingState}, finalStatesFor(steppingState))
					if err != nil {
						// Known issue where the previous CommissioningState was reported prematurely by the API, so we reattempt up to maxRetries times
						if lastKnownState != nil && latestState != nil && *latestState == *lastKnownState && retries < maxRetries {
							retries++
							i--
							log.Printf("[DEBUG] Retrying %d of %d times to set CommissionedState field to %q (current state: %q) for %s..", retries, maxRetries, steppingState, *latestState, id)
							continue
						}

						return lastKnownState, err
					}

					// Update the lastKnownState, so we can monitor for retries on the next iteration
					lastKnownState = latestState
				}

				return r.waitForCommissionedState(ctx, id, transitioningStatesFor(desiredState), commissionedStates{desiredState})
			}
		}
	} else {
		return nil, fmt.Errorf("internal-error: unsupported state %q", desiredState)
	}

	return nil, fmt.Errorf("internal-error: could not transition CommissionedState to %q", desiredState)
}

// setCommissionedState sends a PUT request to effect a transition to a different CommissionedState. The provided
// desiredState should always be a contextual transition state rather than the desired end state (i.e. procedural).
func (r CustomIpPrefixResource) setCommissionedState(ctx context.Context, id customipprefixes.CustomIPPrefixId, desiredState customipprefixes.CommissionedState, noInternetAdvertise *bool) error {
	existing, err := r.client.Get(ctx, id, customipprefixes.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving existing %s: %+v", id, err)
	}
	if existing.Model == nil {
		return fmt.Errorf("retrieving existing %s: `model` was nil", id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving existing %s: `properties` was nil", id)
	}

	existing.Model.Properties.CommissionedState = pointer.To(desiredState)
	existing.Model.Properties.NoInternetAdvertise = noInternetAdvertise

	log.Printf("[DEBUG] Updating the CommissionedState field to %q for %s..", desiredState, id)
	if err := r.client.CreateOrUpdateThenPoll(ctx, id, *existing.Model); err != nil {
		return fmt.Errorf("updating CommissionedState to %q for %s: %+v", desiredState, id, err)
	}

	return nil
}

// waitForCommissionedState polls the resource and returns when one of the targetStates is reached and seen for 3
// consecutive polls, also returning an error if a state is reached that isn't in pendingStates or targetStates. Waits
// for 10 minutes before polling to account for delays in the service reporting the actual latest state, since this
// method is usually called soon after setting a new CommissionedState (known service bug).
func (r CustomIpPrefixResource) waitForCommissionedState(ctx context.Context, id customipprefixes.CustomIPPrefixId, pendingStates, targetStates commissionedStates) (*customipprefixes.CommissionedState, error) {
	log.Printf("[DEBUG] Polling for the CommissionedState field for %s..", id)
	timeout, ok := ctx.Deadline()
	if !ok {
		return nil, fmt.Errorf("internal-error: context has no deadline")
	}

	stateConf := &pluginsdk.StateChangeConf{
		Delay:        10 * time.Minute,
		Pending:      pendingStates.strings(),
		Target:       targetStates.strings(),
		Refresh:      r.commissionedStateRefreshFunc(ctx, id),
		PollInterval: 5 * time.Minute,
		Timeout:      time.Until(timeout),

		// `Provisioned` is known to flip-flop
		ContinuousTargetOccurence: 3,
	}

	result, err := stateConf.WaitForStateContext(ctx)

	if result == nil {
		return nil, fmt.Errorf("retrieving %s: response was nil", id)
	}

	prefix, ok := result.(customipprefixes.CustomIPPrefix)
	if !ok {
		return nil, fmt.Errorf("retrieving %s: response was not a valid Custom IP Prefix", id)
	}

	if prefix.Properties == nil {
		return prefix.Properties.CommissionedState, fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	if err != nil {
		return prefix.Properties.CommissionedState, fmt.Errorf("waiting for CommissionedState of %s: %+v", id, err)
	}

	return prefix.Properties.CommissionedState, nil
}

func (r CustomIpPrefixResource) commissionedStateRefreshFunc(ctx context.Context, id customipprefixes.CustomIPPrefixId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := r.client.Get(ctx, id, customipprefixes.DefaultGetOperationOptions())
		if err != nil {
			return nil, "", fmt.Errorf("polling for %s: %+v", id, err)
		}

		if res.Model == nil {
			return nil, "", fmt.Errorf("polling for %s: `model` was nil", id)
		}
		if res.Model.Properties == nil {
			return nil, "", fmt.Errorf("polling for %s: `properties` was nil", id)
		}

		return res, string(pointer.From(res.Model.Properties.CommissionedState)), nil
	}
}

func (r CustomIpPrefixResource) provisioningStateRefreshFunc(ctx context.Context, id customipprefixes.CustomIPPrefixId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := r.client.Get(ctx, id, customipprefixes.DefaultGetOperationOptions())
		if err != nil {
			return nil, "", fmt.Errorf("polling for %s: %+v", id, err)
		}

		if res.Model == nil {
			return nil, "", fmt.Errorf("polling for %s: `model` was nil", id)
		}
		if res.Model.Properties == nil {
			return nil, "", fmt.Errorf("polling for %s: `properties` was nil", id)
		}

		return res, string(pointer.From(res.Model.Properties.ProvisioningState)), nil
	}
}
