package error

import (
	"errors"
	"net/http"
)

var (
	ErrDuplicateTrx        = errors.New("duplicate transaction")
	ErrBalanceInsufficient = errors.New("balance insufficient")
)

var mapErrorCode = map[error]int{
	ErrDuplicateTrx:        http.StatusUnprocessableEntity,
	ErrBalanceInsufficient: http.StatusUnprocessableEntity,
}

func ParseError(err error) (code int, errMsg string) {
	if err != nil {
		if errCode, ok := mapErrorCode[err]; ok {
			code = errCode
			return code, err.Error()
		} else {
			return http.StatusInternalServerError, err.Error()
		}
	}
	return http.StatusOK, ""
}
