package event

import "testing"

func TestDequeueNoWait(t *testing.T) {
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

func TestDequeueWait(t *testing.T) {
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
