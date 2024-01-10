package apis

var ProductListApi = "https://time.geekbang.org/serv/v4/pvip/product_list"

type ProductListResp struct {
	Data struct {
		Page struct {
			More  bool
			Total int
		}
		Products []struct {
			Id    int
			Title string
		}
	} `json:"data"`
	Code int `json:"code"`
}
