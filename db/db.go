package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type OpenPlatformInfo struct {
	Name                      string `db:"name" desc:"小程序名称"`
	ShopAppID                 string `db:"shop_app_id" desc:"小程序appid"`
	ShopBindPlatformAppID     string `db:"open_app_id" desc:"小程序绑定的虚拟开放平台appid"`
	ExternalAppID             string `db:"corp_id" desc:"企业微信corp_id"`
	ExternalBindPlatformAppID string `db:"open_app_id" desc:"企业微信绑定的开放平台id"`
	CreatedAt                 int    `db:"created_at" desc:"创建时间"`
	UpdatedAt                 int    `db:"updated_at" desc:"更新时间"`
	GzhBindStatus             int    `db:"bind_status" desc:"公众号解绑状态：绑定状态：0-未知，1-绑定，2-解绑"`
	GzhName                   string `db:"open_name" desc:"公众号名称"`
	GzhAppID                  string `db:"open_appid" desc:"公众号appid"`
	GzhBindCreatedAt          int    `db:"gzh_bind_created_at" desc:"公众号绑定创建时间"`
	GzhBindUpdatedAt          int    `db:"gzh_bind_updated_at" desc:"公众号绑定更新时间"`
}

type ToolsDB struct {
	db *sqlx.DB
}

// DBMicroScrmIndex 索引数据库
var DBMicroScrmIndex *ToolsDB

var DBBossIndex *ToolsDB

var MicroScrmDbName string

var MicroScrmDbOldName string

var BossDbName string

func LoadEnterpriseApp() {
	EnterpriseAppInfo = make(map[int]AppInfo)
	EnterpriseAppInfo[96519191699584] = AppInfo{OfficialOpenAppID: ""}
}

// LoadMicroScrmIndex 加载scrm索引
func LoadMicroScrmIndex(dsn string) {
	DBMicroScrmIndex = NewDB(dsn)
	MicroScrmDbName = "micro_scrm_enterprise_%d"
	MicroScrmDbOldName = "micro_scrm"
}

// LoadBossIndex 加载scrm索引
func LoadBossIndex(dsn string) {
	DBBossIndex = NewDB(dsn)
	BossDbName = "crs_test_boss"
}

// NewAdDB 数据库连接
func NewDB(dsn string) *ToolsDB {
	db := &ToolsDB{}
	db.db = sqlx.MustConnect("mysql", dsn)
	return db
}

