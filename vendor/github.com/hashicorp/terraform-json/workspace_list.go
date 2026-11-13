// Copyright IBM Corp. 2019, 2026
// SPDX-License-Identifier: MPL-2.0

package tfjson

// WorkspaceListOutput represents JSON output from terraform workspace list
// (available from 1.16 onwards)
type WorkspaceListOutput struct {
	FormatVersion string `json:"format_version"`

	Workspaces  []WorkspaceListEntry `json:"workspaces"`
	Diagnostics []Diagnostic         `json:"diagnostics"`
}

// WorkspaceListEntry represents a single workspace entry in the list of workspaces
type WorkspaceListEntry struct {
	Name      string `json:"name"`
	IsCurrent bool   `json:"is_current"`
}

// CurrentWorkspace will return details about the workspace that's currently selected.
// The result may be nil if:
// 1) the command output is empty, i.e. no workspaces exist yet.
// 2) or, if the output doesn't specify a selected workspace.
// Number 2 may happen if the default workspace is selected but not created yet.
//
// This method doesn't return any errors when nil is returned, as there are valid conditions where
// nil is returned. Calling code should look for diagnostics in the receiver WorkspaceListOutput to
// see if an explicit error has occurred or not.
func (wlo *WorkspaceListOutput) CurrentWorkspace() *WorkspaceListEntry {
	if wlo == nil || len(wlo.Workspaces) == 0 {
		return nil
	}

	for _, ws := range wlo.Workspaces {
		if ws.IsCurrent {
			return &ws
		}
	}
	return nil
}
