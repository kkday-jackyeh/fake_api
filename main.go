package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB
var err error

type Api struct {
	ID             uint   `gorm:"AUTO_INCREMENT" json:"id"`
	Protocol       string `gorm:"unique_index:idx_complex" json:"protocol"`
	URL            string `gorm:"unique_index:idx_complex" json:"url"`
	ResponseStatus uint   `json:"responseStatus"`
	ResponseStr    string `gorm:"size:4096" json:"responseStr"`
}

func getResposneFromDb(url, protocol string) ([]byte, uint) {
	content := Api{}
	db.Where("protocol = ? and url = ? ", protocol, url).First(&content)
	return []byte(content.ResponseStr), content.ResponseStatus
}

func getResposne(c *gin.Context, protocol string) {
	url := c.Request.URL.String()
	str, status := getResposneFromDb(url, protocol)

	if len(str) == 0 {
		c.String(http.StatusNotFound, "Page Not Found")
	}

	c.Data(int(status), "application/json; charset=utf-8", str)
}

func getting(c *gin.Context) {
	getResposne(c, "GET")
}

func posting(c *gin.Context) {
	getResposne(c, "POST")
}

func deleting(c *gin.Context) {
	getResposne(c, "DELETE")
}

func putting(c *gin.Context) {
	getResposne(c, "PUT")
}

func getRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/*url", getting)
	router.POST("/*url", posting)
	router.PUT("/*url", putting)
	router.DELETE("/*url", deleting)

	return router
}

func checkErr(err error, c *gin.Context) {
	if err != nil {
		c.String(http.StatusInternalServerError, "?", err)
	}
}

func getContents(c *gin.Context) {

	content := []Api{}
	err := db.Find(&content).Error
	checkErr(err, c)

	res, e := json.Marshal(content)
	checkErr(e, c)

	c.Data(http.StatusOK, "application/json; charset=utf-8", res)
}

func createContent(c *gin.Context) {

	httpStatus, e := strconv.Atoi(c.PostForm("responseStatus"))
	checkErr(e, c)

	api := Api{Protocol: c.PostForm("protocol"),
		URL:            c.PostForm("url"),
		ResponseStr:    c.PostForm("responseStr"),
		ResponseStatus: uint(httpStatus)}

	db.NewRecord(api)

	err := db.Create(&api).Error
	checkErr(err, c)

	c.String(http.StatusOK, "")
}

func updateContent(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("index"))
	checkErr(err, c)

	httpStatus, e := strconv.Atoi(c.PostForm("responseStatus"))
	checkErr(e, c)

	content := Api{ID: uint(id)}
	db.First(&content)

	content.Protocol = c.PostForm("protocol")
	content.URL = c.PostForm("url")
	content.ResponseStatus = uint(httpStatus)
	content.ResponseStr = c.PostForm("responseStr")

	err = db.Save(&content).Error
	checkErr(err, c)

	c.String(http.StatusOK, "true")
}

func deleteContent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("index"))
	checkErr(err, c)

	content := Api{ID: uint(id)}

	err = db.Delete(&content).Error
	checkErr(err, c)

	c.String(http.StatusOK, "true")
}

func getConsoleRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/apicontent", getContents)
	router.POST("/apicontent", createContent)
	router.PUT("/apicontent/:index", updateContent)
	router.DELETE("/apicontent/:index", deleteContent)
	router.OPTIONS("/apicontent/:index", getContents)
	router.Static("/static", "./public/static")
	router.StaticFile("/", "./public/index.html")
	return router
}

func initDb() {
	db, err = gorm.Open("sqlite3", "./gorm.db")
	if err != nil {
		fmt.Println(err)
	}
	// db.DropTable(&Api{})
	db.AutoMigrate(&Api{})
}

func main() {
	defer db.Close()
	initDb()

	fakeAPIRouter := getRouter()
	go fakeAPIRouter.Run(":8880") // listen and serve on 0.0.0.0:8080

	consoleRouter := getConsoleRouter()
	consoleRouter.Run(":8881")
}
