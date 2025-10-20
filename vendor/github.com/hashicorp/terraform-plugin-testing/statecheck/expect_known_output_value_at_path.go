// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package statecheck

import (
	"context"
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// Resource State Check
var _ StateCheck = expectKnownOutputValueAtPath{}

type expectKnownOutputValueAtPath struct {
	outputAddress string
	outputPath    tfjsonpath.Path
	knownValue    knownvalue.Check
}

// CheckState implements the state check logic.
func (e expectKnownOutputValueAtPath) CheckState(ctx context.Context, req CheckStateRequest, resp *CheckStateResponse) {
	var output *tfjson.StateOutput

	if req.State == nil {
		resp.Error = fmt.Errorf("state is nil")

		return
	}

	if req.State.Values == nil {
		resp.Error = fmt.Errorf("state does not contain any state values")

		return
	}

	for address, oc := range req.State.Values.Outputs {
		if e.outputAddress == address {
			output = oc

			break
		}
	}

	if output == nil {
		resp.Error = fmt.Errorf("%s - Output not found in state", e.outputAddress)

		return
	}

	result, err := tfjsonpath.Traverse(output.Value, e.outputPath)

	if err != nil {
		resp.Error = err

		return
	}

	if err := e.knownValue.CheckValue(result); err != nil {
		resp.Error = fmt.Errorf("error checking value for output at path: %s.%s, err: %s", e.outputAddress, e.outputPath.String(), err)

		return
	}
}

// ExpectKnownOutputValueAtPath returns a state check that asserts that the specified output at the given path
// has a known type and value.
func ExpectKnownOutputValueAtPath(outputAddress string, outputPath tfjsonpath.Path, knownValue knownvalue.Check) StateCheck {
	return expectKnownOutputValueAtPath{
		outputAddress: outputAddress,
		outputPath:    outputPath,
		knownValue:    knownValue,
	}
}
