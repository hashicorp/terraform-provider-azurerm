package validate

import (
	"fmt"
	"reflect"
	"testing"
)

const (
	FieldName = "max_bid_price"
)

func TestMaxBidPrice(t *testing.T) {
	type args struct {
		value interface{}
		key   string
	}
	tests := []struct {
		name       string
		args       args
		wantErrors []error
	}{
		{
			name: fmt.Sprintf("%q = 0.055", FieldName),
			args: args{
				value: interface{}(float64(0.055)),
				key:   FieldName,
			},
			wantErrors: nil,
		},
		{
			name: fmt.Sprintf("%q = 0.00001", FieldName),
			args: args{
				value: interface{}(float64(0.00001)),
				key:   FieldName,
			},
			wantErrors: nil,
		},
		{
			name: fmt.Sprintf("%q = 1002.3333", FieldName),
			args: args{
				value: interface{}(float64(1002.3333)),
				key:   FieldName,
			},
			wantErrors: nil,
		},
		{
			name: fmt.Sprintf("%q = 0.00007", FieldName),
			args: args{
				value: interface{}(float64(0.00007)),
				key:   FieldName,
			},
			wantErrors: nil,
		},
		{
			name: fmt.Sprintf("%q = -2", FieldName),
			args: args{
				value: interface{}(float64(-2.0)),
				key:   FieldName,
			},
			wantErrors: []error{
				fmt.Errorf("%q must be between 0 (exclusive) and `math.MaxFloat64` (inclusive) or -1 (special value), got %f",
					FieldName, float64(-2.0)),
			},
		},
		{
			name: fmt.Sprintf("%q = 0.000009", FieldName),
			args: args{
				value: interface{}(float64(0.000009)),
				key:   FieldName,
			},
			wantErrors: []error{
				fmt.Errorf("%q can only include up to 5 digits after the radix point, got %g",
					FieldName, float64(0.000009)),
			},
		},
		{
			name: fmt.Sprintf("%q = 0.000000001", FieldName),
			args: args{
				value: interface{}(float64(0.000000001)),
				key:   FieldName,
			},
			wantErrors: []error{
				fmt.Errorf("%q can only include up to 5 digits after the radix point, got %g",
					FieldName, float64(0.000000001)),
			},
		},
		{
			name: fmt.Sprintf("%q = -2.000000001", FieldName),
			args: args{
				value: interface{}(float64(-2.000000001)),
				key:   FieldName,
			},
			wantErrors: []error{
				fmt.Errorf("%q must be between 0 (exclusive) and `math.MaxFloat64` (inclusive) or -1 (special value), got %f",
					FieldName, float64(-2.000000001)),
				fmt.Errorf("%q can only include up to 5 digits after the radix point, got %g",
					FieldName, float64(-2.000000001)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWarns, gotErrors := MaxBidPrice(tt.args.value, tt.args.key)

			if !reflect.DeepEqual(gotErrors, tt.wantErrors) {
				t.Errorf("SpotPrice() gotErrors = %v, want %v", gotErrors, tt.wantErrors)
			}

			// MaxBidPrice doesn't emit warnings
			if !reflect.DeepEqual(gotWarns, []string(nil)) {
				t.Errorf("SpotPrice() gotWarns = %v, want %v", gotWarns, []string(nil))
			}
		})
	}
}
