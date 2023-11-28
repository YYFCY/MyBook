package sms

import "context"

type Service interface {
	Send(ctx context.Context, tpl string, args []string, numbers ...string) error
}

//type Request struct {
//	tpl string
//	args []string
//	numbers []string
//}
