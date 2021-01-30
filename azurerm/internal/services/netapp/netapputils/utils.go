package netapputils

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

const (
	netAppResourceProviderName string = "Microsoft.NetApp"
)

// GetResourceValue returns the name of a resource from resource id/uri based on resource type name.
func GetResourceValue(resourceURI string, resourceName string) string {

	if len(strings.TrimSpace(resourceURI)) == 0 {
		return ""
	}

	if len(strings.TrimSpace(resourceName)) == 0 {
		return ""
	}

	if !strings.HasPrefix(resourceURI, "/") {
		resourceURI = fmt.Sprintf("/%v", resourceURI)
	}

	if !strings.HasPrefix(resourceName, "/") {
		resourceName = fmt.Sprintf("/%v", resourceName)
	}

	// Checks to see if the ResourceName and ResourceGroup is the same name and if so handles it specially.
	rgResourceName := fmt.Sprintf("/resourceGroups%v", resourceName)
	rgIndex := strings.Index(strings.ToLower(resourceURI), strings.ToLower(rgResourceName))

	// Dealing with case where resource name is the same as resource group
	if rgIndex > -1 {
		removedSameRgName := strings.Split(strings.ToLower(resourceURI), strings.ToLower(resourceName))
		return strings.Split(removedSameRgName[len(removedSameRgName)-1], "/")[1]
	}

	// Dealing with regular cases
	index := strings.Index(strings.ToLower(resourceURI), strings.ToLower(resourceName))
	if index > -1 {
		resource := strings.Split(resourceURI[index+len(resourceName):], "/")
		if len(resource) > 1 {
			return resource[1]
		}
	}

	return ""
}

// GetResourceName gets the resource name from resource id/uri
func GetResourceName(resourceURI string) string {

	if len(strings.TrimSpace(resourceURI)) == 0 {
		return ""
	}

	position := strings.LastIndex(resourceURI, "/")
	return resourceURI[position+1:]
}

// GetSubscription gets he subscription id from resource id/uri
func GetSubscription(resourceURI string) string {

	if len(strings.TrimSpace(resourceURI)) == 0 {
		return ""
	}

	subscriptionID := GetResourceValue(resourceURI, "/subscriptions")
	if subscriptionID == "" {
		return ""
	}

	return subscriptionID
}

// GetResourceGroup gets the resource group name from resource id/uri
func GetResourceGroup(resourceURI string) string {

	if len(strings.TrimSpace(resourceURI)) == 0 {
		return ""
	}

	resourceGroupName := GetResourceValue(resourceURI, "/resourceGroups")
	if resourceGroupName == "" {
		return ""
	}

	return resourceGroupName
}

// GetAnfAccount gets an account name from resource id/uri
func GetAnfAccount(resourceURI string) string {

	if len(strings.TrimSpace(resourceURI)) == 0 {
		return ""
	}

	accountName := GetResourceValue(resourceURI, "/netAppAccounts")
	if accountName == "" {
		return ""
	}

	return accountName
}

// GetAnfCapacityPool gets pool name from resource id/uri
func GetAnfCapacityPool(resourceURI string) string {

	if len(strings.TrimSpace(resourceURI)) == 0 {
		return ""
	}

	accountName := GetResourceValue(resourceURI, "/capacityPools")
	if accountName == "" {
		return ""
	}

	return accountName
}

// GetAnfVolume gets volume name from resource id/uri
func GetAnfVolume(resourceURI string) string {

	if len(strings.TrimSpace(resourceURI)) == 0 {
		return ""
	}

	volumeName := GetResourceValue(resourceURI, "/volumes")
	if volumeName == "" {
		return ""
	}

	return volumeName
}

// GetAnfSnapshot gets snapshot name from resource id/uri
func GetAnfSnapshot(resourceURI string) string {

	if len(strings.TrimSpace(resourceURI)) == 0 {
		return ""
	}

	snapshotName := GetResourceValue(resourceURI, "/snapshots")
	if snapshotName == "" {
		return ""
	}

	return snapshotName
}

