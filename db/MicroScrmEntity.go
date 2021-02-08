package db

// wx_app_config 当前企业接入小程序、公众号、企业微信等
type WxAppConfig struct {
	AppID          string `json:"appid" db:"appid" desc:"小程序appid"`
	AppType        int    `json:"app_type" db:"appid" desc:"小程序app_type，1小程序 2企业微信 3公众号"`
	OpenPlatformID string `json:"open_platform_id" db:"open_platform_id" desc:"虚拟开放平台appid"`
	BindStatus     int    `json:"bind_status" db:"bind_status" desc:"绑定状态 0 未绑定 1已绑定"`
	Name           string `json:"name" db:"name" desc:"虚拟开放平台名称name"`
	CreatedAt      int    `json:"created_at" db:"created_at" desc:"创建时间"`
	UpdatedAt      int    `json:"updated_at" db:"updated_at" desc:"更新时间"`
}
