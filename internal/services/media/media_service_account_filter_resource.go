// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package media

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-08-01/accountfilters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var _ sdk.ResourceWithUpdate = AccountFilterResource{}
var _ sdk.ResourceWithDeprecationAndNoReplacement = AccountFilterResource{}

type AccountFilterModel struct {
	Name                     string                  `tfschema:"name"`
	ResourceGroupName        string                  `tfschema:"resource_group_name"`
	MediaServicesAccountName string                  `tfschema:"media_services_account_name"`
	FirstQualityBitrate      int64                   `tfschema:"first_quality_bitrate"`
	PresentationTimeRange    []PresentationTimeRange `tfschema:"presentation_time_range"`
	TrackSelection           []TrackSelection        `tfschema:"track_selection"`
}

type PresentationTimeRange struct {
	EndInUnits                  int64 `tfschema:"end_in_units"`
	ForceEnd                    bool  `tfschema:"force_end"`
	LiveBackoffInUnits          int64 `tfschema:"live_backoff_in_units"`
	PresentationWindowInUnits   int64 `tfschema:"presentation_window_in_units"`
	StartInUnits                int64 `tfschema:"start_in_units"`
	UnitTimescaleInMillisceonds int64 `tfschema:"unit_timescale_in_milliseconds"`
}

type TrackSelection struct {
	Conditions []Condition `tfschema:"condition"`
}

type Condition struct {
	Operation string `tfschema:"operation"`
	Property  string `tfschema:"property"`
	Value     string `tfschema:"value"`
}

type AccountFilterResource struct{}

func (r AccountFilterResource) DeprecationMessage() string {
	return azureMediaRetirementMessage
}

func (r AccountFilterResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"media_services_account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"first_quality_bitrate": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"presentation_time_range": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"end_in_units": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(0),
						AtLeastOneOf: []string{
							"presentation_time_range.0.end_in_units", "presentation_time_range.0.force_end", "presentation_time_range.0.live_backoff_in_units",
							"presentation_time_range.0.presentation_window_in_units", "presentation_time_range.0.start_in_units",
						},
					},

					"force_end": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						AtLeastOneOf: []string{
							"presentation_time_range.0.end_in_units", "presentation_time_range.0.force_end", "presentation_time_range.0.live_backoff_in_units",
							"presentation_time_range.0.presentation_window_in_units", "presentation_time_range.0.start_in_units",
						},
					},

					"live_backoff_in_units": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(0),
						AtLeastOneOf: []string{
							"presentation_time_range.0.end_in_units", "presentation_time_range.0.force_end", "presentation_time_range.0.live_backoff_in_units",
							"presentation_time_range.0.presentation_window_in_units", "presentation_time_range.0.start_in_units",
						},
					},

					"presentation_window_in_units": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(0),
						AtLeastOneOf: []string{
							"presentation_time_range.0.end_in_units", "presentation_time_range.0.force_end", "presentation_time_range.0.live_backoff_in_units",
							"presentation_time_range.0.presentation_window_in_units", "presentation_time_range.0.start_in_units",
						},
					},

					"start_in_units": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(0),
						AtLeastOneOf: []string{
							"presentation_time_range.0.end_in_units", "presentation_time_range.0.force_end", "presentation_time_range.0.live_backoff_in_units",
							"presentation_time_range.0.presentation_window_in_units", "presentation_time_range.0.start_in_units",
						},
					},

					"unit_timescale_in_milliseconds": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntAtLeast(1),
					},
				},
			},
		},

		"track_selection": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"condition": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"operation": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(accountfilters.FilterTrackPropertyCompareOperationEqual),
										string(accountfilters.FilterTrackPropertyCompareOperationNotEqual),
									}, false),
								},

								"property": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(accountfilters.FilterTrackPropertyTypeBitrate),
										string(accountfilters.FilterTrackPropertyTypeFourCC),
										string(accountfilters.FilterTrackPropertyTypeLanguage),
										string(accountfilters.FilterTrackPropertyTypeName),
										string(accountfilters.FilterTrackPropertyTypeType),
									}, false),
								},

								"value": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r AccountFilterResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r AccountFilterResource) ResourceType() string {
	return "azurerm_media_services_account_filter"
}

