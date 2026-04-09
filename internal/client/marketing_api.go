package client

import (
	"context"
	"encoding/json"
	"fmt"
)

// MarketingAPI provides marketing-related API methods.
type MarketingAPI struct {
	client *Client
}

// NewMarketingAPI creates a new MarketingAPI.
func NewMarketingAPI(c *Client) *MarketingAPI {
	return &MarketingAPI{client: c}
}

type CouponStatusResponse struct {
	Status       int    `json:"status"`
	UseOrderCode string `json:"useOrderCode"`
	UseStoreID   string `json:"useStoreId"`
	UseStoreName string `json:"useStoreName"`
	UseStoreNo   string `json:"useStoreNo"`
	UseTime      string `json:"useTime"`
}

type CouponTemplateSummary struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	SellerID     int64  `json:"sellerId"`
	Status       int    `json:"status"`
	Type         int    `json:"type"`
	ThirdBizCode string `json:"thirdBizCode"`
}

type CouponDetailResponse struct {
	ID             int64                 `json:"id"`
	Status         int                   `json:"status"`
	CustomerID     string                `json:"customerId"`
	CustomerName   string                `json:"customerName"`
	CustomerPhone  string                `json:"customerPhone"`
	TemplateID     int64                 `json:"templateId"`
	TemplateName   string                `json:"templateName"`
	TemplateType   int                   `json:"templateType"`
	UseStatus      int                   `json:"status"`
	UseOrderCode   string                `json:"useOrderCode"`
	UseStoreName   string                `json:"useStoreName"`
	UseTime        string                `json:"useTime"`
	ValidBeginDate string                `json:"validBeginDate"`
	ValidEndDate   string                `json:"validEndDate"`
	TemplateInfo   CouponTemplateSummary `json:"templateInfo"`
}

type AnonymousCouponResponse struct {
	Amount           float64 `json:"amount"`
	Code             string  `json:"code"`
	IssueBatchCode   string  `json:"issueBatchCode"`
	IssueBatchID     int64   `json:"issueBatchId"`
	ProduceBatchCode string  `json:"produceBatchCode"`
	ProduceBatchID   int64   `json:"produceBatchId"`
	Status           int     `json:"status"`
	TemplateID       int64   `json:"templateId"`
	TemplateName     string  `json:"templateName"`
	TemplateType     int     `json:"templateType"`
}

type ThirdCodeTemplateResponse struct {
	ID int64 `json:"id"`
}

type CouponChooseItem struct {
	Attended               int     `json:"attended"`
	AvailableQuantity      int     `json:"availableQuantity"`
	CardID                 int64   `json:"cardId"`
	CouponCenterType       int     `json:"couponCenterType"`
	CouponID               int64   `json:"couponId"`
	CouponType             int     `json:"couponType"`
	EndAt                  string  `json:"endAt"`
	ExpireDesc             string  `json:"expireDesc"`
	Origin                 int     `json:"origin"`
	OriginDesc             string  `json:"originDesc"`
	Title                  string  `json:"title"`
	UnUseDesc              string  `json:"unUseDesc"`
	UnUseStatus            int     `json:"unUseStatus"`
	UseStatus              int     `json:"useStatus"`
	UserCouponCanDisAmount float64 `json:"userCouponCanDisAmount"`
}

type CouponActivityItem struct {
	ActivityCode     string `json:"activityCode"`
	ActivityID       int64  `json:"activityId"`
	ActivityName     string `json:"activityName"`
	EndDate          string `json:"endDate"`
	RemainderNum     int    `json:"remainderNum"`
	StartDate        string `json:"startDate"`
	StockLimitSwitch int    `json:"stockLimitSwitch"`
}

type CouponActivityPageResponse struct {
	CurrentPage int                  `json:"currentPage"`
	List        []CouponActivityItem `json:"list"`
	PageSize    int                  `json:"pageSize"`
	Total       int                  `json:"total"`
	TotalPage   int                  `json:"totalPage"`
}

type ExchangeCodeStatusResponse struct {
	ActivityCode   string `json:"activityCode"`
	ActivityName   string `json:"activityName"`
	ExchangeStatus int    `json:"exchangeStatus"`
}

