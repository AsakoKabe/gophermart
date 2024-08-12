package server

import (
	"errors"
)

var ErrCreateDBPoll = errors.New("error to create db poll")
var ErrCreateStorages = errors.New("error to create storages")
var ErrRegisterEndpoints = errors.New("error to register endpoints")
var ErrConnectToDB = errors.New("error connect to db")
