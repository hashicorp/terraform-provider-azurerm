// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package service_plan

import (
	"strings"
)

func PlanSupportsScaleOut(plan string) bool {
	return strings.HasPrefix(plan, "EP") || strings.HasPrefix(plan, "WS")
}
