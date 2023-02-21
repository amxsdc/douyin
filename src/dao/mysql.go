package dao

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// DRIVER 指定驱动
const DRIVER = "mysql"

var SqlSession *gorm.DB

var COS_SECRET_ID = "AKIDzoAOVXVSH9yfFOPL8lFb4PhmgYElBmFY"
var COS_SECRET_KEY = "EKKCwHew64wdEjAWkrpeKpCUitQwWBnI"
var COS_BUCKET_NAME = "douyin-1314435734"
var COS_APP_ID = "1314435734"
var COS_REGION = "https://douyin-1314435734.cos.ap-chongqing.myqcloud.com"
var COS_URL_FORMAT = ""

// 配置参数映射结构体
type conf struct {
	Url      string `yaml:"url"`
	UserName string `yaml:"userName"`
	PassWord string `yaml:"passWord"`
	DbName   string `yaml:"dbName"`
	Port     string `yaml:"port"`
}

// 获取配置参数数据
func (c *conf) getConf() *conf {
	// 读取resources/application.yaml文件
	yamlFile, err := ioutil.ReadFile("resources/application.yaml")
	// 若出现错误，打印错误提示
	if err != nil {
		fmt.Println(err.Error())
	}
	// 将读取的字符串转换成结构体conf
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Println(err.Error())
	}
	return c
}

// InitMySql 初始化连接数据库
func InitMySql() (err error) {
	var c conf
	// 获取yaml配置参数
	conf := c.getConf()
	// 将yaml配置参数拼接成连接数据库的url
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.UserName,
		conf.PassWord,
		conf.Url,
		conf.Port,
		conf.DbName,
	)
	// 连接数据库
	SqlSession, err = gorm.Open(DRIVER, dsn)
	if err != nil {
		panic(err)
	}
	// 验证数据库连接是否成功，若成功，则无异常
	return SqlSession.DB().Ping()
}

// Close 关闭数据库连接
func Close() {
	err := SqlSession.Close()
	if err != nil {
		return
	}
}
