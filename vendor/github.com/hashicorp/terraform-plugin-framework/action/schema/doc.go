// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package schema contains all available schema functionality for actions.
// Action schemas define the structure and value types for configuration data.
// Schemas are implemented via the action.Action type Schema method.
//
// There are three different types of action schemas, which define how a practitioner can trigger an action,
// as well as what effect the action can have on the state.
//   - [UnlinkedSchema] actions are actions that cannot cause changes to resource states.
//   - [LifecycleSchema] actions are actions that can cause changes to exactly one resource state.
//   - [LinkedSchema] actions are actions that can cause changes to one or more resource states.
package schema
