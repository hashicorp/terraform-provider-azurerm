// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_policy_elastic_san_volume

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/basebackuppolicyresources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func TestBackupPolicyElasticSanVolumeSchema_NoTimeZone(t *testing.T) {
	if _, ok := (Resource{}).Arguments()["time_zone"]; ok {
		t.Fatal("expected `time_zone` to be absent from the schema")
	}
}

func TestBackupPolicyElasticSanVolumePolicyRules_RoundTrip(t *testing.T) {
	tests := map[string]BackupPolicyElasticSanVolumeModel{
		"daily": {
			DailyBackupSchedule:      "R/2026-02-08T13:00:00Z/P1D",
			DefaultRetentionDuration: "P7D",
		},
		"weekly": {
			WeeklyBackupSchedule:     "R/2026-02-08T13:00:00Z/P1W",
			DefaultRetentionDuration: "P7D",
		},
		"retention rules with absolute criteria": {
			DailyBackupSchedule:      "R/2026-02-08T13:00:00Z/P1D",
			DefaultRetentionDuration: "P7D",
			DailyRetentionDuration:   "P14D",
			MonthlyRetentionDuration: "P3M",
			WeeklyRetentionDuration:  "P28D",
		},
		"weekly retention only": {
			DailyBackupSchedule:      "R/2026-02-08T13:00:00Z/P1D",
			DefaultRetentionDuration: "P7D",
			WeeklyRetentionDuration:  "P28D",
		},
		"monthly retention only": {
			WeeklyBackupSchedule:     "R/2026-02-08T13:00:00Z/P1W",
			DefaultRetentionDuration: "P7D",
			MonthlyRetentionDuration: "P3M",
		},
		"daily retention only": {
			DailyBackupSchedule:      "R/2026-02-08T13:00:00Z/P1D",
			DefaultRetentionDuration: "P7D",
			DailyRetentionDuration:   "P14D",
		},
	}

	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			policy := basebackuppolicyresources.BackupPolicy{
				DatasourceTypes: []string{elasticSanVolumeDatasourceType},
				PolicyRules:     expandBackupPolicyElasticSanVolumePolicyRules(input),
			}
			encoded, err := json.Marshal(policy)
			if err != nil {
				t.Fatal(err)
			}
			var decoded basebackuppolicyresources.BackupPolicy
			if err := json.Unmarshal(encoded, &decoded); err != nil {
				t.Fatal(err)
			}
			if !backupPolicySupportsElasticSanVolumes(decoded) {
				t.Fatal("expected an Elastic SAN policy")
			}

			var actual BackupPolicyElasticSanVolumeModel
			if err := flattenBackupPolicyElasticSanVolumePolicyRules(decoded.PolicyRules, &actual); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(input, actual) {
				t.Fatalf("expected %+v, got %+v", input, actual)
			}

			for _, item := range decoded.PolicyRules {
				switch rule := item.(type) {
				case basebackuppolicyresources.AzureBackupRule:
					if rule.DataStore.DataStoreType != basebackuppolicyresources.DataStoreTypesOperationalStore {
						t.Fatalf("unexpected backup datastore: %s", rule.DataStore.DataStoreType)
					}
					parameters, ok := rule.BackupParameters.(basebackuppolicyresources.AzureBackupParams)
					if !ok || parameters.BackupType != "Incremental" {
						t.Fatalf("unexpected backup parameters: %+v", rule.BackupParameters)
					}
					trigger, ok := rule.Trigger.(basebackuppolicyresources.ScheduleBasedTriggerContext)
					if !ok || len(trigger.TaggingCriteria) == 0 {
						t.Fatalf("unexpected trigger: %+v", rule.Trigger)
					}
					expectedSchedule := input.DailyBackupSchedule
					if expectedSchedule == "" {
						expectedSchedule = input.WeeklyBackupSchedule
					}
					if !reflect.DeepEqual(trigger.Schedule.RepeatingTimeIntervals, []string{expectedSchedule}) {
						t.Fatalf("unexpected backup schedule: %+v", trigger.Schedule.RepeatingTimeIntervals)
					}
					if trigger.Schedule.TimeZone != nil {
						t.Fatalf("expected no schedule time zone, got %s", *trigger.Schedule.TimeZone)
					}
					if !trigger.TaggingCriteria[0].IsDefault || trigger.TaggingCriteria[0].TaggingPriority != 99 {
						t.Fatalf("unexpected default tagging criteria: %+v", trigger.TaggingCriteria[0])
					}
					for _, expected := range []struct {
						name     string
						duration string
						marker   basebackuppolicyresources.AbsoluteMarker
						priority int64
					}{
						{"Daily", input.DailyRetentionDuration, basebackuppolicyresources.AbsoluteMarkerFirstOfDay, 3},
						{"Weekly", input.WeeklyRetentionDuration, basebackuppolicyresources.AbsoluteMarkerFirstOfWeek, 2},
						{"Monthly", input.MonthlyRetentionDuration, basebackuppolicyresources.AbsoluteMarkerFirstOfMonth, 1},
					} {
						tagFound := false
						for _, tagging := range trigger.TaggingCriteria {
							if tagging.TagInfo.TagName != expected.name {
								continue
							}
							tagFound = true
							criteria := pointer.From(tagging.Criteria)
							if len(criteria) != 1 {
								t.Fatalf("unexpected %s tagging criteria: %+v", expected.name, criteria)
							}
							scheduleCriteria, ok := criteria[0].(basebackuppolicyresources.ScheduleBasedBackupCriteria)
							if !ok || !reflect.DeepEqual(pointer.From(scheduleCriteria.AbsoluteCriteria), []basebackuppolicyresources.AbsoluteMarker{expected.marker}) {
								t.Fatalf("expected %s criteria, got %+v", expected.marker, criteria[0])
							}
							if tagging.TaggingPriority != expected.priority {
								t.Fatalf("expected %s priority %d, got %d", expected.name, expected.priority, tagging.TaggingPriority)
							}
						}
						if tagFound != (expected.duration != "") {
							t.Fatalf("unexpected %s tagging presence: %t", expected.name, tagFound)
						}
					}
				case basebackuppolicyresources.AzureRetentionRule:
					if rule.Lifecycles[0].SourceDataStore.DataStoreType != basebackuppolicyresources.DataStoreTypesOperationalStore {
						t.Fatalf("unexpected retention datastore: %s", rule.Lifecycles[0].SourceDataStore.DataStoreType)
					}
				}
			}
		})
	}
}

