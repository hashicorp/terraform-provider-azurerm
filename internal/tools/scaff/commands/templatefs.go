// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package commands

import "embed"

// Templatedir previously embedded the text/template assets used by the legacy
// template-based commands (config, document, servicepackage). Those templates
// have been removed as part of the move to the gen/ir packages; this empty FS
// keeps the remaining references compiling until those commands are ported or
// removed. The generate command does not use it.
var Templatedir embed.FS
