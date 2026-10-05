// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package windows_virtual_machine_scale_set_test

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type WindowsVirtualMachineResource struct{}

func (WindowsVirtualMachineResource) generalizeVirtualMachine(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) error {
	id, err := virtualmachines.ParseVirtualMachineID(state.ID)
	if err != nil {
		return err
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 15*time.Minute)
		defer cancel()
	}

	command := []string{
		"$cmd = \"$Env:SystemRoot\\system32\\sysprep\\sysprep.exe\"",
		"$args = \"/generalize /oobe /mode:vm /quit\"",
		"Start-Process powershell -Argument \"$cmd $args\" -Wait",
	}
	runCommand := virtualmachines.RunCommandInput{
		CommandId: "RunPowerShellScript",
		Script:    &command,
	}

	if err := client.Compute.VirtualMachinesClient.RunCommandThenPoll(ctx, *id, runCommand); err != nil {
		return fmt.Errorf("running sysprep for %s: %+v", *id, err)
	}

	if err := client.Compute.VirtualMachinesClient.DeallocateThenPoll(ctx, *id, virtualmachines.DefaultDeallocateOperationOptions()); err != nil {
		return fmt.Errorf("deallocating %s: %+v", *id, err)
	}

	if _, err = client.Compute.VirtualMachinesClient.Generalize(ctx, *id); err != nil {
		return fmt.Errorf("generalizing %s: %+v", *id, err)
	}

	return nil
}
