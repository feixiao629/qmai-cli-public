package client

import (
	"context"
	"encoding/json"
	"fmt"
)

// InventoryAPI provides SCM inventory APIs.
type InventoryAPI struct {
	client *Client
}

// NewInventoryAPI creates a new InventoryAPI.
func NewInventoryAPI(c *Client) *InventoryAPI {
	return &InventoryAPI{client: c}
}

type BusinessIDResponse int64
type BoolResultResponse bool
type UnitCodeResponse struct {
	UnitCode string `json:"unitCode"`
}

type DeclareProductItem struct {
	ProductCode     string  `json:"productCode"`
	ProductName     string  `json:"productName"`
	ProductSpec     string  `json:"productSpec"`
	ProductUnit     string  `json:"productUnit"`
	ProductNum      float64 `json:"productNum"`
	Price           float64 `json:"price"`
	DiscountPrice   float64 `json:"discountPrice"`
	RequireAmount   float64 `json:"requireAmount"`
	AuditAmount     float64 `json:"auditAmount"`
	ExamineNum      float64 `json:"examineNum"`
	TagName         string  `json:"tagName"`
	PerformanceName string  `json:"performanceName"`
}

type DeclareDetailResponse struct {
	DeclareNo          string               `json:"declareNo"`
	CreatedAt          string               `json:"createdAt"`
	UpdatedAt          string               `json:"updatedAt"`
	ExpArrivalDate     string               `json:"expArrivalDate"`
	OrderStatus        int                  `json:"orderStatus"`
	PayStatus          int                  `json:"payStatus"`
	Amount             float64              `json:"amount"`
	ActualAmount       float64              `json:"actualAmount"`
	AuditAmount        float64              `json:"auditAmount"`
	DiscountAmount     float64              `json:"discountAmount"`
	ProductCateNum     int                  `json:"productCateNum"`
	ProductNum         int                  `json:"productNum"`
	Remark             string               `json:"remark"`
	RejectionReason    string               `json:"rejectionReason"`
	DeclareProductList []DeclareProductItem `json:"declareProductList"`
}

type TransferOrderProductItem struct {
	ProductCode       string  `json:"productCode"`
	ProductName       string  `json:"productName"`
	ProductSpec       string  `json:"productSpec"`
	StockUnit         string  `json:"stockUnit"`
	TransferNum       float64 `json:"transferNum"`
	TransferPrice     float64 `json:"transferPrice"`
	TransferAmount    float64 `json:"transferAmount"`
	Quantity          float64 `json:"quantity"`
	AvailableQuantity float64 `json:"availableQuantity"`
}

type TransferOrderItem struct {
	TransferNo               string                     `json:"transferNo"`
	CreatedAt                string                     `json:"createdAt"`
	UpdatedAt                string                     `json:"updatedAt"`
	TransferAt               string                     `json:"transferAt"`
	Creator                  string                     `json:"creator"`
	OutWareNo                string                     `json:"outWareNo"`
	OutWareName              string                     `json:"outWareName"`
	InWareNo                 string                     `json:"inWareNo"`
	InWareName               string                     `json:"inWareName"`
	Status                   int                        `json:"status"`
	TransferType             int                        `json:"transferType"`
	ProductCateNum           int                        `json:"productCateNum"`
	ProductNum               float64                    `json:"productNum"`
	PriceAmount              float64                    `json:"priceAmount"`
	Remark                   string                     `json:"remark"`
	TransferOrderProductList []TransferOrderProductItem `json:"transferOrderProductList"`
}

type TransferOrderListResponse struct {
	Data  []TransferOrderItem `json:"data"`
	Total int                 `json:"total"`
}

type ReturnProductItem struct {
	ProductCode string  `json:"productCode"`
	ProductName string  `json:"productName"`
	ProductSpec string  `json:"productSpec"`
	ProductNum  float64 `json:"productNum"`
	ExamineNum  float64 `json:"examineNum"`
	ReceiveNum  float64 `json:"receiveNum"`
	Price       float64 `json:"price"`
	TotalPrice  float64 `json:"totalPrice"`
	RequireNo   string  `json:"requireNo"`
}

