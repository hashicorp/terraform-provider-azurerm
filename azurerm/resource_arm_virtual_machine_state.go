package azurerm

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2017-12-01/compute"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmVirtualMachineState() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmVirtualMachineStateCreate,
		Read:   resourceArmVirtualMachineStateRead,
		Update: resourceArmVirtualMachineStateCreate,
		Delete: resourceArmVirtualMachineStateRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"resource_group_name": resourceGroupNameSchema(),

			"virtual_machine_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"state": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"started",
					"poweredoff",
					"deallocated",
				}, true),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},
		},
	}
}

func resourceArmVirtualMachineStateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).vmClient
	ctx := meta.(*ArmClient).StopContext

	vmName := d.Get("virtual_machine_name").(string)
	resGroup := d.Get("resource_group_name").(string)
	vmState := d.Get("state").(string)

	var err error

	if vmState == "deallocated" {
		err = resourceArmVirtualMachineStateDeallocate(resGroup, vmName, vmState, meta)
	} else if vmState == "poweredoff" {
		err = resourceArmVirtualMachineStatePowerOff(resGroup, vmName, vmState, meta)
	} else if vmState == "started" {
		err = resourceArmVirtualMachineStateStart(resGroup, vmName, vmState, meta)
	}

	if err != nil {
		return err
	}

	vm, err := client.Get(ctx, resGroup, vmName, compute.InstanceView)
	if err != nil {
		return err
	}

	if vm.InstanceView != nil {
		for _, statGet := range *vm.InstanceView.Statuses {
			fmt.Println(statGet.Level, *statGet.DisplayStatus)
		}
	} else {
		return fmt.Errorf("Cannot read Virtual Machine State %s (resource group %s) ID", vmName, resGroup)
	}

	d.SetId(*vm.ID)

	return resourceArmVirtualMachineStateRead(d, meta)
}

func resourceArmVirtualMachineStateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).vmClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	vmName := id.Path["virtualMachines"]

	//vmName := d.Get("virtual_machine_name").(string)
	//resGroup := d.Get("resource_group_name").(string)

	//resp, err := client.InstanceView(ctx, resGroup, vmName)

	vm, err := client.Get(ctx, resGroup, vmName, compute.InstanceView)
	if err != nil {
		if utils.ResponseWasNotFound(vm.Response) {
			d.SetId("")
			return nil
		}
	}

	if vm.InstanceView != nil {
		for _, statGet := range *vm.InstanceView.Statuses {
			//d.Set("state", *statGet.Code)
			//d.Set("state", strings.Split("/", *statGet.Code))
			statusSplit := strings.Split(*statGet.Code, "/")
			d.Set("state", statusSplit[1])
		}
	} else {
		return fmt.Errorf("Cannot read Virtual Machine State %s (resource group %s) ID", vmName, resGroup)
	}

	d.Set("virtual_machine_name", vmName)
	d.Set("resource_group_name", resGroup)

	return nil
}

func resourceArmVirtualMachineStateDeallocate(resGroup string, vmName string, vmState string, meta interface{}) error {
	client := meta.(*ArmClient).vmClient
	ctx := meta.(*ArmClient).StopContext

	future, err := client.Deallocate(ctx, resGroup, vmName)

	if err != nil {
		return err
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return err
	}

	return nil
}

func resourceArmVirtualMachineStateStart(resGroup string, vmName string, vmState string, meta interface{}) error {
	client := meta.(*ArmClient).vmClient
	ctx := meta.(*ArmClient).StopContext

	future, err := client.Start(ctx, resGroup, vmName)

	if err != nil {
		return err
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return err
	}

	return nil
}

func resourceArmVirtualMachineStatePowerOff(resGroup string, vmName string, vmState string, meta interface{}) error {
	client := meta.(*ArmClient).vmClient
	ctx := meta.(*ArmClient).StopContext

	future, err := client.PowerOff(ctx, resGroup, vmName)

	if err != nil {
		return err
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return err
	}

	return nil
}
