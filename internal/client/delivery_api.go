package client

import (
	"context"
	"fmt"
)

// DeliveryAPI provides aggregated delivery API methods.
type DeliveryAPI struct {
	client *Client
}

// NewDeliveryAPI creates a new DeliveryAPI.
func NewDeliveryAPI(c *Client) *DeliveryAPI {
	return &DeliveryAPI{client: c}
}

type DeliveryOrderLog struct {
	LogStr  string `json:"logStr"`
	LogTime string `json:"logTime"`
}

type DeliveryOrderInfoResponse struct {
	BizID             string             `json:"bizId"`
	CancelDesc        string             `json:"cancelDesc"`
	CancelReasonCode  string             `json:"cancelReasonCode"`
	DeliveryCost      string             `json:"deliveryCost"`
	DeliveryOrderLogs []DeliveryOrderLog `json:"deliveryOrderLogs"`
	DriverName        string             `json:"driverName"`
	DriverPhone       string             `json:"driverPhone"`
	ErrorDesc         string             `json:"errorDesc"`
	OrderStatus       string             `json:"orderStatus"`
	PushTime          string             `json:"pushTime"`
}

type RiderLocationResponse struct {
	CarrierName        string `json:"carrierName"`
	CarrierPhone       string `json:"carrierPhone"`
	DeliveryOrderNo    string `json:"deliveryOrderNo"`
	DeliveryStatus     string `json:"deliveryStatus"`
	Latitude           string `json:"latitude"`
	Longitude          string `json:"longitude"`
	UploadDeliveryType string `json:"uploadDeliveryType"`
}

type CancelAllDeliveryOrderItem struct {
	DeliveryID    string `json:"deliveryId"`
	OriginOrderNo string `json:"originOrderNo"`
	Result        bool   `json:"result"`
}

func (a *DeliveryAPI) CreateDeliveryOrder(ctx context.Context, params map[string]interface{}) error {
	_, err := a.client.Call(ctx, "v3/delivery/createDeliveryOrder", params)
	if err != nil {
		return fmt.Errorf("create delivery order: %w", err)
	}
	return nil
}

func (a *DeliveryAPI) CancelDeliveryOrder(ctx context.Context, params map[string]interface{}) error {
	_, err := a.client.Call(ctx, "v3/delivery/cancelDeliveryOrder", params)
	if err != nil {
		return fmt.Errorf("cancel delivery order: %w", err)
	}
	return nil
}

func (a *DeliveryAPI) GetDeliveryOrderInfo(ctx context.Context, params map[string]interface{}) (*DeliveryOrderInfoResponse, error) {
	resp, err := a.client.Call(ctx, "v3/delivery/getDeliveryOrderInfo", params)
	if err != nil {
		return nil, fmt.Errorf("get delivery order info: %w", err)
	}
	return decodeData[DeliveryOrderInfoResponse](resp, "delivery order info")
}

func (a *DeliveryAPI) GetRiderLocation(ctx context.Context, params map[string]interface{}) (*RiderLocationResponse, error) {
	resp, err := a.client.Call(ctx, "v3/delivery/getRiderLocation", params)
	if err != nil {
		return nil, fmt.Errorf("get rider location: %w", err)
	}
	return decodeData[RiderLocationResponse](resp, "rider location")
}

func (a *DeliveryAPI) UpdateSelfOrderStatus(ctx context.Context, params map[string]interface{}) error {
	_, err := a.client.Call(ctx, "v3/delivery/updateSelfOrderStatus", params)
	if err != nil {
		return fmt.Errorf("update self order status: %w", err)
	}
	return nil
}

func (a *DeliveryAPI) CancelAllDeliveryOrder(ctx context.Context, params map[string]interface{}) (*[]CancelAllDeliveryOrderItem, error) {
	resp, err := a.client.Call(ctx, "v3/delivery/cancelAllDeliveryOrder", params)
	if err != nil {
		return nil, fmt.Errorf("cancel all delivery orders: %w", err)
	}
	return decodeData[[]CancelAllDeliveryOrderItem](resp, "cancel all delivery orders")
}

func (a *DeliveryAPI) CreateAndUpdateDeliveryStatus(ctx context.Context, params map[string]interface{}) error {
	_, err := a.client.Call(ctx, "v3/delivery/updateDeliveryStatus", params)
	if err != nil {
		return fmt.Errorf("create and update delivery status: %w", err)
	}
	return nil
}
