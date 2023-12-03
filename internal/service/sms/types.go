package sms

import "context"

type Service interface {
	Send(ctx context.Context, tpl string, args []NamedArg, numbers ...string) error
	// Send(ctx context.Context, tpl string, args []string, numbers ...string) error
}

//type Request struct {
//	tpl string

type NamedArg struct {
	Val  string
	Name string
}