type ActivityTaskRecord struct {
	ActivityID        string `json:"activityId"`
	ActivityName      string `json:"activityName"`
	ActivityStatus    int    `json:"activityStatus"`
	CurrentJoinNum    int    `json:"currentJoinNum"`
	CycleJoinNum      int    `json:"cycleJoinNum"`
	OrderTaskCanGet   int    `json:"orderTaskCanGet"`
	ResidueCanJoinNum int    `json:"residueCanJoinNum"`
	TotalJoinNum      int    `json:"totalJoinNum"`
}

type GiftCardTemplateSummary struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	SellerID     int64   `json:"sellerId"`
	Status       int     `json:"status"`
	ThirdBizCode string  `json:"thirdBizCode"`
	Type         int     `json:"type"`
	SalePrice    float64 `json:"salePrice"`
}

type GiftCardInfoResponse struct {
	Amount          float64                 `json:"amount"`
	CanUse          bool                    `json:"canUse"`
	CardNo          string                  `json:"cardNo"`
	CustomerID      int64                   `json:"customerId"`
	CustomerPhone   string                  `json:"customerPhone"`
	ID              int64                   `json:"id"`
	Name            string                  `json:"name"`
	RemainingAmount float64                 `json:"remainingAmount"`
	Status          int                     `json:"status"`
	TemplateID      int64                   `json:"templateId"`
	TemplateInfo    GiftCardTemplateSummary `json:"templateInfo"`
	Type            int                     `json:"type"`
	ValidBeginDate  string                  `json:"validBeginDate"`
	ValidEndDate    string                  `json:"validEndDate"`
}

type GiftCardFlowItem struct {
	ActionType   int     `json:"actionType"`
	Balance      float64 `json:"balance"`
	BizID        string  `json:"bizId"`
	CardNo       string  `json:"cardNo"`
	CreatedAt    string  `json:"createdAt"`
	OrderAmount  float64 `json:"orderAmount"`
	Reason       string  `json:"reason"`
	UseAmount    float64 `json:"useAmount"`
	UseChannel   int     `json:"useChannel"`
	UseOrderCode string  `json:"useOrderCode"`
	UseStoreName string  `json:"useStoreName"`
	UseStoreNo   string  `json:"useStoreNo"`
}

type GiftCardFlowPageResponse struct {
	CurrentPageSize int                `json:"currentPageSize"`
	List            []GiftCardFlowItem `json:"list"`
	PageCount       int                `json:"pageCount"`
	PageNum         int                `json:"pageNum"`
	PageSize        int                `json:"pageSize"`
	TotalCount      int                `json:"totalCount"`
}

type CustomerGiftCardItem struct {
	Balance        float64                 `json:"balance"`
	BeginAt        string                  `json:"beginAt"`
	CanPresent     int                     `json:"canPresent"`
	CardNo         string                  `json:"cardNo"`
	CardPrice      float64                 `json:"cardPrice"`
	CardStatus     int                     `json:"cardStatus"`
	CardTemplateID int64                   `json:"cardTemplateId"`
	CardType       int                     `json:"cardType"`
	CustomerID     int64                   `json:"customerId"`
	EndAt          string                  `json:"endAt"`
	ID             int64                   `json:"id"`
	Name           string                  `json:"name"`
	Template       GiftCardTemplateSummary `json:"template"`
	TotalBalance   float64                 `json:"totalBalance"`
}

type PromotionActivityItem struct {
	ActivityCode   string `json:"activityCode"`
	ActivityName   string `json:"activityName"`
	ActivityStatus int    `json:"activityStatus"`
	ActivityTitle  string `json:"activityTitle"`
	ActivityType   string `json:"activityType"`
	StoreID        int64  `json:"storeId"`
}

type ConfirmOrderFee struct {
	ActualAmount   float64 `json:"actualAmount"`
	DiscountAmount float64 `json:"discountAmount"`
	GoodsAmount    float64 `json:"goodsAmount"`
	GoodsSubTotal  float64 `json:"goodsSubTotal"`
	MealBoxAmount  float64 `json:"mealBoxAmount"`
	PackingAmount  float64 `json:"packingAmount"`
	TotalAmount    float64 `json:"totalAmount"`
}

type ConfirmPricingResponse struct {
	Discounted    bool            `json:"discounted"`
	NewDiscounted int             `json:"newDiscounted"`
	OrderFee      ConfirmOrderFee `json:"orderFee"`
}

type ActivityRealtimeGrantParams map[string]interface{}

