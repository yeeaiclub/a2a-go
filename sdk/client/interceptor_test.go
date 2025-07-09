package auth

import "testing"

func TestIntercept(t *testing.T) {
	testcases := []struct {
		name string
	}{
		{
			name: "",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
		})
	}
}