type ReturnOrderItem struct {
	ReturnNo           string              `json:"returnNo"`
	RequireNo          string              `json:"requireNo"`
	StoreName          string              `json:"storeName"`
	WarehouseName      string              `json:"warehouseName"`
	CreatedAt          string              `json:"createdAt"`
	UpdatedAt          string              `json:"updatedAt"`
	ReturnStatus       int                 `json:"returnStatus"`
	RefundStatus       int                 `json:"refundStatus"`
	ReturnNum          float64             `json:"returnNum"`
	ReturnAmount       float64             `json:"returnAmount"`
	ActualReturnAmount float64             `json:"actualReturnAmount"`
	ReturnProductList  []ReturnProductItem `json:"returnProductList"`
}

type ReturnOrderListResponse struct {
	Data  []ReturnOrderItem `json:"data"`
	Total int               `json:"total"`
}

type RequireOrderItem struct {
	DeclareNo      string  `json:"declareNo"`
	RequireNo      string  `json:"requireNo"`
	WarehouseNo    string  `json:"warehouseNo"`
	WarehouseName  string  `json:"warehouseName"`
	StoreName      string  `json:"storeName"`
	ContactName    string  `json:"contactName"`
	ContactPhone   string  `json:"contactPhone"`
	OrderAt        string  `json:"orderAt"`
	OrderStatus    int     `json:"orderStatus"`
	ProductCateNum int     `json:"productCateNum"`
	ProductNum     float64 `json:"productNum"`
	Amount         float64 `json:"amount"`
}

type RequireOrderListResponse struct {
	Total int                `json:"total"`
	Data  []RequireOrderItem `json:"data"`
}

type InboundProductItem struct {
	ProductCode   string  `json:"productCode"`
	ProductName   string  `json:"productName"`
	ProductSpec   string  `json:"productSpec"`
	ProductUnit   string  `json:"productUnit"`
	InboundNum    float64 `json:"inboundNum"`
	InboundPrice  float64 `json:"inboundPrice"`
	InboundAmount float64 `json:"inboundAmount"`
}

type InboundOrderItem struct {
	InboundNo          string               `json:"inboundNo"`
	BizNo              string               `json:"bizNo"`
	CreatedAt          string               `json:"createdAt"`
	InboundAt          string               `json:"inboundAt"`
	InboundPerson      string               `json:"inboundPerson"`
	InboundType        int                  `json:"inboundType"`
	ProviderName       string               `json:"providerName"`
	SupplierName       string               `json:"supplierName"`
	WarehouseName      string               `json:"warehouseName"`
	Status             int                  `json:"status"`
	Amount             float64              `json:"amount"`
	ProductAllNum      int                  `json:"productAllNum"`
	ProductTypeNum     int                  `json:"productTypeNum"`
	InboundProductList []InboundProductItem `json:"inboundProductList"`
}

type OutboundProductItem struct {
	ProductCode    string  `json:"productCode"`
	ProductName    string  `json:"productName"`
	ProductSpec    string  `json:"productSpec"`
	ProductUnit    string  `json:"productUnit"`
	OutboundNum    float64 `json:"outboundNum"`
	OutboundPrice  float64 `json:"outboundPrice"`
	OutboundAmount string  `json:"outboundAmount"`
}

type OutboundOrderItem struct {
	OutboundNo          string                `json:"outboundNo"`
	BizNo               string                `json:"bizNo"`
	CreatedAt           string                `json:"createdAt"`
	OutboundAt          string                `json:"outboundAt"`
	OutboundPerson      string                `json:"outboundPerson"`
	OutboundType        int                   `json:"outboundType"`
	ReceiptName         string                `json:"receiptName"`
	SupplierName        string                `json:"supplierName"`
	WarehouseName       string                `json:"warehouseName"`
	Status              int                   `json:"status"`
	Amount              float64               `json:"amount"`
	DeliveryNo          string                `json:"deliveryNo"`
	OutboundProductList []OutboundProductItem `json:"outboundProductList"`
}

