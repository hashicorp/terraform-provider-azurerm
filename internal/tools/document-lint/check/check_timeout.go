// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package check

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/model"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/util"
)

type TimeoutType int

const (
	TimeoutMissed TimeoutType = iota // no timeout part in document
	TimeoutCreate
	TimeoutRead
	TimeoutUpdate
	TimeoutDelete
)

func (t TimeoutType) String() string {
	return []string{"", "create", "read", "update", "delete"}[t]
}

func (t TimeoutType) IngString() string {
	return []string{"", "creating", "retrieving", "updating", "deleting"}[t]
}

// GenLine generate a line for timeout if not exists
func (t *TimeoutDiffItem) GenLine(rt string) string {
	return fmt.Sprintf("* `%s` - (Defaults to %s) Used when %s the %s.",
		t.Type.String(),
		t.ValueString(),
		t.Type.IngString(),
		rt,
	)
}

type TimeoutDiffItem struct {
	Line int         // which line
	Type TimeoutType //
	Want int64       // correct timeout value in second
}

func (t *TimeoutDiffItem) ValueString() string {
	val, suf := t.Want/60, "minute"
	if t.Want > 60*60 && (t.Want%(60*60)) == 0 {
		val, suf = t.Want/(60*60), "hour"
	}
	if val > 1 {
		suf += "s" // add a 's' suffix
	}
	return fmt.Sprintf("%d %s", val, suf)
}

func (t *TimeoutDiffItem) FixLine(line string) string {
	if t.Want <= 0 {
		return line
	}
	// find place to replace
	start, end := util.TimeoutValueIdx(line)
	if end <= start {
		return line
	}
	res := fmt.Sprintf("%s%s%s", line[:start], t.ValueString(), line[end:])
	return res
}

func NewTimeoutDiffItem(line int, typ TimeoutType, want int64) TimeoutDiffItem {
	return TimeoutDiffItem{
		Line: line,
		Type: typ,
		Want: want,
	}
}

type timeoutDiff struct {
	checkBase
	TimeoutDiff []TimeoutDiffItem
}

func newTimeoutDiff(checkBase checkBase, items []TimeoutDiffItem) *timeoutDiff {
	return &timeoutDiff{checkBase: checkBase, TimeoutDiff: items}
}

func (t timeoutDiff) String() string {
	var bs strings.Builder
	bs.WriteString(fmt.Sprintf("%d Document has incorrect or missing Timeouts. These should be:", t.checkBase.Line()))
	for _, item := range t.TimeoutDiff {
		bs.WriteString(fmt.Sprintf("\n* `%s` - (Defaults to %s) Used when %s the xxx", item.Type, item.ValueString(), item.Type.IngString()))
	}
	return bs.String()
}

func (t timeoutDiff) Fix(line string) (result string, err error) {
	// cannot fix timeout by line
	return line, nil
}

var _ Checker = (*timeoutDiff)(nil)

func diffTimeout(r *schema.Resource, md *model.ResourceDoc) (res []Checker) {
	to := r.Schema.Timeouts
	if to == nil {
		return
	}
	var items []TimeoutDiffItem
	if md.Timeouts == nil {
		items = append(items, NewTimeoutDiffItem(0, TimeoutMissed, 0))
		md.Timeouts = &model.Timeouts{} // use an empty timeouts object
	}

	if ptr := to.Read; ptr != nil {
		val := int64((*ptr) / time.Second)
		if mdVal := md.Timeouts.Read; val != mdVal.Value {
			items = append(items, NewTimeoutDiffItem(mdVal.Line, TimeoutRead, val))
		}
	}

	if ptr := to.Create; ptr != nil {
		val := int64((*ptr) / time.Second)
		if mdVal := md.Timeouts.Create; val != mdVal.Value {
			items = append(items, NewTimeoutDiffItem(mdVal.Line, TimeoutCreate, val))
		}
	}

	if ptr := to.Update; ptr != nil {
		val := int64((*ptr) / time.Second)
		if mdVal := md.Timeouts.Update; val != mdVal.Value {
			items = append(items, NewTimeoutDiffItem(mdVal.Line, TimeoutUpdate, val))
		}
	}
	if ptr := to.Delete; ptr != nil {
		val := int64((*ptr) / time.Second)
		if mdVal := md.Timeouts.Delete; val != mdVal.Value {
			items = append(items, NewTimeoutDiffItem(mdVal.Line, TimeoutDelete, val))
		}
	}

	if len(items) > 0 {
		res = append(res, newTimeoutDiff(newCheckBase(md.Timeouts.Read.Line, "", nil), items))
	}
	return res
}
