package utils

func GetIP(ForwardHeader string) string {
	if len(ForwardHeader) > 0 {
		return ForwardHeader
	}

	return ""
}
