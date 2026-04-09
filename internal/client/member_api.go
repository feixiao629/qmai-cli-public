package client

import (
	"context"
	"encoding/json"
	"fmt"
)

// MemberAPI provides member-related API methods for the open platform.
type MemberAPI struct {
	client *Client
}

// NewMemberAPI creates a new MemberAPI.
func NewMemberAPI(c *Client) *MemberAPI {
	return &MemberAPI{client: c}
}

type MemberPointsResponse struct {
	ExpiredTime       string  `json:"expiredTime"`
	ForeverPoints     float64 `json:"foreverPoints"`
	SoonExpiredPoints float64 `json:"soonExpiredPoints"`
	TotalPoints       float64 `json:"totalPoints"`
}

type CustomerBalanceResponse struct {
	CustomerID      string  `json:"customerId"`
	FreezeBalance   float64 `json:"freezeBalance"`
	GiftBalance     float64 `json:"giftBalance"`
	RechargeBalance float64 `json:"rechargeBalance"`
	TotalBalance    float64 `json:"totalBalance"`
}

type MemberDiscountResponse struct {
	CustomerID string  `json:"customerId"`
	Discount   float64 `json:"discount"`
}

type LevelExperienceResponse struct {
	CustomerID         string  `json:"customerId"`
	Level              int     `json:"level"`
	LevelName          string  `json:"levelName"`
	TotalExperienceNum float64 `json:"totalExperienceNum"`
}

type MemberRiskLevelResponse struct {
	Code      string `json:"code"`
	Type      int    `json:"type"`
	RiskLevel string `json:"riskLevel"`
}

type MemberLevelResponse struct {
	CardNo      string `json:"cardNo"`
	Level       int    `json:"level"`
	LevelName   string `json:"levelName"`
	MobilePhone string `json:"mobilePhone"`
}

type CustomerCode struct {
	Code string `json:"code"`
	Type int    `json:"type"`
}

type CustomerIDByCodeResponse struct {
	CustomerID string `json:"customerId"`
}

type CustomerIDByPhoneItem struct {
	CustomerID string `json:"customerId"`
	Phone      string `json:"phone"`
}

type CustomerIDsByPhoneResponse struct {
	List []CustomerIDByPhoneItem `json:"list"`
}

type MemberCoupon struct {
	CardID       string `json:"cardId"`
	CouponNo     string `json:"couponNo"`
	CouponType   int    `json:"couponType"`
	EndTime      string `json:"endTime"`
	TemplateID   string `json:"templateId"`
	TemplateName string `json:"templateName"`
	UseStatus    int    `json:"useStatus"`
}

type MemberCouponListResponse struct {
	List  []MemberCoupon `json:"list"`
	Total int            `json:"total"`
}

type CouponDetailItem struct {
	CardID         string  `json:"cardId"`
	CustomerID     string  `json:"customerId"`
	DiscountAmount float64 `json:"discountAmount"`
	EndAt          string  `json:"endAt"`
	FaceAmount     float64 `json:"faceAmount"`
	TemplateID     string  `json:"templateId"`
	Title          string  `json:"title"`
	UseStatus      int     `json:"useStatus"`
}

type CouponDetailListResponse struct {
	FollowPageNum int                `json:"followPageNum"`
	List          []CouponDetailItem `json:"list"`
	PageNum       int                `json:"pageNum"`
	Total         int                `json:"total"`
}

type BalanceDecreaseDetail struct {
	Amount     float64 `json:"amount"`
	BizID      string  `json:"bizId"`
	ChangeTime string  `json:"changeTime"`
	CustomerID string  `json:"customerId"`
	OrderNo    string  `json:"orderNo"`
}

type AccountFlowItem struct {
	AfterAmount float64 `json:"afterAmount"`
	Amount      float64 `json:"amount"`
	ChangeTime  string  `json:"changeTime"`
	ChangeType  string  `json:"changeType"`
	OrderNo     string  `json:"orderNo"`
	Remark      string  `json:"remark"`
}

type AccountFlowResponse struct {
	List  []AccountFlowItem `json:"list"`
	Total int               `json:"total"`
}

type PointsFlowItem struct {
	AfterPoints  float64 `json:"afterPoints"`
	ChangePoints float64 `json:"changePoints"`
	ChangeTime   string  `json:"changeTime"`
	ChangeType   string  `json:"changeType"`
	OrgType      string  `json:"orgType"`
	SourceType   string  `json:"sourceType"`
}

type PointsFlowResponse struct {
	List  []PointsFlowItem `json:"list"`
	Total int              `json:"total"`
}

type PersonalAssetResponse struct {
	Balance       float64 `json:"balance"`
	CouponCount   int     `json:"couponCount"`
	CustomerID    string  `json:"customerId"`
	Points        float64 `json:"points"`
	StoredBalance float64 `json:"storedBalance"`
}

type DepositRuleCoupon struct {
	CouponName       string `json:"couponName"`
	CouponNum        int    `json:"couponNum"`
	CouponTemplateID string `json:"couponTemplateId"`
}

