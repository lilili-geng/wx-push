package modules

type Sentence struct {
	Content     string `json:"content"`
	Note        string `json:"note"`
	Translation string `json:"translation"`
}

type WeatherInfo struct {
	Status   string     `json:"status"`
	Count    string     `json:"count"`
	Info     string     `json:"info"`
	InfoCode string     `json:"infocode"`
	Lives    []LiveInfo `json:"lives"`
}

type LiveInfo struct {
	Province      string `json:"province"`
	City          string `json:"city"` //城市名
	Adcode        string `json:"adcode"`
	Weather       string `json:"weather"`     //天气现象（汉字描述）
	Temperature   string `json:"temperature"` //实时气温，单位：摄氏度
	WindDirection string `json:"winddirection"`
	WindPower     string `json:"windpower"` //风力级别，单位：级
	Humidity      string `json:"humidity"`  //空气湿度
	ReportTime    string `json:"reporttime"`
}

type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}
