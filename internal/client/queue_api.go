package client

import (
	"context"
	"fmt"
)

// QueueAPI provides queuing APIs.
type QueueAPI struct {
	client *Client
}

// NewQueueAPI creates a new QueueAPI.
func NewQueueAPI(c *Client) *QueueAPI {
	return &QueueAPI{client: c}
}

type ShopQueueProgressItem struct {
	CupTotal int    `json:"cupTotal"`
	MakeTime int    `json:"makeTime"`
	OrderNum int    `json:"orderNum"`
	ShopID   string `json:"shopId"`
}

type OrderQueueProgressResponse struct {
	CupTotal    int `json:"cupTotal"`
	MakeTime    int `json:"makeTime"`
	OrderNum    int `json:"orderNum"`
	OrderStatus int `json:"orderStatus"`
}

type QueueNoItem struct {
	OrderNo            string `json:"orderNo"`
	QueueNo            string `json:"queueNo"`
	QueueNoOrderSource int    `json:"queueNoOrderSource"`
	QueueNoOrderType   int    `json:"queueNoOrderType"`
	QueueNoStatus      int    `json:"queueNoStatus"`
	UserID             int64  `json:"userId"`
}

type ShopQueueNoListResponse struct {
	List  []QueueNoItem `json:"list"`
	Page  int           `json:"page"`
	Size  int           `json:"size"`
	Total int           `json:"total"`
}

func (a *QueueAPI) QueryShopQueueProgress(ctx context.Context, shopIDList []int64, shopType int) (*[]ShopQueueProgressItem, error) {
	resp, err := a.client.Call(ctx, "v3/queuing/queryShopQueueCup", map[string]interface{}{
		"shopIdList": shopIDList,
		"shopType":   shopType,
	})
	if err != nil {
		return nil, fmt.Errorf("query shop queue progress: %w", err)
	}
	return decodeData[[]ShopQueueProgressItem](resp, "shop queue progress")
}

func (a *QueueAPI) QueryOrderQueueProgress(ctx context.Context, orderNo, sourceNo string) (*OrderQueueProgressResponse, error) {
	params := map[string]interface{}{}
	if orderNo != "" {
		params["orderNo"] = orderNo
	}
	if sourceNo != "" {
		params["sourceNo"] = sourceNo
	}
	resp, err := a.client.Call(ctx, "v3/queuing/queryOrderQueueCup", params)
	if err != nil {
		return nil, fmt.Errorf("query order queue progress: %w", err)
	}
	return decodeData[OrderQueueProgressResponse](resp, "order queue progress")
}

func (a *QueueAPI) QueryShopQueueNoList(ctx context.Context, shopCode string, page, size int, statusList []int) (*ShopQueueNoListResponse, error) {
	params := map[string]interface{}{
		"shopCode": shopCode,
		"page":     page,
		"size":     size,
	}
	if len(statusList) > 0 {
		params["queueNoStatusList"] = statusList
	}
	resp, err := a.client.Call(ctx, "v3/queuing/queueNo/queryShopQueueNoList", params)
	if err != nil {
		return nil, fmt.Errorf("query shop queue no list: %w", err)
	}
	return decodeData[ShopQueueNoListResponse](resp, "shop queue no list")
}
