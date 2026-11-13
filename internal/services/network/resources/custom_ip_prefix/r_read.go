// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package custom_ip_prefix

import (
	"context"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/customipprefixes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			r.client = metadata.Client.Network.CustomIPPrefixes

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

type commissionedStates []customipprefixes.CommissionedState

func (t commissionedStates) contains(i customipprefixes.CommissionedState) bool {
	return slices.Contains(t, i)
}

func (t commissionedStates) strings() (out []string) {
	for _, s := range t {
		out = append(out, string(s))
	}
	return
}

// updateCommissionedState implements a state machine to coordinate transitions between different values of CommissionedState for both v4 and v6 prefixes.
// The provided desiredState should be the sought after end state, and the method will work out a path to achieving that state and walk the resource to get there.
func (r Resource) updateCommissionedState(ctx context.Context, id customipprefixes.CustomIPPrefixId, desiredState customipprefixes.CommissionedState) (*customipprefixes.CommissionedState, error) {
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
func (r Resource) setCommissionedState(ctx context.Context, id customipprefixes.CustomIPPrefixId, desiredState customipprefixes.CommissionedState, noInternetAdvertise *bool) error {
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

	if err := r.client.CreateOrUpdateThenPoll(ctx, id, *existing.Model); err != nil {
		return fmt.Errorf("updating CommissionedState to %q for %s: %+v", desiredState, id, err)
	}

	return nil
}

// waitForCommissionedState polls the resource and returns when one of the targetStates is reached and seen for 3
// consecutive polls, also returning an error if a state is reached that isn't in pendingStates or targetStates. Waits
// for 10 minutes before polling to account for delays in the service reporting the actual latest state, since this
// method is usually called soon after setting a new CommissionedState (known service bug).
func (r Resource) waitForCommissionedState(ctx context.Context, id customipprefixes.CustomIPPrefixId, pendingStates, targetStates commissionedStates) (*customipprefixes.CommissionedState, error) {
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

	resp, ok := result.(customipprefixes.GetOperationResponse)
	if !ok {
		return nil, fmt.Errorf("retrieving %s: response was not a valid Custom IP Prefix", id)
	}

	if resp.Model == nil {
		return nil, fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	if resp.Model.Properties == nil {
		return nil, fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	if err != nil {
		return resp.Model.Properties.CommissionedState, fmt.Errorf("waiting for CommissionedState of %s: %+v", id, err)
	}

	return resp.Model.Properties.CommissionedState, nil
}

func (r Resource) commissionedStateRefreshFunc(ctx context.Context, id customipprefixes.CustomIPPrefixId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
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

		return res, pointer.FromEnum(res.Model.Properties.CommissionedState), nil
	}
}

func (r Resource) provisioningStateRefreshFunc(ctx context.Context, id customipprefixes.CustomIPPrefixId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
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

		return res, pointer.FromEnum(res.Model.Properties.ProvisioningState), nil
	}
}