type ActivityAsyncGrantParams map[string]interface{}

type TemplateRealtimeGrantParams map[string]interface{}

type CouponRecoveryParams map[string]interface{}

type ActivityRealtimeGrantItem struct {
	ID         int64 `json:"id"`
	TemplateID int64 `json:"templateId"`
}

type TemplateRealtimeGrantResponse struct {
	CardID     string   `json:"cardId"`
	CardIDList []string `json:"cardIdList"`
	Num        int      `json:"num"`
}

type CouponRecoveryResponse struct {
	CouponIDs      []interface{} `json:"couponIds"`
	RecoveryIDs    []interface{} `json:"recoveryIds"`
	RecoveryStatus int           `json:"recoveryStatus"`
}

type GiftCardConsumeParams map[string]interface{}

type GiftCardConsumeReverseParams map[string]interface{}

type GiftCardConsumePartReverseParams map[string]interface{}

type GiftCardBindParams map[string]interface{}

type GiftCardConsumeInfo struct {
	AfterGiftBalance        float64 `json:"afterGiftBalance"`
	AfterRechargeBalance    float64 `json:"afterRechargeBalance"`
	BeforeGiftBalance       float64 `json:"beforeGiftBalance"`
	BeforeRechargeBalance   float64 `json:"beforeRechargeBalance"`
	BizID                   string  `json:"bizId"`
	CardNo                  string  `json:"cardNo"`
	DecreaseGiftBalance     float64 `json:"decreaseGiftBalance"`
	DecreaseRechargeBalance float64 `json:"decreaseRechargeBalance"`
	DecreaseTotalBalance    float64 `json:"decreaseTotalBalance"`
	SingleCostAmount        float64 `json:"singleCostAmount"`
}

type GiftCardConsumeResponse struct {
	Bind            bool                `json:"bind"`
	CardConsumeInfo GiftCardConsumeInfo `json:"cardConsumeInfo"`
	CardType        int                 `json:"cardType"`
	CustomerID      int64               `json:"customerId"`
	Success         bool                `json:"success"`
}

type GiftCardPartReverseResponse struct {
	AfterGiftBalance       float64 `json:"afterGiftBalance"`
	AfterRechargeBalance   float64 `json:"afterRechargeBalance"`
	AfterTotalBalance      float64 `json:"afterTotalBalance"`
	BeforeGiftBalance      float64 `json:"beforeGiftBalance"`
	BeforeRechargeBalance  float64 `json:"beforeRechargeBalance"`
	BeforeTotalBalance     float64 `json:"beforeTotalBalance"`
	BizID                  string  `json:"bizId"`
	CardNo                 string  `json:"cardNo"`
	ReverseGiftBalance     float64 `json:"reverseGiftBalance"`
	ReverseRechargeBalance float64 `json:"reverseRechargeBalance"`
	ReverseTotalBalance    float64 `json:"reverseTotalBalance"`
	SubBizID               string  `json:"subBizId"`
}

type ExchangeCodeDispatchParams map[string]interface{}

type ExchangeCodeDisableParams map[string]interface{}

type ActivityTaskClaimParams map[string]interface{}

type CouponRecycleParams map[string]interface{}

type ExchangeCodeItem struct {
	Code string `json:"code"`
}

type ExchangeCodeDisableItem struct {
	Code       string `json:"code"`
	Disabled   bool   `json:"disabled"`
	FailReason string `json:"failReason"`
}

type CouponTemplateEnableResponse struct {
	Data bool `json:"data"`
}

type CouponTemplateEnableParams map[string]interface{}

type GiftCardReportLossParams map[string]interface{}

type GiftCardRelieveLossParams map[string]interface{}

type GiftCardRecycleParams map[string]interface{}

type GiftCardExchangeParams map[string]interface{}

type GiftCardBatchGrantParams map[string]interface{}

type GiftCardGrantItem struct {
	Amount         float64 `json:"amount"`
	BindValidUnit  int     `json:"bindValidUnit"`
	CardNo         string  `json:"cardNo"`
	CustomerID     int64   `json:"customerId"`
	GiveCount      int     `json:"giveCount"`
	ImageURL       string  `json:"imageUrl"`
	Name           string  `json:"name"`
	SalePrice      float64 `json:"salePrice"`
	TemplateID     int64   `json:"templateId"`
	ThirdCardCode  string  `json:"thirdCardCode"`
	ThirdCardNo    string  `json:"thirdCardNo"`
	Type           int     `json:"type"`
	ValidBeginDate string  `json:"validBeginDate"`
	ValidEndDate   string  `json:"validEndDate"`
	ValidTime      string  `json:"validTime"`
	ValidType      int     `json:"validType"`
}

