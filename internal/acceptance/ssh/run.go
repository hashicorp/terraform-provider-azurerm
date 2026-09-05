// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package ssh

type Runner struct {
	Hostname      string
	Port          int
	Username      string
	Password      string
	CommandsToRun []string
}