func (r AccountFilterResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return accountfilters.ValidateAccountFilterID
}

func (r AccountFilterResource) ModelObject() interface{} {
	return &AccountFilterModel{}
}

func (r AccountFilterResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			metadata.Logger.Info("Decoding state..")
			var state AccountFilterModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			client := metadata.Client.Media.V20220801Client.AccountFilters
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := accountfilters.NewAccountFilterID(subscriptionId, state.ResourceGroupName, state.MediaServicesAccountName, state.Name)
			metadata.Logger.Infof("creating %s", id)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			input := accountfilters.AccountFilter{
				Name: utils.String(state.Name),
				Properties: &accountfilters.MediaFilterProperties{
					PresentationTimeRange: expandAccountFilterPresentationTimeRange(state.PresentationTimeRange),
					Tracks:                expandAccountFilterTrackSelections(state.TrackSelection),
				},
			}
			if state.FirstQualityBitrate != 0 {
				input.Properties.FirstQuality = &accountfilters.FirstQuality{
					Bitrate: state.FirstQualityBitrate,
				}
			}

			if _, err = client.CreateOrUpdate(ctx, id, input); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r AccountFilterResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Media.V20220801Client.AccountFilters
			id, err := accountfilters.ParseAccountFilterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("retrieving %s", *id)
			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					metadata.Logger.Infof("%s was not found - removing from state!", *id)
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil || resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: got empty Model", *id)
			}
			prop := resp.Model.Properties
			var firstQualityBitrate int64
			if prop.FirstQuality != nil {
				firstQualityBitrate = prop.FirstQuality.Bitrate
			}

			return metadata.Encode(&AccountFilterModel{
				Name:                     id.AccountFilterName,
				ResourceGroupName:        id.ResourceGroupName,
				MediaServicesAccountName: id.MediaServiceName,
				FirstQualityBitrate:      firstQualityBitrate,
				PresentationTimeRange:    flattenAccountFilterPresentationTimeRange(prop.PresentationTimeRange),
				TrackSelection:           flattenAccountFilterTracks(prop.Tracks),
			})
		},
		Timeout: 5 * time.Minute,
	}
}

func (r AccountFilterResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := accountfilters.ParseAccountFilterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("updating %s..", *id)
			client := metadata.Client.Media.V20220801Client.AccountFilters

			var state AccountFilterModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if resp.Model == nil || resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: got empty Model", *id)
			}
			existing := resp.Model

			if metadata.ResourceData.HasChange("first_quality_bitrate") {
				existing.Properties.FirstQuality = &accountfilters.FirstQuality{
					Bitrate: state.FirstQualityBitrate,
				}
			}

			if metadata.ResourceData.HasChange("presentation_time_range") {
				existing.Properties.PresentationTimeRange = expandAccountFilterPresentationTimeRange(state.PresentationTimeRange)
			}

			if metadata.ResourceData.HasChange("track_selection") {
				existing.Properties.Tracks = expandAccountFilterTrackSelections(state.TrackSelection)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *existing); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r AccountFilterResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Media.V20220801Client.AccountFilters
			id, err := accountfilters.ParseAccountFilterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting %s..", *id)
			resp, err := client.Delete(ctx, *id)
			if err != nil && !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func expandAccountFilterPresentationTimeRange(input []PresentationTimeRange) *accountfilters.PresentationTimeRange {
	if len(input) == 0 {
		return nil
	}

	timeRange := input[0]
	presentationTimeRange := &accountfilters.PresentationTimeRange{
		ForceEndTimestamp: utils.Bool(timeRange.ForceEnd),
		Timescale:         utils.Int64((incrementsInASecond * nanoSecondsInAIncrement) / milliSecondsInASecond / timeRange.UnitTimescaleInMillisceonds),
	}

	if timeRange.EndInUnits != 0 {
		presentationTimeRange.EndTimestamp = utils.Int64(timeRange.EndInUnits * milliSecondsInASecond)
	}

	if timeRange.LiveBackoffInUnits != 0 {
		presentationTimeRange.LiveBackoffDuration = utils.Int64(timeRange.LiveBackoffInUnits * milliSecondsInASecond)
	}

	if timeRange.PresentationWindowInUnits != 0 {
		presentationTimeRange.PresentationWindowDuration = utils.Int64(timeRange.PresentationWindowInUnits * milliSecondsInASecond)
	}

	if timeRange.StartInUnits != 0 {
		presentationTimeRange.StartTimestamp = utils.Int64(timeRange.StartInUnits * milliSecondsInASecond)
	}

	return presentationTimeRange
}

