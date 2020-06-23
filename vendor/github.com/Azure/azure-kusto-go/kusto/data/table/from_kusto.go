package table

// value.go provides methods for converting a row to a *struct and for converting KustoValue into Go types
// or in the reverse.

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/Azure/azure-kusto-go/kusto/data/types"
	"github.com/Azure/azure-kusto-go/kusto/data/value"

	"github.com/google/uuid"
)

// decodeToStruct takes a list of columns and a row to decode into "p" which will be a pointer
// to a struct (enforce in the decoder).
func decodeToStruct(cols Columns, row value.Values, p interface{}) error {
	t := reflect.TypeOf(p)
	v := reflect.ValueOf(p)
	fields := newFields(cols, t)

	for i, col := range cols {
		if err := fields.convert(col, row[i], t, v); err != nil {
			return err
		}
	}
	return nil
}

// fields represents the fields inside a struct.
type fields struct {
	colNameToFieldName map[string]string
}

// newFields takes in the Columns from our row and the reflect.Type of our *struct.
func newFields(cols Columns, ptr reflect.Type) fields {
	nFields := fields{colNameToFieldName: map[string]string{}}
	for i := 0; i < ptr.Elem().NumField(); i++ {
		field := ptr.Elem().Field(i)
		if tag := field.Tag.Get("kusto"); strings.TrimSpace(tag) != "" {
			nFields.colNameToFieldName[tag] = field.Name
		} else {
			nFields.colNameToFieldName[field.Name] = field.Name
		}
	}

	return nFields
}

// match returns the name of the field in the struct that matches col. Empty string indicates there
// is no match.
func (f fields) match(col Column) string {
	return f.colNameToFieldName[col.Name]
}

// convert converts a KustoValue that is for Column col into "v" reflect.Value with reflect.Type "t".
func (f fields) convert(col Column, k value.Kusto, t reflect.Type, v reflect.Value) error {
	fieldName, ok := f.colNameToFieldName[col.Name]
	if !ok {
		return nil
	}

	if fieldName == "-" {
		return nil
	}

	switch col.Type {
	case types.Bool:
		return boolConvert(fieldName, col, k, t, v)
	case types.DateTime:
		return dateTimeConvert(fieldName, col, k, t, v)
	case types.Dynamic:
		return dynamicConvert(fieldName, col, k, t, v)
	case types.GUID:
		return guidConvert(fieldName, col, k, t, v)
	case types.Int:
		return intConvert(fieldName, col, k, t, v)
	case types.Long:
		return longConvert(fieldName, col, k, t, v)
	case types.Real:
		return realConvert(fieldName, col, k, t, v)
	case types.String:
		return stringConvert(fieldName, col, k, t, v)
	case types.Timespan:
		return timespanConvert(fieldName, col, k, t, v)
	case types.Decimal:
		return decimalConvert(fieldName, col, k, t, v)
	}
	return fmt.Errorf("received a field type %q we don't recognize", col.Type)
}

/*
The next section has conversion types that allow us to change our value.Values into the underlying types or
Go representations of them, and vice versus. Conversion from a value.Values will be <type>Convert and
conversion to a value.Values will be convert<Type>.  We support conversion to/from value.Values to
mulitple Go types.  For example, our types.Bool could be converted to types.Bool, *types.Bool,
bool or *bool.  Others work similar to this.
*/

func boolConvert(fieldName string, col Column, k value.Kusto, t reflect.Type, v reflect.Value) error {
	val, ok := k.(value.Bool)
	if !ok {
		return fmt.Errorf("Column %s is type %s was trying to store a KustoValue type of %T", col.Name, col.Type, k)
	}
	sf, _ := t.Elem().FieldByName(fieldName)
	switch {
	case sf.Type.Kind() == reflect.Bool:
		if val.Valid {
			v.Elem().FieldByName(fieldName).SetBool(val.Value)
		}
		return nil
	case sf.Type.ConvertibleTo(reflect.TypeOf(new(bool))):
		if val.Valid {
			b := new(bool)
			if val.Value {
				*b = true
			}
			v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(b))
		}
		return nil
	case sf.Type.ConvertibleTo(reflect.TypeOf(value.Bool{})):
		v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(val))
		return nil
	case sf.Type.ConvertibleTo(reflect.TypeOf(&value.Bool{})):
		v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(&val))
		return nil
	}
	return fmt.Errorf("column %s could not store in struct.%s: column was type Kusto.Bool, struct had base Kind %s ", col.Name, fieldName, sf.Type.Kind())
}

