package lis

import "net/http/cookiejar"

type instance struct {
	endpoint string
	userid   uint64
	username string
	password string
	cookie   cookiejar.Jar
	faketime *string
}

func NewInstance(endpoint string) *instance {
	inst := instance{
		endpoint: endpoint,
	}
	return &inst
}

func (inst *instance) GetEndpoint() string {
	return inst.endpoint
}

func (inst *instance) SetFaketime(new_time string) {
	faketime := string(new_time)
	inst.faketime = &faketime
}

func (inst *instance) GetFakeTime() string {
	return *(inst.faketime)
}
