package feishu

type Record interface {
	GetRecordID() string
	GetCoverImage() string
	GetHasDownload() bool
	SetHasDownload(bool)
	SetImageDir(string)
	GetImageDir() string
	ToOutputInfo() interface{}
}

// RecordType 表示记录类型
type RecordType string

const (
	BookRecordType  RecordType = "book"
	AnimeRecordType RecordType = "anime"
	MovieRecordType RecordType = "movie"
)

type MovieRecord struct {
	Title       string `json:"影名,omitempty"`
	Description string `json:"简介,omitempty"`
	Author      string `json:"作者,omitempty"`
	Comment     string `json:"影评,omitempty"`
	Grade       int    `json:"评价,omitempty"`
	Date        string `json:"日期,omitempty"`
	ConverImage string `json:"封面,omitempty"`
	RecordID    string `json:"record_id"`
	ImageDir    string
	HasDownload bool
}

func (m *MovieRecord) GetRecordID() string            { return m.RecordID }
func (m *MovieRecord) GetCoverImage() string          { return m.ConverImage }
func (m *MovieRecord) GetHasDownload() bool           { return m.HasDownload }
func (m *MovieRecord) SetHasDownload(downloaded bool) { m.HasDownload = downloaded }
func (m *MovieRecord) SetImageDir(dir string)         { m.ImageDir = dir }
func (m *MovieRecord) GetImageDir() string            { return m.ImageDir }
func (m *MovieRecord) ToOutputInfo() interface{} {
	return &MovieInfo{
		Title:       m.Title,
		Description: m.Description,
		Author:      m.Author,
		Comment:     m.Comment,
		Grade:       m.Grade,
		Date:        m.Date,
		ImageDir:    m.ImageDir,
	}
}

type AnimeRecord struct {
	Title       string `json:"番名,omitempty"`
	Description string `json:"简介,omitempty"`
	Grade       int    `json:"评分,omitempty"`
	Comment     string `json:"评价,omitempty"`
	ConverImage string `json:"封面,omitempty"`
	RecordID    string `json:"record_id"`
	ImageDir    string
	HasDownload bool
}

// 实现 Record 接口
func (a *AnimeRecord) GetRecordID() string            { return a.RecordID }
func (a *AnimeRecord) GetCoverImage() string          { return a.ConverImage }
func (a *AnimeRecord) GetHasDownload() bool           { return a.HasDownload }
func (a *AnimeRecord) SetHasDownload(downloaded bool) { a.HasDownload = downloaded }
func (a *AnimeRecord) SetImageDir(dir string)         { a.ImageDir = dir }
func (a *AnimeRecord) GetImageDir() string            { return a.ImageDir }
func (a *AnimeRecord) ToOutputInfo() interface{} {
	return &AnimeInfo{
		Title:       a.Title,
		Description: a.Description,
		Grade:       a.Grade,
		Comment:     a.Comment,
		ImageDir:    a.ImageDir,
	}
}

type BookRecord struct {
	Title           string `json:"书名,omitempty"`
	Author          string `json:"作者,omitempty"`
	ConverImage     string `json:"书籍封面,omitempty"`
	ReadDate        string `json:"完成阅读时期,omitempty"`
	RecommendStatus string `json:"推荐状态,omitempty"`
	Tag             string `json:"标签,omitempty"`
	Description     string `json:"简介,omitempty"`
	Comment         string `json:"书评,omitempty"`
	RecordID        string `json:"record_id"`
	ImageDir        string
	HasDownload     bool
}

// 实现 Record 接口
func (b *BookRecord) GetRecordID() string            { return b.RecordID }
func (b *BookRecord) GetCoverImage() string          { return b.ConverImage }
func (b *BookRecord) GetHasDownload() bool           { return b.HasDownload }
func (b *BookRecord) SetHasDownload(downloaded bool) { b.HasDownload = downloaded }
func (b *BookRecord) SetImageDir(dir string)         { b.ImageDir = dir }
func (b *BookRecord) GetImageDir() string            { return b.ImageDir }
func (b *BookRecord) ToOutputInfo() interface{} {
	return &BookInfo{
		Title:           b.Title,
		Author:          b.Author,
		ReadDate:        b.ReadDate,
		RecommendStatus: b.RecommendStatus,
		Tag:             b.Tag,
		Description:     b.Description,
		Comment:         b.Comment,
		ImageDir:        b.ImageDir,
	}
}

type Account struct {
	APPID     string `json:"app_id"`
	APPSECRET string `json:"app_secret"`
}

type TokenResponse struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Token  string `json:"tenant_access_token"`
	Expire int
}

// 输出信息结构体
type BookInfo struct {
	Title           string
	Author          string
	ReadDate        string
	RecommendStatus string
	Tag             string
	Description     string
	Comment         string
	ImageDir        string
}

type AnimeInfo struct {
	Title       string
	Description string
	Grade       int
	Comment     string
	ImageDir    string
}

type MovieInfo struct {
	Title       string
	Description string
	Author      string
	Comment     string
	Grade       int
	Date        string
	ImageDir    string
}
