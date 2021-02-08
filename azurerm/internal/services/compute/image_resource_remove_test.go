package compute_test

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"golang.org/x/crypto/ssh"
)

// TODO: separating this out so this is easier to remove

// nolint unparam
func testGeneralizeVMImage(resourceGroup string, vmName string, userName string, password string, hostName string, port string, location string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		armClient := acceptance.AzureProvider.Meta().(*clients.Client)
		vmClient := armClient.Compute.VMClient
		ctx := armClient.StopContext

		normalizedLocation := azure.NormalizeLocation(location)
		suffix := armClient.Account.Environment.ResourceManagerVMDNSSuffix
		dnsName := fmt.Sprintf("%s.%s.%s", hostName, normalizedLocation, suffix)

		if err := deprovisionVM(userName, password, dnsName, port); err != nil {
			return fmt.Errorf("Bad: Deprovisioning error %+v", err)
		}

		future, err := vmClient.Deallocate(ctx, resourceGroup, vmName)
		if err != nil {
			return fmt.Errorf("Bad: Deallocating error %+v", err)
		}

		if err = future.WaitForCompletionRef(ctx, vmClient.Client); err != nil {
			return fmt.Errorf("Bad: Deallocating error %+v", err)
		}

		if _, err = vmClient.Generalize(ctx, resourceGroup, vmName); err != nil {
			return fmt.Errorf("Bad: Generalizing error %+v", err)
		}

		return nil
	}
}

func deprovisionVM(userName string, password string, hostName string, port string) error {
	// SSH into the machine and execute a waagent deprovisioning command
	var b bytes.Buffer
	cmd := "sudo waagent -verbose -deprovision+user -force"

	config := &ssh.ClientConfig{
		User: userName,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	log.Printf("[INFO] Connecting to %s:%v remote server...", hostName, port)

	hostAddress := strings.Join([]string{hostName, port}, ":")
	client, err := ssh.Dial("tcp", hostAddress, config)
	if err != nil {
		return fmt.Errorf("Bad: deprovisioning error %+v", err)
	}

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("Bad: deprovisioning error, failure creating session %+v", err)
	}
	defer session.Close()

	session.Stdout = &b
	if err := session.Run(cmd); err != nil {
		return fmt.Errorf("Bad: deprovisioning error, failure running command %+v", err)
	}

	return nil
}

// nolint unparam
func testCheckAzureVMExists(sourceVM string, shouldExist bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.VMClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		log.Printf("[INFO] testing MANAGED IMAGE VM EXISTS - BEGIN.")

		vmRs, vmOk := s.RootModule().Resources[sourceVM]
		if !vmOk {
			return fmt.Errorf("VM Not found: %s", sourceVM)
		}
		vmName := vmRs.Primary.Attributes["name"]

		resourceGroup, hasResourceGroup := vmRs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for VM: %s", vmName)
		}

		resp, err := client.Get(ctx, resourceGroup, vmName, "")
		if err != nil {
			return fmt.Errorf("Bad: Get on client: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound && shouldExist {
			return fmt.Errorf("Bad: VM %q (resource group %q) does not exist", vmName, resourceGroup)
		}
		if resp.StatusCode != http.StatusNotFound && !shouldExist {
			return fmt.Errorf("Bad: VM %q (resource group %q) still exists", vmName, resourceGroup)
		}

		log.Printf("[INFO] testing MANAGED IMAGE VM EXISTS - END.")

		return nil
	}
}

func testCheckAzureVMSSExists(sourceVMSS string, shouldExist bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		vmssClient := acceptance.AzureProvider.Meta().(*clients.Client).Compute.VMScaleSetClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		log.Printf("[INFO] testing MANAGED IMAGE VMSS EXISTS - BEGIN.")

		vmRs, vmOk := s.RootModule().Resources[sourceVMSS]
		if !vmOk {
			return fmt.Errorf("VMSS Not found: %s", sourceVMSS)
		}
		vmssName := vmRs.Primary.Attributes["name"]

		resourceGroup, hasResourceGroup := vmRs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for VMSS: %s", vmssName)
		}

		resp, err := vmssClient.Get(ctx, resourceGroup, vmssName)
		if err != nil {
			return fmt.Errorf("Bad: Get on vmssClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound && shouldExist {
			return fmt.Errorf("Bad: VMSS %q (resource group %q) does not exist", vmssName, resourceGroup)
		}
		if resp.StatusCode != http.StatusNotFound && !shouldExist {
			return fmt.Errorf("Bad: VMSS %q (resource group %q) still exists", vmssName, resourceGroup)
		}

		log.Printf("[INFO] testing MANAGED IMAGE VMSS EXISTS - END.")

		return nil
	}
}
