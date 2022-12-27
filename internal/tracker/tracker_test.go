package tracker

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetSrcAndDst(t *testing.T) {
	type input struct {
		flights Flights
	}

	type want struct {
		flight []string
		err    error
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"handles empty flight data",
			input{
				flights: [][]string{},
			},
			want{
				err: errors.New("failed to find any flights to track"),
			},
		},
		{
			"handles invalid flight data format",
			input{
				flights: [][]string{
					{"IND", "EWR", "NYC"},
					{"SFO", "ATL"},
					{"GSO", "IND"},
					{"ATL", "GSO"},
				},
			},
			want{
				err: errors.New("invalid flight data format"),
			},
		},
		{
			"handles disconnected flights",
			input{
				flights: [][]string{
					{"IND", "EWR"},
					{"SFO", "ATL"},
					{"GSO", "IND"},
					{"ATL", "NYC"},
				},
			},
			want{
				err: errors.New("failed to find connecting flights: broken flight chain"),
			},
		},
		{
			"is successful with no layovers",
			input{
				flights: [][]string{
					{"SFO", "EWR"},
				},
			},
			want{
				flight: []string{"SFO", "EWR"},
			},
		},
		{
			"is successful with one layover",
			input{
				flights: [][]string{
					{"SFO", "ATL"},
					{"ATL", "EWR"},
				},
			},
			want{
				flight: []string{"SFO", "EWR"},
			},
		},
		{
			"is successful with multiple layovers",
			input{
				flights: [][]string{
					{"IND", "EWR"},
					{"SFO", "ATL"},
					{"GSO", "IND"},
					{"ATL", "GSO"},
				},
			},
			want{
				flight: []string{"SFO", "EWR"},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			tracker := New()
			flight, err := tracker.GetSrcAndDst(tt.input.flights)
			if tt.want.err != nil {
				if !assert.NotNil(t, err) {
					t.Fatalf("expected an error: %s", tt.want.err.Error())
				}

				assert.Contains(t, tt.want.err.Error(), err.Error())
			}

			assert.Equal(t, tt.want.flight, flight)
		})
	}
}
