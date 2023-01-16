package handler

import (
	"context"
	common "github.com/bufengmobuganhuo/micro-service-common"
	"github.com/bufengmobuganhuo/order/domain/model"
	"github.com/bufengmobuganhuo/order/domain/service"
	order "github.com/bufengmobuganhuo/order/proto/order"
)

type Order struct {
	OrderDataService service.IOrderDataService
}

func (o Order) GetOrderByID(ctx context.Context, req *order.OrderID, resp *order.OrderInfo) error {
	order, err := o.OrderDataService.FindOrderByID(req.OrderId)
	if err != nil {
		return err
	}
	if err := common.SwapTo(order, resp); err != nil {
		return err
	}
	return nil
}

func (o Order) GetAllOrder(ctx context.Context, req *order.AllOrderRequest, resp *order.AllOrder) error {
	orders, err := o.OrderDataService.FindAllOrder()
	if err != nil {
		return err
	}
	for _, v := range orders {
		order := &order.OrderInfo{}
		if err := common.SwapTo(v, orders); err != nil {
			return err
		}
		resp.OrderInfo = append(resp.OrderInfo, order)
	}
	return nil
}

func (o Order) CreateOrder(ctx context.Context, req *order.OrderInfo, resp *order.OrderID) error {
	orderAdd := &model.Order{}
	if err := common.SwapTo(req, orderAdd); err != nil {
		return err
	}
	orderID, err := o.OrderDataService.AddOrder(orderAdd)
	if err != nil {
		return err
	}
	resp.OrderId = orderID
	return nil
}

func (o Order) DeleteOrderByID(ctx context.Context, req *order.OrderID, response *order.Response) error {
	return o.OrderDataService.DeleteOrder(req.OrderId)
}

func (o Order) UpdateOrderPayStatus(ctx context.Context, req *order.PayStatus, response *order.Response) error {
	err := o.OrderDataService.UpdatePayStatus(req.OrderId, req.PayStatus)
	if err != nil {
		return err
	}
	response.Msg = "支付状态更新成功"
	return nil
}

func (o Order) UpdateOrderShipStatus(ctx context.Context, req *order.ShipStatus, resp *order.Response) error {
	err := o.OrderDataService.UpdateShipStatus(req.OrderId, req.ShipStatus)
	if err != nil {
		return err
	}
	resp.Msg = "发货状态更新成功"
	return nil
}

func (o Order) UpdateOrder(ctx context.Context, req *order.OrderInfo, resp *order.Response) error {
	order := &model.Order{}
	if err := common.SwapTo(order, req); err != nil {
		return err
	}
	if err := o.OrderDataService.UpdateOrder(order); err != nil {
		return err
	}
	resp.Msg = "订单更新成功"
	return nil
}