type MultiGiftCardConsumeItem struct {
	Bind            bool                `json:"bind"`
	CardConsumeInfo GiftCardConsumeInfo `json:"cardConsumeInfo"`
	CardType        int                 `json:"cardType"`
	CustomerID      int64               `json:"customerId"`
	Success         bool                `json:"success"`
}

type MultiGiftCardConsumeResponse struct {
	BizID           string                     `json:"bizId"`
	CardConsumeInfo []MultiGiftCardConsumeItem `json:"cardConsumeInfo"`
	Success         bool                       `json:"success"`
}

func (a *MarketingAPI) GetCouponStatus(ctx context.Context, id string, codeType int) (*CouponStatusResponse, error) {
	params := map[string]interface{}{"id": id}
	if codeType != 0 {
		params["type"] = codeType
	}
	resp, err := a.client.Call(ctx, "v3/coupon/status", params)
	if err != nil {
		return nil, fmt.Errorf("get coupon status: %w", err)
	}
	return decodeData[CouponStatusResponse](resp, "coupon status")
}

func (a *MarketingAPI) GrantCouponsByActivity(ctx context.Context, params ActivityRealtimeGrantParams) (*[]ActivityRealtimeGrantItem, error) {
	resp, err := a.client.Call(ctx, "v3/crm/marketing/couponGrant", params)
	if err != nil {
		return nil, fmt.Errorf("grant coupons by activity: %w", err)
	}
	return decodeData[[]ActivityRealtimeGrantItem](resp, "grant coupons by activity")
}

func (a *MarketingAPI) GrantCouponsByActivityAsync(ctx context.Context, params ActivityAsyncGrantParams) error {
	_, err := a.client.Call(ctx, "v3/crm/marketing/couponGrantSync", params)
	if err != nil {
		return fmt.Errorf("grant coupons by activity async: %w", err)
	}
	return nil
}

func (a *MarketingAPI) GrantCouponsByTemplate(ctx context.Context, params TemplateRealtimeGrantParams) (*TemplateRealtimeGrantResponse, error) {
	resp, err := a.client.Call(ctx, "v3/coupon/issueCouponsToCustomer", params)
	if err != nil {
		return nil, fmt.Errorf("grant coupons by template: %w", err)
	}
	return decodeData[TemplateRealtimeGrantResponse](resp, "grant coupons by template")
}

func (a *MarketingAPI) RecoverGrantedCoupons(ctx context.Context, params CouponRecoveryParams) (*CouponRecoveryResponse, error) {
	resp, err := a.client.Call(ctx, "v3/crm/marketing/couponRecovery", params)
	if err != nil {
		return nil, fmt.Errorf("recover granted coupons: %w", err)
	}
	return decodeData[CouponRecoveryResponse](resp, "recover granted coupons")
}

func (a *MarketingAPI) GetCouponDetail(ctx context.Context, userCouponCode string) (*CouponDetailResponse, error) {
	resp, err := a.client.Call(ctx, "v3/coupon/getCouponDetail", map[string]string{"userCouponCode": userCouponCode})
	if err != nil {
		return nil, fmt.Errorf("get coupon detail: %w", err)
	}
	return decodeData[CouponDetailResponse](resp, "coupon detail")
}

func (a *MarketingAPI) GetCouponTemplateDetail(ctx context.Context, id int64) (*CouponTemplateSummary, error) {
	resp, err := a.client.Call(ctx, "v3/coupon/getTemplateDetail", map[string]int64{"id": id})
	if err != nil {
		return nil, fmt.Errorf("get coupon template detail: %w", err)
	}
	return decodeData[CouponTemplateSummary](resp, "coupon template detail")
}

func (a *MarketingAPI) EnableCouponTemplate(ctx context.Context, params CouponTemplateEnableParams) (*CouponTemplateEnableResponse, error) {
	resp, err := a.client.Call(ctx, "v3/coupon/templateEnable", params)
	if err != nil {
		return nil, fmt.Errorf("enable coupon template: %w", err)
	}
	return decodeData[CouponTemplateEnableResponse](resp, "enable coupon template")
}

