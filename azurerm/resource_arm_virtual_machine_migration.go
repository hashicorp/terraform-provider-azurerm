package azurerm

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/terraform"
)

func resourceAzureRMVirtualMachineMigrateState(
	v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	switch v {
	case 0:
		log.Println("[INFO] Found AzureRM Virtual Machine State v0; migrating to v1")
		return migrateAzureRMVirtualMachineStateV0toV1(is)
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func migrateAzureRMVirtualMachineStateV0toV1(is *terraform.InstanceState) (*terraform.InstanceState, error) {
	if is.Empty() {
		log.Println("[DEBUG] Empty InstanceState; nothing to migrate.")
		return is, nil
	}

	log.Printf("[DEBUG] ARM Virtual Machine Attributes before Migration: %#v", is.Attributes)

	if is.Attributes["os_profile_windows_config.#"] == "1" {
		hash, err := computeAzureRMVirtualMachineStateV1WinConfigHash(is, "UTC")
		if err != nil {
			return is, fmt.Errorf("Failed to calculate the new hash code for os_profile_windows_config field: %#v", err)
		}
		replaceAzureRMVirtualMachineStateV0WinConfigWithNewHash(is, hash)
		tzKey := fmt.Sprintf("%s.%d.%s", "os_profile_windows_config", hash, "timezone")
		is.Attributes[tzKey] = "UTC"
	}

	log.Printf("[DEBUG] ARM Virtual Machine Attributes after State Migration: %#v", is.Attributes)

	return is, nil
}

func computeAzureRMVirtualMachineStateV1WinConfigHash(is *terraform.InstanceState, timezone string) (int, error) {
	data := make(map[string]interface{})
	for k, v := range is.Attributes {
		if !strings.HasPrefix(k, "os_profile_windows_config.") {
			continue
		}
		paths := strings.Split(k, ".")
		if len(paths) == 3 && (paths[2] == "enable_automatic_upgrades" || paths[2] == "provision_vm_agent") {
			b, err := strconv.ParseBool(v)
			if err != nil {
				return 0, err
			}
			data[paths[2]] = b
		}
	}
	data["timezone"] = timezone
	return resourceArmVirtualMachineStorageOsProfileWindowsConfigHash(data), nil
}

func replaceAzureRMVirtualMachineStateV0WinConfigWithNewHash(is *terraform.InstanceState, hash int) {
	toDel := make([]string, 0)
	toAdd := make(map[string]string)
	for k, v := range is.Attributes {
		if !strings.HasPrefix(k, "os_profile_windows_config.") {
			continue
		}
		paths := strings.Split(k, ".")
		if len(paths) > 2 {
			toDel = append(toDel, k)
			paths[1] = fmt.Sprintf("%d", hash)
			toAdd[strings.Join(paths, ".")] = v
		}
	}
	for _, k := range toDel {
		delete(is.Attributes, k)
	}
	for k, v := range toAdd {
		is.Attributes[k] = v
	}
}
