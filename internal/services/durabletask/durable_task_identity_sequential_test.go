// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package durabletask_test

import "testing"

func TestAccDurableTaskScheduler_resourceIdentity(t *testing.T) {
	testAccDurableTaskScheduler_resourceIdentity(t)
}

func TestAccDurableTaskHub_resourceIdentity(t *testing.T) {
	testAccDurableTaskHub_resourceIdentity(t)
}

func TestAccDurableTaskRetentionPolicy_resourceIdentity(t *testing.T) {
	testAccDurableTaskRetentionPolicy_resourceIdentity(t)
}