func (db *ToolsDB) CreateEnterpriseDataBase(e int) {
	databaseName := fmt.Sprintf("micro_scrm_enterprise_%d", e)
	//创建数据库
	createDatabase := "CREATE DATABASE IF NOT EXISTS " + databaseName
	db.db.Exec(createDatabase)
	//创建表
	var schemaEnterpriseUserTable = `
CREATE TABLE IF NOT EXISTS `+databaseName+`.enterprise_user (
    id bigint not null primary key,
    wx_union text not null comment '关联用户列表',
    wx_shop  text null comment '商城信息集合',
    wx_external  text null comment '外部联系人集合',
    wx_open  text null comment '公众号用户信息集合',
    tags  text null comment '标签列表',
    groups text null comment 'rfm/aipl分组',
    name  varchar(255) default '' not null comment '姓名',
    nick  varchar(255) default '' not null comment '昵称',
    avatar  varchar(255) default '' not null comment '头像',
    mobile  varchar(20)  default '' null comment '手机',
    birthday varchar(20)  default '' null comment '生日',
    sex  smallint(4)  default 0  not null comment '用户的性别，值为1时是男性，值为2时是女性，值为0时是未知',
    language varchar(20)  default '' null comment '语言',
    country varchar(255) default '' not null comment '国家',
    province varchar(255) default '' not null comment '省份',
    city varchar(255) default '' not null comment '城市',
    remark  text  null comment '备注',
    version  int  default 0  not null,
    pay_order  int default 0  not null comment '已支付订单数',
    total_fee  int default 0  not null comment '实付金额总数',
    refund_order  int  default 0  not null comment '退款订单数',
    refund_fee  int default 0  not null comment '退款金额数',
    per_customer_transaction int default 0  not null comment '客单价',
    latest_market_time  int  default 0  not null comment '最新营销时间',
    market_count int default 0  not null comment '营销次数',
    bind_to bigint default 0  not null comment '绑定对象',
    created_at  int  default 0  not null,
    updated_at  int  default 0  not null
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='企业用户聚合表';
`
	fmt.Println(schemaEnterpriseUserTable)
	db.db.MustExec(schemaEnterpriseUserTable)

	var schemaWorkRelationAddRecordTable = `
CREATE TABLE IF NOT EXISTS  `+databaseName+`.work_relation_add_record (
    corp_id  varchar(32)  default '' not null comment '企业ID',
    add_time  int unsigned default 0  not null comment '消息创建时间 （整型）',
    user_id  varchar(32)  default '' not null comment '企业服务人员的UserID',
    external_user_id varchar(32)  default '' not null comment '外部联系人的userid',
    add_channel  bigint  default 0  not null comment '添加时存在的渠道',
    primary key (corp_id, external_user_id, user_id, add_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='企业微信外部联系人添加记录表';
`
	fmt.Println(schemaWorkRelationAddRecordTable)
	db.db.MustExec(schemaWorkRelationAddRecordTable)

	var schemaWxAppConfig = `
CREATE TABLE IF NOT EXISTS `+databaseName+`.wx_app_config (
    appid  char(64) default '' not null comment '小程序或公众号appid 或 企业微信corpid' primary key,
    app_type  smallint  default 0  not null comment '1小程序 2企业微信 3公众号',
    open_platform_id char(64) default '' not null comment '当前虚拟开放平台id',
    bind_status  smallint(4)  default 0  not null comment '绑定状态：0-未知，1-绑定，2-解绑',
    name  varchar(255) default '' not null comment '名称',
    created_at  int  default 0  not null,
    updated_at  int  default 0  not null,
    deleted_at  int  default 0  not null
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='当前企业接入小程序、公众号、企业微信等';
`
	fmt.Println(schemaWxAppConfig)
	db.db.MustExec(schemaWxAppConfig)

	var schemaWxUserCorpExternal = `
CREATE TABLE IF NOT EXISTS `+databaseName+`.wx_user_corp_external (
    id  bigint  default 0  not null primary key,
    appid  char(64)  default '' not null comment '对应的是corp_id',
    openid  char(64) default '' not null comment '对应的是external_user_id',
    unionid  char(64) default '' not null comment '用户unionid',
    euid  bigint   default 0  not null comment '聚合id',
    name  varchar(255) default '' not null comment '用户名称',
    position  varchar(255) default '' not null comment '用户职称',
    avatar varchar(255) default '' not null comment '头像',
    corp_name  varchar(255) default '' not null comment '企业简称',
    corp_full_name   varchar(255) default '' not null comment '企业全称',
    type  smallint(4)  default 0  not null comment '外部联系人的类型，1表示该外部联系人是微信用户，2表示该外部联系人是企业微信用户',
    sex  smallint(4)  default 0  not null comment '外部联系人性别 0-未知 1-男性 2-女性',
    external_profile text  not null comment '外部联系人的自定义展示信息,JSON',
    created_at  int  default 0  not null comment '创建时间',
    updated_at  int  default 0  not null comment '更新时间',
    deleted_at  int  default 0  not null comment '删除时间',
    constraint uin unique (appid, openid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='企业微信外部联系人';
`
	fmt.Println(schemaWxUserCorpExternal)
	db.db.MustExec(schemaWxUserCorpExternal)

	var schemaWxUserOpen = `
CREATE TABLE IF NOT EXISTS `+databaseName+`.wx_user_open (
    id  bigint  default 0  not null comment '自增id' primary key,
    appid char(64) default '' not null comment '公众号appid',
    openid  char(64)  default '' not null comment '公众号openid',
    unionid  char(64)  default '' not null comment '当前开放平台Unionid',
    euid  bigint default 0  not null comment '聚合id',
    subscribe_at     int default 0  not null comment '订阅时间',
    subscribe_status int  default 0  not null comment '订阅状态',
    subscribe_scene  int  default 0  not null comment '订阅场景值',
    remark  varchar(255) default '' not null comment '备注',
    qr_scene  varchar(255) default '' not null comment '二维码扫码场景',
    qr_scene_str  varchar(255) default '' not null comment '二维码扫码场景描述',
    sex  tinyint unsigned default 0  not null comment '用户的性别，值为1时是男性，值为2时是女性，值为0时是未知',
    avatar varchar(255) default '' not null comment '用户的头像',
    language varchar(255) default '' not null comment '用户的语言，简体中文为zh_CN',
    city varchar(255) default '' not null comment '用户所在城市',
    province varchar(255) default '' not null comment '用户所在省份',
    country varchar(255) default '' not null comment '用户所在国家',
    nick_name varchar(255) default '' not null comment '用户昵称',
    created_at  int(10) default 0  not null comment '创建时间',
    updated_at  int(10) default 0  not null comment '更新时间',
    deleted_at  int(10) default 0  not null comment '删除时间',
    constraint uin unique (appid, openid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='公众号用户';
`
	fmt.Println(schemaWxUserOpen)
	db.db.MustExec(schemaWxUserOpen)

	var schemaWxUserShop = `
CREATE TABLE IF NOT EXISTS `+databaseName+`.wx_user_shop (
    id    bigint  default 0  not null comment '自增id' primary key,
    appid  char(64)  default '' not null comment '小程序appid',
    openid  char(64)  default '' not null comment '小程序openid',
    third_user_id varchar(64)  default '' not null comment '第三方用户id',
    unionid  varchar(64)  default '' not null comment '当前开放平台unionid',
    euid  bigint  default 0  not null comment '企业聚合id',
    nick_name  varchar(64)  default '' not null comment '用户眤称',
    name  varchar(64)  default '' not null comment '用户名称',
    avatar  varchar(255) default '' not null comment '用户头像',
    mobile  char(16) default '' not null comment '手机号',
    birthday  varchar(30)  default '' not null comment '生日',
    sex  smallint(4)  default 0  not null comment '用户的性别，值为1时是男性，值为2时是女性，值为0时是未知',
    country  varchar(255) default '' not null comment '用户所在国家',
    province  varchar(255) default '' not null comment '用户所在省份',
    city  varchar(255) default '' not null comment '用户所在城市',
    created_at  int default 0  not null comment '创建时间',
    updated_at  int default 0  not null comment '更新时间',
    deleted_at  int default 0  not null comment '删除时间',
    constraint uin unique (appid, openid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='交易用户（小程序用户）';
`
	fmt.Println(schemaWxUserShop)
	db.db.MustExec(schemaWxUserShop)

	var schemaWxUserUnion = `
CREATE TABLE IF NOT EXISTS `+databaseName+`.wx_user_union (
    unionid  char(64)  default '' not null comment '微信开放平台unionid' primary key,
    open_platform_id  char(64) default '' not null comment 'unionid所属开放平台',
    open_platform_source smallint(4) default 0  not null comment 'unionid归属开放平台是否确认1确认 0当前配置',
    euid  bigint  default 0  not null comment '企业聚合用户id'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='开放平台unionid映射聚合id'
`
	fmt.Println(schemaWxUserUnion)
	db.db.MustExec(schemaWxUserUnion)
}

