package grpcapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/anypb"
	"io"
	"tenflix/lintang/order-aggregator-service/internal/entity"
	"tenflix/lintang/order-aggregator-service/pb"
	"tenflix/lintang/order-aggregator-service/pkg/grpc"
)

type OrderGrpcAPI struct {
	c *grpc.ServiceClient
}

func NewOrderGrpc(orderGrpc *grpc.ServiceClient) *OrderGrpcAPI {
	return &OrderGrpcAPI{orderGrpc}
}

func (o *OrderGrpcAPI) CreateOrder(ctx context.Context, c entity.CreateOrderRequest, plan pb.PlanDto, userId string) (entity.Order, error) {
	orderClient := o.c.OrderClient

	createOrderGrpcRequest := &pb.CreateOrderGrpcRequest{
		Order: &pb.OrderDto{
			Id:              uuid.NewString(),
			UserId:          userId,
			Price:           plan.Price,
			OrderStatus:     pb.OrderStatus_PENDING,
			PaymentId:       uuid.NewString(),
			FailureMessages: "",
			Plan: &pb.OrderPlanDto{
				Id:       uuid.NewString(),
				PlanId:   plan.PlanId,
				Price:    plan.Price,
				SubTotal: plan.Price,
			},
		},
	}
	res, err := orderClient.CreateOrder(context.Background(), createOrderGrpcRequest)
	if err != nil {
		return entity.Order{}, fmt.Errorf("OrderGrpcAPI - CreateOrder - orderClient.CreateOrder: %w", err)
	}

	order := entity.Order{
		Id:              res.CreatedOrder.Id,
		UserId:          res.CreatedOrder.UserId,
		Price:           int64(res.CreatedOrder.Price),
		OrderStatus:     res.CreatedOrder.OrderStatus.String(),
		PaymentId:       res.CreatedOrder.PaymentId,
		FailureMessages: res.CreatedOrder.FailureMessages,
		Plan: entity.OrderPlan{
			Id:          res.CreatedOrder.Id,
			Name:        plan.Name,
			Description: plan.Description,
			PlanId:      int64(res.CreatedOrder.Plan.PlanId),
			Price:       int64(plan.Price),
			Subtotal:    int64(plan.Price),
		},
	}

	return order, nil
}

func (o *OrderGrpcAPI) ProcessOrderGrpc(ctx context.Context, notificationRes map[string]interface{}) error {
	orderClient := o.c.OrderClient
	var notifProto = map[string]*anypb.Any{}
	for key, val := range notificationRes {
		byteSliceVal, _ := json.Marshal(val)
		notifProto[key] = &anypb.Any{Value: byteSliceVal}
	}
	_, err := orderClient.ProcessOrderSaga(context.Background(), &pb.ProcessOrderRequest{
		PaymentNotification: &pb.PaymentNotification{
			NotificationRes: notifProto,
		},
	})
	if err != nil {
		return fmt.Errorf("OrderGrpcAPI - ProcessOrderGrpc - orderClient.ProcessOrderSaga: %w", err)
	}

	return nil
}

func (o *OrderGrpcAPI) GetUserOrderDetail(ctx context.Context, orderId string, userId string) (entity.Order, error) {
	orderClient := o.c.OrderClient

	res, err := orderClient.GetUserOrderDetail(context.Background(), &pb.GetUserOrderDetailRequest{
		OrderId: orderId,
		UserId:  userId,
	})

	if err != nil {
		return entity.Order{}, fmt.Errorf("OrderGrpcAPI - GetUserOrderDetail - orderClient.GetUserOrderDetail: %w", err)
	}
	order := entity.Order{
		Id:              res.OrderDto.Id,
		UserId:          res.OrderDto.UserId,
		Price:           int64(res.OrderDto.Price),
		OrderStatus:     res.OrderDto.OrderStatus.String(),
		PaymentId:       res.OrderDto.PaymentId,
		FailureMessages: res.OrderDto.FailureMessages,
		Plan: entity.OrderPlan{
			Id:       res.OrderDto.Id,
			PlanId:   int64(res.OrderDto.Plan.PlanId),
			Price:    int64(res.OrderDto.Price),
			Subtotal: int64(res.OrderDto.Price),
		},
	}
	return order, nil
}

func (o *OrderGrpcAPI) GetUserOrderHistory(ctx context.Context, userId string) ([]entity.Order, error) {
	orderClient := o.c.OrderClient

	stream, err := orderClient.GetUserOrderHistory(context.Background(), &pb.GetUserOrderHistoryRequest{
		UserId: userId,
	})

	if err != nil {
		return []entity.Order{}, fmt.Errorf("OrderGrpcAPI - GetUserOrderHistory - orderClient.GetUserOrderHistory: %w", err)
	}
	var orders []entity.Order
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			return orders, nil
		}
		if err != nil {
			return []entity.Order{}, fmt.Errorf("OrderGrpcAPI - GetUserOrderHistory - orderClient.GetUserOrderHistory (stream): %w", err)
		}
		order := entity.Order{
			Id:              res.OrderDto.Id,
			UserId:          res.OrderDto.UserId,
			Price:           int64(res.OrderDto.Price),
			OrderStatus:     res.OrderDto.OrderStatus.String(),
			PaymentId:       res.OrderDto.PaymentId,
			FailureMessages: res.OrderDto.FailureMessages,
			Plan: entity.OrderPlan{
				Id:       res.OrderDto.Id,
				PlanId:   int64(res.OrderDto.Plan.PlanId),
				Price:    int64(res.OrderDto.Price),
				Subtotal: int64(res.OrderDto.Price),
			},
		}
		_ = append(orders, order)

	}

}