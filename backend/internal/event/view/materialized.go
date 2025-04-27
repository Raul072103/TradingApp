package view

import (
	"TradingSimulation/backend/internal/event"
	"errors"
)

var (
	ErrUnknownEventView            = errors.New("unknown event")
	ErrUnknownOrderType            = errors.New("unknown order type")
	ErrNoAccountWithID             = errors.New("there is no registered account with the given ID")
	ErrCannotCancelAnUnplacedOrder = errors.New("trying to cancel an order that wasn't placed")
	ErrTradeHasUnknownUsers        = errors.New("trade has unknwon users")
)

type Stock struct {
	ID    int64   `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type AccountState struct {
	BuyOrders        []event.Order         `json:"buy_orders"`
	SellOrders       []event.Order         `json:"sell_orders"`
	CanceledOrders   map[int64]event.Order `json:"canceled_orders"`
	SuccessfulOrders []event.Order         `json:"successful_orders"`

	Funds float64 `json:"funds"`

	Events []event.Event `json:"events"` // Events represent a slice where the latest events are pushed on it
}

type MaterializedView struct {
	Accounts map[int64]AccountState // Accounts map of all the accounts accessible by each AccountID
	Trades   []event.Trade          `json:"trades"` // Trades represent a slice with all the executed trades
	Orders   []event.Order          `json:"orders"` // Orders represent a slice with all the orders placed
	Stocks   map[int64]Stock        `json:"stocks"` // Stocks represents the current stocks traded by this app
}

// New creates a new Materialized View with all the events that were recorded till that point.
func New(events []event.Event) (MaterializedView, error) {
	var materializedView MaterializedView
	materializedView.Accounts = make(map[int64]AccountState)
	materializedView.Trades = make([]event.Trade, 0)
	materializedView.Orders = make([]event.Order, 0)
	materializedView.Stocks = map[int64]Stock{
		0: {
			ID:    0,
			Name:  "NVDA",
			Price: 111.01,
		},
		1: {
			ID:    1,
			Name:  "TSLA",
			Price: 284.95,
		},
		2: {
			ID:    2,
			Name:  "PLTR",
			Price: 112.78,
		},
		3: {
			ID:    3,
			Name:  "GOOGL",
			Price: 161.96,
		},
		4: {
			ID:    4,
			Name:  "AMZN",
			Price: 188.99,
		},
		5: {
			ID:    5,
			Name:  "CRWD",
			Price: 422.88,
		},
	}

	for _, currEvent := range events {
		err := materializedView.handleEvent(currEvent)
		if err != nil {
			return materializedView, err
		}
	}

	return materializedView, nil
}

// RegisterEvent updates the MaterializedView with a new event.
func (view *MaterializedView) RegisterEvent(event event.Event) error {
	err := view.handleEvent(event)
	return err
}

// GetAccount return the AccountState corresponding to the given ID if that account exists, otherwise an error.
func (view *MaterializedView) GetAccount(accountID int64) (AccountState, error) {
	accountState, exists := view.Accounts[accountID]
	if !exists {
		return accountState, ErrNoAccountWithID
	}

	return accountState, nil
}

func (view *MaterializedView) handleEvent(currEvent event.Event) error {
	switch currEvent.Type() {

	case event.OrdersCanceledEvent:
		orderCanceled := currEvent.(*event.OrderCanceled)
		accountID := orderCanceled.AccountID
		accountState, exists := view.Accounts[accountID]

		if !exists {
			return ErrCannotCancelAnUnplacedOrder
		}

		accountState.CanceledOrders[orderCanceled.Order.ID] = orderCanceled.Order
		accountState.Events = append(accountState.Events, orderCanceled)

		// TODO() delete events related to the order canceled

		view.Orders = append(view.Orders, orderCanceled.Order)
		view.Accounts[accountID] = accountState

	case event.OrdersPlacedEvent:
		orderPlaced := currEvent.(*event.OrderPlaced)
		accountID := orderPlaced.AccountID
		accountState, exists := view.Accounts[accountID]

		if !exists {
			accountState = initializeAccountState()
		}

		switch orderPlaced.Order.Type {
		case event.BuyOrder:
			accountState.BuyOrders = append(accountState.BuyOrders, orderPlaced.Order)
		case event.SellOrder:
			accountState.SellOrders = append(accountState.SellOrders, orderPlaced.Order)
		default:
			return ErrUnknownOrderType
		}

		accountState.Events = append(accountState.Events, orderPlaced)

		view.Accounts[accountID] = accountState
		view.Orders = append(view.Orders, orderPlaced.Order)

	case event.FundsCreditedEvent:
		fundsCredited := currEvent.(*event.FundsCredited)
		accountID := fundsCredited.AccountID
		accountState, exists := view.Accounts[accountID]

		if !exists {
			accountState = initializeAccountState()
		}

		accountState.Funds += fundsCredited.Sum
		accountState.Events = append(accountState.Events, fundsCredited)

		view.Accounts[accountID] = accountState

	case event.FundsDebitedEvent:
		fundsDebited := currEvent.(*event.FundsDebited)
		accountID := fundsDebited.AccountID
		accountState, exists := view.Accounts[accountID]

		if !exists {
			accountState = initializeAccountState()
		}

		accountState.Funds -= fundsDebited.Sum
		accountState.Events = append(accountState.Events, fundsDebited)

		view.Accounts[accountID] = accountState

	case event.TradeExecutedEvent:
		tradeExecuted := currEvent.(*event.TradeExecuted)

		accountID1 := tradeExecuted.Trade.AccountIDs[0]
		accountID2 := tradeExecuted.Trade.AccountIDs[1]

		order1 := tradeExecuted.Trade.Orders[0]
		order2 := tradeExecuted.Trade.Orders[1]

		accountState1, exists1 := view.Accounts[accountID1]
		accountState2, exists2 := view.Accounts[accountID2]

		if !exists1 || !exists2 {
			return ErrTradeHasUnknownUsers
		}

		accountState1.SuccessfulOrders = append(accountState1.SuccessfulOrders, order1)
		accountState1.Events = append(accountState1.Events, tradeExecuted)

		accountState2.SuccessfulOrders = append(accountState2.SuccessfulOrders, order2)
		accountState2.Events = append(accountState2.Events, tradeExecuted)

		view.Accounts[accountID1] = accountState1
		view.Accounts[accountID2] = accountState2

		view.Trades = append(view.Trades, tradeExecuted.Trade)

	default:
		return ErrUnknownEventView
	}

	return nil
}

func initializeAccountState() AccountState {
	return AccountState{
		BuyOrders:        make([]event.Order, 0),
		SellOrders:       make([]event.Order, 0),
		CanceledOrders:   make(map[int64]event.Order),
		SuccessfulOrders: make([]event.Order, 0),
		Funds:            0,
		Events:           make([]event.Event, 0),
	}
}
