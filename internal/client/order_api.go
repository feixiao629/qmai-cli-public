package client

import (
	"context"
	"encoding/json"
	"fmt"
)

// OrderAPI provides order-related API methods.
type OrderAPI struct {
	client *Client
}

// NewOrderAPI creates a new OrderAPI.
func NewOrderAPI(c *Client) *OrderAPI {
	return &OrderAPI{client: c}
}

type UserOrderItem struct {
	BizID               string `json:"bizId"`
	ActualAmount        int    `json:"actualAmount"`
	BizType             int    `json:"bizType"`
	CreatedAt           string `json:"createdAt"`
	OrderNo             string `json:"orderNo"`
	OrderType           int    `json:"orderType"`
	OrderSubType        int    `json:"orderSubType"`
	PayAt               string `json:"payAt"`
	PayStatus           int    `json:"payStatus"`
	PerformancePlatform string `json:"performancePlatform"`
	PerformanceType     int    `json:"performanceType"`
	ShopCode            string `json:"shopCode"`
	Source              int    `json:"source"`
	SourceNo            string `json:"sourceNo"`
	Status              int    `json:"status"`
	StoreID             int64  `json:"storeId"`
	TotalAmount         int    `json:"totalAmount"`
	UpdatedAt           string `json:"updatedAt"`
	UserID              int64  `json:"userId"`
}

type UserOrderListResponse struct {
	Data  []UserOrderItem `json:"data"`
	Total int             `json:"total"`
}

type OrderDetailItem struct {
	ItemName   string `json:"itemName"`
	ItemSign   string `json:"itemSign"`
	ItemType   int    `json:"itemType"`
	ItemUnit   string `json:"itemUnit"`
	Num        string `json:"num"`
	ItemPrice  int    `json:"itemPrice"`
	PackAmount int    `json:"packAmount"`
}

type OrderDiscountItem struct {
	DiscountName    string `json:"discountName"`
	DiscountType    int    `json:"discountType"`
	DiscountSummary string `json:"discountSummary"`
	DiscountLevel   int    `json:"discountLevel"`
	DiscountAmount  int    `json:"discountAmount"`
}

type OrderPayItem struct {
	PayType   int `json:"payType"`
	PayAmount int `json:"payAmount"`
}

type OrderDetailResponse struct {
	OrderType      int                 `json:"orderType"`
	Source         int                 `json:"source"`
	OrderNo        string              `json:"orderNo"`
	PayNo          string              `json:"payNo"`
	ThirdPayNo     string              `json:"thirdPayNo"`
	StoreOrderNo   string              `json:"storeOrderNo"`
	CreatedAt      string              `json:"createdAt"`
	CompletedAt    string              `json:"completedAt"`
	TotalAmount    int                 `json:"totalAmount"`
	ActualAmount   int                 `json:"actualAmount"`
	DiscountAmount int                 `json:"discountAmount"`
	ItemAmount     int                 `json:"itemAmount"`
	PackAmount     int                 `json:"packAmount"`
	FreightAmount  int                 `json:"freightAmount"`
	BuyerRemarks   string              `json:"buyerRemarks"`
	ShopCode       string              `json:"shopCode"`
	ShopName       string              `json:"shopName"`
	ItemList       []OrderDetailItem   `json:"itemList"`
	PayList        []OrderPayItem      `json:"payList"`
	DiscountList   []OrderDiscountItem `json:"discountList"`
	ContactTel     string              `json:"contactTel"`
	ContactName    string              `json:"contactName"`
}

type RechargeOrderItem struct {
	OrderType       int    `json:"orderType"`
	OrderSubType    int    `json:"orderSubType"`
	PerformanceType int    `json:"performanceType"`
	Source          int    `json:"source"`
	OrderNo         string `json:"orderNo"`
	ShopName        string `json:"shopName"`
	ShopCode        string `json:"shopCode"`
	ShopID          int64  `json:"shopId"`
	CreatedAt       string `json:"createdAt"`
	CompletedAt     string `json:"completedAt"`
	Status          int    `json:"status"`
	TotalAmount     int    `json:"totalAmount"`
	ActualAmount    int    `json:"actualAmount"`
	DiscountAmount  int    `json:"discountAmount"`
	ItemAmount      int    `json:"itemAmount"`
	PackAmount      int    `json:"packAmount"`
	FreightAmount   int    `json:"freightAmount"`
}

