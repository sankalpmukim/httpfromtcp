package headers

import "net/http"

func isValidHeaderChars(s string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !((c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') ||
			c == '!' || c == '#' || c == '$' || c == '%' ||
			c == '&' || c == '\'' || c == '*' || c == '+' ||
			c == '-' || c == '.' || c == '^' || c == '_' ||
			c == '`' || c == '|' || c == '~') {
			return false
		}
	}
	return len(s) > 0 // token = 1*tchar
}

// utility to convert map[string]string[] to map[string]string
// to only take the last value in case of multiple.
func ConvertInbuiltHeadersToOurHeaders(from http.Header) Headers {
	to := NewHeaders()
	for k, v := range from {
		to[k] = v[0]
	}
	return to
}
