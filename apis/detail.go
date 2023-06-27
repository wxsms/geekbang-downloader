package apis

// ArticleDetailApi {"id":"426265","include_neighbors":true,"is_freelyread":true}
var ArticleDetailApi = "https://time.geekbang.org/serv/v1/article"

type ArticleDetail struct {
	ArticleContent string `json:"article_content"`
	ArticleTitle   string `json:"article_title"`
	Id             int    `json:"id"`
}

type ArticleDetailResp struct {
	Data ArticleDetail `json:"data"`
	Code int           `json:"code"`
}
