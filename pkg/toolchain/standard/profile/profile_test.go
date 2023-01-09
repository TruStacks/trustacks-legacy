package profile

import (
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		Profile *Profile
		Valid   bool
	}{

		{
			&Profile{},
			false,
		},
		{
			&Profile{Domain: "local.gd"},
			true,
		},
	}
	for _, tc := range tests {
		if err := tc.Profile.Validate(); err != nil && tc.Valid {
			t.Fatal(err)
		} else if err == nil && !tc.Valid {
			t.Fatal("expected validation to fail")
		}
	}
}
