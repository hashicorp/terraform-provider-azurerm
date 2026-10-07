// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2024-07-10/machines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/configurationassignments"
)

func TestMaintenanceAssignmentArcMachineID(t *testing.T) {
	machineId := machines.NewMachineID("00000000-0000-0000-0000-000000000000", "ExampleRG", "ExampleMachine")
	id := NewMaintenanceAssignmentArcMachineID(machineId.SubscriptionId, machineId.ResourceGroupName, machineId.MachineName, "ExampleAssignment")
	sdkId := configurationassignments.NewScopedConfigurationAssignmentID(machineId.ID(), "ExampleAssignment")
	if id.ID() != sdkId.ID() {
		t.Fatalf("expected the SDK ARM ID %q, got %q", sdkId.ID(), id.ID())
	}

	parsed, err := MaintenanceAssignmentArcMachineID(id.ID())
	if err != nil {
		t.Fatal(err)
	}
	if *parsed != id {
		t.Fatalf("expected %#v, got %#v", id, *parsed)
	}
}

func TestMaintenanceAssignmentArcMachineIDInsensitively(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "lowercase service response",
			input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/example/providers/microsoft.hybridcompute/machines/machine/providers/microsoft.maintenance/configurationassignments/assignment",
			want:  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.HybridCompute/machines/machine/providers/Microsoft.Maintenance/configurationAssignments/assignment",
		},
		{
			name:  "preserve resource name casing",
			input: "/subscriptions/00000000-0000-0000-0000-000000000000/RESOURCEGROUPS/ExampleRG/providers/MICROSOFT.HYBRIDCOMPUTE/MACHINES/ExampleMachine/providers/MICROSOFT.MAINTENANCE/CONFIGURATIONASSIGNMENTS/ExampleAssignment",
			want:  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/ExampleRG/providers/Microsoft.HybridCompute/machines/ExampleMachine/providers/Microsoft.Maintenance/configurationAssignments/ExampleAssignment",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			id, err := MaintenanceAssignmentArcMachineIDInsensitively(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if got := id.ID(); got != tc.want {
				t.Fatalf("expected %q, got %q", tc.want, got)
			}
		})
	}
}

func TestMaintenanceAssignmentArcMachineIDRejectsOtherScopes(t *testing.T) {
	for _, scope := range []string{
		"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Compute/virtualMachines/machine",
		"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Compute/hostGroups/group/hosts/host",
		"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example",
	} {
		t.Run(scope, func(t *testing.T) {
			id := configurationassignments.NewScopedConfigurationAssignmentID(scope, "assignment")
			if _, err := MaintenanceAssignmentArcMachineID(id.ID()); err == nil {
				t.Fatalf("expected to reject non-Arc scope %q", scope)
			}
			if _, err := MaintenanceAssignmentArcMachineIDInsensitively(strings.ToLower(id.ID())); err == nil {
				t.Fatalf("expected to reject lowercase non-Arc scope %q", scope)
			}
		})
	}
}
