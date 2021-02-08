package db

type AppInfo struct {
	AppID        string  //小程序的appid
	OpenAppID    string  //公众号名称
	ExternalName string  //企业微信名称
}

var EnterpriseAppInfo map[int]AppInfo
