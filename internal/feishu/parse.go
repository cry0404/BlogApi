package feishu

import (
	"strconv"
	"strings"
	"time"
)

func parseTextField(fields map[string]any, key string) string {
	v, ok := fields[key]
	if !ok || v == nil {
		return ""
	}
	switch arr := v.(type) {
	case []any:
		var b strings.Builder
		for _, it := range arr {
			if m, ok := it.(map[string]any); ok {
				if s, ok := m["text"].(string); ok {
					b.WriteString(s)
				}
			}
		}
		return b.String()
	case string:
		return arr
	}
	return ""
}

func parseStringField(fields map[string]any, key string) string {
	if v, ok := fields[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func parseFirstFileURL(fields map[string]any, key string) string {
	v, ok := fields[key]
	if !ok || v == nil {
		return ""
	}
	arr, ok := v.([]any)
	if !ok || len(arr) == 0 {
		return ""
	}
	m, ok := arr[0].(map[string]any)
	if !ok {
		return ""
	}
	if s, ok := m["url"].(string); ok && s != "" {
		return s
	}
	if s, ok := m["tmp_url"].(string); ok && s != "" {
		return s
	}
	return ""
}

func parseUnixFieldRFC3339(fields map[string]any, key string) string {
	v, ok := fields[key]
	if !ok || v == nil {
		return ""
	}
	var ts int64
	switch t := v.(type) {
	case float64:
		ts = int64(t)
	case int64:
		ts = t
	case int:
		ts = int64(t)
	case string:
		if p, err := strconv.ParseInt(t, 10, 64); err == nil {
			ts = p
		} else {
			return ""
		}
	default:
		return ""
	}
	sec, nsec := normalizeUnix(ts)
	return time.Unix(sec, nsec).UTC().Format(time.RFC3339)
}

func parseIntField(fields map[string]any, key string) int {
	v, ok := fields[key]
	if !ok || v == nil {
		return 0
	}
	switch t := v.(type) {
	case float64:
		return int(t)
	case int:
		return t
	case int64:
		return int(t)
	default:
		return 0
	}
}


func normalizeUnix(ts int64) (int64, int64) {
	switch {
	case ts >= 1e16: // 纳秒 ns（~1.7e18）
		return ts / 1e9, ts % 1e9
	case ts >= 1e13: // 微秒 μs（~1.7e15）
		return ts / 1e6, (ts % 1e6) * 1e3
	case ts >= 1e10: // 毫秒 ms（~1.7e12）
		return ts / 1e3, (ts % 1e3) * 1e6
	default: // 秒 s（~1.7e9）
		return ts, 0
	}
}