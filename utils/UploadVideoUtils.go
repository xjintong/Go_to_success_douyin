package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go"
)

// 上传文件到本地的中间件
// 返回上传文件的文件绝对路径和文件名
func SingleFileUploadMidWare(c *gin.Context) (filePath string, fileName string) {
	projectPath, _ := os.Getwd()
	//直接从formfile中获得文件
	file, _ := c.FormFile("data")
	filename := file.Filename
	extename := filename[strings.Index(filename, "."):]
	uuid := uuid.New()
	filename = uuid.String() + extename
	log.Println(filename)
	//"E:/webServer/golang/file_test/upload"+"/images/"绝对路径测试
	c.SaveUploadedFile(file, projectPath+"/upload/videos/"+filename)
	return projectPath + "/upload/videos/" + filename, filename
}

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

// 获取 minio 容器的ip
func resolveContainerIP(containerName string) (string, error) {
	// 使用 net.LookupIP 函数来解析容器名称获取IP地址
	ips, err := net.LookupIP(containerName)
	if err != nil {
		return "", err
	}

	// 获取第一个解析出的IP地址（通常只有一个）
	ip := ips[0].String()

	return ip, nil
}

// 上传视频到服务器并返回视频的地址
// 最终返回服务器可访问的地址
func UploadVideo(fileInProjectPath string, filename string) string {
	// 指定要查找的另一个容器的名称
	targetContainerName := "minio"

	// 尝试解析另一个容器的名称以获取IP地址
	minioIp, err := resolveContainerIP(targetContainerName)
	if err != nil {
		fmt.Printf("无法获取容器 %s 的IP地址: %s\n", targetContainerName, err)
	}

	fmt.Printf("容器 %s 的IP地址是: %s\n", targetContainerName, minioIp)

	endpoint := minioIp + ":9000"
	accessKeyID := "minioUser"
	secretAccessKey := "minioPassword"
	useSSL := false

	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Fatalln("创建 MinIO 客户端失败", err)
	}
	log.Printf("创建 MinIO 客户端成功")

	// 创建一个叫 mybucket 的存储桶。
	bucketName := "mybucket"
	location := "beijing"

	err = minioClient.MakeBucket(bucketName, location)
	if err != nil {
		// 检查存储桶是否已经存在。
		exists, err := minioClient.BucketExists(bucketName)
		if err == nil && exists {
			log.Printf("存储桶 %s 已经存在", bucketName)
		} else {
			log.Fatalln("查询存储桶状态异常", err)
		}
	}
	log.Printf("创建存储桶 %s 成功", bucketName)

	// 指定上传文件为 test.txt
	objectName := filename
	// // 指定上传文件路径
	filePath := fileInProjectPath
	// // 指定上传文件类型
	contentType := "video/mp4"

	// 调用 FPutObject 接口上传文件。
	n, err := minioClient.FPutObject(bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln("上传文件失败", err)
	}

	log.Printf("上传文件 %s 成功，成功上传 %d 个对象\n", objectName, n)

	// 获取宿主机公网IP
	serverIp, errIp := getIP()
	if errIp != nil {
		log.Println("获取公网ip失败")
	}

	wailian := "http://" + serverIp + ":9000" + "/" + bucketName + "/" + objectName
	return wailian
}
