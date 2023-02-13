package main

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isSpace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\r' || ch == '\v'
}

func isAlpha(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isAlphaNumeric(ch byte) bool {
	return isAlpha(ch) || isDigit(ch)
}
