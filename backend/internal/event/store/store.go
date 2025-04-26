package store

import (
	"TradingSimulation/common"
	"bufio"
	"encoding/json"
	"errors"
	"os"
)

var (
	ErrEventFormatWrong = errors.New("the line format isn't separated by commas")
	ErrUnknownEvent     = errors.New("the evnet number is unknown")
)

type Store struct {
	file *os.File
	path string
}

// New creates an instance of Store struct an initializes the file based events store.
func New() (*Store, error) {
	var store Store
	store.path = "./container/events.log"

	err := store.createStore()
	return &store, err
}

// Close frees the resources of the Store struct.
func (s *Store) Close() error {
	err := s.file.Close()
	return err
}

// AppendEvent appends an event to the log file of the vent store.
func (s *Store) AppendEvent(event common.Event) error {
	eventStr, err := s.marshalEvent(event)
	if err != nil {
		return err
	}

	_, err = s.file.WriteString(eventStr)
	if err != nil {
		return err
	}
	return nil
}

// GetAllEvents returns a list with all the events saved in the store logs
func (s *Store) GetAllEvents() ([]common.Event, error) {
	var events []common.Event

	scanner := bufio.NewScanner(s.file)

	// read line by line from the file
	var lines [][]byte
	for scanner.Scan() {
		lines = append(lines, scanner.Bytes())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	for _, line := range lines {
		event, err := s.unmarshalEvent(line)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

// createStore creates the file that stores all the logs, only if that file doesn't exist at the moment of calling
// this method.
func (s *Store) createStore() error {
	file, err := os.OpenFile(s.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	s.file = file
	return nil
}

// marshalEvnet helper function to write an event as a string.
func (s *Store) marshalEvent(event common.Event) (string, error) {
	data, err := json.Marshal(event)
	return string(data), err
}

// unmarshalEvent helper function to parse a line of bytes read by the store.
func (s *Store) unmarshalEvent(data []byte) (common.Event, error) {
	eventType, err := common.UnmarshalEventTypeJSON(data)
	if err != nil {
		return nil, err
	}

	switch eventType {

	case common.OrdersCanceledEvent:
		var ordersCanceled common.OrderCanceled
		err := unmarshalEvent(data, &ordersCanceled)
		return &ordersCanceled, err

	case common.OrdersPlacedEvent:
		var ordersPlaced common.OrderPlaced
		err := unmarshalEvent(data, &ordersPlaced)
		return &ordersPlaced, err

	case common.FundsCreditedEvent:
		var fundsCredited common.FundsCredited
		err := unmarshalEvent(data, &fundsCredited)
		return &fundsCredited, err

	case common.FundsDebitedEvent:
		var fundsDebited common.FundsDebited
		err := unmarshalEvent(data, &fundsDebited)
		return &fundsDebited, err

	case common.TradeExecutedEvent:
		var tradeExecuted common.TradeExecuted
		err := unmarshalEvent(data, &tradeExecuted)
		return &tradeExecuted, err

	default:
		return nil, ErrUnknownEvent

	}
}

func unmarshalEvent(data []byte, event common.Event) error {
	err := json.Unmarshal(data, event)
	return err
}

func (s *Store) writeEvent() {

}
