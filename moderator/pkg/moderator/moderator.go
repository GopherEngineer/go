// Пакет "Moderator" для тестирования текста на запрещенные слова.
package moderator

import (
	"regexp"
)

// Moderate принимает строку и проверяет её на запрещенные слова.
// Если в строке есть запрещенные слова, то функция вернет false.
func Moderate(str string) bool {
	// в данном случае просто используется ругулярное выражение
	reg := regexp.MustCompile("(?i)qwerty|йцукен|zxvbnm")
	res := reg.FindAllString(str, -1)
	return len(res) == 0
}
