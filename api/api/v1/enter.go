package v1

import (
	"github.com/cngamesdk/live-chat/api/service"
)

type ApiGroup struct {
	ProductApi ProductApi
	ChatApi    ChatApi
	FaqApi     FaqApi
	UploadApi  UploadApi
	AdminApi   AdminApi
	SessionApi SessionApi
}

var ServiceGroupApp = new(service.ServiceGroup)