type StoreInventorySummaryItem struct {
	TheDate              string  `json:"thedate"`
	WarehouseNo          string  `json:"warehouseNo"`
	StoreName            string  `json:"storeName"`
	ProductCode          string  `json:"productCode"`
	ProductName          string  `json:"productName"`
	Spec                 string  `json:"spec"`
	CategoryName         string  `json:"categoryName"`
	FinancialSubjectName string  `json:"financialSubjectName"`
	BeforeNum            float64 `json:"beforeNum"`
	BeforeAmount         float64 `json:"beforeAmount"`
	AfterNum             float64 `json:"afterNum"`
	AfterAmount          float64 `json:"afterAmount"`
	UnitName             string  `json:"unitName"`
}

type StoreInventorySummaryResponse struct {
	PageNo     int                         `json:"pageNo"`
	PageSize   int                         `json:"pageSize"`
	TotalCount int                         `json:"totalCount"`
	ResultList []StoreInventorySummaryItem `json:"resultList"`
}

type WarehouseProductItem struct {
	WarehouseNo        string  `json:"warehouseNo"`
	WarehouseName      string  `json:"warehouseName"`
	ProductCode        string  `json:"productCode"`
	ProductName        string  `json:"productName"`
	ProductSpec        string  `json:"productSpec"`
	CategoryName       string  `json:"categoryName"`
	StockUnit          string  `json:"stockUnit"`
	Quantity           float64 `json:"quantity"`
	AvailableQuantity  float64 `json:"availableQuantity"`
	CurrentAmount      float64 `json:"currentAmount"`
	EstimateNum        float64 `json:"estimateNum"`
	OccupyQuantity     float64 `json:"occupyQuantity"`
	CostPrice          float64 `json:"costPrice"`
	EstimateInboundNum float64 `json:"estimateInboundNum"`
}

type WarehouseProductListResponse struct {
	Total int                    `json:"total"`
	Data  []WarehouseProductItem `json:"data"`
}

func decodeBusinessID(resp *APIResponse) (*BusinessIDResponse, error) {
	var result BusinessIDResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("parse business id: %w", err)
	}
	return &result, nil
}

func decodeBoolResult(resp *APIResponse) (*BoolResultResponse, error) {
	var result BoolResultResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("parse bool result: %w", err)
	}
	return &result, nil
}

func (a *InventoryAPI) ReceiveDeclareOrder(ctx context.Context, params map[string]interface{}) (*BusinessIDResponse, error) {
	resp, err := a.client.Call(ctx, "v3/v3/scm/order/declare/order/receive", params)
	if err != nil {
		return nil, fmt.Errorf("receive declare order: %w", err)
	}
	return decodeBusinessID(resp)
}

func (a *InventoryAPI) GetDeclareOrderDetail(ctx context.Context, declareNo, shopCode string) (*DeclareDetailResponse, error) {
	params := map[string]interface{}{"declareNo": declareNo}
	if shopCode != "" {
		params["shopCode"] = shopCode
	}
	resp, err := a.client.Call(ctx, "v3/scm/order/declare/order/detail", params)
	if err != nil {
		return nil, fmt.Errorf("get declare order detail: %w", err)
	}
	return decodeData[DeclareDetailResponse](resp, "declare order detail")
}

func (a *InventoryAPI) CompleteRequireOrder(ctx context.Context, params map[string]interface{}) (*BusinessIDResponse, error) {
	resp, err := a.client.Call(ctx, "v3/scm/order/require/order/complete", params)
	if err != nil {
		return nil, fmt.Errorf("complete require order: %w", err)
	}
	return decodeBusinessID(resp)
}

