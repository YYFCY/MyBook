package tencent

import (
	"context"
	"fmt"
	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/slice"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	mysms "github/yyfzy/mybook/internal/service/sms"
)

type Service struct {
	appId    *string
	signName *string
	client   *sms.Client
}

func NewService(client *sms.Client, appId string, signName string) *Service {
	return &Service{
		appId:    ekit.ToPtr(appId),
		signName: ekit.ToPtr(signName),
		client:   client,
	}
}

func (s *Service) Send(ctx context.Context, tpl string, args []mysms.NamedArg, numbers ...string) error {
	req := sms.NewSendSmsRequest()
	req.SmsSdkAppId = s.appId
	req.SignName = s.signName
	req.TemplateId = ekit.ToPtr(tpl)
	req.PhoneNumberSet = s.toStringPtrSlice(numbers)
	req.TemplateParamSet = slice.Map[mysms.NamedArg, *string](args, func(idx int, src mysms.NamedArg) *string {
		return &src.Val
	})
	resp, err := s.client.SendSmsWithContext(ctx, req)
	if err != nil {
		return err
	}
	for _, status := range resp.Response.SendStatusSet {
		if status.Code == nil || *(status.Code) != "Ok" {
			return fmt.Errorf("发送失败 code: %s, 原因: %s", *status.Code, *status.Message)
		}
	}
	return nil
}

func (s *Service) toStringPtrSlice(src []string) []*string {
	return slice.Map[string, *string](src, func(idx int, src string) *string {
		return &src
	})
}
