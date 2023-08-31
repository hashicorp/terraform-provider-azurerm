// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package suppress

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func RFC3339Time(_, old, new string, _ *schema.ResourceData) bool {
	ot, oerr := time.Parse(time.RFC3339, old)
	nt, nerr := time.Parse(time.RFC3339, new)

	if oerr != nil || nerr != nil {
		return false
	}

	return nt.Equal(ot)
}

func RFC3339MinuteTime(_, old, new string, _ *schema.ResourceData) bool {
	ot, oerr := time.Parse(time.RFC3339, old)
	nt, nerr := time.Parse(time.RFC3339, new)

	if oerr != nil || nerr != nil {
		return false
	}

	return nt.Unix()-int64(nt.Second()) == ot.Unix()-int64(ot.Second())
}
