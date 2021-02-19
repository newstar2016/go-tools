package db

type AppInfo struct {
	OfficialOpenAppID    string  //公众号绑定的开放平台的appid
}

var EnterpriseAppInfo map[int]AppInfo
