package ioc

import (
	"github/yyfzy/mybook/internal/service/sms"
	"github/yyfzy/mybook/internal/service/sms/memory"
)

func InitSMSService() sms.Service {
	// 此处可更换其他实现
	return memory.NewService()
}