func (a *MarketingAPI) ChooseCouponsAsOne(ctx context.Context, params map[string]interface{}) (*[]CouponChooseItem, error) {
	resp, err := a.client.Call(ctx, "v3/bsns/customer/chooseCouponAsOne", params)
	if err != nil {
		return nil, fmt.Errorf("choose coupons as one: %w", err)
	}
	return decodeData[[]CouponChooseItem](resp, "choose coupons")
}

func (a *MarketingAPI) GetAnonymousCoupon(ctx context.Context, code string) (*AnonymousCouponResponse, error) {
	resp, err := a.client.Call(ctx, "v3/coupon/anonymous-coupon/get", map[string]string{"code": code})
	if err != nil {
		return nil, fmt.Errorf("get anonymous coupon: %w", err)
	}
	return decodeData[AnonymousCouponResponse](resp, "anonymous coupon")
}

func (a *MarketingAPI) GetCouponTemplateListByIDs(ctx context.Context, ids []int64) (*[]CouponTemplateSummary, error) {
	resp, err := a.client.Call(ctx, "v3/coupon/getCouponTemplateListByIds", map[string][]int64{"templateIds": ids})
	if err != nil {
		return nil, fmt.Errorf("get coupon template list by ids: %w", err)
	}
	return decodeData[[]CouponTemplateSummary](resp, "coupon template list")
}

func (a *MarketingAPI) GetCouponTemplateByThirdCode(ctx context.Context, sellerType int, thirdBizCode string) (*ThirdCodeTemplateResponse, error) {
	resp, err := a.client.Call(ctx, "v3/newPattern/couponCenter/get/template/third-code/get", map[string]interface{}{
		"sellerType":   sellerType,
		"thirdBizCode": thirdBizCode,
	})
	if err != nil {
		return nil, fmt.Errorf("get template by third code: %w", err)
	}
	return decodeData[ThirdCodeTemplateResponse](resp, "template by third code")
}

func (a *MarketingAPI) ListCouponActivities(ctx context.Context, channelID, pageNo, pageSize int) (*CouponActivityPageResponse, error) {
	resp, err := a.client.Call(ctx, "v3/crm/marketing/couponActivityList/page", map[string]int{
		"channelId": channelID,
		"pageNo":    pageNo,
		"pageSize":  pageSize,
	})
	if err != nil {
		return nil, fmt.Errorf("list coupon activities: %w", err)
	}
	return decodeData[CouponActivityPageResponse](resp, "coupon activities")
}

func (a *MarketingAPI) QueryExchangeCodeStatus(ctx context.Context, code string) (*ExchangeCodeStatusResponse, error) {
	resp, err := a.client.Call(ctx, "v3/coupon/exchange/queryCodeStatus", map[string]string{"code": code})
	if err != nil {
		return nil, fmt.Errorf("query exchange code status: %w", err)
	}
	return decodeData[ExchangeCodeStatusResponse](resp, "exchange code status")
}

func (a *MarketingAPI) DispatchExchangeCodes(ctx context.Context, params ExchangeCodeDispatchParams) (*[]ExchangeCodeItem, error) {
	resp, err := a.client.Call(ctx, "v3/coupon/exchange/codeDispatch", params)
	if err != nil {
		return nil, fmt.Errorf("dispatch exchange codes: %w", err)
	}
	return decodeData[[]ExchangeCodeItem](resp, "dispatch exchange codes")
}

func (a *MarketingAPI) DisableExchangeCodes(ctx context.Context, params ExchangeCodeDisableParams) (*[]ExchangeCodeDisableItem, error) {
	resp, err := a.client.Call(ctx, "v3/coupon/exchange/codeDisableBatch", params)
	if err != nil {
		return nil, fmt.Errorf("disable exchange codes: %w", err)
	}
	return decodeData[[]ExchangeCodeDisableItem](resp, "disable exchange codes")
}

func (a *MarketingAPI) GetActivityTasks(ctx context.Context, activityIDList []int64, customerID int64, mobilePhone string) (*[]ActivityTaskRecord, error) {
	params := map[string]interface{}{
		"activityIdList": activityIDList,
		"customerId":     customerID,
	}
	if mobilePhone != "" {
		params["mobilePhone"] = mobilePhone
	}
	resp, err := a.client.Call(ctx, "v3/marketing/activity/activityTasks", params)
	if err != nil {
		return nil, fmt.Errorf("get activity tasks: %w", err)
	}
	return decodeData[[]ActivityTaskRecord](resp, "activity tasks")
}

