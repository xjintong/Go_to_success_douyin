package dao

import (
	"douyin/config"
	"douyin/models"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

// 获取宿主机的公网ip
func getIP() (string, error) {
	// 使用一个公网IP查询服务获取宿主机的公网IP
	resp, err := http.Get("https://ipinfo.io/ip")
	if err != nil {
		fmt.Println("无法获取公网IP:", err)
		return "", errors.New("无法获取公网IP")
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("无法读取响应:", err)
		return "", errors.New("无法读取响应")
	}

	// 清理和输出公网IP
	publicIP := strings.TrimSpace(string(body))
	return publicIP, nil
}


// SetupDB 初始化数据库和 ORM
func SetupDB() {

	// 获取数据库配置
	config, err := config.GetConfig("db")
	if err != nil {
		panic("获取数据库配置失败")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		config.GetString("mysql.username"),
		config.GetString("mysql.password"),
		host,
		config.GetInt("mysql.port"),
		config.GetString("mysql.dbname"))

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("链接数据库失败, error=" + err.Error())
	}

	db.AutoMigrate(&models.User{}, &models.Video{}, &models.Comment{}, &models.FavoriteVideoRelation{}, &models.FollowRelation{}, &models.Message{}, &models.FavoriteCommentRelation{})

}
