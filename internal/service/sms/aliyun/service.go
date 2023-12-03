package aliyun

import (
	"context"
	"encoding/json"
	"fmt"
	_ "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	sms "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	"github.com/ecodeclub/ekit"
	mysms "github/yyfzy/mybook/internal/service/sms"
)

type Service struct {
	signName string
	client   *sms.Client
}

func NewService(client *sms.Client, signName string) *Service {
	return &Service{signName: signName, client: client}
}

func (s *Service) Send(ctx context.Context, tpl string, args []mysms.NamedArg, numbers ...string) error {
	argsMap := make(map[string]string, len(args))
	for _, arg := range args {
		argsMap[arg.Name] = arg.Val
	}
	bCode, err := json.Marshal(argsMap)
	if err != nil {
		return err
	}
	for _, number := range numbers {
		smsRequest := &sms.SendSmsRequest{
			PhoneNumbers:  ekit.ToPtr(number),
			SignName:      ekit.ToPtr(s.signName),
			TemplateCode:  ekit.ToPtr(tpl),
			TemplateParam: ekit.ToPtr(string(bCode)),
		}
		smsResponse, err := s.client.SendSms(smsRequest)
		if err != nil {
			return err
		}
		if smsResponse.StatusCode == nil || *smsResponse.Body.Code != "OK" {
			return fmt.Errorf("send sms failed: %s", *smsResponse.Body.Message)
		}
	}
	return nil
}
