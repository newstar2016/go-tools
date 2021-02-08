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
	ExternalBindPlatformAppID string `db:"app_id" desc:"企业微信绑定的开放平台id"`
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
	EnterpriseAppInfo[96519191699584] = AppInfo{AppID: "", OpenAppID: "", ExternalName: ""}
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

	externalQuery := fmt.Sprintf("select a.corp_id,b.app_id,a.name,a.created_at,a.updated_at from %s.%s a, %s.%s b where a.open_platform_id = b.id and a.id = ?", dbname, "enterprise", dbname, "enterprise_open_platform_config")

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

//获取公众号的绑定信息
func (db *ToolsDB) GetOfficiallOpenPlatformInfo(enterpriseID int, dbname string, openPlatformInfo *OpenPlatformInfo) error {
	//查询公众号的虚拟开放平台appid
	openQuery := fmt.Sprintf("select name as open_name,appid as open_appid,bind_status,created_at as gzh_bind_created_at, updated_at as gzh_bind_updated_at from %s.%s where enterprise_id= ? order by created_at desc limit 1", dbname, "official_account")

	//查询公众号开放平台信息
	if err := db.db.Get(openPlatformInfo, openQuery, enterpriseID); err != nil && err != sql.ErrNoRows {
		fmt.Println(err)
		return err
	}
	return nil
}

//获取当前scrm库中有哪些企业已经接入了用户
//func (db *ToolsDB) GetWaitMigrateEnterpriseID(enterpriseID int, dbname string) (enterpriseIDs []int ,err error) {
//	openQuery := fmt.Sprintf("select DISTINCT(enterprise_id) as enterprise_id from %s.%s where id>0", dbname, "user")
//
//	//查询公众号开放平台信息
//	if err := db.db.Select(&enterpriseIDs, openQuery, enterpriseID); err != nil && err != sql.ErrNoRows {
//		fmt.Println(err)
//		return enterpriseIDs,err
//	}
//	return enterpriseIDs,nil
//}