func TestBackupPolicySupportsElasticSanVolumes(t *testing.T) {
	tests := []struct {
		name     string
		input    basebackuppolicyresources.BaseBackupPolicy
		expected bool
	}{
		{name: "nil"},
		{name: "missing datasource", input: basebackuppolicyresources.BackupPolicy{}},
		{
			name: "cosmosdb",
			input: basebackuppolicyresources.BackupPolicy{
				DatasourceTypes: []string{"Microsoft.DocumentDB/databaseAccounts"},
			},
		},
		{
			name: "elastic san",
			input: basebackuppolicyresources.BackupPolicy{
				DatasourceTypes: []string{elasticSanVolumeDatasourceType},
			},
			expected: true,
		},
		{
			name: "case insensitive",
			input: basebackuppolicyresources.BackupPolicy{
				DatasourceTypes: []string{strings.ToUpper(elasticSanVolumeDatasourceType)},
			},
			expected: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if actual := backupPolicySupportsElasticSanVolumes(test.input); actual != test.expected {
				t.Fatalf("expected %t, got %t", test.expected, actual)
			}
		})
	}
}

func TestBackupPolicyElasticSanVolumePolicyRules_InvalidShape(t *testing.T) {
	tests := map[string][]basebackuppolicyresources.BasePolicyRule{
		"missing lifecycle": {
			basebackuppolicyresources.AzureRetentionRule{Name: "Default", IsDefault: pointer.To(true)},
		},
		"missing delete option": {
			basebackuppolicyresources.AzureRetentionRule{
				Name:       "Default",
				Lifecycles: []basebackuppolicyresources.SourceLifeCycle{{}},
			},
		},
		"unsupported retention": {
			expandBackupPolicyElasticSanVolumeRetentionRule("Yearly", "P1Y", false),
		},
		"duplicate weekly retention": {
			basebackuppolicyresources.AzureBackupRule{
				Name: "Backup",
				Trigger: basebackuppolicyresources.ScheduleBasedTriggerContext{
					Schedule: basebackuppolicyresources.BackupSchedule{
						RepeatingTimeIntervals: []string{"R/2026-02-08T13:00:00Z/P1D"},
					},
				},
			},
			expandBackupPolicyElasticSanVolumeRetentionRule("Weekly", "P7D", false),
			expandBackupPolicyElasticSanVolumeRetentionRule("Weekly", "P14D", false),
		},
		"unsupported trigger": {
			basebackuppolicyresources.AzureBackupRule{Name: "Backup"},
		},
		"custom weekly weekdays": {
			basebackuppolicyresources.AzureBackupRule{
				Name: "Backup",
				Trigger: basebackuppolicyresources.ScheduleBasedTriggerContext{
					Schedule: basebackuppolicyresources.BackupSchedule{
						RepeatingTimeIntervals: []string{"R/2026-02-08T13:00:00Z/P1D"},
					},
					TaggingCriteria: []basebackuppolicyresources.TaggingCriteria{
						expandBackupPolicyElasticSanVolumeTaggingCriteria("Weekly", 1, basebackuppolicyresources.ScheduleBasedBackupCriteria{
							DaysOfTheWeek: pointer.To([]basebackuppolicyresources.DayOfWeek{basebackuppolicyresources.DayOfWeekMonday}),
						}),
					},
				},
			},
			expandBackupPolicyElasticSanVolumeRetentionRule("Weekly", "P28D", false),
		},
		"missing backup rule": {},
		"missing schedule": {
			basebackuppolicyresources.AzureBackupRule{
				Name:    "Backup",
				Trigger: basebackuppolicyresources.ScheduleBasedTriggerContext{},
			},
		},
		"unsupported schedule cadence": {
			basebackuppolicyresources.AzureBackupRule{
				Name: "Backup",
				Trigger: basebackuppolicyresources.ScheduleBasedTriggerContext{
					Schedule: basebackuppolicyresources.BackupSchedule{
						RepeatingTimeIntervals: []string{"R/2026-02-08T13:00:00Z/PT1H"},
					},
				},
			},
		},
		"multiple schedules": {
			basebackuppolicyresources.AzureBackupRule{
				Name: "Backup",
				Trigger: basebackuppolicyresources.ScheduleBasedTriggerContext{
					Schedule: basebackuppolicyresources.BackupSchedule{
						RepeatingTimeIntervals: []string{
							"R/2026-02-08T13:00:00Z/P1D",
							"R/2026-02-08T14:00:00Z/P1D",
						},
					},
				},
			},
		},
		"multiple backup rules": {
			basebackuppolicyresources.AzureBackupRule{
				Name: "First",
				Trigger: basebackuppolicyresources.ScheduleBasedTriggerContext{
					Schedule: basebackuppolicyresources.BackupSchedule{
						RepeatingTimeIntervals: []string{"R/2026-02-08T13:00:00Z/P1D"},
					},
				},
			},
			basebackuppolicyresources.AzureBackupRule{
				Name: "Second",
				Trigger: basebackuppolicyresources.ScheduleBasedTriggerContext{
					Schedule: basebackuppolicyresources.BackupSchedule{
						RepeatingTimeIntervals: []string{"R/2026-02-08T14:00:00Z/P1D"},
					},
				},
			},
		},
	}
	for name, rules := range tests {
		t.Run(name, func(t *testing.T) {
			var state BackupPolicyElasticSanVolumeModel
			if err := flattenBackupPolicyElasticSanVolumePolicyRules(rules, &state); err == nil {
				t.Fatal("expected an error for an unsupported policy shape")
			}
		})
	}
}

