package client

import (
	"context"
	"encoding/json"
	"fmt"
)

// ProductAPI provides product-related API methods for the open platform.
type ProductAPI struct {
	client *Client
}

// NewProductAPI creates a new ProductAPI.
func NewProductAPI(c *Client) *ProductAPI {
	return &ProductAPI{client: c}
}

// --- Request/Response types aligned to open platform ---

// ProductListParams maps to /v3/goods/item/getItemList params.
type ProductListParams struct {
	ShopCode    string `json:"shopCode"`
	Name        string `json:"name,omitempty"`
	SaleChannel int    `json:"saleChannel,omitempty"` // 1=堂食, 2=外卖, etc.
	SaleType    int    `json:"saleType,omitempty"`    // 1=普通, 2=套餐
	PageNo      int    `json:"pageNo,omitempty"`
	PageSize    int    `json:"pageSize,omitempty"`
}

// GoodsSku represents a SKU within a product.
type GoodsSku struct {
	SkuId     string  `json:"skuId"`
	TradeMark string  `json:"tradeMark"`
	SalePrice float64 `json:"salePrice"`
	Inventory float64 `json:"inventory"` // API 返回可能是 int 或 float
	Barcode   string  `json:"barcode,omitempty"`
}

// OpenProduct represents a product from the open platform list response.
type OpenProduct struct {
	ID               int64      `json:"id"`
	GoodsId          int64      `json:"goodsId"` //nolint:revive
	Name             string     `json:"name"`
	Status           int        `json:"status"` // 10=上架, 20=下架
	SaleChannel      int        `json:"saleChannel"`
	SaleType         int        `json:"saleType"`
	ShowPriceLow     int        `json:"showPriceLow"`  // 价格（分）
	ShowPriceHigh    int        `json:"showPriceHigh"` // 价格（分）
	CategoryNameList []string   `json:"categoryNameList"`
	GoodsSkuList     []GoodsSku `json:"goodsSkuList"`
	Unit             string     `json:"unit,omitempty"`
}

// ProductListResponse is the response from getItemList.
type ProductListResponse struct {
	Data  []OpenProduct `json:"data"`
	Total int           `json:"total"`
}

// List queries the store product list.
// POST /v3/goods/item/getItemList
func (a *ProductAPI) List(ctx context.Context, shopCode, name string, saleChannel, saleType, pageNo, pageSize int) (*ProductListResponse, error) {
	params := ProductListParams{
		ShopCode:    shopCode,
		Name:        name,
		SaleChannel: saleChannel,
		SaleType:    saleType,
		PageNo:      pageNo,
		PageSize:    pageSize,
	}

	resp, err := a.client.Call(ctx, "v3/goods/item/getItemList", params)
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}

	var result ProductListResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("parse product list: %w", err)
	}
	return &result, nil
}

// --- Sync (batch create/update) ---

// SyncCategory is the category info for the sync API.
type SyncCategory struct {
	CategoryName string `json:"categoryName"`
}

// SyncSku is the SKU info for the sync API.
type SyncSku struct {
	SalePrice float64 `json:"salePrice"`
	Inventory int     `json:"inventory,omitempty"`
}

// ShopGoods is the user-facing input for creating/updating products.
// Sync() will merge these into the full API payload with all required fields.
type ShopGoods struct {
	TradeName  string  // 商品名称（必填）
	TradePrice float64 // 价格（必填）
	TradeNo    string  // 商户自定义编号
	ClassName  string  // 分类名称
	Stock      int     // 库存
}