func flattenAccountFilterPresentationTimeRange(input *accountfilters.PresentationTimeRange) []PresentationTimeRange {
	if input == nil {
		return make([]PresentationTimeRange, 0)
	}

	var unitTimescaleInMillisceonds, endInUnits, startInUnits, liveBackoffInUnits, presentationWindowInUnits int64
	var forceEnd bool
	if input.Timescale != nil {
		unitTimescaleInMillisceonds = (incrementsInASecond * nanoSecondsInAIncrement) / milliSecondsInASecond / *input.Timescale
	}

	if input.EndTimestamp != nil {
		endInUnits = *input.EndTimestamp / milliSecondsInASecond
	}

	if input.ForceEndTimestamp != nil {
		forceEnd = *input.ForceEndTimestamp
	}

	if input.LiveBackoffDuration != nil {
		liveBackoffInUnits = *input.LiveBackoffDuration / milliSecondsInASecond
	}

	if input.PresentationWindowDuration != nil {
		presentationWindowInUnits = *input.PresentationWindowDuration / milliSecondsInASecond
	}

	if input.StartTimestamp != nil {
		startInUnits = *input.StartTimestamp / milliSecondsInASecond
	}
	return []PresentationTimeRange{
		{
			EndInUnits:                  endInUnits,
			ForceEnd:                    forceEnd,
			LiveBackoffInUnits:          liveBackoffInUnits,
			PresentationWindowInUnits:   presentationWindowInUnits,
			StartInUnits:                startInUnits,
			UnitTimescaleInMillisceonds: unitTimescaleInMillisceonds,
		},
	}
}

func expandAccountFilterTrackSelections(input []TrackSelection) *[]accountfilters.FilterTrackSelection {
	results := make([]accountfilters.FilterTrackSelection, 0)

	for _, track := range input {
		conditions := make([]accountfilters.FilterTrackPropertyCondition, 0)
		for _, condition := range track.Conditions {
			conditions = append(conditions, accountfilters.FilterTrackPropertyCondition{
				Operation: accountfilters.FilterTrackPropertyCompareOperation(condition.Operation),
				Property:  accountfilters.FilterTrackPropertyType(condition.Property),
				Value:     condition.Value,
			})
		}
		results = append(results, accountfilters.FilterTrackSelection{
			TrackSelections: conditions,
		})
	}

	return &results
}

func flattenAccountFilterTracks(input *[]accountfilters.FilterTrackSelection) []TrackSelection {
	if input == nil {
		return make([]TrackSelection, 0)
	}

	result := make([]TrackSelection, 0)
	for _, track := range *input {
		conditions := make([]Condition, 0)
		for _, condition := range track.TrackSelections {
			conditions = append(conditions, Condition{
				Operation: string(condition.Operation),
				Property:  string(condition.Property),
				Value:     condition.Value,
			})
		}
		result = append(result, TrackSelection{
			Conditions: conditions,
		})
	}

	return result
}
