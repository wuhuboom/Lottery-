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

// FuCai3D 福彩3D采集
func FuCai3D(db *gorm.DB) {
	for {
		now := time.Now()
		// 计算下一个运行时间（21:15）
		nextRunTime := time.Date(now.Year(), now.Month(), now.Day(), 21, 15, 0, 0, now.Location())
		//nextRunTime := time.Date(now.Year(), now.Month(), now.Day(), 22, 19, 0, 0, now.Location())
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
			result, _ := GetTodayResult(time.Now().Format("2006-01-02"), db)
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

func GetTodayResult(data string, db *gorm.DB) (bool, error) {
	url := "https://www.cwl.gov.cn/cwl_admin/front/cwlkj/search/kjxx/findDrawNotice?name=3d&issueCount=&issueStart=&issueEnd=&dayStart=" + data + "&dayEnd=" + data + "&pageNo=1&pageSize=30&week=&systemType=PC"
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
	var data1 model.FCR
	// 解析 JSON 数据
	err = json.Unmarshal(body, &data1)
	if err != nil {
		zap.L().Debug("GetTodayResult 74 -->" + string(body))
		return false, err
	}

	if data1.State != 0 {
		zap.L().Debug("GetTodayResult 80 -->" + string(body))
		return false, eeor.OtherError("is not 0")
	}
	//入库
	db.Save(&model.CaiJi{
		QiShu:     data1.Result[0].Code,
		CaiJiTime: data1.Result[0].Date,
		Created:   time.Now().Unix(),
		Result:    data1.Result[0].Red,
		Type:      1, Kind: 2,
		Date: time.Now().Format("2006-01-02"),
		Md5:  data1.Result[0].Code + strconv.Itoa(1),
	})
	return true, nil
}
