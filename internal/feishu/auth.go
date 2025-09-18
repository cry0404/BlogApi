package feishu

import (
	"BlogApi/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)




func getAccessToken(cfg *config.Config) (string, error) {
	body := Account{
		APPID: cfg.FeiShu.FeiShuAppID,
		APPSECRET: cfg.FeiShu.FeiShuAppSecret,
	}
	bs, _ := json.Marshal(body)
	//调试
	//fmt.Println(string(bs))
	req, err := http.NewRequest("POST", "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal", bytes.NewReader(bs))
	if err != nil {
		return "", fmt.Errorf("设置请求时失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)

	if err !=nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("请求失败，未能获取到 token: %v", err)
	}

	data, _ := io.ReadAll(resp.Body)
	var tokenResponse TokenResponse
	err = json.Unmarshal(data, &tokenResponse)
	if err != nil {
		return "", fmt.Errorf("解析时发生错误: %v", err)
	}


	return tokenResponse.Token, nil 
}