func dateTimeConvert(fieldName string, col Column, k value.Kusto, t reflect.Type, v reflect.Value) error {
	val, ok := k.(value.DateTime)
	if !ok {
		return fmt.Errorf("Column %s is type %s was trying to store a KustoValue type of %T", col.Name, col.Type, k)
	}
	sf, _ := t.Elem().FieldByName(fieldName)
	switch {
	case sf.Type.AssignableTo(reflect.TypeOf(time.Time{})):
		if val.Valid {
			v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(val.Value))
		}
		return nil
	case sf.Type.ConvertibleTo(reflect.TypeOf(new(time.Time))):
		if val.Valid {
			t := &val.Value
			v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(t))
		}
		return nil
	case sf.Type.ConvertibleTo(reflect.TypeOf(value.DateTime{})):
		v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(val))
		return nil
	case sf.Type.ConvertibleTo(reflect.TypeOf(&value.DateTime{})):
		v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(&val))
		return nil
	}
	return fmt.Errorf("column %s could not store in struct.%s: column was type Kusto.DateTime, struct had type %s ", col.Name, fieldName, sf.Type.Name())
}

func timespanConvert(fieldName string, col Column, k value.Kusto, t reflect.Type, v reflect.Value) error {
	val, ok := k.(value.Timespan)
	if !ok {
		return fmt.Errorf("Column %s is type %s was trying to store a KustoValue type of %T", col.Name, col.Type, k)
	}
	sf, _ := t.Elem().FieldByName(fieldName)
	switch {
	case sf.Type.AssignableTo(reflect.TypeOf(time.Duration(0))):
		if val.Valid {
			v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(val.Value))
		}
		return nil
	case sf.Type.ConvertibleTo(reflect.TypeOf(new(time.Duration))):
		if val.Valid {
			t := &val.Value
			v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(t))
		}
		return nil
	case sf.Type.ConvertibleTo(reflect.TypeOf(value.Timespan{})):
		v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(val))
		return nil
	case sf.Type.ConvertibleTo(reflect.TypeOf(&value.Timespan{})):
		v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(&val))
		return nil
	}
	return fmt.Errorf("column %s could not store in struct.%s: column was type Kusto.Timespan, struct had type %s ", col.Name, fieldName, sf.Type.Name())
}

