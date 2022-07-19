package validate

import (
	"strconv"
	"testing"
)

func TestFluidRelayServerName(t *testing.T) {
	type args struct {
		Input string
		Valid bool
	}
	tests := []args{
		{
			Input: "",
			Valid: false,
		},
		{
			Input: "Ahc_dfs",
			Valid: false,
		},
		{
			Input: "fds##$%#",
			Valid: false,
		},
		{
			Input: "SUlwnodfs",
			Valid: true,
		},
		{
			Input: "09810-",
			Valid: true,
		},
		{
			Input: "-u432948230",
			Valid: true,
		},
		{
			Input: "jkfdsj_de",
			Valid: false,
		},
		{
			Input: "njkjoinjknkjnakjvdsjaneoihroiwjioionb6789233jn",
			Valid: true,
		},
		{
			Input: "njk joinjknkjnakjvdsjaneoihroiwjioionb6789233jn",
			Valid: false,
		},
		{
			Input: "njk1234567890sdfghjklcvbnmkjhgfcvbnjoinjknkjnakjvdsjaneoihroiwjioionb6789233jn",
			Valid: false,
		},
	}
	for idx, tt := range tests {
		t.Run(strconv.FormatInt(int64(idx), 10), func(t *testing.T) {
			_, gotErrors := FluidRelayServerName(tt.Input, "test")
			if (len(gotErrors) == 0) != tt.Valid {
				t.Fatalf("server name validate `%s` expect %#v but got %#v", tt.Input, tt.Valid, gotErrors)
			}
		})
	}
}