type RechargeRefundOrderItem struct {
	OrderType    int    `json:"orderType"`
	OrderSubType int    `json:"orderSubType"`
	Source       int    `json:"source"`
	OrderNo      string `json:"orderNo"`
	RefundNo     string `json:"refundNo"`
	ShopName     string `json:"shopName"`
	ShopCode     string `json:"shopCode"`
	CreatedAt    string `json:"createdAt"`
	DealAt       string `json:"dealAt"`
	ApplyReason  string `json:"applyReason"`
	ApplyAt      string `json:"applyAt"`
	RefundStatus int    `json:"refundStatus"`
	RefundAmount int    `json:"refundAmount"`
}

type OrderStatusResponse struct {
	Status       string `json:"status"`
	StoreOrderNo string `json:"storeOrderNo"`
	UserID       int64  `json:"userId"`
}

type MemberOrderCheckResponse int

type OrderUploadResponse struct {
	OrderNo string `json:"orderNo"`
}

type RefundOrderUpResponse struct {
	OrderID       int64  `json:"orderId"`
	OrderNo       string `json:"orderNo"`
	RefundOrderNo string `json:"refundOrderNo"`
}

type ProductionRecordItem struct {
	BizType           int                    `json:"bizType"`
	ExtraMap          map[string]interface{} `json:"extraMap"`
	ID                int64                  `json:"id"`
	IsDeleted         int                    `json:"isDeleted"`
	Num               int                    `json:"num"`
	OperateAt         string                 `json:"operateAt"`
	OrderItemID       int64                  `json:"orderItemId"`
	OrderNo           string                 `json:"orderNo"`
	PdcType           int                    `json:"pdcType"`
	ProductionAtEnd   string                 `json:"productionAtEnd"`
	ProductionAtStart string                 `json:"productionAtStart"`
	ProductionSign    string                 `json:"productionSign"`
	SellerID          int64                  `json:"sellerId"`
	SkuID             string                 `json:"skuId"`
	SpuID             int64                  `json:"spuId"`
	Status            int                    `json:"status"`
	StoreID           int64                  `json:"storeId"`
	SuccessNum        int                    `json:"successNum"`
	UserID            int64                  `json:"userId"`
}

type BatchOrderUploadResponse struct {
	Error        int      `json:"error"`
	ErrorOrderNo []string `json:"errorOrderNo"`
	Success      int      `json:"success"`
}

func (a *OrderAPI) GetUserOrderList(ctx context.Context, params map[string]interface{}) (*UserOrderListResponse, error) {
	resp, err := a.client.Call(ctx, "v3/order/getUserOrderList", params)
	if err != nil {
		return nil, fmt.Errorf("get user order list: %w", err)
	}
	return decodeData[UserOrderListResponse](resp, "user order list")
}

func (a *OrderAPI) GetOrderDetail(ctx context.Context, bizType int, orderNo string, userID int64) (*OrderDetailResponse, error) {
	params := map[string]interface{}{"bizType": bizType, "orderNo": orderNo}
	if userID != 0 {
		params["userId"] = userID
	}
	resp, err := a.client.Call(ctx, "v3/order/getDetail", params)
	if err != nil {
		return nil, fmt.Errorf("get order detail: %w", err)
	}
	return decodeData[OrderDetailResponse](resp, "order detail")
}

func (a *OrderAPI) ListRechargeOrders(ctx context.Context, params map[string]interface{}) (*[]RechargeOrderItem, error) {
	resp, err := a.client.Call(ctx, "v3/order/standard/recharge/order", params)
	if err != nil {
		return nil, fmt.Errorf("list recharge orders: %w", err)
	}
	return decodeData[[]RechargeOrderItem](resp, "recharge orders")
}