func (db *ToolsDB) AddWxAppConfig(wxAppConfig *WxAppConfig, dbname string) error {
	wxAppConfigQuery := fmt.Sprintf("insert into %s.%s (appid,app_type,open_platform_id,bind_status,name,created_at,updated_at) values (?,?,?,?,?,?,?)", dbname, "wx_app_config")
	_, err := db.db.Exec(wxAppConfigQuery, wxAppConfig.AppID, wxAppConfig.AppType, wxAppConfig.OpenPlatformID, wxAppConfig.BindStatus, wxAppConfig.Name, wxAppConfig.CreatedAt, wxAppConfig.UpdatedAt)
	if err != nil {
		return err
	}
	return err
}

func (db *ToolsDB) GetOpenPlatformInfo(enterpriseID int, dbname string) (*OpenPlatformInfo, error) {
	openPlatformInfo := &OpenPlatformInfo{}
	//查询小程序的虚拟开放平台appid
	shopQuery := fmt.Sprintf("select open_app_id from %s.%s where enterprise_id = ?", dbname, "enterprise_virtual_open_platform")

	externalQuery := fmt.Sprintf("select a.corp_id,b.app_id as shop_app_id,a.name,a.created_at,a.updated_at from %s.%s a, %s.%s b where a.id = b.enterprise_id and a.id = ?", dbname, "enterprise", dbname, "enterprise_third_part_account")

	//查询小程序企业开放平台信息
	if err := db.db.Get(openPlatformInfo, shopQuery, enterpriseID); err != nil {
		fmt.Println(err)
		return openPlatformInfo, err
	}
	//查询企业微信开放平台信息
	if err := db.db.Get(openPlatformInfo, externalQuery, enterpriseID); err != nil {
		fmt.Println(err)
		return openPlatformInfo, err
	}
	//MicroScrmDbName = fmt.Sprintf(MicroScrmDbName, enterpriseID)
	err := db.GetOfficiallOpenPlatformInfo(enterpriseID, MicroScrmDbOldName, openPlatformInfo)
	if err != nil {
		fmt.Println(err)
	}
	return openPlatformInfo, nil
}

//获取公众号的信息
func (db *ToolsDB) GetOfficiallOpenPlatformInfo(enterpriseID int, dbname string, openPlatformInfo *OpenPlatformInfo) error {
	openQuery := fmt.Sprintf("select name as open_name,appid as open_appid,bind_status,created_at as gzh_bind_created_at, updated_at as gzh_bind_updated_at from %s.%s where enterprise_id= ? order by created_at desc limit 1", dbname, "official_account")

	//查询公众号开放平台信息
	if err := db.db.Get(openPlatformInfo, openQuery, enterpriseID); err != nil && err != sql.ErrNoRows {
		fmt.Println(err)
		return err
	}
	return nil
}
