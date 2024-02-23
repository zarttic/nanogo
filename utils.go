package nanogo

import "strings"

// SubstringLast 返回 str 中最后一次出现 substr 后面的子串
func SubstringLast(str string, substr string) string {
	index := strings.Index(str, substr) // 获取 substr 在 str 中的索引
	if index < 0 {                      // 如果索引小于 0，表示 substr 不存在于 str 中
		return "" // 返回空字符串
	}
	return str[index+len(substr):] // 返回从 substr 后面开始的子串
}
