package tracker

import (
	"errors"
)

// Flights represents an unordered itinerary of person's flights including any layovers.
// Example: [["IND", "EWR"], ["SFO", "ATL"], ["GSO", "IND"], ["ATL", "GSO"]]
type Flights [][]string

func (f Flights) validate() error {
	if len(f) < 1 {
		return errors.New("failed to find any flights to track")
	}

	for _, flight := range f {
		if len(flight) != 2 {
			return errors.New("invalid flight data format")
		}
	}

	return nil
}

// Tracker is a flight tracker that exposes funtions that operate over a person's Flights data.
type Tracker struct{}

func New() *Tracker {
	return &Tracker{}
}

type flightMap map[currentAirport]flightHops

type currentAirport string

type flightHops struct {
	LastHop string
	NextHop string
}

// GetSrcAndDst will take a list of Flights (unordered) and determine the source and destination flight
func (t *Tracker) GetSrcAndDst(flights Flights) (flight []string, err error) {
	if err := flights.validate(); err != nil {
		return nil, err
	}

	var firstFlight, lastFlight []string
	for key, val := range t.mapFlights(flights) {
		if val.NextHop == "" {
			lastFlight = append(lastFlight, string(key))
		}

		if val.LastHop == "" {
			firstFlight = append(firstFlight, string(key))
		}
	}

	if len(firstFlight) != 1 && len(lastFlight) != 1 {
		return nil, errors.New("failed to find connecting flights: broken flight chain")
	}

	return []string{firstFlight[0], lastFlight[0]}, nil
}

func (t *Tracker) mapFlights(flights Flights) flightMap {
	fm := flightMap{}
	for _, flight := range flights {
		fst := flight[0]
		lst := flight[1]

		flightPair, found := fm[currentAirport(lst)]
		if !found {
			fm[currentAirport(lst)] = flightHops{
				LastHop: fst,
			}
		} else {
			flightPair.LastHop = fst
			fm[currentAirport(lst)] = flightPair

		}

		flightPair1, found := fm[currentAirport(fst)]
		if !found {
			fm[currentAirport(fst)] = flightHops{
				NextHop: lst,
			}
		} else {
			flightPair1.NextHop = lst
			fm[currentAirport(fst)] = flightPair1
		}
	}

	return fm
}
