package store

import (
	"TradingSimulation/backend/internal/event"
	"bufio"
	"encoding/json"
	"errors"
	"os"
)

var (
	ErrEventFormatWrong = errors.New("the line format isn't separated by commas")
	ErrUnknownEvent     = errors.New("the evnet number is unknown")
)

type logStruct struct {
	Event     any   `json:"event"`
	EventType int64 `json:"type"`
}

type Store struct {
	file *os.File
	path string
}

// New creates an instance of Store struct an initializes the file based events store.
func New() (*Store, error) {
	var store Store
	store.path = "./events.log"

	err := store.createStore()
	return &store, err
}

// Close frees the resources of the Store struct.
func (s *Store) Close() error {
	err := s.file.Close()
	return err
}

// AppendEvent appends an event to the log file of the vent store.
func (s *Store) AppendEvent(event event.Event) error {
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
func (s *Store) GetAllEvents() ([]event.Event, error) {
	var events []event.Event

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
		unmarshalledEvent, err := s.unmarshalEvent(line)
		if err != nil {
			return nil, err
		}
		events = append(events, unmarshalledEvent)
	}

	return events, nil
}

// createStore creates the file that stores all the logs, only if that file doesn't exist at the moment of calling
// this method.
func (s *Store) createStore() error {
	file, err := os.OpenFile(s.path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return err
	}

	s.file = file
	return nil
}

// marshalEvnet helper function to write an event as a string.
func (s *Store) marshalEvent(event event.Event) (string, error) {
	logStruct := logStruct{
		Event:     event,
		EventType: event.Type(),
	}
	data, err := json.Marshal(logStruct)
	return string(data), err
}

// unmarshalEvent helper function to parse a line of bytes read by the store.
func (s *Store) unmarshalEvent(data []byte) (event.Event, error) {
	var logStruct logStruct
	err := unmarshalLogStruct(data, &logStruct)
	if err != nil {
		return nil, err
	}

	// Marshal the Event (map[string]any) back to JSON
	eventBytes, err := json.Marshal(logStruct.Event)
	if err != nil {
		return nil, err
	}

	switch logStruct.EventType {
	case event.OrdersCanceledEvent:
		var ordersCanceled event.OrderCanceled
		err = json.Unmarshal(eventBytes, &ordersCanceled)
		return &ordersCanceled, err

	case event.OrdersPlacedEvent:
		var ordersPlaced event.OrderPlaced
		err = json.Unmarshal(eventBytes, &ordersPlaced)
		return &ordersPlaced, err

	case event.FundsCreditedEvent:
		var fundsCredited event.FundsCredited
		err = json.Unmarshal(eventBytes, &fundsCredited)
		return &fundsCredited, err

	case event.FundsDebitedEvent:
		var fundsDebited event.FundsDebited
		err = json.Unmarshal(eventBytes, &fundsDebited)
		return &fundsDebited, err

	case event.TradeExecutedEvent:
		var tradeExecuted event.TradeExecuted
		err = json.Unmarshal(eventBytes, &tradeExecuted)
		return &tradeExecuted, err

	default:
		return nil, ErrUnknownEvent
	}
}

func unmarshalLogStruct(data []byte, logStruct *logStruct) error {
	err := json.Unmarshal(data, logStruct)
	return err
}

func (s *Store) writeEvent() {

}
