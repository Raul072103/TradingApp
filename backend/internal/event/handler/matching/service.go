package matching

import (
	"TradingSimulation/backend/internal/event"
	"TradingSimulation/backend/internal/event/view"
	"errors"
)

var (
	ErrUnknownOrderType = errors.New("the order can be either sell or buy")
)

type OrderBook struct {
	BuyOrders  map[int64][]event.Order
	SellOrders map[int64][]event.Order
}

type Service struct {
	MainChannel chan event.Event
	// MatchingChannel channel trough which buy or sell orders are received to be matched
	OrdersChannel         chan event.Order
	CancelOrdersChannel   chan event.Order
	Stocks                map[int64]view.Stock
	ActiveOrders          map[int64]OrderBook
	TradesExecutedCounter int64
}

func New(
	mainChannel chan event.Event,
	ordersChannel chan event.Order,
	cancelOrdersChannel chan event.Order,
	stocks map[int64]view.Stock) Service {

	var service Service

	service.MainChannel = mainChannel
	service.OrdersChannel = ordersChannel
	service.CancelOrdersChannel = cancelOrdersChannel
	service.Stocks = stocks
	service.ActiveOrders = make(map[int64]OrderBook)
	service.TradesExecutedCounter = 0

	for stockID := range stocks {
		service.ActiveOrders[stockID] = OrderBook{
			BuyOrders:  make(map[int64][]event.Order),
			SellOrders: make(map[int64][]event.Order),
		}
	}

	return service
}
func (service *Service) Run() error {
	for {
		select {
		case order := <-service.OrdersChannel:
			orderType := order.Type
			stockID := order.Stock
			count := order.Count

			orderBook := service.ActiveOrders[stockID]

			switch orderType {

			case event.BuyOrder:
				sellOrders, exists := orderBook.SellOrders[count]
				if exists && len(sellOrders) > 0 {
					sellOrder := sellOrders[len(sellOrders)-1]
					orderBook.SellOrders[count] = sellOrders[:len(sellOrders)-1]

					if len(orderBook.SellOrders[count]) == 0 {
						delete(orderBook.SellOrders, count)
					}

					tradeExecuted := event.TradeExecuted{
						Trade: event.Trade{
							ID:         service.TradesExecutedCounter,
							AccountIDs: []int64{order.AccountID, sellOrder.AccountID},
							Orders:     []event.Order{order, sellOrder},
							Status:     event.ActiveTrade,
						},
					}

					service.TradesExecutedCounter++
					service.MainChannel <- &tradeExecuted
				} else {
					orderBook.BuyOrders[count] = append(orderBook.BuyOrders[count], order)
				}

			case event.SellOrder:
				buyOrders, exists := orderBook.BuyOrders[count]
				if exists && len(buyOrders) > 0 {
					buyOrder := buyOrders[len(buyOrders)-1]
					orderBook.BuyOrders[count] = buyOrders[:len(buyOrders)-1]

					if len(orderBook.BuyOrders[count]) == 0 {
						delete(orderBook.BuyOrders, count)
					}

					tradeExecuted := event.TradeExecuted{
						Trade: event.Trade{
							ID:         service.TradesExecutedCounter,
							AccountIDs: []int64{buyOrder.AccountID, order.AccountID},
							Orders:     []event.Order{buyOrder, order},
							Status:     event.ActiveTrade,
						},
					}

					service.TradesExecutedCounter++
					service.MainChannel <- &tradeExecuted
				} else {
					orderBook.SellOrders[count] = append(orderBook.SellOrders[count], order)
				}

			default:
				return ErrUnknownOrderType
			}

			service.ActiveOrders[stockID] = orderBook

		case cancelOrder := <-service.CancelOrdersChannel:
			orderType := cancelOrder.Type
			stockID := cancelOrder.Stock
			count := cancelOrder.Count

			orderBook := service.ActiveOrders[stockID]

			switch orderType {
			case event.BuyOrder:
				orders := orderBook.BuyOrders[count]
				newOrders := removeOrderByID(orders, cancelOrder.ID)
				if len(newOrders) == 0 {
					delete(orderBook.BuyOrders, count)
				} else {
					orderBook.BuyOrders[count] = newOrders
				}

			case event.SellOrder:
				orders := orderBook.SellOrders[count]
				newOrders := removeOrderByID(orders, cancelOrder.ID)
				if len(newOrders) == 0 {
					delete(orderBook.SellOrders, count)
				} else {
					orderBook.SellOrders[count] = newOrders
				}

			default:
				return ErrUnknownOrderType
			}

			service.ActiveOrders[stockID] = orderBook
		}
	}
}

func removeOrderByID(orders []event.Order, id int64) []event.Order {
	for i, order := range orders {
		if order.ID == id {
			return append(orders[:i], orders[i+1:]...)
		}
	}
	return orders
}
