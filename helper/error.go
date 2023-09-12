package helper

import "fmt"

var (
	Forbidden      = fmt.Errorf("Forbidden")
	InvalidToken   = fmt.Errorf("Invalid Token")
	NotFound       = fmt.Errorf("Data not found")
	InternalServer = fmt.Errorf("Internal Server Error")
	InvalidChiper  = fmt.Errorf("invalid ciphertext block size")
	BadGateway     = fmt.Errorf("Bad Gateway")
	Conflict       = fmt.Errorf("Conflicts")
	AccessDenied   = fmt.Errorf("Access Denied")
)

func PanicIfError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
