// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func Timestamp(in time.Time) *timestamppb.Timestamp {
	if in.IsZero() {
		return nil
	}

	return timestamppb.New(in)
}
