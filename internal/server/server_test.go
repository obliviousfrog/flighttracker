package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/obliviousfrog/flighttracker/internal/tracker"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_GetSrcAndDst(t *testing.T) {
	type input struct {
		requestFn func() ([]string, error)
	}

	type want struct {
		flights []string
		err     error
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"handles wrong endpoint",
			input{
				requestFn: func() ([]string, error) {
					flights := "garbage"

					v, err := json.Marshal(&flights)
					if !assert.Nil(t, err) {
						t.Fatal("failed to setup test")
					}

					resp, err := http.Post("http://127.0.0.1:8080/garbage", "application/json", bytes.NewBuffer(v))
					if err != nil {
						return nil, err
					}

					if resp.StatusCode == http.StatusOK {
						v, err = ioutil.ReadAll(resp.Body)
						if !assert.Nil(t, err) {
							t.Fatal("failed to setup test")
						}

						var respFlights []string
						err = json.Unmarshal(v, &respFlights)
						if !assert.Nil(t, err) {
							t.Fatal("failed to setup test")
						}

						return respFlights, nil

					}

					return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
				},
			},
			want{
				err: errors.New("unexpected status code: 404"),
			},
		},
		{
			"handles garbage data",
			input{
				requestFn: func() ([]string, error) {
					flights := "garbage"

					v, err := json.Marshal(&flights)
					if !assert.Nil(t, err) {
						t.Fatal("failed to setup test")
					}

					resp, err := http.Post("http://127.0.0.1:8080/calculate", "application/json", bytes.NewBuffer(v))
					if err != nil {
						return nil, err
					}

					if resp.StatusCode == http.StatusOK {
						v, err = ioutil.ReadAll(resp.Body)
						if !assert.Nil(t, err) {
							t.Fatal("failed to setup test")
						}

						var respFlights []string
						err = json.Unmarshal(v, &respFlights)
						if !assert.Nil(t, err) {
							t.Fatal("failed to setup test")
						}

						return respFlights, nil

					}

					return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
				},
			},
			want{
				err: errors.New("unexpected status code: 400"),
			},
		},
		{
			"handles illy formatted data",
			input{
				requestFn: func() ([]string, error) {
					flights := [][]string{
						{"IND", "EWR", "NYC"},
						{"SFO", "ATL"},
						{"GSO", "IND"},
						{"ATL", "GSO"},
					}

					v, err := json.Marshal(&flights)
					if !assert.Nil(t, err) {
						t.Fatal("failed to setup test")
					}

					resp, err := http.Post("http://127.0.0.1:8080/calculate", "application/json", bytes.NewBuffer(v))
					if err != nil {
						return nil, err
					}

					if resp.StatusCode == http.StatusOK {
						v, err = ioutil.ReadAll(resp.Body)
						if !assert.Nil(t, err) {
							t.Fatal("failed to setup test")
						}

						var respFlights []string
						err = json.Unmarshal(v, &respFlights)
						if !assert.Nil(t, err) {
							t.Fatal("failed to setup test")
						}

						return respFlights, nil

					}

					return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
				},
			},
			want{
				err: errors.New("unexpected status code: 500"),
			},
		},
		{
			"is successful",
			input{
				requestFn: func() ([]string, error) {
					flights := [][]string{
						{"IND", "EWR"},
						{"SFO", "ATL"},
						{"GSO", "IND"},
						{"ATL", "GSO"},
					}

					v, err := json.Marshal(&flights)
					if !assert.Nil(t, err) {
						t.Fatal("failed to setup test")
					}

					resp, err := http.Post("http://127.0.0.1:8080/calculate", "application/json", bytes.NewBuffer(v))
					if err != nil {
						return nil, err
					}

					if resp.StatusCode == http.StatusOK {
						v, err = ioutil.ReadAll(resp.Body)
						if !assert.Nil(t, err) {
							t.Fatal("failed to setup test")
						}

						var respFlights []string
						err = json.Unmarshal(v, &respFlights)
						if !assert.Nil(t, err) {
							t.Fatal("failed to setup test")
						}

						return respFlights, nil

					}

					return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
				},
			},
			want{
				flights: []string{"SFO", "EWR"},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := New(Config{
				Host:    "127.0.0.1",
				Port:    8080,
				Tracker: tracker.New(),
				Log:     zap.NewNop(),
			})

			go server.Start()

			flights, err := tt.input.requestFn()
			if tt.want.err != nil {
				if !assert.NotNil(t, err) {
					t.Fatalf("expected an error: %s", tt.want.err.Error())
				}

				assert.Contains(t, tt.want.err.Error(), err.Error())
			}

			assert.Equal(t, tt.want.flights, flights)
		})
	}
}
