// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azurecli

var cache *cachedCliData

type cachedCliData struct {
	data map[string][]byte
}

func (c *cachedCliData) Set(index string, data []byte) {
	c.data[index] = data
}

func (c *cachedCliData) Get(index string) ([]byte, bool) {
	if data, ok := c.data[index]; ok {
		return data, true
	}
	return nil, false
}

func init() {
	cache = &cachedCliData{
		data: make(map[string][]byte),
	}
}