func (a *InventoryAPI) DeliverRequireOrder(ctx context.Context, params map[string]interface{}) (*BusinessIDResponse, error) {
	resp, err := a.client.Call(ctx, "v3/scm/order/require/order/doDeliver", params)
	if err != nil {
		return nil, fmt.Errorf("deliver require order: %w", err)
	}
	return decodeBusinessID(resp)
}

func (a *InventoryAPI) CancelReturnOrder(ctx context.Context, params map[string]interface{}) (*BusinessIDResponse, error) {
	resp, err := a.client.Call(ctx, "v3/scm/order/return/order/cancel", params)
	if err != nil {
		return nil, fmt.Errorf("cancel return order: %w", err)
	}
	return decodeBusinessID(resp)
}

func (a *InventoryAPI) ExamineReturnOrder(ctx context.Context, params map[string]interface{}) (*BusinessIDResponse, error) {
	resp, err := a.client.Call(ctx, "v3/scm/order/return/order/examine", params)
	if err != nil {
		return nil, fmt.Errorf("examine return order: %w", err)
	}
	return decodeBusinessID(resp)
}

func (a *InventoryAPI) ReceiptReturnOrder(ctx context.Context, params map[string]interface{}) (*BusinessIDResponse, error) {
	resp, err := a.client.Call(ctx, "v3/scm/order/return/order/doReceipt", params)
	if err != nil {
		return nil, fmt.Errorf("receipt return order: %w", err)
	}
	return decodeBusinessID(resp)
}

func (a *InventoryAPI) ConfirmDeliveryArrive(ctx context.Context, params map[string]interface{}) (*BusinessIDResponse, error) {
	resp, err := a.client.Call(ctx, "v3/scm/delivery/performance/deliveryArrive", params)
	if err != nil {
		return nil, fmt.Errorf("confirm delivery arrive: %w", err)
	}
	return decodeBusinessID(resp)
}

func (a *InventoryAPI) ListTransferOrders(ctx context.Context, params map[string]interface{}) (*TransferOrderListResponse, error) {
	resp, err := a.client.Call(ctx, "v3/scm/order/transfer/list", params)
	if err != nil {
		return nil, fmt.Errorf("list transfer orders: %w", err)
	}
	return decodeData[TransferOrderListResponse](resp, "transfer orders")
}

func (a *InventoryAPI) ListReturnOrders(ctx context.Context, params map[string]interface{}) (*ReturnOrderListResponse, error) {
	resp, err := a.client.Call(ctx, "v3/scm/order/returnQuery", params)
	if err != nil {
		return nil, fmt.Errorf("list return orders: %w", err)
	}
	return decodeData[ReturnOrderListResponse](resp, "return orders")
}

func (a *InventoryAPI) ListRequireOrders(ctx context.Context, params map[string]interface{}) (*RequireOrderListResponse, error) {
	resp, err := a.client.Call(ctx, "v3/newPattern/scmApiserver/post/require/order/list", params)
	if err != nil {
		return nil, fmt.Errorf("list require orders: %w", err)
	}
	return decodeData[RequireOrderListResponse](resp, "require orders")
}

func (a *InventoryAPI) CreateRequireOrder(ctx context.Context, params map[string]interface{}) error {
	_, err := a.client.Call(ctx, "v3/newPattern/scmApiserver/post/require/order/submitRequireOrder", params)
	if err != nil {
		return fmt.Errorf("create require order: %w", err)
	}
	return nil
}

func (a *InventoryAPI) UpdateRequireOrderDetails(ctx context.Context, params map[string]interface{}) error {
	_, err := a.client.Call(ctx, "v3/newPattern/scmApiserver/post/require/order/detail-update", params)
	if err != nil {
		return fmt.Errorf("update require order details: %w", err)
	}
	return nil
}

func (a *InventoryAPI) GetTransferOrderDetail(ctx context.Context, transferNo string) (*TransferOrderItem, error) {
	resp, err := a.client.Call(ctx, "v3/newPattern/scmApiserver/post/transfer/order/detail", map[string]interface{}{"transferNo": transferNo})
	if err != nil {
		return nil, fmt.Errorf("get transfer order detail: %w", err)
	}
	return decodeData[TransferOrderItem](resp, "transfer order detail")
}

