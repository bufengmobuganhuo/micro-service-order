package repository

import (
	"errors"
	"github.com/bufengmobuganhuo/order/domain/model"
	"github.com/jinzhu/gorm"
)

type IOrderRepository interface {
	InitTable() error
	FindOrderByID(int64) (*model.Order, error)
	CreateOrder(*model.Order) (int64, error)
	DeleteOrderByID(int64) error
	UpdateOrder(*model.Order) error
	FindAll() ([]model.Order, error)
	UpdateShipStatus(int64, int32) error
	UpdatePayStatus(int64, int32) error
}

// NewOrderRepository 创建orderRepository
func NewOrderRepository(db *gorm.DB) IOrderRepository {
	return &OrderRepository{mysqlDb: db}
}

type OrderRepository struct {
	mysqlDb *gorm.DB
}

func (u *OrderRepository) UpdateShipStatus(orderId int64, shipStatus int32) error {
	db := u.mysqlDb.Model(&model.Order{}).Where("id = ?", orderId).UpdateColumn("ship_status", shipStatus)
	if db.Error != nil {
		return db.Error
	}
	if db.RowsAffected == 0 {
		return errors.New("更新失败")
	}
	return nil
}

func (u *OrderRepository) UpdatePayStatus(orderId int64, payStatus int32) error {
	db := u.mysqlDb.Model(&model.Order{}).Where("id = ?", orderId).UpdateColumn("pay_status", payStatus)
	if db.Error != nil {
		return db.Error
	}
	if db.RowsAffected == 0 {
		return errors.New("更新失败")
	}
	return nil
}

// InitTable 初始化表
func (u *OrderRepository) InitTable() error {
	return u.mysqlDb.CreateTable(&model.Order{}).Error
}

// FindOrderByID 根据ID查找Order信息
func (u *OrderRepository) FindOrderByID(orderID int64) (order *model.Order, err error) {
	order = &model.Order{}
	return order, u.mysqlDb.Preloads("OrderDetail").First(order, orderID).Error
}

// CreateOrder 创建Order信息
func (u *OrderRepository) CreateOrder(order *model.Order) (int64, error) {
	return order.ID, u.mysqlDb.Create(order).Error
}

// DeleteOrderByID 根据ID删除Order信息
func (u *OrderRepository) DeleteOrderByID(orderID int64) error {
	tx := u.mysqlDb.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}

	// 删除order信息
	if err := tx.Unscoped().Where("id = ?", orderID).Delete(&model.Order{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除order detail信息
	if err := tx.Unscoped().Where("order_id = ?", orderID).Delete(&model.OrderDetail{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// UpdateOrder 更新Order信息
func (u *OrderRepository) UpdateOrder(order *model.Order) error {
	return u.mysqlDb.Model(order).Update(order).Error
}

// FindAll 获取结果集
func (u *OrderRepository) FindAll() (orderAll []model.Order, err error) {
	return orderAll, u.mysqlDb.Find(&orderAll).Error
}
