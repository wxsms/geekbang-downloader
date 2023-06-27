package apis

// ArticleListApi `{"cid":"100093501","size":500,"prev":0,"order":"earliest","sample":false}`
var ArticleListApi = "https://time.geekbang.org/serv/v1/column/articles"

type ArticleListItem struct {
	Title string `json:"article_title"`
	Id    int    `json:"id"`
}

type ArticleListResp struct {
	Data struct {
		List []ArticleListItem `json:"list"`
	} `json:"data"`
	Code int `json:"code"`
}
