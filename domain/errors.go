package domain

import "fmt"

// CheckError is a simple error check that panics if the error is not nil
func CheckError(err error) {
	if err != nil {
		panic(fmt.Sprintf("err=%#+v", err))
	}
}
