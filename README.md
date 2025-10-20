# BlogApi

> GitHub Action 配置请参考 config 目录下的 sync-feishu.yml

---

> 你可以参考 https://ycnne3fqukbq.feishu.cn/base/XGk3bupgbaCEAPsRxAqc47ytnyg?table=tbllVnipZfa9cHoX&view=vewHYeOSGR 的多维表格模板 

一个用于个人博客数据同步的 Go 应用程序，支持从飞书多维表格和 Steam 平台自动同步数据，并生成静态资源文件。

> TODO: 为 bilibili 豆瓣等平台实现抓取功能, 配置对应的 webhook, 可以通过配置手动调整时更新

## 功能特性

### 核心功能

- **飞书数据同步**: 从飞书多维表格同步书籍、电影、动漫记录
- **Steam 游戏同步**: 获取 Steam 游戏库和最近游戏记录
- **图片资源管理**: 自动下载并本地化存储图片资源
- **数据格式化**: 生成 NDJSON 格式的结构化数据文件

### 支持的数据类型

1. **书籍记录** (`bookcase/`)
   - 从飞书表格同步书籍信息
   - 自动下载封面图片
   - 生成 `record.ndjson` 数据文件

2. **电影记录** (`movie/`)
   - 从飞书表格同步电影信息
   - 自动下载海报图片
   - 生成 `record.ndjson` 数据文件

3. **动漫记录** (`anime/`)
   - 从飞书表格同步动漫信息
   - 自动下载封面图片
   - 生成 `record.ndjson` 数据文件

4. **Steam 游戏** (`steam/`)
   - 同步 Steam 游戏库
   - 获取最近游戏记录（最多10个）
   - 自动补充随机游戏至10个
   - 下载游戏封面图片

## 项目结构

```
BlogApi/
├── cmd/                    # 主程序入口
│   └── main.go
├── config/                 # 配置文件
│   ├── config.go          # 配置结构定义
│   ├── config.yaml        # 主配置文件
│   └── config.yaml.example # 配置示例文件
├── internal/              # 内部模块
│   ├── feishu/           # 飞书相关功能
│   │   ├── client.go     # 飞书客户端
│   │   ├── auth.go       # 认证逻辑
│   │   ├── download.go   # 图片下载
│   │   ├── models.go     # 数据模型
│   │   ├── parse.go      # 数据解析
│   │   ├── save.go       # 数据保存
│   │   ├── animeRecords.go
│   │   ├── bookRecords.go
│   │   └── movieRecords.go
│   └── steam/            # Steam 相关功能
│       ├── client.go     # Steam 客户端
│       ├── download.go   # 图片下载
│       ├── models.go     # 数据模型
│       ├── save.go       # 数据保存
│       └── steam_api.go  # Steam API 调用
├── public/               # 生成的静态资源
│   ├── anime/           # 动漫相关文件
│   ├── bookcase/        # 书籍相关文件
│   ├── movie/           # 电影相关文件
│   └── steam/           # Steam 游戏文件
├── utils/               # 工具函数
│   └── update.go       # 更新逻辑
├── go.mod              # Go 模块定义
├── go.sum              # 依赖校验
└── README.md           # 项目说明
```

## 快速开始

### 环境要求

- Go 1.25.0 或更高版本
- 有效的飞书应用凭据
- 有效的 Steam Web API Key

### 安装依赖

```bash
go mod download
```

### 配置设置

1. 复制配置示例文件：
```bash
cp config/config.yaml.example config/config.yaml
```

2. 编辑 `config/config.yaml`，填入实际的配置信息：

```yaml
# 飞书配置
feishu:
  app_id: "your_feishu_app_id"
  app_secret: "your_feishu_app_secret"
  app_token: "your_feishu_app_token"
  book_table_id: "your_book_table_id"
  movie_table_id: "your_movie_table_id"
  anime_table_id: "your_anime_table_id"

# Steam 配置
steam:
  steam_key: "your_steam_web_api_key"
  steam_id: "your_steam_64bit_id"
```

### 运行程序

```bash
# 运行数据同步
go run ./cmd

# 或者构建后运行
go build -o blogapi ./cmd
./blogapi
```

## 配置说明

### 飞书配置获取

1. **创建飞书应用**：
   - 访问 [飞书开放平台](https://open.feishu.cn/)
   - 创建企业自建应用
   - 获取 `app_id` 和 `app_secret`

2. **获取表格信息**：
   - 打开目标多维表格
   - 从 URL 中提取 `app_token` 和各个 `table_id`

3. **权限设置**：
   - 确保应用有读取表格的权限
   - 添加必要的 API 权限范围

### Steam 配置获取

1. **Steam Web API Key**：
   - 访问 [Steam Web API](https://steamcommunity.com/dev/apikey)
   - 使用 Steam 账号登录并申请 API Key

2. **Steam ID**：
   - 获取你的 64 位 Steam ID
   - 可以使用 [SteamID.io](https://steamid.io/) 等工具转换

## 环境变量支持

程序支持通过环境变量配置，环境变量优先级高于配置文件：

```bash
# 飞书配置
export BLOGAPI_FEISHU_APP_ID="your_app_id"
export BLOGAPI_FEISHU_APP_SECRET="your_app_secret"
export BLOGAPI_FEISHU_APP_TOKEN="your_app_token"
export BLOGAPI_FEISHU_BOOK_TABLE_ID="your_book_table_id"
export BLOGAPI_FEISHU_MOVIE_TABLE_ID="your_movie_table_id"
export BLOGAPI_FEISHU_ANIME_TABLE_ID="your_anime_table_id"

# Steam 配置
export BLOGAPI_STEAM_KEY="your_steam_key"
export BLOGAPI_STEAM_ID="your_steam_id"
```

## 输出文件格式

### NDJSON 数据文件

每个类别都会生成一个 `record.ndjson` 文件，包含记录 ID：

```json
{"record_id":"rec123456"}
{"record_id":"rec789012"}
```

### 图片文件

- **格式**: WebP
- **命名**: 使用记录 ID 或游戏 App ID
- **存储**: 按类别分目录存储

## 依赖项

- `github.com/larksuite/oapi-sdk-go/v3` - 飞书 Open API SDK
- `github.com/spf13/viper` - 配置管理
- `github.com/chai2010/webp` - WebP 图片处理

## 开发说明

### 添加新的数据源

1. 在 `internal/` 目录下创建新的包
2. 实现相应的客户端和数据处理逻辑
3. 在 `utils/update.go` 中添加更新调用
4. 更新配置结构和文档

### 自定义数据处理

可以通过修改各个模块的 `save.go` 文件来自定义数据输出格式和处理逻辑。

## 注意事项

- 确保网络连接稳定，程序需要访问外部 API
- 首次运行可能需要较长时间下载图片资源
- 建议定期备份 `public/` 目录下的数据文件
- Steam API 有访问频率限制，避免频繁调用

## 贡献

欢迎提交 Issue 和 Pull Request 来改进这个项目。

---

**注意**: 请妥善保管你的 API 密钥和配置信息，不要将其提交到公共代码仓库中。