func dynamicConvert(fieldName string, col Column, k value.Kusto, t reflect.Type, v reflect.Value) error {
	val, ok := k.(value.Dynamic)
	if !ok {
		return fmt.Errorf("Column %s is type %s was tryihng to store a KustoValue type of %T", col.Name, col.Type, k)
	}
	sf, _ := t.Elem().FieldByName(fieldName)
	switch {
	case sf.Type.ConvertibleTo(reflect.TypeOf(value.Dynamic{})):
		v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(val))
		return nil
	case sf.Type.ConvertibleTo(reflect.TypeOf(&value.Dynamic{})):
		v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(&val))
		return nil
	case sf.Type.ConvertibleTo(reflect.TypeOf([]byte{})):
		if sf.Type.Kind() == reflect.String {
			s := string(val.Value)
			v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(s))
			return nil
		}
		v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(val.Value))
		return nil
	case sf.Type.Kind() == reflect.Map:
		if !val.Valid {
			return nil
		}
		if sf.Type.Key().Kind() != reflect.String {
			return fmt.Errorf("Column %s is type dymanic and can only be stored in a string, *string, map[string]interface{}, *map[string]interface{}, struct or *struct", col.Name)
		}
		if sf.Type.Elem().Kind() != reflect.Interface {
			return fmt.Errorf("Column %s is type dymanic and can only be stored in a string, *string, map[string]interface{}, *map[string]interface{}, struct or *struct", col.Name)
		}

		m := map[string]interface{}{}
		if err := json.Unmarshal([]byte(val.Value), &m); err != nil {
			return fmt.Errorf("Column %s is type dymanic and had an error tyring to marshal into a map[string]interface{}: %s", col.Name, err)
		}

		v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(m))
		return nil
	case sf.Type.Kind() == reflect.Struct:
		structPtr := reflect.New(sf.Type)

		if err := json.Unmarshal([]byte(val.Value), structPtr.Interface()); err != nil {
			return fmt.Errorf("Column %s of type dynamic could not unmarshal into the passed struct: %s", col.Name, err)
		}

		v.Elem().FieldByName(fieldName).Set(structPtr.Elem())
		return nil
	case sf.Type.Kind() == reflect.Ptr:
		if !val.Valid {
			return nil
		}

		switch {
		case sf.Type.Elem().Kind() == reflect.String:
			str := string(val.Value)
			v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(&str))
			return nil
		case sf.Type.Elem().Kind() == reflect.Struct:
			store := reflect.New(sf.Type.Elem())

			if err := json.Unmarshal([]byte(val.Value), store.Interface()); err != nil {
				return fmt.Errorf("Column %s of type dynamic could not unmarshal into the passed *struct: %s", col.Name, err)
			}
			v.Elem().FieldByName(fieldName).Set(store)
			return nil
		case sf.Type.Elem().Kind() == reflect.Map:
			if sf.Type.Elem().Key().Kind() != reflect.String {
				return fmt.Errorf("Column %s is type dymanic and can only be stored in a map of type map[string]interface{}", col.Name)
			}
			if sf.Type.Elem().Elem().Kind() != reflect.Interface {
				return fmt.Errorf("Column %s is type dymanic and can only be stored in a of type map[string]interface{}", col.Name)
			}

			m := map[string]interface{}{}
			if err := json.Unmarshal([]byte(val.Value), &m); err != nil {
				return fmt.Errorf("Column %s is type dymanic and had an error tyring to marshal into a map[string]interface{}: %s", col.Name, err)
			}
			v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(&m))
			return nil
		}
	}
	return fmt.Errorf("column %s could not store in struct.%s: column was type Kusto.Dynamic, struct had base Kind %s ", col.Name, fieldName, sf.Type.Kind())
}

func guidConvert(fieldName string, col Column, k value.Kusto, t reflect.Type, v reflect.Value) error {
	val, ok := k.(value.GUID)
	if !ok {
		return fmt.Errorf("Column %s is type %s was trying to store a KustoValue type of %T", col.Name, col.Type, k)
	}
	sf, _ := t.Elem().FieldByName(fieldName)
	switch {
	case sf.Type.AssignableTo(reflect.TypeOf(uuid.UUID{})):
		if val.Valid {
			v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(val.Value))
		}
		return nil
	case sf.Type.ConvertibleTo(reflect.TypeOf(new(uuid.UUID))):
		if val.Valid {
			t := &val.Value
			v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(t))
		}
		return nil
	case sf.Type.ConvertibleTo(reflect.TypeOf(value.GUID{})):
		v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(val))
		return nil
	case sf.Type.ConvertibleTo(reflect.TypeOf(&value.GUID{})):
		v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(&val))
		return nil
	}
	return fmt.Errorf("column %s could not store in struct.%s: column was type Kusto.GUID, struct had type %s ", col.Name, fieldName, sf.Type.Name())
}

func intConvert(fieldName string, col Column, k value.Kusto, t reflect.Type, v reflect.Value) error {
	val, ok := k.(value.Int)
	if !ok {
		return fmt.Errorf("Column %s is type %s was tryihng to store a KustoValue type of %T", col.Name, col.Type, k)
	}
	sf, _ := t.Elem().FieldByName(fieldName)
	switch {
	case sf.Type.Kind() == reflect.Int32:
		if val.Valid {
			v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(val.Value))
		}
		return nil
	case sf.Type.ConvertibleTo(reflect.TypeOf(new(int32))):
		if val.Valid {
			i := &val.Value
			v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(i))
		}
		return nil
	case sf.Type.ConvertibleTo(reflect.TypeOf(value.Int{})):
		v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(val))
		return nil
	case sf.Type.ConvertibleTo(reflect.TypeOf(&value.Int{})):
		v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(&val))
		return nil
	}
	return fmt.Errorf("column %s could not store in struct.%s: column was type Kusto.Int, struct had base Kind %s ", col.Name, fieldName, sf.Type.Kind())
}

