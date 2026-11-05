// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/hostpool"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/msiximage"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/method"
)

func GetMsixImageProperties(ctx context.Context, metadata sdk.ResourceMetaData, hostPoolIds []string, storageShareFileId string, msixPackageName string) (*msiximage.ExpandMsixImageProperties, error) {
	msixImageUri := msiximage.MSIXImageURI{
		Uri: pointer.To(storageShareFileId),
	}

	availableMsixPackageNames := make(map[string][]string)
	errs := make([]string, 0)

	for _, hostPoolReference := range hostPoolIds {
		availableMsixPackageNames[hostPoolReference] = make([]string, 0)
		hostPoolId, err := hostpool.ParseHostPoolIDInsensitively(hostPoolReference)
		if err != nil {
			return nil, fmt.Errorf("parsing host pool ID %s: %+v", hostPoolReference, err)
		}

		msixImageHostPoolId := msiximage.NewHostPoolID(hostPoolId.SubscriptionId, hostPoolId.ResourceGroupName, hostPoolId.HostPoolName)

		// Replicate portal behavior where MSIX image properties required for app attach package creation are retrieved with API
		// `Expand` method from imported `msiximage` package is not used as it lacks codes to marshal request body
		result, err := method.ExpandCompleteMsixImage(ctx, metadata, msixImageHostPoolId, msixImageUri)
		if err != nil {
			// Continue to check next host pool if there is error with expanding MSIX image of current host pool
			errs = append(errs, fmt.Sprintf("expanding MSIX image for host pool %s: %+v", hostPoolReference, err))
			continue
		}

		msixImages := result.Items
		for _, msixImage := range msixImages {
			if properties := msixImage.Properties; properties != nil {
				// Return MSIX image with matched `PackageFullName`
				if properties.PackageFullName != nil && strings.EqualFold(pointer.From(properties.PackageFullName), msixPackageName) {
					return properties, nil
				}
			}
		}

		// Store available `PackageFullName` for user reference if no matched MSIX image is found after checking all host pools
		for _, msixImage := range msixImages {
			if properties := msixImage.Properties; properties != nil && properties.PackageFullName != nil {
				availableMsixPackageNames[hostPoolReference] = append(availableMsixPackageNames[hostPoolReference], pointer.From(properties.PackageFullName))
			}
		}
	}

	concatenatedAvailableMsixPackageNames := ""
	for hostPoolReference, packageFullNames := range availableMsixPackageNames {
		concatenatedAvailableMsixPackageNames += fmt.Sprintf("%v from %s, ", packageFullNames, hostPoolReference)
	}

	concatenatedAvailableMsixPackageNames = strings.TrimSuffix(concatenatedAvailableMsixPackageNames, ", ")

	return nil, fmt.Errorf("no matched MSIX image with package name %s was found. The available package names are %s. The list of errors encountered is: %v", msixPackageName, concatenatedAvailableMsixPackageNames, errs)
}
