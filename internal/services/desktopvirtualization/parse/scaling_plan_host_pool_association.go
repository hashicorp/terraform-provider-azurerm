// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2022-02-10-preview/scalingplan"
)

var _ resourceids.Id = ScalingPlanHostPoolAssociationId{}

type ScalingPlanHostPoolAssociationId struct {
	ScalingPlan scalingplan.ScalingPlanId
	HostPool    scalingplan.HostPoolId
}

func (id ScalingPlanHostPoolAssociationId) String() string {
	components := []string{
		fmt.Sprintf("Scaling Plan %s", id.ScalingPlan.String()),
		fmt.Sprintf("Host Pool %s", id.HostPool.String()),
	}
	return fmt.Sprintf("Scaling Plan Host Pool Association %s", strings.Join(components, " / "))
}

func (id ScalingPlanHostPoolAssociationId) ID() string {
	scalingplanId := id.ScalingPlan.ID()
	hostPoolId := id.HostPool.ID()
	return fmt.Sprintf("%s|%s", scalingplanId, hostPoolId)
}

func NewScalingPlanHostPoolAssociationId(scalingplan scalingplan.ScalingPlanId, hostPool scalingplan.HostPoolId) ScalingPlanHostPoolAssociationId {
	return ScalingPlanHostPoolAssociationId{
		ScalingPlan: scalingplan,
		HostPool:    hostPool,
	}
}

func ScalingPlanHostPoolAssociationID(input string) (*ScalingPlanHostPoolAssociationId, error) {
	segments := strings.Split(input, "|")
	if len(segments) != 2 {
		return nil, fmt.Errorf("expected an ID in the format {scalingplanID}|{hostPoolID} but got %q", input)
	}

	scalingplanId, err := scalingplan.ParseScalingPlanID(segments[0])
	if err != nil {
		return nil, fmt.Errorf("parsing Scaling Plan ID for Scaling Plan/Host Pool Association %q: %+v", segments[0], err)
	}

	hostPoolId, err := scalingplan.ParseHostPoolID(segments[1])
	if err != nil {
		return nil, fmt.Errorf("parsing Host Pool ID for Scaling Plan/Host Pool Association %q: %+v", segments[1], err)
	}

	return &ScalingPlanHostPoolAssociationId{
		ScalingPlan: *scalingplanId,
		HostPool:    *hostPoolId,
	}, nil
}
