package client

import (
	"context"
	"encoding/json"
	"fmt"
)

// StoreAPI provides org/store related API methods.
type StoreAPI struct {
	client *Client
}

// NewStoreAPI creates a new StoreAPI.
func NewStoreAPI(c *Client) *StoreAPI {
	return &StoreAPI{client: c}
}

type ShopImage struct {
	URL string `json:"url"`
}

type ShopWorktime struct {
	Time []string `json:"time"`
}

type ShopOpentime struct {
	Worktime []ShopWorktime `json:"worktime"`
	Workweek []int          `json:"workweek"`
}

type ShopChannelOpenTime struct {
	ChannelID int            `json:"channelId"`
	Opentime  []ShopOpentime `json:"opentime"`
	Status    int            `json:"status"`
}

type ShopDetail struct {
	ID             int64                 `json:"id"`
	Code           string                `json:"code"`
	Name           string                `json:"name"`
	Address        string                `json:"address"`
	FullAddress    string                `json:"fullAddress"`
	ContactName    string                `json:"contactName"`
	ContactPhone   string                `json:"contactPhone"`
	ContactTel     string                `json:"contactTel"`
	ManagerPhone   string                `json:"managerPhone"`
	ProvinceName   string                `json:"provinceName"`
	CityName       string                `json:"cityName"`
	DistrictName   string                `json:"districtName"`
	OpenStatus     int                   `json:"openStatus"`
	OperateStatus  string                `json:"operateStatus"`
	ManagerStatus  int                   `json:"managerStatus"`
	IsAppletShow   int                   `json:"isAppletShow"`
	IsEat          int                   `json:"isEat"`
	IsTakeout      int                   `json:"isTakeout"`
	IsAppoint      int                   `json:"isAppoint"`
	Lat            string                `json:"lat"`
	Lng            string                `json:"lng"`
	Images         []ShopImage           `json:"images"`
	Opentimes      []ShopChannelOpenTime `json:"opentimes"`
	LabelIDs       []int64               `json:"lableIds"`
	AliPayShopCode string                `json:"aliPayShopCode"`
}

type ShopListParams struct {
	Keyfield         string  `json:"keyfield,omitempty"`
	Keyword          string  `json:"keyword,omitempty"`
	LabelID          int64   `json:"lableId,omitempty"`
	PageNum          int     `json:"pageNum"`
	PageSize         int     `json:"pageSize"`
	Search           string  `json:"search,omitempty"`
	ShopIDs          []int64 `json:"shopIds,omitempty"`
	Type             int     `json:"type,omitempty"`
	TypeID           int64   `json:"typeId,omitempty"`
	ContainCloseFlag int     `json:"containCloseFlag,omitempty"`
}

type ShopListResponse struct {
	List  []ShopDetail `json:"list"`
	Total int          `json:"total"`
}

type TakeoutShopMapping struct {
	ExternalShopID string `json:"externalShopId"`
	PlatformType   int    `json:"platformType"`
	ShopCode       string `json:"shopCode"`
	ShopID         int64  `json:"shopId"`
}

type TakeoutShopMappingResponse struct {
	List  []TakeoutShopMapping `json:"list"`
	Total int                  `json:"total"`
}

type TakeoutShopMappingParams struct {
	PageNum      int `json:"pageNum,omitempty"`
	PageSize     int `json:"pageSize,omitempty"`
	PlatformType int `json:"platformType"`
}

