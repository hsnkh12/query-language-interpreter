package query_parser

func IsLogicalAndOr(t TokenType) bool {
	if t == LOGICAL_OR || t == LOGICAL_AND {
		return true
	}

	return false
}

func IsLogicalOperation(t TokenType) bool {

	if t == LOGICAL_BIGGER || t == LOGICAL_EQUAL || t == LOGICAL_SMALLER {
		return true
	}
	return false
}
