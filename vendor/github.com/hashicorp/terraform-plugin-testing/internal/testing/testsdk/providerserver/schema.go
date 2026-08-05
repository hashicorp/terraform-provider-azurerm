// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package providerserver

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func SchemaAttributeAtPath(schema *tfprotov6.Schema, path *tftypes.AttributePath) *tfprotov6.SchemaAttribute {
	if schema == nil || schema.Block == nil || path == nil || len(path.Steps()) == 0 {
		return nil
	}

	steps := path.Steps()
	nextStep := steps[0]
	remainingSteps := steps[1:]

	switch nextStep := nextStep.(type) {
	case tftypes.AttributeName:
		for _, attribute := range schema.Block.Attributes {
			if attribute == nil {
				continue
			}

			if attribute.Name != string(nextStep) {
				continue
			}

			if len(remainingSteps) == 0 {
				return attribute
			}

			// If needed, recursive attribute.NestedType handling would go here.
		}

		for _, block := range schema.Block.BlockTypes {
			if block == nil {
				continue
			}

			if block.TypeName != string(nextStep) {
				continue
			}

			// Blocks cannot be computed.
			if len(remainingSteps) == 0 {
				return nil
			}

			// If needed, recursive block handling would go here.
		}
	}

	return nil
}
