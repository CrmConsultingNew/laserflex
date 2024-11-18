package endpoints

import "os"

var (
	LocalBitrixAppID     = os.Getenv("LOCAL_BITRIX_APP_ID")
	LocalBitrixAppSecret = os.Getenv("LOCAL_BITRIX_APP_SECRET")
	AccessToken          = "YOUR_ACCESS_TOKEN"
	//BitrixDomain         = "onviz.bitrix24.ru"
)

var BitrixDomain string
var NewBitrixDomain string
