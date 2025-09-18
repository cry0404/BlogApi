package feishu

// 从理论上来讲这里只需要解析 bookcase 的有关信息即可，但 ai 写数据结构比较方便



type FeishuResponse[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}


type TableListResponse struct {
	HasMore bool          `json:"has_more"`
	Items   []BookRecord  `json:"items"`
	Total   int           `json:"total"`
}


type FieldValue struct {
	Text string `json:"text"`
	Type string `json:"type"`
}


type FileValue struct {
	FileToken string `json:"file_token"`
	Name      string `json:"name"`
	Size      int64  `json:"size"`
	TmpURL    string `json:"tmp_url"`
	Type      string `json:"type"`
	URL       string `json:"url"`
}


type BookFields struct {
	Title           []FieldValue `json:"书名,omitempty"`
	Author          []FieldValue `json:"作者,omitempty"`
	CoverImage      []FileValue  `json:"书籍封面,omitempty"`
	ReadDate        int64        `json:"完成阅读时期,omitempty"`
	RecommendStatus string       `json:"推荐状态,omitempty"`
	Tag             string       `json:"标签,omitempty"`
	Description     []FieldValue `json:"简介,omitempty"`
	Comment         []FieldValue `json:"书评,omitempty"`
}


type BookRecord struct {
	Fields  BookFields `json:"fields"`
	RecordID string    `json:"record_id"`
}


type ProcessedBook struct {
	RecordID        string     `json:"record_id"`
	Title           string     `json:"title"`
	Author          string     `json:"author"`
	CoverImage      *FileValue `json:"cover_image,omitempty"`
	ReadDate        int64      `json:"read_date"`
	RecommendStatus string     `json:"recommend_status"`
	Tag             string     `json:"tag"`
	Description     string     `json:"description"`
	Comment         string     `json:"comment"`
} //将处理后的信息写入 bookcase.json 文件