type ShopSyncDistrict struct {
	Code        string `json:"code"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Name        string `json:"name"`
	Seq         int    `json:"seq,omitempty"`
}

type ShopSyncExtData struct {
	ID    int64  `json:"id"`
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type ShopSyncImage struct {
	URL string `json:"url"`
}

type ShopSyncWorktime struct {
	Time []string `json:"time"`
}

type ShopSyncOpenTime struct {
	Worktime []ShopSyncWorktime `json:"worktime"`
	Workweek []int              `json:"workweek"`
}

type ShopSyncChannelOpenTime struct {
	ChannelID int                `json:"channelId"`
	Status    int                `json:"status"`
	Opentime  []ShopSyncOpenTime `json:"opentime"`
}

// ShopSyncRequest mirrors the open platform payload for shopInfoSync.
type ShopSyncRequest struct {
	Name          string                    `json:"name"`
	Address       string                    `json:"address,omitempty"`
	AppletCodeURL string                    `json:"appletCodeUrl,omitempty"`
	IsAppletShow  *int                      `json:"isAppletShow,omitempty"`
	ChannelIDs    []int                     `json:"channelIds,omitempty"`
	CompName      string                    `json:"compName,omitempty"`
	ContactName   string                    `json:"contactName,omitempty"`
	ContactPhone  string                    `json:"contactPhone,omitempty"`
	ContactTel    string                    `json:"contactTel,omitempty"`
	IDCard        string                    `json:"idCard,omitempty"`
	CityID        int64                     `json:"cityId"`
	Code          string                    `json:"code,omitempty"`
	DistrictID    int64                     `json:"districtId"`
	FacilityIDs   []int64                   `json:"facilityIds,omitempty"`
	FoodLicense   []string                  `json:"foodLicense,omitempty"`
	IsAppoint     *int                      `json:"isAppoint,omitempty"`
	IsEat         *int                      `json:"isEat,omitempty"`
	IsTakeout     *int                      `json:"isTakeout,omitempty"`
	IsTest        int                       `json:"isTest,omitempty"`
	LabelNames    []string                  `json:"lableNames,omitempty"`
	Lat           string                    `json:"lat"`
	Lng           string                    `json:"lng"`
	ManagerPhone  string                    `json:"managerPhone,omitempty"`
	ManagerStatus int                       `json:"managerStatus"`
	Notice        string                    `json:"notice,omitempty"`
	Notices       []string                  `json:"notices,omitempty"`
	OpenTime      string                    `json:"openTime,omitempty"`
	OrgID         int64                     `json:"orgId,omitempty"`
	OutHotline    string                    `json:"outHotline,omitempty"`
	PerPrice      string                    `json:"perPrice,omitempty"`
	ProvinceID    int64                     `json:"provinceId"`
	SapCode       string                    `json:"sapCode,omitempty"`
	ShareContent  string                    `json:"shareContent,omitempty"`
	ShareImage    string                    `json:"shareImage,omitempty"`
	Districts     []ShopSyncDistrict        `json:"districts,omitempty"`
	ExtData       []ShopSyncExtData         `json:"extData,omitempty"`
	Images        []ShopSyncImage           `json:"images,omitempty"`
	Opentimes     []ShopSyncChannelOpenTime `json:"opentimes,omitempty"`
}

type ShopCodeToIDResponse struct {
	ID int64 `json:"id"`
}

type BoolResult struct {
	Data bool `json:"data"`
}

type BatchConfigItem struct {
	FieldCode  string `json:"fieldCode"`
	FieldValue string `json:"fieldValue"`
}

type StoreConfigBatchResponse struct {
	ConfigList []BatchConfigItem `json:"configList"`
	StoreID    int64             `json:"storeId"`
}

type SellerConfigBatchResponse struct {
	ConfigList []BatchConfigItem `json:"configList"`
}

type ShopExtDataItem struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ShopLabel struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type OrgTreeLabel struct {
	ID      int64  `json:"id"`
	BrandID int64  `json:"brandId"`
	Name    string `json:"name"`
}

type OrgTreeShop struct {
	ID      int64          `json:"id"`
	BrandID int64          `json:"brandId"`
	Code    string         `json:"code"`
	Name    string         `json:"name"`
	IsShow  int            `json:"isShow"`
	Labels  []OrgTreeLabel `json:"lables"`
}

type OrgTreeNode struct {
	BrandID          int64         `json:"brandId"`
	BrandName        string        `json:"brandName"`
	IsShow           int           `json:"isShow"`
	OrgList          []interface{} `json:"orgList"`
	UnBoundShopCount int           `json:"unBoundShopCount"`
	UnBoundShopList  []OrgTreeShop `json:"unBoundShopList"`
}

type ShopTeam struct {
	ID                 int64      `json:"id"`
	ManagerAccountID   int64      `json:"managerAccountId"`
	ManagerAccountName string     `json:"managerAccountName"`
	Name               string     `json:"name"`
	Num                int        `json:"num"`
	Path               string     `json:"path"`
	Pid                int64      `json:"pid"`
	Sort               int        `json:"sort"`
	Children           []ShopTeam `json:"children"`
}

type ShopDeptNode struct {
	ID          int64          `json:"id"`
	LeaderID    int64          `json:"leaderId"`
	LeaderName  string         `json:"leaderName"`
	Name        string         `json:"name"`
	Path        string         `json:"path"`
	Pid         int64          `json:"pid"`
	Seq         int            `json:"seq"`
	Type        int            `json:"type"`
	SubDeptTree []ShopDeptNode `json:"subDeptTree"`
}

func (a *StoreAPI) GetShopDetailByCode(ctx context.Context, code string) (*ShopDetail, error) {
	resp, err := a.client.Call(ctx, "v3/org/shop/getShopDetail", map[string]string{"code": code})
	if err != nil {
		return nil, fmt.Errorf("get shop detail by code: %w", err)
	}
	var result ShopDetail
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("parse shop detail by code: %w", err)
	}
	return &result, nil
}

func (a *StoreAPI) ListShops(ctx context.Context, params ShopListParams) (*ShopListResponse, error) {
	resp, err := a.client.Call(ctx, "v3/org/shop/getShopList", params)
	if err != nil {
		return nil, fmt.Errorf("list shops: %w", err)
	}
	var result ShopListResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("parse shop list: %w", err)
	}
	return &result, nil
}

func (a *StoreAPI) ListTakeoutMappings(ctx context.Context, params TakeoutShopMappingParams) (*TakeoutShopMappingResponse, error) {
	resp, err := a.client.Call(ctx, "v3/dist/meTakeoutShopMapPage", params)
	if err != nil {
		return nil, fmt.Errorf("list takeout mappings: %w", err)
	}
	var result TakeoutShopMappingResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("parse takeout mappings: %w", err)
	}
	return &result, nil
}

func (a *StoreAPI) ChangeShopStatus(ctx context.Context, code string, status int) error {
	_, err := a.client.Call(ctx, "v3/org/shop/changeStatus", map[string]interface{}{
		"multiMark": code,
		"status":    status,
	})
	if err != nil {
		return fmt.Errorf("change shop status: %w", err)
	}
	return nil
}

func (a *StoreAPI) SyncShopInfo(ctx context.Context, req ShopSyncRequest) error {
	_, err := a.client.Call(ctx, "v3/org/shop/shopInfoSync", req)
	if err != nil {
		return fmt.Errorf("sync shop info: %w", err)
	}
	return nil
}

func (a *StoreAPI) ShopCodeToID(ctx context.Context, code string) (*ShopCodeToIDResponse, error) {
	resp, err := a.client.Call(ctx, "v3/org/shop/shopCode2Id", map[string]string{"multiMark": code})
	if err != nil {
		return nil, fmt.Errorf("shop code to id: %w", err)
	}
	var result ShopCodeToIDResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("parse shop code to id: %w", err)
	}
	return &result, nil
}

func (a *StoreAPI) GetShopDetailByID(ctx context.Context, shopID int64) (*ShopDetail, error) {
	resp, err := a.client.Call(ctx, "v3/org/shop/getShopDetailById", map[string]int64{"shopId": shopID})
	if err != nil {
		return nil, fmt.Errorf("get shop detail by id: %w", err)
	}
	var result ShopDetail
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("parse shop detail by id: %w", err)
	}
	return &result, nil
}

func (a *StoreAPI) ModifyShopTeam(ctx context.Context, shopID, teamID int64) (*BoolResult, error) {
	resp, err := a.client.Call(ctx, "v3/org/shop/modifyShopTeam", map[string]int64{
		"shopId": shopID,
		"teamId": teamID,
	})
	if err != nil {
		return nil, fmt.Errorf("modify shop team: %w", err)
	}
	var result BoolResult
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("parse modify shop team: %w", err)
	}
	return &result, nil
}

func (a *StoreAPI) QueryStoreConfigBatch(ctx context.Context, storeID int64, fieldCodes []string) (*StoreConfigBatchResponse, error) {
	resp, err := a.client.Call(ctx, "v3/storeConfig/queryStoreConfigBatch", map[string]interface{}{
		"storeId":    storeID,
		"fieldCodes": fieldCodes,
	})
	if err != nil {
		return nil, fmt.Errorf("query store config batch: %w", err)
	}
	var result StoreConfigBatchResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("parse store config batch: %w", err)
	}
	return &result, nil
}

func (a *StoreAPI) QuerySellerConfigBatch(ctx context.Context, fieldCodes []string) (*SellerConfigBatchResponse, error) {
	resp, err := a.client.Call(ctx, "v3/sellerConfig/querySellerConfigBatch", map[string]interface{}{
		"fieldCodes": fieldCodes,
	})
	if err != nil {
		return nil, fmt.Errorf("query seller config batch: %w", err)
	}
	var result SellerConfigBatchResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("parse seller config batch: %w", err)
	}
	return &result, nil
}

func (a *StoreAPI) GetShopExtData(ctx context.Context, shopID int64) ([]ShopExtDataItem, error) {
	resp, err := a.client.Call(ctx, "v3/newPattern/orgCenter/post/shop/get-shop-ext-data", map[string]int64{"shopId": shopID})
	if err != nil {
		return nil, fmt.Errorf("get shop ext data: %w", err)
	}
	var result []ShopExtDataItem
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("parse shop ext data: %w", err)
	}
	return result, nil
}

func (a *StoreAPI) ListShopLabels(ctx context.Context, shopID int64) ([]ShopLabel, error) {
	resp, err := a.client.Call(ctx, "v3/newPattern/orgCenter/get/shop/label-by-id", map[string]int64{"shopId": shopID})
	if err != nil {
		return nil, fmt.Errorf("list shop labels: %w", err)
	}
	var result []ShopLabel
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("parse shop labels: %w", err)
	}
	return result, nil
}

func (a *StoreAPI) GetOrgTree(ctx context.Context, containCloseFlag int) ([]OrgTreeNode, error) {
	params := map[string]int{}
	if containCloseFlag > 0 {
		params["containCloseFlag"] = containCloseFlag
	}
	resp, err := a.client.Call(ctx, "v3/org/shop/getOrgTree", params)
	if err != nil {
		return nil, fmt.Errorf("get org tree: %w", err)
	}
	var result []OrgTreeNode
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("parse org tree: %w", err)
	}
	return result, nil
}

func (a *StoreAPI) ListShopTeams(ctx context.Context, name string) ([]ShopTeam, error) {
	params := map[string]string{}
	if name != "" {
		params["name"] = name
	}
	resp, err := a.client.Call(ctx, "v3/org/shop/shopTeamList", params)
	if err != nil {
		return nil, fmt.Errorf("list shop teams: %w", err)
	}
	var result []ShopTeam
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("parse shop teams: %w", err)
	}
	return result, nil
}

func (a *StoreAPI) GetShopDeptTree(ctx context.Context) (*ShopDeptNode, error) {
	resp, err := a.client.Call(ctx, "v3/org/shop/shopDeptTree", map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("get shop dept tree: %w", err)
	}
	var result ShopDeptNode
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("parse shop dept tree: %w", err)
	}
	return &result, nil
}
