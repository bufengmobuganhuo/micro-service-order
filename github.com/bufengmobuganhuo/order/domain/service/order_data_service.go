package service

import (
	"github.com/bufengmobuganhuo/order/domain/model"
	"github.com/bufengmobuganhuo/order/domain/repository"
)

type IOrderDataService interface {
	AddOrder(*model.Order) (int64, error)
	DeleteOrder(int64) error
	UpdateOrder(*model.Order) error
	FindOrderByID(int64) (*model.Order, error)
	FindAllOrder() ([]model.Order, error)
	UpdateShipStatus(int64, int32) error
	UpdatePayStatus(int64, int32) error
}

//创建
func NewOrderDataService(orderRepository repository.IOrderRepository) IOrderDataService {
	return &OrderDataService{orderRepository}
}

type OrderDataService struct {
	OrderRepository repository.IOrderRepository
}

func (u *OrderDataService) UpdateShipStatus(orderId int64, shipStatus int32) error {
	return u.OrderRepository.UpdateShipStatus(orderId, shipStatus)
}

func (u *OrderDataService) UpdatePayStatus(orderId int64, payStatus int32) error {
	return u.OrderRepository.UpdatePayStatus(orderId, payStatus)
}

// AddOrder 插入
func (u *OrderDataService) AddOrder(order *model.Order) (int64, error) {
	return u.OrderRepository.CreateOrder(order)
}

// DeleteOrder 删除
func (u *OrderDataService) DeleteOrder(orderID int64) error {
	return u.OrderRepository.DeleteOrderByID(orderID)
}

// UpdateOrder 更新
func (u *OrderDataService) UpdateOrder(order *model.Order) error {
	return u.OrderRepository.UpdateOrder(order)
}

// FindOrderByID 查找
func (u *OrderDataService) FindOrderByID(orderID int64) (*model.Order, error) {
	return u.OrderRepository.FindOrderByID(orderID)
}

// FindAllOrder 查找
func (u *OrderDataService) FindAllOrder() ([]model.Order, error) {
	return u.OrderRepository.FindAll()
}
