package store

import (
	"TradingSimulation/common"
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
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

func (s *Store) AppendEvent(event common.Event) {

}

func (s *Store) GetAllEvents() ([]common.Event, error) {
	var events []common.Event

	scanner := bufio.NewScanner(s.file)

	// read line by line from the file
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return nil, nil
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

// parseEvent helper function to parse a line of bytes read by the store.
func (s *Store) parseEvent(line string) (common.Event, error) {
	splitStr := strings.Split(line, ",")

	if len(splitStr) <= 2 {
		return nil, ErrEventFormatWrong
	}

	eventType, err := strconv.ParseInt(splitStr[0], 10, 64)
	if err != nil {
		return nil, err
	}

	eventID, err := strconv.ParseInt(splitStr[1], 10, 64)
	if err != nil {
		return nil, err
	}

	accountID, err := strconv.ParseInt(splitStr[2], 10, 64)
	if err != nil {
		return nil, err
	}

	switch eventType {

	case common.OrdersCanceledEvent:
		order := common.Order{
			ID:        eventID,
			AccountID: accountID,
			Type:      eventType,
		}

		return nil, err
	case common.OrdersPlacedEvent:

	case common.FundsCreditedEvent:

	case common.FundsDebitedEvent:

	case common.TradeExecutedEvent:

	default:
		return nil, ErrUnknownEvent

	}

	return nil, nil
}

func (s *Store) writeEvent() {

}