// Sync batch creates/updates products.
// POST /v3/goods/sync/shopGoodsSync
func (a *ProductAPI) Sync(ctx context.Context, storeCode string, saleType int, goodsList []ShopGoods) error {
	apiGoods := make([]map[string]interface{}, len(goodsList))
	for i, g := range goodsList {
		name := g.TradeName
		price := g.TradePrice
		stock := g.Stock
		if stock == 0 {
			stock = 9999
		}

		catList := []map[string]interface{}{}
		if g.ClassName != "" {
			catList = []map[string]interface{}{
				{
					"name":         g.ClassName,
					"categoryName": g.ClassName,
					"mark":         g.ClassName,
					"isRequired":   0,
					"isBackend":    0,
					"isFront":      0,
					"sort":         0,
					"type":         0,
				},
			}
		}

		apiGoods[i] = map[string]interface{}{
			// 基本信息
			"name":       name,
			"tradeName":  name,
			"tradeNo":    g.TradeNo,
			"tradePrice": price,
			// 分类
			"categoryList": catList,
			// SKU
			"skuList": []map[string]interface{}{
				{
					"tradeMark":        g.TradeNo,
					"salePrice":        price,
					"marketPrice":      price,
					"costPrice":        0,
					"inventory":        stock,
					"stock":            stock,
					"weight":           "0.000",
					"barcode":          "",
					"specCode":         "",
					"packingFee":       0,
					"clearStatus":      1,
					"specValueList":    []string{},
					"pictureUrlList":   []string{},
					"skuImageList":     []string{},
					"skuItemList":      []string{},
					"categoryList":     []string{},
					"categoryNameList": []string{},
				},
			},
			// 必填默认值
			"packingFee":    0,
			"minBuyNum":     1,
			"saleChannel":   1,
			"isAttach":      0,
			"isMemberPrice": 0,
			"isMultiSpec":   0,
			"isWeigh":       0,
			"isPractice":    1,
			"isShow":        0,
			"type":          1,
			"subType":       11,
			"status":        10,
			"sort":          0,
			"unit":          "",
			"barcode":       "",
			"moreBarcode":   "",
			"coverUrl":      "",
			"videoCoverUrl": "",
			"videoUrl":      "",
			"baseSaleNum":   0,
			"saleTime": map[string]interface{}{
				"dateStart":    "",
				"dateEnd":      "",
				"saleTimeList": []string{},
				"weekdayList":  []int{},
			},
			"pictureUrlList": []string{"https://images.qmai.cn/s213951/2025/11/05/3bba455d9c8a42294f.jpg"},
		}
	}

	params := map[string]interface{}{
		"storeCode": storeCode,
		"saleType":  saleType,
		"goodsList": apiGoods,
	}

	_, err := a.client.Call(ctx, "v3/goods/sync/shopGoodsSync", params)
	if err != nil {
		return fmt.Errorf("sync products: %w", err)
	}
	return nil
}

// --- Batch Up/Down ---

// BatchUpParams maps to /v3/goods/item/externalUp params.
type BatchUpParams struct {
	ShopCode      string   `json:"shopCode"`
	TradeMarkList []string `json:"tradeMarkList"`
	SaleChannel   int      `json:"saleChannel,omitempty"`
}

// BatchUp sets products as on-sale (上架).
// POST /v3/goods/item/externalUp
func (a *ProductAPI) BatchUp(ctx context.Context, shopCode string, tradeMarkList []string, saleChannel int) error {
	params := BatchUpParams{
		ShopCode:      shopCode,
		TradeMarkList: tradeMarkList,
		SaleChannel:   saleChannel,
	}

	_, err := a.client.Call(ctx, "v3/goods/item/externalUp", params)
	if err != nil {
		return fmt.Errorf("batch up: %w", err)
	}
	return nil
}

// BatchDownParams maps to /v3/goods/item/externalDown params.
type BatchDownParams struct {
	ShopCode      string   `json:"shopCode"`
	TradeMarkList []string `json:"tradeMarkList"`
	SaleChannel   int      `json:"saleChannel,omitempty"`
}

// BatchDown sets products as off-sale (下架).
// POST /v3/goods/item/externalDown
func (a *ProductAPI) BatchDown(ctx context.Context, shopCode string, tradeMarkList []string, saleChannel int) error {
	params := BatchDownParams{
		ShopCode:      shopCode,
		TradeMarkList: tradeMarkList,
		SaleChannel:   saleChannel,
	}

	_, err := a.client.Call(ctx, "v3/goods/item/externalDown", params)
	if err != nil {
		return fmt.Errorf("batch down: %w", err)
	}
	return nil
}

// --- Sell Out (估清) ---

// SellOutParams maps to /v3/goods/item/sellOut params.
type SellOutParams struct {
	StoreCode string `json:"storeCode"`
	TradeMark string `json:"tradeMark"`
	IsSellOut int    `json:"isSellOut"` // 1=估清, 0=取消估清
}

// SellOutResponse is the response from sellOut.
type SellOutResponse struct {
	TradeMark string `json:"tradeMark"`
	IsSellOut int    `json:"isSellOut"`
}

// SellOut sets a product as sold out or restores it.
// POST /v3/goods/item/sellOut
func (a *ProductAPI) SellOut(ctx context.Context, storeCode, tradeMark string, isSellOut int) (*SellOutResponse, error) {
	params := SellOutParams{
		StoreCode: storeCode,
		TradeMark: tradeMark,
		IsSellOut: isSellOut,
	}

	resp, err := a.client.Call(ctx, "v3/goods/item/sellOut", params)
	if err != nil {
		return nil, fmt.Errorf("sell out: %w", err)
	}

	var result SellOutResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("parse sell out response: %w", err)
	}
	return &result, nil
}