func (a *InventoryAPI) ListInboundOrders(ctx context.Context, params map[string]interface{}) (*[]InboundOrderItem, error) {
	resp, err := a.client.Call(ctx, "v3/scm/order/inbound/order/list", params)
	if err != nil {
		return nil, fmt.Errorf("list inbound orders: %w", err)
	}
	return decodeData[[]InboundOrderItem](resp, "inbound orders")
}

func (a *InventoryAPI) CreateInboundOrder(ctx context.Context, params map[string]interface{}) (*BusinessIDResponse, error) {
	resp, err := a.client.Call(ctx, "v3/scm/order/inbound/order/create", params)
	if err != nil {
		return nil, fmt.Errorf("create inbound order: %w", err)
	}
	return decodeBusinessID(resp)
}

func (a *InventoryAPI) UpdateInboundOrder(ctx context.Context, params map[string]interface{}) (*BusinessIDResponse, error) {
	resp, err := a.client.Call(ctx, "v3/scm/order/inbound/order/update", params)
	if err != nil {
		return nil, fmt.Errorf("update inbound order: %w", err)
	}
	return decodeBusinessID(resp)
}

func (a *InventoryAPI) FinishInboundOrders(ctx context.Context, params map[string]interface{}) error {
	_, err := a.client.Call(ctx, "v3/scm/order/inbound/order/finish", params)
	if err != nil {
		return fmt.Errorf("finish inbound orders: %w", err)
	}
	return nil
}

func (a *InventoryAPI) ListOutboundOrders(ctx context.Context, params map[string]interface{}) (*[]OutboundOrderItem, error) {
	resp, err := a.client.Call(ctx, "v3/scm/order/outbound/order/list", params)
	if err != nil {
		return nil, fmt.Errorf("list outbound orders: %w", err)
	}
	return decodeData[[]OutboundOrderItem](resp, "outbound orders")
}

func (a *InventoryAPI) CreateOutboundOrder(ctx context.Context, params map[string]interface{}) (*BusinessIDResponse, error) {
	resp, err := a.client.Call(ctx, "v3/scm/order/outbound/order/create", params)
	if err != nil {
		return nil, fmt.Errorf("create outbound order: %w", err)
	}
	return decodeBusinessID(resp)
}

func (a *InventoryAPI) UpdateOutboundOrder(ctx context.Context, params map[string]interface{}) (*BusinessIDResponse, error) {
	resp, err := a.client.Call(ctx, "v3/scm/order/outbound/order/update", params)
	if err != nil {
		return nil, fmt.Errorf("update outbound order: %w", err)
	}
	return decodeBusinessID(resp)
}

func (a *InventoryAPI) GetStoreInventorySummary(ctx context.Context, params map[string]interface{}) (*StoreInventorySummaryResponse, error) {
	resp, err := a.client.Call(ctx, "v3/newPattern/dataoneApiserver/post/scm/pos/storeInventory/summary", params)
	if err != nil {
		return nil, fmt.Errorf("get store inventory summary: %w", err)
	}
	return decodeData[StoreInventorySummaryResponse](resp, "store inventory summary")
}

func (a *InventoryAPI) ListWarehouseProducts(ctx context.Context, params map[string]interface{}) (*WarehouseProductListResponse, error) {
	resp, err := a.client.Call(ctx, "v3/newPattern/scmApiserver/post/warehouse-product/list", params)
	if err != nil {
		return nil, fmt.Errorf("list warehouse products: %w", err)
	}
	return decodeData[WarehouseProductListResponse](resp, "warehouse products")
}

