package storage

import "errors"

var AccountNotFoundError = errors.New("account didn't exist")
var InternalError = errors.New("internal error")
