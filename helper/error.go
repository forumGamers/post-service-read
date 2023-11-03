package helper

import (
	"fmt"

	"github.com/olivere/elastic/v7"
)

var (
	Forbidden      = fmt.Errorf("Forbidden")
	InvalidToken   = fmt.Errorf("Invalid Token")
	NotFound       = fmt.Errorf("Data not found")
	InternalServer = fmt.Errorf("Internal Server Error")
	InvalidChiper  = fmt.Errorf("invalid ciphertext block size")
	BadGateway     = fmt.Errorf("Bad Gateway")
	Conflict       = fmt.Errorf("Conflicts")
	AccessDenied   = fmt.Errorf("Access Denied")
	TimeOut        = fmt.Errorf("Timeout")
	Limit          = fmt.Errorf("Limit")
)

func PanicIfError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func ElasticError(err error) error {
	switch true {
	case elastic.IsConflict(err):
		return Conflict
	case elastic.IsConnErr(err):
		return BadGateway
	case elastic.IsForbidden(err):
		return Forbidden
	case elastic.IsNotFound(err):
		return NotFound
	case elastic.IsTimeout(err):
		return TimeOut
	case elastic.IsUnauthorized(err):
		return AccessDenied
	case elastic.IsStatusCode(err, 429):
		return Limit
	default:
		return InternalServer
	}
}
