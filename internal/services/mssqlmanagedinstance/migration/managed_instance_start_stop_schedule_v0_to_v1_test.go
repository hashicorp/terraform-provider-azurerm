// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"encoding/json"
	"testing"
)

func TestMsSqlManagedInstanceStartStopScheduleV0ToV1(t *testing.T) {
	input := map[string]interface{}{
		"managed_instance_id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Sql/managedInstances/instance1",
		"schedule": []interface{}{
			map[string]interface{}{
				"start_day":  "Wednesday",
				"start_time": "11:00",
				"stop_day":   "Wednesday",
				"stop_time":  "23:00",
			},
		},
		"timezone_id": "UTC",
	}

	actual, err := MsSqlManagedInstanceStartStopScheduleV0ToV1{}.UpgradeFunc()(context.TODO(), input, nil)
	if err != nil {
		t.Fatalf("expected no error but got: %+v", err)
	}

	if _, ok := actual["schedule"].([]interface{}); !ok {
		t.Fatalf("expected schedule to remain a JSON-serializable slice, got %T", actual["schedule"])
	}

	if _, err := json.Marshal(actual); err != nil {
		t.Fatalf("expected upgraded state to be JSON serializable, got: %+v", err)
	}
}
