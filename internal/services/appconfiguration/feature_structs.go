// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appconfiguration

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"
)

const (
	TargetingFilterName  = "Microsoft.Targeting"
	TimewindowFilterName = "Microsoft.TimeWindow"
	PercentageFilterName = "Microsoft.Percentage"
)

type ClientFilter struct {
	Filters []interface{}
}

func (p *ClientFilter) UnmarshalJSON(b []byte) error {
	var tempIntf []interface{}

	if err := json.Unmarshal(b, &tempIntf); err != nil {
		return err
	}

	filtersOut := make([]interface{}, 0)
	for _, filterRawIntf := range tempIntf {
		filterRaw, ok := filterRawIntf.(map[string]interface{})
		if !ok {
			return fmt.Errorf("wtf")
		}
		nameRaw, ok := filterRaw["name"]
		if !ok {
			return fmt.Errorf("missing name ...")
		}

		name := nameRaw.(string)
		switch strings.ToLower(name) {
		case "microsoft.targeting":
			{
				var out TargetingFeatureFilter
				mpc := mapstructure.DecoderConfig{TagName: "json", Result: &out}
				mpd, err := mapstructure.NewDecoder(&mpc)
				if err != nil {
					return err
				}
				err = mpd.Decode(filterRaw)
				if err != nil {
					return err
				}
				filtersOut = append(filtersOut, out)
			}
		case "microsoft.timewindow":
			{
				var out TimewindowFeatureFilter
				mpc := mapstructure.DecoderConfig{TagName: "json", Result: &out}
				mpd, err := mapstructure.NewDecoder(&mpc)
				if err != nil {
					return err
				}
				err = mpd.Decode(filterRaw)
				if err != nil {
					return err
				}
				filtersOut = append(filtersOut, out)
			}
		case "microsoft.percentage":
			{
				var out PercentageFeatureFilter
				mpc := mapstructure.DecoderConfig{TagName: "json", Result: &out}
				mpd, err := mapstructure.NewDecoder(&mpc)
				if err != nil {
					return err
				}
				err = mpd.Decode(filterRaw)
				if err != nil {
					return err
				}
				filtersOut = append(filtersOut, out)
			}

		default:
			return fmt.Errorf("unknown type %q", name)
		}
	}

	p.Filters = filtersOut
	return nil
}

func (p ClientFilter) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Filters)
}

type PercentageFilterParameters struct {
	Value int `json:"Value"`
}

type PercentageFeatureFilter struct {
	Name       string                     `json:"name"`
	Parameters PercentageFilterParameters `json:"parameters"`
}

type TargetingGroupParameter struct {
	Name              string `json:"Name" tfschema:"name"`
	RolloutPercentage int    `json:"RolloutPercentage" tfschema:"rollout_percentage"`
}

type TargetingFilterParameters struct {
	Audience TargetingFilterAudience `json:"Audience"`
}

type TargetingFilterAudience struct {
	DefaultRolloutPercentage int                       `json:"DefaultRolloutPercentage" tfschema:"default_rollout_percentage"`
	Users                    []string                  `json:"Users" tfschema:"users"`
	Groups                   []TargetingGroupParameter `json:"Groups" tfschema:"groups"`
}

type TargetingFeatureFilter struct {
	Name       string                    `json:"name"`
	Parameters TargetingFilterParameters `json:"parameters"`
}

type TimewindowFilterParameters struct {
	Start string `json:"Start" tfschema:"start"`
	End   string `json:"End" tfschema:"end"`
}

type TimewindowFeatureFilter struct {
	Name       string                     `json:"name"`
	Parameters TimewindowFilterParameters `json:"parameters"`
}

type Conditions struct {
	ClientFilters ClientFilter `json:"client_filters"`
}

type FeatureValue struct {
	ID          string     `json:"id"`
	Description string     `json:"description"`
	Enabled     bool       `json:"enabled"`
	Conditions  Conditions `json:"conditions"`
}