func longConvert(fieldName string, col Column, k value.Kusto, t reflect.Type, v reflect.Value) error {
	val, ok := k.(value.Long)
	if !ok {
		return fmt.Errorf("Column %s is type %s was tryihng to store a KustoValue type of %T", col.Name, col.Type, k)
	}
	sf, _ := t.Elem().FieldByName(fieldName)
	switch {
	case sf.Type.Kind() == reflect.Int64:
		if val.Valid {
			v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(val.Value))
		}
		return nil
	case sf.Type.ConvertibleTo(reflect.TypeOf(new(int64))):
		if val.Valid {
			i := &val.Value
			v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(i))
		}
		return nil
	case sf.Type.ConvertibleTo(reflect.TypeOf(value.Long{})):
		v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(val))
		return nil
	case sf.Type.ConvertibleTo(reflect.TypeOf(&value.Long{})):
		v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(&val))
		return nil
	}
	return fmt.Errorf("column %s could not store in struct.%s: column was type Kusto.Long, struct had base Kind %s ", col.Name, fieldName, sf.Type.Kind())
}

func realConvert(fieldName string, col Column, k value.Kusto, t reflect.Type, v reflect.Value) error {
	val, ok := k.(value.Real)
	if !ok {
		return fmt.Errorf("Column %s is type %s was tryihng to store a KustoValue type of %T", col.Name, col.Type, k)
	}
	sf, _ := t.Elem().FieldByName(fieldName)
	switch {
	case sf.Type.Kind() == reflect.Float64:
		if val.Valid {
			v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(val.Value))
		}
		return nil
	case sf.Type.ConvertibleTo(reflect.TypeOf(new(float64))):
		if val.Valid {
			i := &val.Value
			v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(i))
		}
		return nil
	case sf.Type.ConvertibleTo(reflect.TypeOf(value.Real{})):
		v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(val))
		return nil
	case sf.Type.ConvertibleTo(reflect.TypeOf(&value.Real{})):
		v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(&val))
		return nil
	}
	return fmt.Errorf("column %s could not store in struct.%s: column was type Kusto.Real, struct had base Kind %s ", col.Name, fieldName, sf.Type.Kind())
}

func stringConvert(fieldName string, col Column, k value.Kusto, t reflect.Type, v reflect.Value) error {
	val, ok := k.(value.String)
	if !ok {
		return fmt.Errorf("Column %s is type %s was trying to store a KustoValue type of %T", col.Name, col.Type, k)
	}
	sf, _ := t.Elem().FieldByName(fieldName)
	switch {
	case sf.Type.Kind() == reflect.String:
		if val.Valid {
			v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(val.Value))
		}
		return nil
	case sf.Type.ConvertibleTo(reflect.TypeOf(new(string))):
		if val.Valid {
			i := &val.Value
			v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(i))
		}
		return nil
	case sf.Type.ConvertibleTo(reflect.TypeOf(value.String{})):
		v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(val))
		return nil
	case sf.Type.ConvertibleTo(reflect.TypeOf(&value.String{})):
		v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(&val))
		return nil
	}
	return fmt.Errorf("column %s could not store in struct.%s: column was type Kusto.String, struct had base Kind %s ", col.Name, fieldName, sf.Type.Kind())
}

func decimalConvert(fieldName string, col Column, k value.Kusto, t reflect.Type, v reflect.Value) error {
	val, ok := k.(value.Decimal)
	if !ok {
		return fmt.Errorf("Column %s is type %s was trying to store a KustoValue type of %T", col.Name, col.Type, k)
	}
	sf, _ := t.Elem().FieldByName(fieldName)
	switch {
	case sf.Type.Kind() == reflect.String:
		if val.Valid {
			v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(val.Value))
		}
		return nil
	case sf.Type.ConvertibleTo(reflect.TypeOf(new(string))):
		if val.Valid {
			i := &val.Value
			v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(i))
		}
		return nil
	case sf.Type.ConvertibleTo(reflect.TypeOf(value.Decimal{})):
		v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(val))
		return nil
	case sf.Type.ConvertibleTo(reflect.TypeOf(&value.Decimal{})):
		v.Elem().FieldByName(fieldName).Set(reflect.ValueOf(&val))
		return nil
	}
	return fmt.Errorf("column %s could not store in struct.%s: column was type Kusto.Decimal, struct had base Kind %s ", col.Name, fieldName, sf.Type.Kind())
}
