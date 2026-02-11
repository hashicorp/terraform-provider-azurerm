// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"log"
	"reflect"
	"strings"
	"sync"
)

type IdentityData struct {
	// raw identity data will be stored internally
	raw    map[string]string
	schema map[string]*Schema

	// Don't set
	once        sync.Once
	multiReader *MultiLevelFieldReader
	setWriter   *MapFieldWriter

	panicOnError bool
}

// Reading/writing data will be similar to the *schema.ResourceData flatmap
func (d *IdentityData) Get(key string) interface{} {
	v, _ := d.GetOk(key)
	return v
}

func (d *IdentityData) GetOk(key string) (interface{}, bool) {
	r := d.getRaw(key)
	exists := r.Exists
	if exists {
		// If it exists, we also want to verify it is not the zero-value.
		value := r.Value
		zero := r.Schema.Type.Zero()

		if eq, ok := value.(Equal); ok {
			exists = !eq.Equal(zero)
		} else {
			exists = !reflect.DeepEqual(value, zero)
		}
	}

	return r.Value, exists
}

func (d *IdentityData) Set(key string, value interface{}) error {
	d.once.Do(d.init)

	// If the value is a pointer to a non-struct, get its value and
	// use that. This allows Set to take a pointer to primitives to
	// simplify the interface.
	reflectVal := reflect.ValueOf(value)
	if reflectVal.Kind() == reflect.Ptr {
		if reflectVal.IsNil() {
			// If the pointer is nil, then the value is just nil
			value = nil
		} else {
			// Otherwise, we dereference the pointer as long as its not
			// a pointer to a struct, since struct pointers are allowed.
			reflectVal = reflect.Indirect(reflectVal)
			if reflectVal.Kind() != reflect.Struct {
				value = reflectVal.Interface()
			}
		}
	}

	err := d.setWriter.WriteField(strings.Split(key, "."), value)
	if err != nil {
		if d.panicOnError {
			panic(err)
		} else {
			log.Printf("[ERROR] setting identity state: %s", err)
		}
	}
	return err
}

func (d *IdentityData) init() {
	// Initialize the map for storing data set by the user
	d.setWriter = &MapFieldWriter{Schema: d.schema}

	// Initialize the reader for getting data from the
	// underlying sources (config, diff, etc.)
	readers := make(map[string]FieldReader)
	if d.raw != nil {
		readers["raw"] = &MapFieldReader{
			Schema: d.schema,
			Map:    BasicMapReader(d.raw),
		}
	}
	readers["set"] = &MapFieldReader{
		Schema: d.schema,
		Map:    BasicMapReader(d.setWriter.Map()),
	}
	d.multiReader = &MultiLevelFieldReader{
		Levels:  []string{"raw", "set"},
		Readers: readers,
	}
}

func (d *IdentityData) getRaw(key string) getResult {
	var parts []string
	if key != "" {
		parts = strings.Split(key, ".")
	}

	return d.get(parts)
}

func (d *IdentityData) get(addr []string) getResult {
	d.once.Do(d.init)

	result, err := d.multiReader.ReadFieldMerge(addr, "set")

	if err != nil {
		panic(err)
	}

	// If the result doesn't exist, then we set the value to the zero value
	var schema *Schema
	if schemaL := addrToSchema(addr, d.schema); len(schemaL) > 0 {
		schema = schemaL[len(schemaL)-1]
	}

	if result.Value == nil && schema != nil {
		result.Value = result.ValueOrZero(schema)
	}

	// Transform the FieldReadResult into a getResult. It might be worth
	// merging these two structures one day.
	return getResult{
		Value:          result.Value,
		ValueProcessed: result.ValueProcessed,
		Computed:       result.Computed,
		Exists:         result.Exists,
		Schema:         schema,
	}
}
