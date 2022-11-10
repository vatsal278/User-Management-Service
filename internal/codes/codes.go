package codes

import (
	"fmt"
)

type errCode int

const (
	ErrUnauthorized errCode = iota + 1000
	ErrTokenExpired
	ErrAssertClaims
	ErrMatchingToken
	ErrAssertUserid
	ErrUnauthorizedAgent
	ErrUnauthorizedUrl
	ErrKeyNotFound
	ErrEncodingFile
	ErrConvertingToPdf
	ErrDecodingData
	ErrIdNeeded
)

var errCodes = map[errCode]string{
	ErrUnauthorized:      "UnAuthorized",
	ErrTokenExpired:      "Token is expired",
	ErrMatchingToken:     "Compared literals are not same",
	ErrAssertClaims:      "unable to assert claims",
	ErrAssertUserid:      "unable to assert userid",
	ErrUnauthorizedAgent: "UnAuthorized user agent",
	ErrUnauthorizedUrl:   "UnAuthorized url",
	ErrKeyNotFound:       "unable to find this Uuid",
	ErrEncodingFile:      "unable to json encode the data",
	ErrConvertingToPdf:   "unable to convert to pdf format",
	ErrIdNeeded:          "id needed",
	ErrDecodingData:      "unable to decode the data",
}

func GetErr(code errCode) string {
	x, ok := errCodes[code]
	if !ok {
		return ""
	}
	return fmt.Sprintf("%d: %s", code, x)
}
