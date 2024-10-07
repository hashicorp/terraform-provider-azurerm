// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package check

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/md"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/schema"
)

// function to call diff all resource provided by expose_schema.go

type DiffResult struct {
	start  time.Time
	end    time.Time
	result []*ResourceDiff
	byCate map[string][]*ResourceDiff
	mux    sync.Mutex
}

func NewDiffResult() *DiffResult {
	return &DiffResult{
		byCate: map[string][]*ResourceDiff{},
		start:  time.Now(),
	}
}

func (d *DiffResult) FixDocuments() (err error) {
	for _, r := range d.result {
		fix := NewFixer(r)
		if err = fix.TryFix(); err != nil {
			return fmt.Errorf("when try fix document: %v", err)
		}
		if err = fix.WriteBack(); err != nil {
			return fmt.Errorf("when write back: %v", err)
		}
	}
	return nil
}

func (d *DiffResult) GetResult() []*ResourceDiff {
	return d.result
}

func (d *DiffResult) HasDiff() bool {
	result := d.GetResult()
	for _, r := range result {
		for _, f := range r.Diff {
			if !f.ShouldSkip() {
				return true
			}
		}
	}
	return false
}

func (d *DiffResult) ToString() string {
	var bs strings.Builder
	var count int
	var possiblevalueMiss int
	var crossCount, resourceCount int
	var reqCount, defaultCount, timeoutCount, forceNewCount int
	var skipCount int
	for _, diff := range d.result {
		if len(diff.Diffs()) > 0 {
			resourceCount++
			count += len(diff.Diffs())
			bs.WriteString(diff.ToString())
			bs.WriteString("\n")

			for _, df := range diff.Diffs() {
				if df.ShouldSkip() {
					skipCount++
					continue
				}
				switch df.(type) {
				case propertyMissDiff:
					crossCount++
				case possibleValueDiff:
					possiblevalueMiss++
				case requireDiff:
					reqCount++
				case defaultDiff:
					defaultCount++
				case timeoutDiff:
					timeoutCount++
				case forceNewDiff:
					forceNewCount++

				}
			}
		}
	}
	fg := color.FgGreen
	if count > 0 {
		fg = color.FgYellow
	}
	bs.WriteString(
		color.New(color.Bold, fg).Sprintf(
			`------
%d issues found in %d resources
------`,
			count, resourceCount))
	return bs.String()
}

func (d *DiffResult) CostTime() time.Duration {
	return d.end.Sub(d.start)
}

func DiffAll(regs Resources, dryRun bool) *DiffResult {
	var dr = NewDiffResult()

	// can not split to package in different goroutine which may cause data-race and mix shared pointer up
	// register may repeat in typed and untyped, so use a map to remove the repeat entry
	for _, res := range regs.resources {
		// one register exists in both typed register and untyped register
		var rds []*ResourceDiff
		var catName string

		sch := schema.NewResource(res.schema, res.name)
		rd := NewResourceDiff(sch)
		if !dryRun {
			md.FixFileNormalize(rd.MDFile)
		}
		rd.DiffAll()

		if len(rd.Diffs()) > 0 {
			rds = append(rds, rd)
		}

		if len(rds) > 0 {
			dr.mux.Lock()
			dr.result = append(dr.result, rds...)
			dr.byCate[catName] = rds
			dr.mux.Unlock()
		}
	}
	dr.end = time.Now()
	return dr
}
