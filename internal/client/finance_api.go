package client

import (
	"context"
	"encoding/json"
	"fmt"
)

// FinanceAPI provides finance-related API methods.
type FinanceAPI struct {
	client *Client
}

// NewFinanceAPI creates a new FinanceAPI.
func NewFinanceAPI(c *Client) *FinanceAPI {
	return &FinanceAPI{client: c}
}

type SplitOrderFlowItem struct {
	ActualAmount          float64 `json:"actualAmount"`
	ActualDonationAmount  float64 `json:"actualDonationAmount"`
	ActualPrincipalAmount float64 `json:"actualPrincipalAmount"`
	BatchNo               string  `json:"batchNo"`
	BizID                 string  `json:"bizId"`
	CreatedAt             string  `json:"createdAt"`
	DeliveryAmount        string  `json:"deliveryAmount"`
	HandlingFee           float64 `json:"handlingFee"`
	OriginOrderNo         string  `json:"originOrderNo"`
	RefundOrderNo         string  `json:"refundOrderNo"`
	ShopCode              string  `json:"shopCode"`
	ShopID                int64   `json:"shopId"`
	ShopName              string  `json:"shopName"`
	SplitAmount           string  `json:"splitAmount"`
	SplitRate             string  `json:"splitRate"`
	SplitStatus           string  `json:"splitStatus"`
	SplitType             string  `json:"splitType"`
	TotalAmount           string  `json:"totalAmount"`
}

type SplitOrderFlowListResponse struct {
	Data  []SplitOrderFlowItem `json:"data"`
	Total int                  `json:"total"`
}

type WechatBillExtra struct {
	BankType           string  `json:"bankType"`
	CurrencyType       string  `json:"currencyType"`
	DiscountAmount     float64 `json:"discountAmount"`
	RefundDiscount     float64 `json:"refundDiscount"`
	ServiceFeeInfo     string  `json:"serviceFeeInfo"`
	SettlementTotalFee float64 `json:"settlementTotalFee"`
}

type WechatBillItem struct {
	ChannelAppID           string          `json:"channelAppId"`
	ChannelDiscountAmount  float64         `json:"channelDiscountAmount"`
	ChannelMchID           string          `json:"channelMchId"`
	ChannelNo              string          `json:"channelNo"`
	ChannelRefundNo        string          `json:"channelRefundNo"`
	ChannelSupMchID        string          `json:"channelSupMchId"`
	CustomerID             string          `json:"customerId"`
	CustomerPaymentAmount  float64         `json:"customerPaymentAmount"`
	CustomerRefundAmount   float64         `json:"customerRefundAmount"`
	DeviceID               string          `json:"deviceId"`
	DiscountAmount         float64         `json:"discountAmount"`
	ID                     int64           `json:"id"`
	MerchantDiscountAmount float64         `json:"merchantDiscountAmount"`
	MerchantReceiptAmount  float64         `json:"merchantReceiptAmount"`
	MerchantRefundAmount   float64         `json:"merchantRefundAmount"`
	RefundAmount           float64         `json:"refundAmount"`
	RefundNo               string          `json:"refundNo"`
	ServiceFee             float64         `json:"serviceFee"`
	ServiceFeeRate         float64         `json:"serviceFeeRate"`
	ShopCode               string          `json:"shopCode"`
	ShopName               string          `json:"shopName"`
	TradeAmount            float64         `json:"tradeAmount"`
	TradeChannel           int             `json:"tradeChannel"`
	TradeNo                string          `json:"tradeNo"`
	TradeSource            int             `json:"tradeSource"`
	TradeTime              string          `json:"tradeTime"`
	TradeType              int             `json:"tradeType"`
	WechatData             WechatBillExtra `json:"wechatData"`
}

type WechatBillListResponse struct {
	Data  []WechatBillItem `json:"data"`
	Total int              `json:"total"`
}

type AlipayBillListResponse struct {
	Data      []json.RawMessage `json:"data"`
	Total     int64             `json:"total"`
	TotalPage int64             `json:"totalPage"`
}

type BusinessSummaryItem struct {
	BusinessAmt             float64 `json:"businessAmt"`
	BusinessNo              float64 `json:"businessNo"`
	ExpensesAmt             float64 `json:"expensesAmt"`
	InServiceAmt            float64 `json:"inServiceAmt"`
	IncomeAmt               float64 `json:"incomeAmt"`
	Income                  float64 `json:"income"`
	MerchantAmt             float64 `json:"merchantAmt"`
	OnsiteReceivableAmount  float64 `json:"onsiteReceivableAmount"`
	OnsiteReceivedAmount    float64 `json:"onsiteReceivedAmount"`
	OutsiteOrderCnt         int     `json:"outsiteOrderCnt"`
	OutsiteReceivableAmount float64 `json:"outsiteReceivableAmount"`
	OutsiteReceivedAmount   float64 `json:"outsiteReceivedAmount"`
	PeopleCnt               int     `json:"peopleCnt"`
	PlatformServiceAmt      int     `json:"platformServiceAmt"`
	ProcessDate             string  `json:"processDate"`
	RecordTime              string  `json:"recordTime"`
	RefundAmt               float64 `json:"refundAmt"`
	RefundNum               int     `json:"refundNum"`
	ShopCode                string  `json:"shopCode"`
	ShopName                string  `json:"shopName"`
}