func (a *OrderAPI) ListRechargeRefundOrders(ctx context.Context, params map[string]interface{}) (*[]RechargeRefundOrderItem, error) {
	resp, err := a.client.Call(ctx, "v3/order/standard/recharge/refundOrder", params)
	if err != nil {
		return nil, fmt.Errorf("list recharge refund orders: %w", err)
	}
	return decodeData[[]RechargeRefundOrderItem](resp, "recharge refund orders")
}

func (a *OrderAPI) GetOrderStatus(ctx context.Context, orderNo string) (*OrderStatusResponse, error) {
	resp, err := a.client.Call(ctx, "v3/order/status", map[string]string{"orderNo": orderNo})
	if err != nil {
		return nil, fmt.Errorf("get order status: %w", err)
	}
	return decodeData[OrderStatusResponse](resp, "order status")
}

func (a *OrderAPI) CheckMemberOrder(ctx context.Context, params map[string]interface{}) (*MemberOrderCheckResponse, error) {
	resp, err := a.client.Call(ctx, "v3/order/checkMemberOrder", params)
	if err != nil {
		return nil, fmt.Errorf("check member order: %w", err)
	}
	var result MemberOrderCheckResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("parse member order check: %w", err)
	}
	return &result, nil
}

func (a *OrderAPI) ReplyUserComment(ctx context.Context, orderNo, replyAt, sellerReplyInfo string) error {
	_, err := a.client.Call(ctx, "v3/order/comment/storeReplyUserComment", map[string]string{
		"orderNo":         orderNo,
		"replyAt":         replyAt,
		"sellerReplyInfo": sellerReplyInfo,
	})
	if err != nil {
		return fmt.Errorf("reply user comment: %w", err)
	}
	return nil
}

func (a *OrderAPI) OrderUpload(ctx context.Context, params map[string]interface{}) (*OrderUploadResponse, error) {
	resp, err := a.client.Call(ctx, "v3/bsns/order/orderUpload", params)
	if err != nil {
		return nil, fmt.Errorf("order upload: %w", err)
	}
	return decodeData[OrderUploadResponse](resp, "order upload")
}

func (a *OrderAPI) RefundOrderUp(ctx context.Context, params map[string]interface{}) (*RefundOrderUpResponse, error) {
	resp, err := a.client.Call(ctx, "v3/cy/order/refundOrderUp", params)
	if err != nil {
		return nil, fmt.Errorf("refund order up: %w", err)
	}
	return decodeData[RefundOrderUpResponse](resp, "refund order up")
}

func (a *OrderAPI) OrderBatchUpload(ctx context.Context, params map[string]interface{}) (*BatchOrderUploadResponse, error) {
	resp, err := a.client.Call(ctx, "v3/bsns/order/orderBatchUpload", params)
	if err != nil {
		return nil, fmt.Errorf("order batch upload: %w", err)
	}
	return decodeData[BatchOrderUploadResponse](resp, "order batch upload")
}

func (a *OrderAPI) RefundOrderBatchUpload(ctx context.Context, params map[string]interface{}) (*BatchOrderUploadResponse, error) {
	resp, err := a.client.Call(ctx, "v3/bsns/order/refundOrderBatchUpload", params)
	if err != nil {
		return nil, fmt.Errorf("refund order batch upload: %w", err)
	}
	return decodeData[BatchOrderUploadResponse](resp, "refund order batch upload")
}

func (a *OrderAPI) ListProductionRecords(ctx context.Context, params map[string]interface{}) (*[]ProductionRecordItem, error) {
	resp, err := a.client.Call(ctx, "v3/newPattern/orderCenter/post/order/item/production/record/list", params)
	if err != nil {
		return nil, fmt.Errorf("list production records: %w", err)
	}
	return decodeData[[]ProductionRecordItem](resp, "production records")
}