func TestBackupPolicyElasticSanVolumeRetentionDurationValidation(t *testing.T) {
	arguments := Resource{}.Arguments()
	for _, name := range []string{"weekly_retention_rule"} {
		if _, ok := arguments[name]; ok {
			t.Fatalf("unexpected legacy `%s` argument", name)
		}
	}
	for _, name := range []string{"daily_retention_duration", "weekly_retention_duration", "monthly_retention_duration"} {
		t.Run(name, func(t *testing.T) {
			duration := arguments[name]
			if duration == nil || duration.Type != pluginsdk.TypeString || !duration.Optional || !duration.ForceNew {
				t.Fatalf("expected `%s` to be an optional, ForceNew string", name)
			}
			for input, expected := range map[string]bool{"P28D": true, "P3M": true, "invalid": false} {
				t.Run(input, func(t *testing.T) {
					_, errors := duration.ValidateFunc(input, name)
					if actual := len(errors) == 0; actual != expected {
						t.Fatalf("expected validation result %t, got %t: %+v", expected, actual, errors)
					}
				})
			}
		})
	}
}

func TestBackupPolicyElasticSanVolumeScheduleValidation(t *testing.T) {
	arguments := Resource{}.Arguments()
	for _, name := range []string{"backup_schedule", "backup_repeating_time_intervals"} {
		if _, ok := arguments[name]; ok {
			t.Fatalf("unexpected legacy `%s` argument", name)
		}
	}
	for _, name := range []string{"daily_backup_schedule", "weekly_backup_schedule"} {
		t.Run(name, func(t *testing.T) {
			schedule := arguments[name]
			if schedule == nil || schedule.Type != pluginsdk.TypeString || !schedule.Optional || !schedule.ForceNew {
				t.Fatalf("expected `%s` to be an optional, ForceNew string", name)
			}
			if !reflect.DeepEqual(schedule.ExactlyOneOf, []string{"daily_backup_schedule", "weekly_backup_schedule"}) {
				t.Fatalf("unexpected ExactlyOneOf constraint: %+v", schedule.ExactlyOneOf)
			}
			tests := map[string]bool{
				"R/2026-02-08T13:00:00Z/P1D":  name == "daily_backup_schedule",
				"R/2026-02-08T13:00:00Z/P1W":  name == "weekly_backup_schedule",
				"R/2026-02-08T13:00:00Z/PT1H": false,
				"R/2026-02-08T13:00:00Z/P1M":  false,
				"invalid/P1D":                 false,
				"invalid/P1W":                 false,
			}
			for input, expected := range tests {
				t.Run(input, func(t *testing.T) {
					_, errors := schedule.ValidateFunc(input, name)
					if actual := len(errors) == 0; actual != expected {
						t.Fatalf("expected validation result %t, got %t: %+v", expected, actual, errors)
					}
				})
			}
		})
	}
}

