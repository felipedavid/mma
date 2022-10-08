package main

func isWhiteSpace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\r' || ch == '\v'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isLower(ch byte) bool {
	return ch >= 'a' && ch <= 'z'
}

func isUpper(ch byte) bool {
	return ch >= 'A' && ch <= 'Z'
}

func isLetter(ch byte) bool {
	return isLower(ch) || isUpper(ch)
}

func isAlphaNumeric(ch byte) bool {
	return isLetter(ch) || isDigit(ch)

}

func toLower(ch byte) byte {
	if isUpper(ch) {
		ch += 32
	}
	return ch
}
