package model

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type CaiJi struct {
	ID     int    `gorm:"primaryKey;comment:'主键'"`
	QiShu  string ` json:"qi_shu"` //期数
	Result string `json:"result"`  //结果
	Type   int    `json:"type"`
	//1福彩3D  2排序3   3河内五分彩   4印尼五分彩  5.腾讯分分彩 6.重庆时时彩 7.台湾三分彩  8.澳洲五分彩  9.
	Kind      int    `json:"kind"`        //1 高频彩  2低频彩
	Date      string `json:"date"`        //日期
	CaiJiTime string `json:"cai_ji_time"` //采集时候网站的时间
	Created   int64  `json:"created"`     //我创建的时间
	Md5       string `json:"md5" gorm:"uniqueIndex"`
}

func CheckIsExistModelAccountChange(db *gorm.DB) {
	if db.Migrator().HasTable(&CaiJi{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&CaiJi{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		db.Migrator().CreateTable(&CaiJi{})
	}
}

// Get100FuCai3D 获取近100期的福彩3D
func Get100FuCai3D(db *gorm.DB) {
	url := "https://www.cwl.gov.cn/cwl_admin/front/cwlkj/search/kjxx/findDrawNotice?name=3d&issueCount=100&issueStart=&issueEnd=&dayStart=&dayEnd=&pageNo=1&pageSize=100&week=&systemType=PC"

	//url := "https://www.cwl.gov.cn/ygkj/kjgg/"
	// 发送 GET 请求

	//resp, err := http.Get(url)

	// 创建一个新的请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	// 添加自定义 Header
	//req.Header.Set("Connection", "keep-alive")

	req.Header.Set("Cookie", "HMF_CI=a154b9d2a701523268c7079a4ebfccbee48a8d0e5892f9963ae94ee0f2721f332ebffe51e7e021599636fd0b056e2a37da28d82fc06bae027f71b7c8d171f2f576'")
	req.Header.Set("Host", "www.cwl.gov.cn")
	req.Header.Set("User-Agent", "PostmanRuntime/7.42.0")
	req.Header.Set("Accept", "*/*")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")

	// 创建一个 HTTP 客户端
	client := &http.Client{}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// 打印响应内容
	fmt.Println("Response Body:")
	fmt.Println(string(body)) // 将字节切片转换为字符串输出
	// 打印返回的状态码
	fmt.Println("Response Status:", resp.Status)

	// 创建一个变量来存储解析后的数据
	var data FCR
	// 解析 JSON 数据
	err = json.Unmarshal(body, &data)

	var fc []CaiJi
	for _, i2 := range data.Result {
		fc = append(fc,
			CaiJi{QiShu: i2.Code,
				CaiJiTime: i2.Date,
				Created:   time.Now().Unix(), Result: i2.Red, Type: 1, Kind: 2,
				Date: time.Now().Format("2006-01-02"), Md5: i2.Code + strconv.Itoa(1),
			})
	}

	db.Save(fc)
}

type FCR struct {
	State    int    `json:"state"`
	Message  string `json:"message"`
	Total    int    `json:"total"`
	PageNum  int    `json:"pageNum"`
	PageNo   int    `json:"pageNo"`
	PageSize int    `json:"pageSize"`
	Tflag    int    `json:"Tflag"`
	Result   []struct {
		Name        string `json:"name"`
		Code        string `json:"code"`
		DetailsLink string `json:"detailsLink"`
		VideoLink   string `json:"videoLink"`
		Date        string `json:"date"`
		Week        string `json:"week"`
		Red         string `json:"red"`
		Blue        string `json:"blue"`
		Blue2       string `json:"blue2"`
		Sales       string `json:"sales"`
		Poolmoney   string `json:"poolmoney"`
		Content     string `json:"content"`
		Addmoney    string `json:"addmoney"`
		Addmoney2   string `json:"addmoney2"`
		Msg         string `json:"msg"`
		Z2Add       string `json:"z2add"`
		M2Add       string `json:"m2add"`
		Prizegrades []struct {
			Type      int    `json:"type"`
			Typenum   string `json:"typenum"`
			Typemoney string `json:"typemoney"`
		} `json:"prizegrades"`
	} `json:"result"`
}

//获取近100期的排列3

func Get100PaiXu3(db *gorm.DB) {
	var p int
	p = 24320
	for i := 0; i < 10; i++ {
		p = p - 10
		url := "https://www.vipc.cn/i/results/pl3/" + strconv.Itoa(p)
		// 发送 GET 请求
		//resp, err := http.Get(url)
		// 创建一个新的请求
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatal(err)
		}

		// 添加自定义 Header
		//req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Cookie", "HMF_CI=a154b9d2a701523268c7079a4ebfccbee48a8d0e5892f9963ae94ee0f2721f332ebffe51e7e021599636fd0b056e2a37da28d82fc06bae027f71b7c8d171f2f576'")
		req.Header.Set("Host", "www.cwl.gov.cn")
		req.Header.Set("User-Agent", "PostmanRuntime/7.42.0")
		req.Header.Set("Accept", "*/*")
		//req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		req.Header.Set("Connection", "keep-alive")
		// 创建一个 HTTP 客户端
		client := &http.Client{}
		// 发送请求
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		// 读取响应体
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		// 打印响应内容
		fmt.Println("Response Body:")
		fmt.Println(string(body)) // 将字节切片转换为字符串输出
		// 打印返回的状态码
		fmt.Println("Response Status:", resp.Status)

		//创建一个变量来存储解析后的数据
		var data Px3
		// 解析 JSON 数据
		err = json.Unmarshal(body, &data)

		var fc []CaiJi
		for _, i2 := range data.List {
			fc = append(fc, CaiJi{
				QiShu:     i2.Issue,
				CaiJiTime: i2.Time,
				Created:   time.Now().Unix(),
				Result:    i2.Numbers[0] + "," + i2.Numbers[1] + "," + i2.Numbers[2],
				Type:      2, //排列三
				Kind:      2,
				Date:      time.Now().Format("2006-01-02"), Md5: i2.Issue + strconv.Itoa(2)})
		}

		db.Save(fc)
	}
}

type Px3 struct {
	Residue int    `json:"residue"`
	Name    string `json:"name"`
	Cycle   string `json:"cycle"`
	List    []struct {
		Id      string        `json:"_id"`
		Issue   string        `json:"issue"`
		Time    string        `json:"time"`
		Numbers []string      `json:"numbers"`
		Sjh     []interface{} `json:"sjh"`
		Pool    string        `json:"pool"`
		Sale    string        `json:"sale"`
		Type    string        `json:"type"`
		Link    string        `json:"link"`
	} `json:"list"`
}