type DepositRuleDetail struct {
	ChargeBalance       float64             `json:"chargeBalance"`
	ChargeGiftBalance   float64             `json:"chargeGiftBalance"`
	DepositPlan         int                 `json:"depositPlan"`
	DepositValue        float64             `json:"depositValue"`
	GiftIntegral        string              `json:"giftIntegral"`
	OrderMultiple       float64             `json:"orderMultiple"`
	PresentDepositType  int                 `json:"presentDepositType"`
	PresentDepositValue float64             `json:"presentDepositValue"`
	CouponList          []DepositRuleCoupon `json:"couponList"`
}

type DepositRule struct {
	ActivityContent string            `json:"activityContent"`
	ActivityName    string            `json:"activityName"`
	DetailInfo      DepositRuleDetail `json:"detailInfo"`
	EndDate         int64             `json:"endDate"`
	ID              int64             `json:"id"`
	StartDate       int64             `json:"startDate"`
}

type InflateStatusResponse struct {
	CustomerID         string `json:"customerId"`
	InflateBeforeLevel int    `json:"inflateBeforeLevel"`
	InflateTime        string `json:"inflateTime"`
	IsInflate          int    `json:"isInflate"`
}

type MemberInfoResponse struct {
	CardNo      string `json:"cardNo"`
	CustomerID  string `json:"customerId"`
	MobilePhone string `json:"mobilePhone"`
	NickName    string `json:"nickName"`
	OpenID      string `json:"openId"`
	RealName    string `json:"realName"`
	UnionID     string `json:"unionId"`
}

type MemberBaseInfoResponse struct {
	CardNo      string `json:"cardNo"`
	CustomerID  string `json:"customerId"`
	MobilePhone string `json:"mobilePhone"`
	NickName    string `json:"nickName"`
	RealName    string `json:"realName"`
}

type MemberSearchItem struct {
	CardNo      string `json:"cardNo"`
	CustomerID  string `json:"customerId"`
	LevelName   string `json:"levelName"`
	MobilePhone string `json:"mobilePhone"`
	NickName    string `json:"nickName"`
	UseStatus   int    `json:"useStatus"`
}

type MemberSearchResponse struct {
	List  []MemberSearchItem `json:"list"`
	Total int                `json:"total"`
}

type MemberTagItem struct {
	ID       int64  `json:"id"`
	LabelID  int64  `json:"labelId"`
	Name     string `json:"name"`
	TagName  string `json:"tagName"`
	TypeName string `json:"typeName"`
}

type MemberTagPageResponse struct {
	List  []MemberTagItem `json:"list"`
	Total int             `json:"total"`
}

type MemberTagDetailResponse struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	TagGroupID   int64  `json:"tagGroupId"`
	TagGroupName string `json:"tagGroupName"`
	TypeName     string `json:"typeName"`
}

type BrandTagListResponse struct {
	List []MemberTagItem `json:"list"`
}

type MemberTagIDsResponse struct {
	LabelIDs []int64 `json:"labelIds"`
}

type PhoneRegisterParams struct {
	Birthday            string                 `json:"birthday,omitempty"`
	ChannelsRecruitCode string                 `json:"channelsRecruitCode,omitempty"`
	CountryCode         string                 `json:"countryCode,omitempty"`
	Email               string                 `json:"email,omitempty"`
	Gender              int                    `json:"gender,omitempty"`
	MobilePhone         string                 `json:"mobilePhone"`
	MultiMark           string                 `json:"multiMark,omitempty"`
	MultiStoreID        int64                  `json:"multiStoreId,omitempty"`
	OperationScene      map[string]interface{} `json:"operationScene,omitempty"`
	RecruitChannel      int                    `json:"recruitChannel,omitempty"`
	RegAppType          int                    `json:"regAppType"`
	Remark              string                 `json:"remark,omitempty"`
	Username            string                 `json:"username,omitempty"`
}

type CustomerOpenIDResponse struct {
	OpenID  string `json:"openId"`
	UnionID string `json:"unionId"`
}

type CustomerCodeResponse struct {
	ClientCode string `json:"clientCode"`
}

type PhoneRegisterResponse struct {
	Address    string `json:"address"`
	Avatar     string `json:"avatar"`
	Birthday   string `json:"birthday"`
	City       string `json:"city"`
	CreatedAt  string `json:"createdAt"`
	CustomerID string `json:"customerId"`
	District   string `json:"district"`
	Gender     int    `json:"gender"`
	GroupID    string `json:"groupId"`
	Nickname   string `json:"nickname"`
	Phone      string `json:"phone"`
	Province   string `json:"province"`
	Username   string `json:"username"`
}

