package service

import (
	"context"
	"fmt"
	"github/yyfzy/mybook/internal/repository"
	"github/yyfzy/mybook/internal/service/sms"
	"math/rand"
)

var codeTplId = "SMS_154950909"

var (
	ErrCodeVerifyTooManyTimes = repository.ErrCodeVerifyTooManyTimes
)

type CodeService struct {
	repo   *repository.CodeRepository
	smsSvc sms.Service
}

func NewCodeService(repo *repository.CodeRepository, smsSvc sms.Service) *CodeService {
	return &CodeService{
		repo:   repo,
		smsSvc: smsSvc,
	}
}

func (svc *CodeService) Send(ctx context.Context, biz string, phone string) error {
	// 生成一个验证码；放到redis；发送；
	code := svc.generateCode()
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	nameArg := sms.NamedArg{
		Val:  code,
		Name: "code",
	}
	err = svc.smsSvc.Send(ctx, codeTplId, []sms.NamedArg{nameArg}, phone)
	//if err != nil {
	//	上面的操作成功了，但是这里err了怎么办
	// 不能删掉redis里的验证码，因为这个err可能是超时的err，无法确定是否发送出去
	// 考虑重试的话，在初始化的时候，传入一个自带重试机制的smsSvc
	//}
	return err
}

func (svc *CodeService) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	return svc.repo.Verify(ctx, biz, phone, inputCode)
}

func (svc *CodeService) generateCode() string {
	num := rand.Intn(1000000)
	return fmt.Sprintf("%06d", num)
}
