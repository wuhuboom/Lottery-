package process

import (
	"encoding/json"
	eeor "example.com/m/error"
	"example.com/m/model"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"io"
	"net/http"
	"strconv"
	"time"
)

// PaiXu3 FuCai3D 排列三采集
func PaiXu3(db *gorm.DB) {
	for {
		now := time.Now()
		// 计算下一个运行时间（21:15）
		//nextRunTime := time.Date(now.Year(), now.Month(), now.Day(), 21, 15, 0, 0, now.Location())
		nextRunTime := time.Date(now.Year(), now.Month(), now.Day(), 20, 30, 0, 0, now.Location())
		// 如果当前时间已经过了 21:15，则将运行时间设为第二天的 21:15
		if now.After(nextRunTime) {
			nextRunTime = nextRunTime.Add(24 * time.Hour)
		}
		// 计算距离下一个 21:15 的时间间隔
		duration := nextRunTime.Sub(now)
		// 输出下一个运行时间
		fmt.Printf("Next run time: %v\n", nextRunTime)
		// 等待直到指定的时间点
		time.Sleep(duration)
		// 任务开始执行
		fmt.Println("Task is running at", time.Now())
		var ISRIGHT bool
		ISRIGHT = false
		for i := 0; i < 30; i++ {
			result, _ := GetTodayResultForPX3(db)
			if result == true {
				ISRIGHT = true
				break
			}
			fmt.Println(result)
			time.Sleep(1 * time.Minute)
		}
		if ISRIGHT == false {
			//警报告诉我
		}

	}
}

func GetTodayResultForPX3(db *gorm.DB) (bool, error) {
	//获取最后一次期数
	p3x := model.CaiJi{}
	db.Where("type=2 and kind=2").Order("qi_shu desc").First(&p3x)
	atoi, err := strconv.Atoi(p3x.QiShu)
	if err != nil {
		return false, err
	}
	atoi = atoi + 2

	url := "https://www.vipc.cn/i/results/pl3/" + strconv.Itoa(atoi)
	// 发送 GET 请求
	//resp, err := http.Get(url)
	// 创建一个新的请求
	// 创建一个新的请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		zap.L().Debug("GetTodayResult 47 -->" + err.Error())
		return false, err
	}
	req.Header.Set("Cookie", "HMF_CI=a154b9d2a701523268c7079a4ebfccbee48a8d0e5892f9963ae94ee0f2721f332ebffe51e7e021599636fd0b056e2a37da28d82fc06bae027f71b7c8d171f2f576'")
	req.Header.Set("Host", "www.cwl.gov.cn")
	req.Header.Set("User-Agent", "PostmanRuntime/7.42.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Connection", "keep-alive")
	// 创建一个 HTTP 客户端
	client := &http.Client{}
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		zap.L().Debug("GetTodayResult 59 -->" + err.Error())
		return false, err
	}
	defer resp.Body.Close()
	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.L().Debug("GetTodayResult 67 -->" + err.Error())
		return false, err
	}
	// 创建一个变量来存储解析后的数据
	var data1 Px3
	// 解析 JSON 数据
	err = json.Unmarshal(body, &data1)

	fmt.Println(string(body))
	if err != nil {
		zap.L().Debug("GetTodayResult 74 -->" + string(body))
		return false, err
	}

	if len(data1.List) < 0 {
		zap.L().Debug("GetTodayResult 80 -->" + string(body))
		return false, eeor.OtherError("is not 0")
	}
	//入库
	//var fc []model.CaiJi
	for _, i2 := range data1.List {
		db.Save(&model.CaiJi{
			QiShu:     i2.Issue,
			CaiJiTime: i2.Time,
			Created:   time.Now().Unix(),
			Result:    i2.Numbers[0] + "," + i2.Numbers[1] + "," + i2.Numbers[2],
			Type:      2, //排列三
			Kind:      2,
			Date:      time.Now().Format("2006-01-02"), Md5: i2.Issue + strconv.Itoa(2)})
	}
	//db.Save(fc)
	return true, nil
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
