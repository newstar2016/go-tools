package server

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"go-tools/db"
	"go-tools/pkg/common"
	"go-tools/pkg/logger"
)

//获取小程序绑定的开放平台信息
func AddWxAppConfig(c ctx.Context, e int) {
	config := &db.WxAppConfig{}
	op, err := db.DBBossIndex.GetOpenPlatformInfo(e, db.BossDbName)

	if err != nil {
		logger.Error.Println(fmt.Sprintf("GetOpenPlatformInfo出错，err=%+v", err))
	}

	config.AppID = db.EnterpriseAppInfo[e].AppID
	config.AppType = common.AppTypeEnum.Shop
	config.Name = op.Name
	config.OpenPlatformID = op.ShopBindPlatformAppID
	config.CreatedAt = op.CreatedAt
	config.UpdatedAt = op.UpdatedAt
	config.BindStatus = common.BindStatusEnum.Yes

	err = db.DBMicroScrmIndex.AddWxAppConfig(config, fmt.Sprintf(db.MicroScrmDbName, e))

	if err != nil {
		logger.Error.Println(fmt.Sprintf("插入小程序AddWxAppConfig出错，err=%+v", err))
	}

	if op.ExternalAppID != "" {
		config.AppID = op.ExternalAppID
		config.AppType = common.AppTypeEnum.External
		corpName := "未知"
		if db.EnterpriseAppInfo[e].ExternalName!=""{
			corpName=db.EnterpriseAppInfo[e].ExternalName
		}
		config.Name = corpName
		config.OpenPlatformID = op.ExternalBindPlatformAppID
		config.CreatedAt = op.CreatedAt
		config.UpdatedAt = op.UpdatedAt
		config.BindStatus = common.BindStatusEnum.Yes

		err = db.DBMicroScrmIndex.AddWxAppConfig(config, db.MicroScrmDbName)
		if err != nil {
			logger.Error.Println(fmt.Sprintf("企业微信插入企业微信AddWxAppConfig出错，err=%+v",err))
		}
	}

	if op.GzhAppID != "" {
		config.AppID = db.EnterpriseAppInfo[e].OpenAppID
		config.AppType = common.AppTypeEnum.Open
		config.Name = op.GzhName
		config.OpenPlatformID = op.GzhAppID
		config.CreatedAt = op.GzhBindCreatedAt
		config.UpdatedAt = op.GzhBindUpdatedAt
		config.BindStatus = op.GzhBindStatus

		err = db.DBMicroScrmIndex.AddWxAppConfig(config, db.MicroScrmDbName)

		if err != nil {
			logger.Error.Println(fmt.Sprintf("公众号插入开放平台AddWxAppConfig出错，err=%+v",err))
		}
	}

}