// FillUpParams maps to /v3/goods/item/fillUp params.
type FillUpParams struct {
	IsAllEmpty    int      `json:"isAllEmpty"`
	SaleChannel   int      `json:"saleChannel,omitempty"`
	SaleType      int      `json:"saleType,omitempty"`
	SpecCodeList  []string `json:"specCodeList,omitempty"`
	StoreCode     string   `json:"storeCode"`
	TradeMark     string   `json:"tradeMark,omitempty"`
	TradeMarkList []string `json:"tradeMarkList,omitempty"`
}

// BatchResultItem contains success/failure item info for some batch APIs.
type BatchResultItem struct {
	Message     string `json:"message"`
	Name        string `json:"name"`
	SaleChannel int    `json:"saleChannel"`
	SaleType    int    `json:"saleType"`
	Status      bool   `json:"status"`
}

// FillUpResponse is the response from fillUp.
type FillUpResponse struct {
	FailList    []BatchResultItem `json:"failList"`
	SuccessList []BatchResultItem `json:"successList"`
}

// FillUp cancels sell-out/estimated clear.
func (a *ProductAPI) FillUp(ctx context.Context, params FillUpParams) (*FillUpResponse, error) {
	resp, err := a.client.Call(ctx, "v3/goods/item/fillUp", params)
	if err != nil {
		return nil, fmt.Errorf("fill up: %w", err)
	}

	var result FillUpResponse
	if len(resp.Data) == 0 || string(resp.Data) == "null" {
		return &result, nil
	}
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("parse fill up response: %w", err)
	}
	return &result, nil
}

type ExternalStockParams struct {
	IsAll           bool     `json:"isAll,omitempty"`
	SaleChannelList []int    `json:"saleChannelList,omitempty"`
	SaleTypeList    []int    `json:"saleTypeList,omitempty"`
	MultiMark       string   `json:"multiMark"`
	TradeMarkList   []string `json:"tradeMarkList"`
	SpecCodeList    []string `json:"specCodeList,omitempty"`
}

// ExternalEmpty marks products as sold out.
func (a *ProductAPI) ExternalEmpty(ctx context.Context, params ExternalStockParams) error {
	_, err := a.client.Call(ctx, "v3/goods/item/externalEmpty", params)
	if err != nil {
		return fmt.Errorf("external empty: %w", err)
	}
	return nil
}

// ExternalFull marks products as fully stocked.
func (a *ProductAPI) ExternalFull(ctx context.Context, params ExternalStockParams) error {
	_, err := a.client.Call(ctx, "v3/goods/item/externalFull", params)
	if err != nil {
		return fmt.Errorf("external full: %w", err)
	}
	return nil
}

type GoodsAttachListParams struct {
	SaleChannel int    `json:"saleChannel,omitempty"`
	SaleType    int    `json:"saleType,omitempty"`
	ShopCode    string `json:"shopCode"`
	StockStatus int    `json:"stockStatus,omitempty"`
}

// GoodsAttach represents a simplified attach item.
type GoodsAttach struct {
	ID          string `json:"id"`
	AttachCode  string `json:"attachCode"`
	Name        string `json:"name"`
	Inventory   string `json:"inventory"`
	ShowPrice   string `json:"showPrice"`
	Status      int    `json:"status"`
	SaleChannel int    `json:"saleChannel"`
	SaleType    int    `json:"saleType"`
}

// ListAttachGoods queries attach goods list.
func (a *ProductAPI) ListAttachGoods(ctx context.Context, params GoodsAttachListParams) ([]GoodsAttach, error) {
	resp, err := a.client.Call(ctx, "v3/goods/item/shopGoodsAttachList", params)
	if err != nil {
		return nil, fmt.Errorf("list attach goods: %w", err)
	}
	var result []GoodsAttach
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("parse attach goods: %w", err)
	}
	return result, nil
}

type ShopGoodsListParams struct {
	FilterDownAttach bool     `json:"filterDownAttach,omitempty"`
	GoodsIDList      []int64  `json:"goodsIdList,omitempty"`
	IncludeProps     []string `json:"includeProperties,omitempty"`
	Name             string   `json:"name,omitempty"`
	PageNo           int      `json:"pageNo,omitempty"`
	PageSize         int      `json:"pageSize,omitempty"`
	SaleChannel      int      `json:"saleChannel"`
	SaleMethod       int      `json:"saleMethod,omitempty"`
	SaleType         int      `json:"saleType"`
	ShopCode         string   `json:"shopCode"`
	Status           string   `json:"status"`
}

type ShopGoodsWithPractice struct {
	ID      int64  `json:"id"`
	GoodsID int64  `json:"goodsId"`
	Name    string `json:"name"`
	Status  int    `json:"status"`
}

type ShopGoodsListResponse struct {
	Data  []ShopGoodsWithPractice `json:"data"`
	Total int                     `json:"total"`
}

