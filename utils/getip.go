package utils

import "strings"

func GetIP(ForwardHeader string) string {
	// The X-Forwarded-For might be like : client, proxy1, proxy2
	// So, the first field is the client IP
	ForwardedFields := strings.Split(ForwardHeader, ",")
	if len(ForwardedFields) > 0 {
		return ForwardedFields[0]
	}

	return ""
}
