package methods

func TruncateTextFromEnd(text string, maxLength int) string {
	if len(text) > maxLength {
		maxLength -= 3
		startIndex := len(text) - maxLength
		return "..." + text[startIndex:]
	}
	return text
}