func (a *MarketingAPI) ClaimActivityTask(ctx context.Context, params ActivityTaskClaimParams) error {
	_, err := a.client.Call(ctx, "v3/marketing/activity/getTask", params)
	if err != nil {
		return fmt.Errorf("claim activity task: %w", err)
	}
	return nil
}

func (a *MarketingAPI) RecycleCoupons(ctx context.Context, params CouponRecycleParams) error {
	_, err := a.client.Call(ctx, "v3/coupon/recycleCoupons", params)
	if err != nil {
		return fmt.Errorf("recycle coupons: %w", err)
	}
	return nil
}

func (a *MarketingAPI) GetCardInfo(ctx context.Context, cardNo string, takeTemplate int) (*GiftCardInfoResponse, error) {
	params := map[string]interface{}{}
	if cardNo != "" {
		params["cardNo"] = cardNo
	}
	if takeTemplate != 0 {
		params["takeCardTemplate"] = takeTemplate
	}
	resp, err := a.client.Call(ctx, "v3/coupon/card/unifyCardInfo", params)
	if err != nil {
		return nil, fmt.Errorf("get card info: %w", err)
	}
	return decodeData[GiftCardInfoResponse](resp, "card info")
}

func (a *MarketingAPI) ConsumeGiftCard(ctx context.Context, params GiftCardConsumeParams) (*GiftCardConsumeResponse, error) {
	resp, err := a.client.Call(ctx, "v3/crm/giftCard/unifyCardConsume", params)
	if err != nil {
		return nil, fmt.Errorf("consume gift card: %w", err)
	}
	return decodeData[GiftCardConsumeResponse](resp, "consume gift card")
}

func (a *MarketingAPI) ConsumeMultiGiftCards(ctx context.Context, params map[string]interface{}) (*MultiGiftCardConsumeResponse, error) {
	resp, err := a.client.Call(ctx, "v3/crm/giftCard/unifyMultiCardConsume", params)
	if err != nil {
		return nil, fmt.Errorf("consume multi gift cards: %w", err)
	}
	return decodeData[MultiGiftCardConsumeResponse](resp, "consume multi gift cards")
}

func (a *MarketingAPI) ReverseGiftCardConsume(ctx context.Context, params GiftCardConsumeReverseParams) (*GiftCardConsumeResponse, error) {
	resp, err := a.client.Call(ctx, "v3/crm/giftCard/unifyCardConsumeReverse", params)
	if err != nil {
		return nil, fmt.Errorf("reverse gift card consume: %w", err)
	}
	return decodeData[GiftCardConsumeResponse](resp, "reverse gift card consume")
}

func (a *MarketingAPI) ReverseMultiGiftCards(ctx context.Context, params map[string]interface{}) error {
	_, err := a.client.Call(ctx, "v3/crm/giftCard/unifyMultiCardConsumeReverse", params)
	if err != nil {
		return fmt.Errorf("reverse multi gift cards: %w", err)
	}
	return nil
}

func (a *MarketingAPI) PartReverseGiftCardConsume(ctx context.Context, params GiftCardConsumePartReverseParams) (*GiftCardPartReverseResponse, error) {
	resp, err := a.client.Call(ctx, "v3/crm/giftCard/unifyCardConsumePartReverse", params)
	if err != nil {
		return nil, fmt.Errorf("part reverse gift card consume: %w", err)
	}
	return decodeData[GiftCardPartReverseResponse](resp, "part reverse gift card consume")
}

func (a *MarketingAPI) BindGiftCard(ctx context.Context, params GiftCardBindParams) error {
	_, err := a.client.Call(ctx, "v3/coupon/card/unifyBindCard", params)
	if err != nil {
		return fmt.Errorf("bind gift card: %w", err)
	}
	return nil
}

func (a *MarketingAPI) ReportGiftCardLoss(ctx context.Context, params GiftCardReportLossParams) error {
	_, err := a.client.Call(ctx, "v3/crm/giftCard/giftCardReportLoss", params)
	if err != nil {
		return fmt.Errorf("report gift card loss: %w", err)
	}
	return nil
}