type BusinessSummaryResponse struct {
	PageNo     int                   `json:"pageNo"`
	PageSize   int                   `json:"pageSize"`
	TotalCount int                   `json:"totalCount"`
	ResultList []BusinessSummaryItem `json:"resultList"`
}

type OrderTypeItem struct {
	OrderType     string `json:"orderType"`
	OrderTypeName string `json:"orderTypeName"`
}

type SettleSceneItem struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	SellerID string `json:"sellerId"`
}

type WechatBillURLResponse struct {
	FileURL string `json:"fileUrl"`
}

type StoreTurnoverItem struct {
	AfterAvgAmount   float64 `json:"afterAvgAmount"`
	BeforAvgAmount   float64 `json:"beforAvgAmount"`
	CategoryName     string  `json:"categoryName"`
	CostAmount       float64 `json:"costAmount"`
	Name             string  `json:"name"`
	Num              float64 `json:"num"`
	PathName         string  `json:"pathName"`
	ProcessDate      string  `json:"processDate"`
	ReceivableAmount float64 `json:"receivableAmount"`
	ReceivedAmount   float64 `json:"receivedAmount"`
	ReceivedRate     float64 `json:"receivedRate"`
	RecordDate       string  `json:"recordDate"`
	RefundAmount     float64 `json:"refundAmount"`
	RefundNum        float64 `json:"refundNum"`
	StoreCode        string  `json:"storeCode"`
	StoreID          string  `json:"storeId"`
	StoreName        string  `json:"storeName"`
}

type StoreTurnoverResponse struct {
	PageNo     int                 `json:"pageNo"`
	PageSize   int                 `json:"pageSize"`
	TotalCount int                 `json:"totalCount"`
	ResultList []StoreTurnoverItem `json:"resultList"`
}

func (a *FinanceAPI) GetSplitOrderFlows(ctx context.Context, params map[string]interface{}) (*SplitOrderFlowListResponse, error) {
	resp, err := a.client.Call(ctx, "v3/payCenter/split/orderFlow", params)
	if err != nil {
		return nil, fmt.Errorf("get split order flows: %w", err)
	}
	return decodeData[SplitOrderFlowListResponse](resp, "split order flows")
}

func (a *FinanceAPI) GetWechatBills(ctx context.Context, params map[string]interface{}) (*WechatBillListResponse, error) {
	resp, err := a.client.Call(ctx, "v3/payCenter/custom/wechatBill", params)
	if err != nil {
		return nil, fmt.Errorf("get wechat bills: %w", err)
	}
	return decodeData[WechatBillListResponse](resp, "wechat bills")
}

func (a *FinanceAPI) GetAlipayBills(ctx context.Context, params map[string]interface{}) (*AlipayBillListResponse, error) {
	resp, err := a.client.Call(ctx, "v3/payCenter/offlinePayments/pullAlipayBillData", params)
	if err != nil {
		return nil, fmt.Errorf("get alipay bills: %w", err)
	}
	return decodeData[AlipayBillListResponse](resp, "alipay bills")
}

func (a *FinanceAPI) GetBusinessSummary(ctx context.Context, params map[string]interface{}) (*BusinessSummaryResponse, error) {
	resp, err := a.client.Call(ctx, "v3/dataone/finance/summary/businessRecord", params)
	if err != nil {
		return nil, fmt.Errorf("get business summary: %w", err)
	}
	return decodeData[BusinessSummaryResponse](resp, "business summary")
}

func (a *FinanceAPI) ListOrderTypes(ctx context.Context) (*[]OrderTypeItem, error) {
	resp, err := a.client.Call(ctx, "v3/dataone/finance/detail/orderType", map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("list order types: %w", err)
	}
	return decodeData[[]OrderTypeItem](resp, "order types")
}

func (a *FinanceAPI) ListSettleScenes(ctx context.Context) (*[]SettleSceneItem, error) {
	resp, err := a.client.Call(ctx, "v3/dataone/finance/detail/settleScene", map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("list settle scenes: %w", err)
	}
	return decodeData[[]SettleSceneItem](resp, "settle scenes")
}

func (a *FinanceAPI) GetWechatBillURL(ctx context.Context, billDate, shopCode string) (*WechatBillURLResponse, error) {
	params := map[string]interface{}{"billDate": billDate}
	if shopCode != "" {
		params["shopCode"] = shopCode
	}
	resp, err := a.client.Call(ctx, "v3/pay/getWechatMerchantBillUrl", params)
	if err != nil {
		return nil, fmt.Errorf("get wechat bill url: %w", err)
	}
	return decodeData[WechatBillURLResponse](resp, "wechat bill url")
}

func (a *FinanceAPI) GetStoreTurnover(ctx context.Context, params map[string]interface{}) (*StoreTurnoverResponse, error) {
	resp, err := a.client.Call(ctx, "v3/dataone/item/store/turnover", params)
	if err != nil {
		return nil, fmt.Errorf("get store turnover: %w", err)
	}
	return decodeData[StoreTurnoverResponse](resp, "store turnover")
}
