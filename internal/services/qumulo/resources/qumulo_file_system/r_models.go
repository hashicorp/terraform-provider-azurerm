// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package qumulo_file_system

func (r Resource) ModelObject() any {
	return &FileSystemResourceSchema{}
}
