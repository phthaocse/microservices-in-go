package domain

import (
	"time"
)

type Order struct {
	ID         int64       `json:"id" bson:"_id"`
	CustomerID int64       `json:"customer_id" bson:"customer_id"`
	Status     string      `json:"status" bson:"status"`
	OrderItems []OrderItem `json:"order_items" bson:"order_items"`
	CreatedAt  int64       `json:"created_at" bson:"created_at"`
}

type OrderItem struct {
	ProductCode string  `json:"product_code" bson:"product_code"`
	UnitPrice   float32 `json:"unit_price" bson:"unit_price"`
	Quantity    int32   `json:"quantity" bson:"quantity"`
}

func NewOrder(customerId int64, orderItems []OrderItem) Order {
	return Order{
		CreatedAt:  time.Now().Unix(),
		Status:     "Pending",
		CustomerID: customerId,
		OrderItems: orderItems,
	}
}

func (o *Order) TotalPrice() float32 {
	var totalPrice float32
	for _, orderItem := range o.OrderItems {
		totalPrice += orderItem.UnitPrice * float32(orderItem.Quantity)
	}
	return totalPrice
}
