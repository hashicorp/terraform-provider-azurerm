// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryapplicationversions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
)

var _ resourceids.Id = VirtualMachineGalleryApplicationAssignmentId{}

type VirtualMachineGalleryApplicationAssignmentId struct {
	VirtualMachineId            virtualmachines.VirtualMachineId
	GalleryApplicationVersionId galleryapplicationversions.ApplicationVersionId
}

func (v VirtualMachineGalleryApplicationAssignmentId) ID() string {
	return fmt.Sprintf("%s|%s", v.VirtualMachineId.ID(), v.GalleryApplicationVersionId.ID())
}

func (v VirtualMachineGalleryApplicationAssignmentId) String() string {
	components := []string{
		fmt.Sprintf("VirtualMachineId %s", v.VirtualMachineId.ID()),
		fmt.Sprintf("GalleryApplicationVersionId %s", v.GalleryApplicationVersionId.ID()),
	}
	return fmt.Sprintf("Virtual Machine Gallery Application Assignment: %s", strings.Join(components, " / "))
}

func NewVirtualMachineGalleryApplicationAssignmentID(virtualMachineId virtualmachines.VirtualMachineId, galleryApplicationVersionId galleryapplicationversions.ApplicationVersionId) VirtualMachineGalleryApplicationAssignmentId {
	return VirtualMachineGalleryApplicationAssignmentId{
		VirtualMachineId:            virtualMachineId,
		GalleryApplicationVersionId: galleryApplicationVersionId,
	}
}

func VirtualMachineGalleryApplicationAssignmentID(input string) (VirtualMachineGalleryApplicationAssignmentId, error) {
	splitId := strings.Split(input, "|")
	if len(splitId) != 2 {
		return VirtualMachineGalleryApplicationAssignmentId{}, fmt.Errorf("expected ID to be in the format {VirtualMachineId}|{GalleryApplicationVersionId} but got %q", input)
	}

	virtualMachineId, err := virtualmachines.ParseVirtualMachineID(splitId[0])
	if err != nil {
		return VirtualMachineGalleryApplicationAssignmentId{}, err
	}

	galleryApplicationVersionId, err := galleryapplicationversions.ParseApplicationVersionID(splitId[1])
	if err != nil {
		return VirtualMachineGalleryApplicationAssignmentId{}, err
	}

	if virtualMachineId == nil || galleryApplicationVersionId == nil {
		return VirtualMachineGalleryApplicationAssignmentId{}, fmt.Errorf("parse error, both VirtualMachineId and GalleryApplicationVersionId should not be nil")
	}

	return VirtualMachineGalleryApplicationAssignmentId{
		VirtualMachineId:            *virtualMachineId,
		GalleryApplicationVersionId: *galleryApplicationVersionId,
	}, nil
}

func VirtualMachineGalleryApplicationAssignmentIDValidation(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := VirtualMachineGalleryApplicationAssignmentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}
