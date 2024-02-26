package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testwx/modules"

	"github.com/spf13/viper"
)

func RunScheduledTasks() {
	InitConfig()

	// 定时
	// cronScheduler := cron.New()
	// cronScheduler.AddFunc("0 10 19 * * *", Weather) // 每天晚上7:10执行Weather函数
	// cronScheduler.Start()

	// fmt.Println("定时任务已启动")
	Weather()
}

func GetLatandLng() {
	resp, err := http.Get("https://ipapi.co/json/")
	if err != nil {
		fmt.Println("HTTP请求失败:", err)
		return
	}

	fmt.Println(resp)
}
func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("config")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("读取配置文件失败:", err)
		return
	}
}

// 封装请求
func getResponseBody(url string, data interface{}) interface{} {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("发送HTTP GET请求失败：%v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应体失败：%v", err)
	}

	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("解析JSON失败:", err)
		return err
	}

	return data
}

// 深圳高德api
func getShenZhen(city string) modules.LiveInfo {
	key := viper.GetString("gd.key")
	url := fmt.Sprintf("https://restapi.amap.com/v3/weather/weatherInfo?city=%s&key=%s", city, key)

	var weatherInfo modules.WeatherInfo

	data := getResponseBody(url, &weatherInfo)

	if err, ok := data.(error); ok {
		fmt.Println("获取天气信息失败:", err)
		return modules.LiveInfo{}
	}

	fmt.Println("weatherInfo.Lives", weatherInfo.Lives)

	return weatherInfo.Lives[0]
}

// 获取token
func Getaccesstoken() string {

	appid := viper.GetString("wx.appid")

	appsecret := viper.GetString("wx.appsecret")

	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%v&secret=%v", appid, appsecret)

	var token modules.Token

	data := getResponseBody(url, &token)

	if err, ok := data.(error); ok {
		fmt.Println("获取天气信息失败:", err)
		return ""
	}

	fmt.Println("token.AccessToken", token.AccessToken)

	return token.AccessToken
}

// 获取关注者列表
func Getflist(access_token string) interface{} {
	url := "https://api.weixin.qq.com/cgi-bin/user/get?access_token=" + access_token + "&next_openid="

	var flist interface{}

	data := getResponseBody(url, &flist)

	if err, ok := data.(error); ok {
		fmt.Println("获取天气信息失败:", err)
		return nil
	}

	fmt.Println("dataGetflist", data)

	return flist
}

// 发送天气预报
func Weather() {
	access_token := Getaccesstoken()

	if access_token == "" {
		return
	}

	flist := Getflist(access_token)

	if flist == nil {
		return
	}

	flistMap, ok := flist.(map[string]interface{})

	if !ok {
		fmt.Println("flist 不是 map[string]interface{} 类型")
		return
	}

	// 遍历 flistMap
	for _, value := range flistMap {
		switch v := value.(type) {
		case string:
			fmt.Println("openid:", v)
			sendweather(access_token, v)
		case map[string]interface{}:
			if openid, ok := v["openid"].(string); ok {
				fmt.Println("openid:", openid)
			}
		}
	}
}

// 发送天气
func sendweather(access_token, openid string) {

	fmt.Println("access_token", access_token)

	fmt.Println("openid", openid)

	data := getShenZhen("440300")

	weatTemplateID := viper.GetString("wx.weatTemplateID")

	fmt.Println("data", data)

	reqdata := "{\"city\":{\"value\":\"" + data.City + "\", \"color\":\"#0000CD\"}, \"day\":{\"value\":\"" + data.ReportTime + "\"}, \"wea\":{\"value\":\"" + data.Weather + "\"}, \"tem1\":{\"value\":\"" + data.Temperature + "°C" + "\"}}"

	fmt.Println(reqdata)

	templatepost(access_token, reqdata, "", weatTemplateID, openid)
}

func templatepost(access_token string, reqdata string, fxurl string, templateid string, openid string) {
	url := "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=" + access_token

	fmt.Println("templatepost", url)

	reqbody := "{\"touser\":\"" + openid + "\", \"template_id\":\"" + templateid + "\", \"url\":\"" + fxurl + "\", \"data\": " + reqdata + "}"

	fmt.Println("reqbody", reqbody)

	resp, err := http.Post(url,
		"application/x-www-form-urlencoded",
		strings.NewReader(string(reqbody)))

	if err != nil {
		fmt.Println("err", err)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("string(body)", string(body))
}
