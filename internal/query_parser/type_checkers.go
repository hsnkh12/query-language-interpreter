package query_parser

func IsLogicalAndOr(t TokenType) bool {
	if t == LOGICAL_OR || t == LOGICAL_AND {
		return true
	}

	return false
}

func IsLogicalOperation(t TokenType) bool {

	if t == LOGICAL_EQUAL || t == LOGICAL_SMALLER || t == LOGICAL_BIGGER || t == LOGICAL_ESMALLER {
		return true
	}
	return false
}

func IsTFN(t TokenType) bool {
	if t == TRUE || t == FALSE || t == NULL {
		return true
	}
	return false
}
