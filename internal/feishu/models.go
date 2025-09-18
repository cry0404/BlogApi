package feishu

type BookRecord struct {
	//理论上都该是 string
	Title           string 		`json:"书名,omitempty"`
	Author          string 		`json:"作者,omitempty"`
	ConverImage		string		`json:"书籍封面,omitempty"`
	ReadDate        string      `json:"完成阅读时期,omitempty"`
	RecommendStatus string      `json:"推荐状态,omitempty"`
	Tag             string      `json:"标签,omitempty"`
	Description     string 		`json:"简介,omitempty"`
	Comment         string 		`json:"书评,omitempty"`
	RecordID		string      `json:"record_id"`
	ImageDir		string
}

type Account struct {
	APPID  		string 	`json:"app_id"`
	APPSECRET	string	`json:"app_secret"`
}

type TokenResponse struct {
	Code 	int  	`json:"code"`
	Msg		string	`json:"msg"`
	Token	string  `json:"tenant_access_token"`
	Expire  int
}

