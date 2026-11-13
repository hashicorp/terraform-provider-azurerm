// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package datadog_monitor_tag_rule

import "testing"

func TestExpandLogRulesEmptyBlock(t *testing.T) {
	// `log {}` has only optional attributes, so the SDK hands over a nil element
	got := expandLogRules([]any{nil})
	if got == nil {
		t.Fatal("expected log rules for an empty block, got nil")
	}
	if got.SendAadLogs == nil || *got.SendAadLogs || got.SendSubscriptionLogs == nil || *got.SendSubscriptionLogs || got.SendResourceLogs == nil || *got.SendResourceLogs {
		t.Fatalf("expected all log toggles to be false, got %+v", got)
	}
	if got.FilteringTags == nil || len(*got.FilteringTags) != 0 {
		t.Fatalf("expected no filtering tags, got %+v", got.FilteringTags)
	}
}

func TestExpandMetricRulesEmptyBlock(t *testing.T) {
	// `metric {}` has only optional attributes, so the SDK hands over a nil element
	got := expandMetricRules([]any{nil})
	if got == nil {
		t.Fatal("expected metric rules for an empty block, got nil")
	}
	if got.FilteringTags == nil || len(*got.FilteringTags) != 0 {
		t.Fatalf("expected no filtering tags, got %+v", got.FilteringTags)
	}
}
