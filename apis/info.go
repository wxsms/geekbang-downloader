package apis

// ColumnInfoApi {"product_id":100093501,"with_recommend_article":true}
var ColumnInfoApi = "https://time.geekbang.org/serv/v3/column/info"

type ColumnInfoResp struct {
	Data struct {
		Title string `json:"title"`
		Id    int    `json:"id"`
	} `json:"data"`
	Code int `json:"code"`
}
