package model

import "time"

type Order struct {
	ID int64 `gorm:"primary_key;not_null;auto_increment"`
	// 用于实现幂等性
	OrderCode   string        `gorm:"unique_index;not_null",json:"order_code"`
	PayStatus   int32         `json:"pay_status"`
	ShipStatus  int32         `json:"ship_status"`
	Price       float64       `json:"price"`
	OrderDetail []OrderDetail `gorm:"ForeignKey:OrderID" json:"order_detail"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