func (a *InventoryAPI) OccupyProductStock(ctx context.Context, params map[string]interface{}) (*BoolResultResponse, error) {
	resp, err := a.client.Call(ctx, "v3/scm/product/stock/occupy", params)
	if err != nil {
		return nil, fmt.Errorf("occupy product stock: %w", err)
	}
	return decodeBoolResult(resp)
}

func (a *InventoryAPI) ReleaseProductStock(ctx context.Context, params map[string]interface{}) (*BoolResultResponse, error) {
	resp, err := a.client.Call(ctx, "v3/scm/product/stock/release", params)
	if err != nil {
		return nil, fmt.Errorf("release product stock: %w", err)
	}
	return decodeBoolResult(resp)
}

func (a *InventoryAPI) CreateInventoryAdjust(ctx context.Context, params map[string]interface{}) error {
	_, err := a.client.Call(ctx, "v3/newPattern/scmApiserver/post/inventory/adjust/create", params)
	if err != nil {
		return fmt.Errorf("create inventory adjust: %w", err)
	}
	return nil
}

func (a *InventoryAPI) AuditTransferOrder(ctx context.Context, params map[string]interface{}) error {
	_, err := a.client.Call(ctx, "v3/newPattern/scmApiserver/post/transfer/order/audit", params)
	if err != nil {
		return fmt.Errorf("audit transfer order: %w", err)
	}
	return nil
}

func (a *InventoryAPI) CreateProduct(ctx context.Context, params map[string]interface{}) (*BusinessIDResponse, error) {
	resp, err := a.client.Call(ctx, "v3/scm/product/changeProduct", params)
	if err != nil {
		return nil, fmt.Errorf("create product: %w", err)
	}
	return decodeBusinessID(resp)
}

func (a *InventoryAPI) CreateCategory(ctx context.Context, params map[string]interface{}) (*BusinessIDResponse, error) {
	resp, err := a.client.Call(ctx, "v3/scm/product/createCategory", params)
	if err != nil {
		return nil, fmt.Errorf("create category: %w", err)
	}
	return decodeBusinessID(resp)
}

func (a *InventoryAPI) CreateUnit(ctx context.Context, params map[string]interface{}) (*UnitCodeResponse, error) {
	resp, err := a.client.Call(ctx, "v3/scm/product/createUnit", params)
	if err != nil {
		return nil, fmt.Errorf("create unit: %w", err)
	}
	return decodeData[UnitCodeResponse](resp, "create unit")
}

func (a *InventoryAPI) UpdateProduct(ctx context.Context, params map[string]interface{}) (*BusinessIDResponse, error) {
	resp, err := a.client.Call(ctx, "v3/scm/product/update", params)
	if err != nil {
		return nil, fmt.Errorf("update product: %w", err)
	}
	return decodeBusinessID(resp)
}

func (a *InventoryAPI) CreateSupplier(ctx context.Context, params map[string]interface{}) error {
	_, err := a.client.Call(ctx, "v3/newPattern/scmApiserver/post/supplier/basic/create", params)
	if err != nil {
		return fmt.Errorf("create supplier: %w", err)
	}
	return nil
}

func (a *InventoryAPI) BatchGroupProduct(ctx context.Context, params map[string]interface{}) error {
	_, err := a.client.Call(ctx, "v3/newPattern/scmApiserver/post/product/batch-group-product", params)
	if err != nil {
		return fmt.Errorf("batch group product: %w", err)
	}
	return nil
}

func (a *InventoryAPI) UpdateSupplier(ctx context.Context, params map[string]interface{}) error {
	_, err := a.client.Call(ctx, "v3/newPattern/scmApiserver/post/supplier/basic/update", params)
	if err != nil {
		return fmt.Errorf("update supplier: %w", err)
	}
	return nil
}

func (a *InventoryAPI) BatchCreateMachiningCards(ctx context.Context, params map[string]interface{}) error {
	_, err := a.client.Call(ctx, "v3/newPattern/scmApiserver/post/machining/card/batchCreate", params)
	if err != nil {
		return fmt.Errorf("batch create machining cards: %w", err)
	}
	return nil
}
