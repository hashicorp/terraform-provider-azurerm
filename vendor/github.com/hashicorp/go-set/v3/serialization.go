// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package set

import "encoding/json"

// marshalJSON will serialize a Serializable[T] into a json byte array
func marshalJSON[T any](s Collection[T]) ([]byte, error) {
	return json.Marshal(s.Slice())
}

// unmarshalJSON will deserialize a json byte array into a Serializable[T]
func unmarshalJSON[T any](s Collection[T], data []byte) error {
	slice := make([]T, 0)
	err := json.Unmarshal(data, &slice)
	if err != nil {
		return err
	}
	s.InsertSlice(slice)
	return nil
}
