// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package dataprotection

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/basebackuppolicyresources"
)

func TestBackupPolicyCosmosdbAccountSchemaWithoutTimeZone(t *testing.T) {
	r := DataProtectionBackupPolicyCosmosdbAccountResource{}
	if _, ok := r.Arguments()["time_zone"]; ok {
		t.Fatal("time_zone must not be exposed as an argument")
	}
	if _, ok := r.Attributes()["time_zone"]; ok {
		t.Fatal("time_zone must not be exposed as an attribute")
	}
}

func TestBackupPolicyCosmosdbAccountBackupRulesWithoutTimeZone(t *testing.T) {
	const fullSchedule = "R/2026-02-08T10:00:00+02:00/P1W"
	incrementalSchedules := []string{
		"R/2026-02-09T10:00:00+02:00/P1W",
		"R/2026-02-10T10:00:00+02:00/P1W",
		"R/2026-02-11T10:00:00+02:00/P1W",
		"R/2026-02-12T10:00:00+02:00/P1W",
		"R/2026-02-13T10:00:00+02:00/P1W",
		"R/2026-02-14T10:00:00+02:00/P1W",
	}

	for _, testCase := range []struct {
		name        string
		dailyBackup bool
		ruleCount   int
	}{
		{name: "full", dailyBackup: false, ruleCount: 1},
		{name: "full_and_incremental", dailyBackup: true, ruleCount: 2},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			rules, err := expandBackupPolicyCosmosdbAccountBackupRules(fullSchedule, testCase.dailyBackup, expandBackupPolicyCosmosdbAccountTaggingCriteria(nil))
			if err != nil {
				t.Fatal(err)
			}
			if len(rules) != testCase.ruleCount {
				t.Fatalf("expected %d rules, got %d", testCase.ruleCount, len(rules))
			}

			for i, item := range rules {
				rule, ok := item.(basebackuppolicyresources.AzureBackupRule)
				if !ok {
					t.Fatalf("unexpected backup rule type %T", item)
				}
				trigger, ok := rule.Trigger.(basebackuppolicyresources.ScheduleBasedTriggerContext)
				if !ok {
					t.Fatalf("unexpected trigger type %T", rule.Trigger)
				}
				if trigger.Schedule.TimeZone != nil {
					t.Fatal("schedule must not specify an API time zone")
				}
				payload, err := json.Marshal(trigger.Schedule)
				if err != nil {
					t.Fatal(err)
				}
				if strings.Contains(string(payload), `"timeZone"`) {
					t.Fatalf("schedule must omit timeZone from its JSON payload: %s", payload)
				}

				expected := []string{fullSchedule}
				if i == 1 {
					expected = incrementalSchedules
				}
				if !reflect.DeepEqual(trigger.Schedule.RepeatingTimeIntervals, expected) {
					t.Fatalf("expected schedules %v, got %v", expected, trigger.Schedule.RepeatingTimeIntervals)
				}
			}

			_, _, actualFullSchedule, actualIncrementalSchedules, dailyBackup := flattenBackupPolicyCosmosdbAccountPolicyRules(rules)
			if actualFullSchedule != fullSchedule {
				t.Fatalf("expected full schedule %q, got %q", fullSchedule, actualFullSchedule)
			}
			if dailyBackup != testCase.dailyBackup {
				t.Fatalf("expected daily backup enabled to be %t, got %t", testCase.dailyBackup, dailyBackup)
			}
			expectedIncrementalSchedules := []string{}
			if testCase.dailyBackup {
				expectedIncrementalSchedules = incrementalSchedules
			}
			if !reflect.DeepEqual(actualIncrementalSchedules, expectedIncrementalSchedules) {
				t.Fatalf("expected incremental schedules %v, got %v", expectedIncrementalSchedules, actualIncrementalSchedules)
			}
		})
	}
}