// ListShopGoodsWithPractice queries product list with optional practice fields.
func (a *ProductAPI) ListShopGoodsWithPractice(ctx context.Context, params ShopGoodsListParams) (*ShopGoodsListResponse, error) {
	resp, err := a.client.Call(ctx, "v3/goods/item/getShopGoodsList", params)
	if err != nil {
		return nil, fmt.Errorf("list shop goods with practice: %w", err)
	}
	var result ShopGoodsListResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("parse shop goods with practice: %w", err)
	}
	return &result, nil
}

type GoodsRealtimeParams struct {
	StoreID     int64   `json:"storeId"`
	GoodsIDList []int64 `json:"goodsIdList"`
	SaleType    int     `json:"saleType"`
	SaleChannel int     `json:"saleChannel"`
}

type GoodsRealtimeSKU struct {
	Inventory     float64 `json:"inventory"`
	InventoryType int     `json:"inventoryType"`
	ItemSkuID     int64   `json:"itemSkuId"`
	PackingFee    int     `json:"packingFee"`
	SalePrice     int     `json:"salePrice"`
	SkuID         int64   `json:"skuId"`
}

type GoodsRealtimeItem struct {
	GoodsID      int64              `json:"goodsId"`
	GoodsSKUList []GoodsRealtimeSKU `json:"goodsSkuList"`
	ID           int64              `json:"id"`
	PackingFee   int                `json:"packingFee"`
	SaleChannel  int                `json:"saleChannel"`
	SaleType     int                `json:"saleType"`
	Status       int                `json:"status"`
}

// ListRealtimeGoods queries realtime goods data.
func (a *ProductAPI) ListRealtimeGoods(ctx context.Context, params GoodsRealtimeParams) ([]GoodsRealtimeItem, error) {
	resp, err := a.client.Call(ctx, "v3/newPattern/goodsCenter/post/v2/item/real-time/list", params)
	if err != nil {
		return nil, fmt.Errorf("list realtime goods: %w", err)
	}
	var result []GoodsRealtimeItem
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("parse realtime goods: %w", err)
	}
	return result, nil
}

type GoodsEnergyItem struct {
	EnergyValue           float64 `json:"energyValue"`
	GoodsID               int64   `json:"goodsId"`
	GradingIdentification string  `json:"gradingIdentification"`
	GradingInstructions   string  `json:"gradingInstructions"`
	PracticeValueIDs      []int64 `json:"practiceValueIds"`
	SkuID                 int64   `json:"skuId"`
	Unit                  string  `json:"unit"`
}

// GetGoodsEnergy queries product energy values.
func (a *ProductAPI) GetGoodsEnergy(ctx context.Context, goodsID, storeID int64) ([]GoodsEnergyItem, error) {
	resp, err := a.client.Call(ctx, "v3/goods/nutritional/energy", map[string]interface{}{
		"goodsId": goodsID,
		"storeId": storeID,
	})
	if err != nil {
		return nil, fmt.Errorf("get goods energy: %w", err)
	}
	var result []GoodsEnergyItem
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("parse goods energy: %w", err)
	}
	return result, nil
}

type PracticeStatusParams struct {
	PracticeValues []string `json:"practiceValues"`
	ShopCode       string   `json:"shopCode"`
	Status         int      `json:"status"`
}

// PracticeOnOff enables or disables practice values.
func (a *ProductAPI) PracticeOnOff(ctx context.Context, params PracticeStatusParams) error {
	_, err := a.client.Call(ctx, "v3/goods/item/practiceOnOff", params)
	if err != nil {
		return fmt.Errorf("practice on off: %w", err)
	}
	return nil
}

type DeleteTaskParams struct {
	SaleChannel   int      `json:"saleChannel"`
	SaleType      int      `json:"saleType"`
	SpecCodeList  []string `json:"specCodeList,omitempty"`
	TradeMarkList []string `json:"tradeMarkList,omitempty"`
	StoreCode     string   `json:"storeCode"`
}

type DeleteTaskResponse struct {
	SaleType  int    `json:"saleType"`
	StoreCode string `json:"storeCode"`
	TaskID    string `json:"taskId"`
}

// SubmitDeleteTask submits a delete task for products.
func (a *ProductAPI) SubmitDeleteTask(ctx context.Context, params DeleteTaskParams) (*DeleteTaskResponse, error) {
	resp, err := a.client.Call(ctx, "v3/goods/sync/tripartiteShopGoodsDel", params)
	if err != nil {
		return nil, fmt.Errorf("submit delete task: %w", err)
	}
	var result DeleteTaskResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("parse delete task response: %w", err)
	}
	return &result, nil
}