func TestBackupPolicyElasticSanVolumeScheduleExactlyOneOf(t *testing.T) {
	resource := &pluginsdk.Resource{Schema: Resource{}.Arguments()}
	tests := []struct {
		name        string
		daily       string
		weekly      string
		expectError bool
	}{
		{name: "neither", expectError: true},
		{name: "daily", daily: "R/2026-02-08T13:00:00Z/P1D"},
		{name: "weekly", weekly: "R/2026-02-08T13:00:00Z/P1W"},
		{
			name:        "both",
			daily:       "R/2026-02-08T13:00:00Z/P1D",
			weekly:      "R/2026-02-08T13:00:00Z/P1W",
			expectError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := map[string]interface{}{
				"name":                            "test-policy",
				"data_protection_backup_vault_id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test/providers/Microsoft.DataProtection/backupVaults/test",
				"default_retention_duration":      "P7D",
			}
			if test.daily != "" {
				config["daily_backup_schedule"] = test.daily
			}
			if test.weekly != "" {
				config["weekly_backup_schedule"] = test.weekly
			}
			diags := resource.Validate(terraform.NewResourceConfigRaw(config))
			if diags.HasError() != test.expectError {
				t.Fatalf("expected validation error %t, got %+v", test.expectError, diags)
			}
		})
	}
}