func (a *MarketingAPI) RelieveGiftCardLoss(ctx context.Context, params GiftCardRelieveLossParams) error {
	_, err := a.client.Call(ctx, "v3/crm/giftCard/giftCardRelieveLoss", params)
	if err != nil {
		return fmt.Errorf("relieve gift card loss: %w", err)
	}
	return nil
}

func (a *MarketingAPI) RecycleGiftCard(ctx context.Context, params GiftCardRecycleParams) error {
	_, err := a.client.Call(ctx, "v3/coupon/card/cardRecycle", params)
	if err != nil {
		return fmt.Errorf("recycle gift card: %w", err)
	}
	return nil
}

func (a *MarketingAPI) ExchangeGiftCard(ctx context.Context, params GiftCardExchangeParams) (*[]GiftCardGrantItem, error) {
	resp, err := a.client.Call(ctx, "v3/coupon/card/exchangeGiftCard", params)
	if err != nil {
		return nil, fmt.Errorf("exchange gift card: %w", err)
	}
	return decodeData[[]GiftCardGrantItem](resp, "exchange gift card")
}

func (a *MarketingAPI) BatchGrantGiftCard(ctx context.Context, params GiftCardBatchGrantParams) (*[]GiftCardGrantItem, error) {
	resp, err := a.client.Call(ctx, "v3/coupon/card/batchGrantCard", params)
	if err != nil {
		return nil, fmt.Errorf("batch grant gift card: %w", err)
	}
	return decodeData[[]GiftCardGrantItem](resp, "batch grant gift card")
}

func (a *MarketingAPI) ListGiftCardFlows(ctx context.Context, params map[string]interface{}) (*GiftCardFlowPageResponse, error) {
	resp, err := a.client.Call(ctx, "v3/coupon/card/consumeFlow/page", params)
	if err != nil {
		return nil, fmt.Errorf("list gift card flows: %w", err)
	}
	return decodeData[GiftCardFlowPageResponse](resp, "gift card flows")
}

func (a *MarketingAPI) GetGiftCardTemplateDetail(ctx context.Context, id int64) (*GiftCardTemplateSummary, error) {
	resp, err := a.client.Call(ctx, "v3/coupon/card/getGiftCardTemplateDetail", map[string]int64{"id": id})
	if err != nil {
		return nil, fmt.Errorf("get gift card template detail: %w", err)
	}
	return decodeData[GiftCardTemplateSummary](resp, "gift card template detail")
}

func (a *MarketingAPI) ListCustomerGiftCards(ctx context.Context, params map[string]interface{}) (*[]CustomerGiftCardItem, error) {
	resp, err := a.client.Call(ctx, "v3/crm/giftCard/getCustomerGiftCardListApplyShop", params)
	if err != nil {
		return nil, fmt.Errorf("list customer gift cards: %w", err)
	}
	return decodeData[[]CustomerGiftCardItem](resp, "customer gift cards")
}

func (a *MarketingAPI) ListGiftCardTemplates(ctx context.Context, ids []int64) (*[]GiftCardTemplateSummary, error) {
	resp, err := a.client.Call(ctx, "v3/coupon/card/listGiftCardTemplates", map[string][]int64{"templateIds": ids})
	if err != nil {
		return nil, fmt.Errorf("list gift card templates: %w", err)
	}
	return decodeData[[]GiftCardTemplateSummary](resp, "gift card templates")
}

func (a *MarketingAPI) ListStorePromotionActivities(ctx context.Context, params map[string]interface{}) (*[]PromotionActivityItem, error) {
	resp, err := a.client.Call(ctx, "v3/newPattern/marketingCenter/post/activity/Baking/goods/activity", params)
	if err != nil {
		return nil, fmt.Errorf("list store promotion activities: %w", err)
	}
	return decodeData[[]PromotionActivityItem](resp, "store promotion activities")
}

func (a *MarketingAPI) ConfirmOrderPricing(ctx context.Context, params map[string]interface{}) (*ConfirmPricingResponse, error) {
	resp, err := a.client.Call(ctx, "v3/newPattern/cateringApiserver/post/order/appConfirm", params)
	if err != nil {
		return nil, fmt.Errorf("confirm order pricing: %w", err)
	}
	var result ConfirmPricingResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("parse confirm pricing: %w", err)
	}
	return &result, nil
}
