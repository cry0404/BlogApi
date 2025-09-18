package utils

import (
	config "BlogApi/config"
	feishu "BlogApi/internal/feishu"
	"fmt"
)

//这里应该依次更新所有客户端的配置
func Update(cfg *config.Config) error {

	err := feishu.UpdateBookCase(cfg)

	if err != nil {
		return fmt.Errorf("更新书架失败，检查错误: %v", err)
	}

	// 别的更新逻辑

// 下载完成后发起 webhook
	return nil
}