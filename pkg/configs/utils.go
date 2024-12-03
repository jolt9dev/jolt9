package configs

import (
	"fmt"
	"strings"
	"unicode"
)

func ParseKeyValue(value string) (string, string, error) {
	kv := value
	index := strings.Index(kv, "=")
	if index == -1 {
		return "", "", fmt.Errorf("expected key=value, got %q", kv)
	}

	key := strings.TrimSpace(kv[:index])
	value = kv[index+1:]
	quoted := false
	quote := "double"
	sb := strings.Builder{}
	for _, c := range value {
		if !quoted && sb.Len() == 0 {
			if unicode.IsSpace(c) {
				continue
			}

			if c == '"' || c == '\'' {
				quoted = true
				quote = string(c)
				continue
			}

			sb.WriteRune(c)
			continue
		}

		if quoted {
			if string(c) == quote {
				quoted = false
				break
			}
		} else {
			if unicode.IsSpace(c) {
				break
			}
		}

		sb.WriteRune(c)
	}

	value = sb.String()
	return key, value, nil
}
