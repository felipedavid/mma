package main

func isSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r' || b == '\v'
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func isLower(b byte) bool {
	return b >= 'a' && b <= 'z'
}

func toLower(b byte) byte {
	if isLower(b) {
		return b
	}
	return b + 32
}