type UpdateMemberParams struct {
	Avatar   string `json:"avatar,omitempty"`
	Birthday string `json:"birthday,omitempty"`
	City     string `json:"city,omitempty"`
	Country  string `json:"country,omitempty"`
	District string `json:"district,omitempty"`
	Gender   int    `json:"gender,omitempty"`
	ID       string `json:"id"`
	IDNumber string `json:"idNumber,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	Province string `json:"province,omitempty"`
	QMFrom   string `json:"qmFrom,omitempty"`
	Username string `json:"username,omitempty"`
}

type ThirdIDRegisterParams map[string]interface{}

type StaticLabelPageItem struct {
	ID               int64  `json:"id"`
	LabelAttributed  int    `json:"labelAttributed"`
	LabelCode        string `json:"labelCode"`
	LabelGroupID     int64  `json:"labelGroupId"`
	LabelGroupType   int    `json:"labelGroupType"`
	LabelName        string `json:"labelName"`
	LabelNum         int64  `json:"labelNum"`
	LabelNumUpdateAt string `json:"labelNumUpdateAt"`
	LabelStatus      int    `json:"labelStatus"`
	LabelUpdateType  int    `json:"labelUpdateType"`
}

type StaticLabelPageResponse struct {
	List     []StaticLabelPageItem `json:"list"`
	PageNum  int                   `json:"pageNum"`
	PageSize int                   `json:"pageSize"`
	Total    int64                 `json:"total"`
}

type StaticLabelGroupItem struct {
	ID              int64  `json:"id"`
	LabelAttributed int    `json:"labelAttributed"`
	LabelGroupName  string `json:"labelGroupName"`
	LabelGroupType  int    `json:"labelGroupType"`
	OperateAccount  string `json:"operateAccount"`
	OperateName     string `json:"operateName"`
	SellerID        int64  `json:"sellerId"`
}

type TagMarkParams struct {
	CustomerIDList []string `json:"customerIdList,omitempty"`
	LabelCode      string   `json:"labelCode"`
	MarkDate       string   `json:"markDate"`
	ThirdMemberID  []string `json:"thirdMemberId,omitempty"`
}

type LabelCreateParams struct {
	LabelCode string `json:"labelCode"`
	LabelName string `json:"labelName"`
}

type LabelDeleteParams struct {
	LabelCode string `json:"labelCode"`
}

type LabelClearMembersParams struct {
	CustomerIDList []string `json:"customerIdList,omitempty"`
	LabelCode      string   `json:"labelCode"`
	ThirdMemberID  []string `json:"thirdMemberId,omitempty"`
}

type DeleteCustomerLabelParams struct {
	CustomerID      string `json:"customerId"`
	PanoramaLabelID int64  `json:"panoramaLabelId"`
}

type BrandInfoLabelItem struct {
	LabelCode  string `json:"labelCode"`
	LabelID    int64  `json:"labelId"`
	LabelLevel int    `json:"labelLevel"`
	LabelName  string `json:"labelName"`
}

type BrandInfoMemberLevelItem struct {
	MemberCardName  string `json:"memberCardName"`
	MemberLevel     int    `json:"memberLevel"`
	MemberLevelName string `json:"memberLevelName"`
}

type BrandInfoPaidCardItem struct {
	BeginAt    string `json:"beginAt"`
	CardName   string `json:"cardName"`
	ExpiredAt  string `json:"expiredAt"`
	PaidCardID string `json:"paidCardId"`
}

type BrandInfoResponse struct {
	LabelList                  []BrandInfoLabelItem       `json:"labelList"`
	MemberLevelList            []BrandInfoMemberLevelItem `json:"memberLevelList"`
	PaidBenefitsCardDetailList []BrandInfoPaidCardItem    `json:"paidBenefitsCardDetailList"`
}

type BlacklistStatusResponse struct {
	CustomerType int    `json:"customerType"`
	RegisterTime string `json:"registerTime"`
}

type FreezeRecordResponse struct {
	CustomerID    string `json:"customerId"`
	FreezeEndTime string `json:"freezeEndTime"`
	FreezeStatus  int    `json:"freezeStatus"`
}

type CanLogoffResponse struct {
	HasOrder bool `json:"hasOrder"`
}

type WeComGroupItem struct {
	ChatID    string `json:"chatId"`
	Count     int    `json:"count"`
	GroupName string `json:"groupName"`
	JoinAt    string `json:"joinAt"`
}

type WeComEmployeeItem struct {
	AddTime      string `json:"addTime"`
	EmpAvatar    string `json:"empAvatar"`
	EmployeeName string `json:"employeeName"`
	UserID       string `json:"userId"`
}

type WeComCustomerInfoResponse struct {
	Avatar                  string              `json:"avatar"`
	BelongCustomerGroupList []WeComGroupItem    `json:"belongCustomerGroupList"`
	BelongEmployees         []WeComEmployeeItem `json:"belongEmployees"`
	CustomerID              string              `json:"customerId"`
	Nickname                string              `json:"nickname"`
}

type ActivitySignResponse struct {
	TodaySign int `json:"todaySign"`
}

type AccountLevelItem struct {
	CustomerID int64 `json:"customerId"`
	Level      int   `json:"level"`
}

type CustomerConditionItem struct {
	ID          string  `json:"id"`
	MobilePhone string  `json:"mobilePhone"`
	Nickname    string  `json:"nickname"`
	Username    string  `json:"username"`
	Level       int     `json:"level"`
	LevelName   string  `json:"levelName"`
	Balance     float64 `json:"balance"`
	PointNum    float64 `json:"pointNum"`
	CouponNum   int     `json:"couponNum"`
	GiftCardNum int     `json:"giftCardNum"`
}

type MemberRechargeParams map[string]interface{}

type MemberBalanceDecreaseParams map[string]interface{}

type MemberBalanceReverseParams map[string]interface{}

type MemberPointsDecreaseParams map[string]interface{}

type MemberPointsReverseParams map[string]interface{}

type MemberPointsAddParams map[string]interface{}

type MemberConsumeParams map[string]interface{}

type MemberConsumeReverseParams map[string]interface{}

type MemberCouponWriteOffParams map[string]interface{}

type MemberCouponReverseParams map[string]interface{}

type OfflineBalanceOperationParams map[string]interface{}

type MemberRechargeReverseParams map[string]interface{}

func decodeData[T any](resp *APIResponse, action string) (*T, error) {
	var result T
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("parse %s: %w", action, err)
	}
	return &result, nil
}

func (a *MemberAPI) GetCustomerPoints(ctx context.Context, customerID, orderTime string) (*MemberPointsResponse, error) {
	resp, err := a.client.Call(ctx, "v3/crm/points/getCustomerPoints", map[string]string{
		"customerId": customerID,
		"orderTime":  orderTime,
	})
	if err != nil {
		return nil, fmt.Errorf("get customer points: %w", err)
	}
	return decodeData[MemberPointsResponse](resp, "customer points")
}

func (a *MemberAPI) RechargeBalance(ctx context.Context, params MemberRechargeParams) error {
	_, err := a.client.Call(ctx, "v3/bsns/customerBalance/userRecharge", params)
	if err != nil {
		return fmt.Errorf("recharge member balance: %w", err)
	}
	return nil
}

func (a *MemberAPI) DecreaseBalance(ctx context.Context, params MemberBalanceDecreaseParams) error {
	_, err := a.client.Call(ctx, "v3/crm/account/decreaseBalance", params)
	if err != nil {
		return fmt.Errorf("decrease member balance: %w", err)
	}
	return nil
}

func (a *MemberAPI) ReverseBalance(ctx context.Context, params MemberBalanceReverseParams) error {
	_, err := a.client.Call(ctx, "v3/crm/account/decreaseBalanceReverse", params)
	if err != nil {
		return fmt.Errorf("reverse member balance: %w", err)
	}
	return nil
}

func (a *MemberAPI) ReducePoints(ctx context.Context, params MemberPointsDecreaseParams) error {
	_, err := a.client.Call(ctx, "v3/crm/points/reducePoints", params)
	if err != nil {
		return fmt.Errorf("reduce member points: %w", err)
	}
	return nil
}

func (a *MemberAPI) ReversePoints(ctx context.Context, params MemberPointsReverseParams) error {
	_, err := a.client.Call(ctx, "v3/crm/points/reversePoints", params)
	if err != nil {
		return fmt.Errorf("reverse member points: %w", err)
	}
	return nil
}

func (a *MemberAPI) AddPoints(ctx context.Context, params MemberPointsAddParams) error {
	_, err := a.client.Call(ctx, "v3/crm/points/addCustomerPoints", params)
	if err != nil {
		return fmt.Errorf("add member points: %w", err)
	}
	return nil
}

func (a *MemberAPI) CustomerConsume(ctx context.Context, params MemberConsumeParams) error {
	_, err := a.client.Call(ctx, "v3/crm/customer/customerConsume", params)
	if err != nil {
		return fmt.Errorf("member consume: %w", err)
	}
	return nil
}

func (a *MemberAPI) ConsumeReverse(ctx context.Context, params MemberConsumeReverseParams) error {
	_, err := a.client.Call(ctx, "v3/crm/customer/consumeReverse", params)
	if err != nil {
		return fmt.Errorf("reverse member consume: %w", err)
	}
	return nil
}

func (a *MemberAPI) WriteOffCoupon(ctx context.Context, params MemberCouponWriteOffParams) error {
	_, err := a.client.Call(ctx, "v3/crm/coupon/writeOffCoupon", params)
	if err != nil {
		return fmt.Errorf("write off member coupon: %w", err)
	}
	return nil
}

func (a *MemberAPI) ReverseCoupon(ctx context.Context, params MemberCouponReverseParams) error {
	_, err := a.client.Call(ctx, "v3/crm/coupon/couponReverse", params)
	if err != nil {
		return fmt.Errorf("reverse member coupon: %w", err)
	}
	return nil
}

func (a *MemberAPI) OfflineOperateBalance(ctx context.Context, params OfflineBalanceOperationParams) error {
	_, err := a.client.Call(ctx, "v3/newPattern/crmCenter/post/customer-account/v2/offline-operation-balance", params)
	if err != nil {
		return fmt.Errorf("offline operate member balance: %w", err)
	}
	return nil
}

func (a *MemberAPI) RechargeReverse(ctx context.Context, params MemberRechargeReverseParams) error {
	_, err := a.client.Call(ctx, "v3/bsns/customerBalance/userRechargeReverse", params)
	if err != nil {
		return fmt.Errorf("reverse member recharge: %w", err)
	}
	return nil
}

func (a *MemberAPI) GetCustomerCoupons(ctx context.Context, customerID, cardID, templateID string, useStatus int) (*MemberCouponListResponse, error) {
	params := map[string]interface{}{"customerId": customerID}
	if cardID != "" {
		params["cardId"] = cardID
	}
	if templateID != "" {
		params["templateId"] = templateID
	}
	if useStatus != 0 {
		params["useStatus"] = useStatus
	}
	resp, err := a.client.Call(ctx, "v3/crm/coupon/getCustomerCoupons", params)
	if err != nil {
		return nil, fmt.Errorf("get customer coupons: %w", err)
	}
	return decodeData[MemberCouponListResponse](resp, "customer coupons")
}

func (a *MemberAPI) GetDecreaseBalanceBiz(ctx context.Context, customerID, bizID string) (*BalanceDecreaseDetail, error) {
	resp, err := a.client.Call(ctx, "v3/crm/account/getDecreaseBalanceBiz", map[string]string{
		"bizId":      bizID,
		"customerId": customerID,
	})
	if err != nil {
		return nil, fmt.Errorf("get balance decrease detail: %w", err)
	}
	return decodeData[BalanceDecreaseDetail](resp, "balance decrease detail")
}

func (a *MemberAPI) QueryAccountFlow(ctx context.Context, params map[string]interface{}) (*AccountFlowResponse, error) {
	resp, err := a.client.Call(ctx, "v3/crm/account/queryAccountFlow", params)
	if err != nil {
		return nil, fmt.Errorf("query account flow: %w", err)
	}
	return decodeData[AccountFlowResponse](resp, "account flow")
}

func (a *MemberAPI) GetLevelConsumeDiscount(ctx context.Context, customerID string) (*MemberDiscountResponse, error) {
	resp, err := a.client.Call(ctx, "v3/crm/customer/getLevelConsumeDiscount", map[string]string{"customerId": customerID})
	if err != nil {
		return nil, fmt.Errorf("get level consume discount: %w", err)
	}
	return decodeData[MemberDiscountResponse](resp, "member discount")
}

func (a *MemberAPI) GetLevelExperience(ctx context.Context, customerID string) (*LevelExperienceResponse, error) {
	resp, err := a.client.Call(ctx, "v3/crm/customer/getLevelExperience", map[string]string{"customerId": customerID})
	if err != nil {
		return nil, fmt.Errorf("get level experience: %w", err)
	}
	return decodeData[LevelExperienceResponse](resp, "level experience")
}

func (a *MemberAPI) GetCustomerBalance(ctx context.Context, customerID string) (*CustomerBalanceResponse, error) {
	resp, err := a.client.Call(ctx, "v3/crm/customer/getCustomerBalance", map[string]string{"customerId": customerID})
	if err != nil {
		return nil, fmt.Errorf("get customer balance: %w", err)
	}
	return decodeData[CustomerBalanceResponse](resp, "customer balance")
}

func (a *MemberAPI) GetCustomerCouponList(ctx context.Context, customerID string, couponType, pageNum, pageSize, useStatus int) (*MemberCouponListResponse, error) {
	params := map[string]interface{}{
		"customerId": customerID,
		"pageNum":    pageNum,
		"pageSize":   pageSize,
	}
	if couponType != 0 {
		params["couponType"] = couponType
	}
	if useStatus != 0 {
		params["useStatus"] = useStatus
	}
	resp, err := a.client.Call(ctx, "v3/crm/coupon/getCustomerCouponList", params)
	if err != nil {
		return nil, fmt.Errorf("get customer coupon list: %w", err)
	}
	return decodeData[MemberCouponListResponse](resp, "customer coupon list")
}

func (a *MemberAPI) GetCouponDetailList(ctx context.Context, params map[string]interface{}) (*CouponDetailListResponse, error) {
	resp, err := a.client.Call(ctx, "v3/newPattern/crmCenter/get/customer-coupons/get_my_coupons_v2", params)
	if err != nil {
		return nil, fmt.Errorf("get coupon detail list: %w", err)
	}
	return decodeData[CouponDetailListResponse](resp, "coupon detail list")
}

func (a *MemberAPI) GetCrmPointsFlow(ctx context.Context, params map[string]interface{}) (*PointsFlowResponse, error) {
	resp, err := a.client.Call(ctx, "v3/crm/customer/getCrmPointsFlow", params)
	if err != nil {
		return nil, fmt.Errorf("get points flow: %w", err)
	}
	return decodeData[PointsFlowResponse](resp, "points flow")
}

func (a *MemberAPI) GetPersonalAsset(ctx context.Context, customerID string) (*PersonalAssetResponse, error) {
	resp, err := a.client.Call(ctx, "v3/crm/customerInfo/personalAsset", map[string]string{"customerId": customerID})
	if err != nil {
		return nil, fmt.Errorf("get personal asset: %w", err)
	}
	return decodeData[PersonalAssetResponse](resp, "personal asset")
}

func (a *MemberAPI) GetStoreDepositRules(ctx context.Context, shopCode string, depositValue float64) (*[]DepositRule, error) {
	params := map[string]interface{}{}
	if shopCode != "" {
		params["shopCode"] = shopCode
	}
	if depositValue != 0 {
		params["depositValue"] = depositValue
	}
	resp, err := a.client.Call(ctx, "v3/marketing/deposit/getStoreOrShopDepositList", params)
	if err != nil {
		return nil, fmt.Errorf("get store deposit rules: %w", err)
	}
	return decodeData[[]DepositRule](resp, "store deposit rules")
}

func (a *MemberAPI) GetAssetInflateDetail(ctx context.Context, customerID string) (*InflateStatusResponse, error) {
	resp, err := a.client.Call(ctx, "v3/crm/customer/assetInflateDetail", map[string]string{"customerId": customerID})
	if err != nil {
		return nil, fmt.Errorf("get asset inflate detail: %w", err)
	}
	return decodeData[InflateStatusResponse](resp, "asset inflate detail")
}

func (a *MemberAPI) GetCustomerInfo(ctx context.Context, params map[string]interface{}) (*MemberInfoResponse, error) {
	resp, err := a.client.Call(ctx, "v3/crm/customer/getCustomerInfo", params)
	if err != nil {
		return nil, fmt.Errorf("get customer info: %w", err)
	}
	return decodeData[MemberInfoResponse](resp, "customer info")
}

func (a *MemberAPI) GetCustomerIDByCode(ctx context.Context, code string, codeType int) (*CustomerIDByCodeResponse, error) {
	resp, err := a.client.Call(ctx, "v3/crm/customer/getCustomerIdByCode", map[string]CustomerCode{
		"customerCode": {Code: code, Type: codeType},
	})
	if err != nil {
		return nil, fmt.Errorf("get customer id by code: %w", err)
	}
	return decodeData[CustomerIDByCodeResponse](resp, "customer id by code")
}

func (a *MemberAPI) GetCustomerIDByPhone(ctx context.Context, phones []string) (*CustomerIDsByPhoneResponse, error) {
	resp, err := a.client.Call(ctx, "v3/crm/customer/getCustomerIdByPhone", map[string][]string{"phone": phones})
	if err != nil {
		return nil, fmt.Errorf("get customer id by phone: %w", err)
	}
	return decodeData[CustomerIDsByPhoneResponse](resp, "customer id by phone")
}

func (a *MemberAPI) GetBaseInfo(ctx context.Context, customerID string, infoType int) (*MemberBaseInfoResponse, error) {
	params := map[string]interface{}{"customerId": customerID}
	if infoType != 0 {
		params["type"] = infoType
	}
	resp, err := a.client.Call(ctx, "v3/crm/customer/getBaseInfo", params)
	if err != nil {
		return nil, fmt.Errorf("get base info: %w", err)
	}
	return decodeData[MemberBaseInfoResponse](resp, "base info")
}

func (a *MemberAPI) GetCustomerOpenID(ctx context.Context, customerID string, infoType int) (*CustomerOpenIDResponse, error) {
	params := map[string]interface{}{"customerId": customerID}
	if infoType != 0 {
		params["type"] = infoType
	}
	resp, err := a.client.Call(ctx, "v3/crm/customer/getCustomerOpenId", params)
	if err != nil {
		return nil, fmt.Errorf("get customer open id: %w", err)
	}
	return decodeData[CustomerOpenIDResponse](resp, "customer open id")
}

func (a *MemberAPI) GetCustomerCode(ctx context.Context, customerID int64) (*CustomerCodeResponse, error) {
	resp, err := a.client.Call(ctx, "v3/crm/customer/getCustomerCode", map[string]int64{"customerId": customerID})
	if err != nil {
		return nil, fmt.Errorf("get customer code: %w", err)
	}
	return decodeData[CustomerCodeResponse](resp, "customer code")
}

func (a *MemberAPI) SearchMembers(ctx context.Context, params map[string]interface{}) (*MemberSearchResponse, error) {
	resp, err := a.client.Call(ctx, "v3/crm/customer/allCustomerInfo", params)
	if err != nil {
		return nil, fmt.Errorf("search members: %w", err)
	}
	return decodeData[MemberSearchResponse](resp, "member search")
}

func (a *MemberAPI) GetRiskLevel(ctx context.Context, code string, codeType int) (*MemberRiskLevelResponse, error) {
	resp, err := a.client.Call(ctx, "v3/crm/customer/riskLevel", map[string]interface{}{"code": code, "type": codeType})
	if err != nil {
		return nil, fmt.Errorf("get risk level: %w", err)
	}
	return decodeData[MemberRiskLevelResponse](resp, "member risk level")
}

func (a *MemberAPI) GetMemberLevel(ctx context.Context, mobilePhone string) (*MemberLevelResponse, error) {
	resp, err := a.client.Call(ctx, "v3/crm/customer/member/detail", map[string]string{"mobilePhone": mobilePhone})
	if err != nil {
		return nil, fmt.Errorf("get member level: %w", err)
	}
	return decodeData[MemberLevelResponse](resp, "member level")
}

func (a *MemberAPI) PhoneRegister(ctx context.Context, params PhoneRegisterParams) (*PhoneRegisterResponse, error) {
	resp, err := a.client.Call(ctx, "v3/crm/customer/phoneRegister", params)
	if err != nil {
		return nil, fmt.Errorf("phone register: %w", err)
	}
	return decodeData[PhoneRegisterResponse](resp, "phone register")
}

func (a *MemberAPI) UpdateMember(ctx context.Context, params UpdateMemberParams) error {
	_, err := a.client.Call(ctx, "v3/crm/customer/updateMember", params)
	if err != nil {
		return fmt.Errorf("update member: %w", err)
	}
	return nil
}

func (a *MemberAPI) RegisterByThirdID(ctx context.Context, params ThirdIDRegisterParams) (*PhoneRegisterResponse, error) {
	resp, err := a.client.Call(ctx, "v3/crm/customer/customerRegister", params)
	if err != nil {
		return nil, fmt.Errorf("register customer by third id: %w", err)
	}
	return decodeData[PhoneRegisterResponse](resp, "customer register")
}

func (a *MemberAPI) SendSmsCaptcha(ctx context.Context, params map[string]interface{}) error {
	_, err := a.client.Call(ctx, "v3/message/commonSendSmsCaptcha", params)
	if err != nil {
		return fmt.Errorf("send sms captcha: %w", err)
	}
	return nil
}

func (a *MemberAPI) CheckSmsCaptcha(ctx context.Context, params map[string]interface{}) error {
	_, err := a.client.Call(ctx, "v3/message/commonCheckSmsCaptcha", params)
	if err != nil {
		return fmt.Errorf("check sms captcha: %w", err)
	}
	return nil
}

func (a *MemberAPI) PageTagsByCustomerID(ctx context.Context, customerID string, pageNum, pageSize int) (*MemberTagPageResponse, error) {
	resp, err := a.client.Call(ctx, "v3/cdp/panoramaLabel/pageByCustomerId", map[string]interface{}{
		"customerId": customerID,
		"pageNum":    pageNum,
		"pageSize":   pageSize,
	})
	if err != nil {
		return nil, fmt.Errorf("list member tags: %w", err)
	}
	return decodeData[MemberTagPageResponse](resp, "member tags")
}

func (a *MemberAPI) ListStaticLabelsByGroup(ctx context.Context, labelGroupID int64, labelName string, pageNum, pageSize int) (*StaticLabelPageResponse, error) {
	params := map[string]interface{}{
		"labelGroupId": labelGroupID,
		"pageNum":      pageNum,
		"pageSize":     pageSize,
	}
	if labelName != "" {
		params["labelName"] = labelName
	}
	resp, err := a.client.Call(ctx, "v3/cdp/panoramaLabel/list", params)
	if err != nil {
		return nil, fmt.Errorf("list static labels by group: %w", err)
	}
	return decodeData[StaticLabelPageResponse](resp, "static labels")
}

func (a *MemberAPI) ListStaticLabelGroups(ctx context.Context, labelAttributed int, labelGroupName string, filterEmpty int) (*[]StaticLabelGroupItem, error) {
	params := map[string]interface{}{
		"labelAttributed": labelAttributed,
	}
	if labelGroupName != "" {
		params["labelGroupName"] = labelGroupName
	}
	if filterEmpty != 0 {
		params["filterEmptyTagGroups"] = filterEmpty
	}
	resp, err := a.client.Call(ctx, "v3/cdp/panoramaLabel/listGroup", params)
	if err != nil {
		return nil, fmt.Errorf("list static label groups: %w", err)
	}
	return decodeData[[]StaticLabelGroupItem](resp, "static label groups")
}

func (a *MemberAPI) GetTagDetail(ctx context.Context, id int64) (*MemberTagDetailResponse, error) {
	resp, err := a.client.Call(ctx, "v3/cdp/panoramaLabel/detail", map[string]int64{"id": id})
	if err != nil {
		return nil, fmt.Errorf("get tag detail: %w", err)
	}
	return decodeData[MemberTagDetailResponse](resp, "tag detail")
}

func (a *MemberAPI) GetBrandTags(ctx context.Context) (*BrandTagListResponse, error) {
	resp, err := a.client.Call(ctx, "v3/cdp/panoramaLabel/panoramaLabelList", map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("get brand tags: %w", err)
	}
	return decodeData[BrandTagListResponse](resp, "brand tags")
}

func (a *MemberAPI) GetMemberTagIDs(ctx context.Context, customerID string) (*MemberTagIDsResponse, error) {
	resp, err := a.client.Call(ctx, "v3/newPattern/cdpCenter/get/panorama-label/list-by-customerId", map[string]string{"customerId": customerID})
	if err != nil {
		return nil, fmt.Errorf("get member tag ids: %w", err)
	}
	return decodeData[MemberTagIDsResponse](resp, "member tag ids")
}

func (a *MemberAPI) MarkPanoramaLabel(ctx context.Context, params TagMarkParams) error {
	_, err := a.client.Call(ctx, "v3/cdp/panoramaLabel/panoramaLabelMark", params)
	if err != nil {
		return fmt.Errorf("mark panorama label: %w", err)
	}
	return nil
}

func (a *MemberAPI) CreatePanoramaLabel(ctx context.Context, params LabelCreateParams) error {
	_, err := a.client.Call(ctx, "v3/cdp/panoramaLabel/panoramaLabelCreate", params)
	if err != nil {
		return fmt.Errorf("create panorama label: %w", err)
	}
	return nil
}

func (a *MemberAPI) DeletePanoramaLabel(ctx context.Context, params LabelDeleteParams) error {
	_, err := a.client.Call(ctx, "v3/cdp/panoramaLabel/panoramaLabelDelete", params)
	if err != nil {
		return fmt.Errorf("delete panorama label: %w", err)
	}
	return nil
}

func (a *MemberAPI) ClearPanoramaLabelMembers(ctx context.Context, params LabelClearMembersParams) error {
	_, err := a.client.Call(ctx, "v3/cdp/panoramaLabel/panoramaLabelDelete", params)
	if err != nil {
		return fmt.Errorf("clear panorama label members: %w", err)
	}
	return nil
}

func (a *MemberAPI) DeleteCustomerLabel(ctx context.Context, params DeleteCustomerLabelParams) error {
	_, err := a.client.Call(ctx, "v3/cdp/panoramaLabel/deleteCustomerLabel", params)
	if err != nil {
		return fmt.Errorf("delete customer label: %w", err)
	}
	return nil
}

func (a *MemberAPI) GetBrandInfo(ctx context.Context) (*BrandInfoResponse, error) {
	resp, err := a.client.Call(ctx, "v3/crm/customer/brandInfo", map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("get brand info: %w", err)
	}
	return decodeData[BrandInfoResponse](resp, "brand info")
}

func (a *MemberAPI) BatchAddBlackList(ctx context.Context, params map[string]interface{}) error {
	_, err := a.client.Call(ctx, "v3/crm/customer/batchAddBlackList", params)
	if err != nil {
		return fmt.Errorf("batch add blacklist: %w", err)
	}
	return nil
}

func (a *MemberAPI) BatchCancelBlackList(ctx context.Context, params map[string]interface{}) error {
	_, err := a.client.Call(ctx, "v3/crm/customer/batchCancelBlackList", params)
	if err != nil {
		return fmt.Errorf("batch cancel blacklist: %w", err)
	}
	return nil
}

func (a *MemberAPI) QueryMemberBlack(ctx context.Context, customerID string) (*BlacklistStatusResponse, error) {
	resp, err := a.client.Call(ctx, "v3/crm/customer/queryMemberBlack", map[string]string{"customerId": customerID})
	if err != nil {
		return nil, fmt.Errorf("query member black: %w", err)
	}
	return decodeData[BlacklistStatusResponse](resp, "member black")
}

func (a *MemberAPI) FreezeCustomer(ctx context.Context, customerID, reason string) error {
	_, err := a.client.Call(ctx, "v3/crm/customer/logoffFreeze", map[string]string{
		"customerId": customerID,
		"reason":     reason,
	})
	if err != nil {
		return fmt.Errorf("freeze customer: %w", err)
	}
	return nil
}

func (a *MemberAPI) UnfreezeCustomer(ctx context.Context, customerID string) error {
	_, err := a.client.Call(ctx, "v3/crm/customer/logoffUnfreeze", map[string]string{"customerId": customerID})
	if err != nil {
		return fmt.Errorf("unfreeze customer: %w", err)
	}
	return nil
}

func (a *MemberAPI) LogoffCustomer(ctx context.Context, customerID, reason string) error {
	_, err := a.client.Call(ctx, "v3/crm/customer/logOff", map[string]string{
		"customerId": customerID,
		"reason":     reason,
	})
	if err != nil {
		return fmt.Errorf("logoff customer: %w", err)
	}
	return nil
}

func (a *MemberAPI) UpdateCustomerPhone(ctx context.Context, params map[string]interface{}) error {
	_, err := a.client.Call(ctx, "v3/crm/customer/updatePhone", params)
	if err != nil {
		return fmt.Errorf("update customer phone: %w", err)
	}
	return nil
}

func (a *MemberAPI) QueryFreezeRecord(ctx context.Context, customerID string) (*FreezeRecordResponse, error) {
	resp, err := a.client.Call(ctx, "v3/crm/customer/queryFreezeRecord", map[string]string{"customerId": customerID})
	if err != nil {
		return nil, fmt.Errorf("query freeze record: %w", err)
	}
	return decodeData[FreezeRecordResponse](resp, "freeze record")
}

func (a *MemberAPI) CanLogoff(ctx context.Context, customerID string, bizType int) (*CanLogoffResponse, error) {
	resp, err := a.client.Call(ctx, "v3/newPattern/accountCenter/post/account/can-logoff", map[string]interface{}{
		"customerId": customerID,
		"bizType":    bizType,
	})
	if err != nil {
		return nil, fmt.Errorf("check can logoff: %w", err)
	}
	return decodeData[CanLogoffResponse](resp, "can logoff")
}

func (a *MemberAPI) GetWeComCustomerInfo(ctx context.Context, customerID int64) (*WeComCustomerInfoResponse, error) {
	resp, err := a.client.Call(ctx, "v3/crm/customer/weComCustomerInfo", map[string]int64{"customerId": customerID})
	if err != nil {
		return nil, fmt.Errorf("get wecom customer info: %w", err)
	}
	return decodeData[WeComCustomerInfoResponse](resp, "wecom customer info")
}

func (a *MemberAPI) QueryActivitySign(ctx context.Context, userID int64, activityID int64) (*ActivitySignResponse, error) {
	params := map[string]interface{}{"userId": userID}
	if activityID != 0 {
		params["activityId"] = activityID
	}
	resp, err := a.client.Call(ctx, "v3/crm/customer/queryActivitySign", params)
	if err != nil {
		return nil, fmt.Errorf("query activity sign: %w", err)
	}
	return decodeData[ActivitySignResponse](resp, "activity sign")
}

func (a *MemberAPI) QueryAccountLevel(ctx context.Context, customerID int64) (*[]AccountLevelItem, error) {
	resp, err := a.client.Call(ctx, "v3/newPattern/crmCenter/get/customer-account-card/query-account-level", map[string]int64{"customerId": customerID})
	if err != nil {
		return nil, fmt.Errorf("query account level: %w", err)
	}
	return decodeData[[]AccountLevelItem](resp, "account level")
}

func (a *MemberAPI) QueryCustomerCondition(ctx context.Context, customerIDs []int64, conditions []int) (*[]CustomerConditionItem, error) {
	resp, err := a.client.Call(ctx, "v3/newPattern/crmCenter/post/customer-condition/query", map[string]interface{}{
		"customerIdList": customerIDs,
		"condition":      conditions,
	})
	if err != nil {
		return nil, fmt.Errorf("query customer condition: %w", err)
	}
	return decodeData[[]CustomerConditionItem](resp, "customer condition")
}