// IsAnfResource checks if resource is an ANF related resource
func IsAnfResource(resourceURI string) bool {

	if len(strings.TrimSpace(resourceURI)) == 0 {
		return false
	}

	return strings.Index(resourceURI, netAppResourceProviderName) > -1
}

// IsAnfSnapshot checks resource is a snapshot
func IsAnfSnapshot(resourceURI string) bool {

	if len(strings.TrimSpace(resourceURI)) == 0 || !IsAnfResource(resourceURI) {
		return false
	}

	return strings.LastIndex(resourceURI, "/snapshots/") > -1
}

// IsAnfVolume checks resource is a volume
func IsAnfVolume(resourceURI string) bool {

	if len(strings.TrimSpace(resourceURI)) == 0 || !IsAnfResource(resourceURI) {
		return false
	}

	return !IsAnfSnapshot(resourceURI) &&
		strings.LastIndex(resourceURI, "/volumes/") > -1
}

// IsAnfCapacityPool checks resource is a capacity pool
func IsAnfCapacityPool(resourceURI string) bool {

	if len(strings.TrimSpace(resourceURI)) == 0 || !IsAnfResource(resourceURI) {
		return false
	}

	return !IsAnfSnapshot(resourceURI) &&
		!IsAnfVolume(resourceURI) &&
		strings.LastIndex(resourceURI, "/capacityPools/") > -1
}

// IsAnfAccount checks resource is an account
func IsAnfAccount(resourceURI string) bool {

	if len(strings.TrimSpace(resourceURI)) == 0 || !IsAnfResource(resourceURI) {
		return false
	}

	return !IsAnfSnapshot(resourceURI) &&
		!IsAnfVolume(resourceURI) &&
		!IsAnfCapacityPool(resourceURI) &&
		strings.LastIndex(resourceURI, "/backupPolicies/") == -1 &&
		strings.LastIndex(resourceURI, "/netAppAccounts/") > -1
}

// WaitForANFResource - Waits for ANF resource to be available from Get perspective
func WaitForANFResource(ctx context.Context, meta interface{}, resourceID string, intervalInSec int, retries int, checkForReplication bool) error {

	var err error

	for i := 0; i < retries; i++ {
		time.Sleep(time.Duration(intervalInSec) * time.Second)
		if IsAnfSnapshot(resourceID) {
			client := meta.(*clients.Client).NetApp.SnapshotClient
			_, err = client.Get(
				ctx,
				GetResourceGroup(resourceID),
				GetAnfAccount(resourceID),
				GetAnfCapacityPool(resourceID),
				GetAnfVolume(resourceID),
				GetAnfSnapshot(resourceID),
			)
		} else if IsAnfVolume(resourceID) {
			client := meta.(*clients.Client).NetApp.VolumeClient
			if checkForReplication == false {
				_, err = client.Get(
					ctx,
					GetResourceGroup(resourceID),
					GetAnfAccount(resourceID),
					GetAnfCapacityPool(resourceID),
					GetAnfVolume(resourceID),
				)
			} else {
				_, err = client.ReplicationStatusMethod(
					ctx,
					GetResourceGroup(resourceID),
					GetAnfAccount(resourceID),
					GetAnfCapacityPool(resourceID),
					GetAnfVolume(resourceID),
				)
			}
		} else if IsAnfCapacityPool(resourceID) {
			client := meta.(*clients.Client).NetApp.PoolClient
			_, err = client.Get(
				ctx,
				GetResourceGroup(resourceID),
				GetAnfAccount(resourceID),
				GetAnfCapacityPool(resourceID),
			)
		} else if IsAnfAccount(resourceID) {
			client := meta.(*clients.Client).NetApp.AccountClient
			_, err = client.Get(
				ctx,
				GetResourceGroup(resourceID),
				GetAnfAccount(resourceID),
			)
		}

		// In this case, we exit when there is no error
		if err == nil {
			return nil
		}
	}

	return fmt.Errorf("resource still not found after number of retries: %v, error: %v", retries, err)
}
