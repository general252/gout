package unet

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/axgle/mahonia"
)

type JsonResult struct {
	Ip      string `json:"ip"`
	Address string `json:"address"`
	Type    string `json:"type"`
	Timeout int64  `json:"timeout"` // ms
}

// GetIpAddress 获取ip的地址
func GetIpAddress(ip string) (*JsonResult, error) {
	result := &JsonResult{
		Ip:      ip,
		Address: "",
		Type:    "",
	}
	var timeStart = time.Now()

	var err error
	for {
		if result.Address, err = getIpBaiduOpenData(ip); err == nil {
			result.Type = "Baidu OpenData"
			break
		}

		if result.Address, err = getIpBaidu(ip); err == nil {
			result.Type = "Baidu"
			break
		}

		if result.Address, err = getIpPCOnline(ip); err == nil {
			result.Type = "PCOnLine"
			break
		}

		if result.Address, err = getIpTaoBao(ip); err == nil {
			result.Type = "TaoBao"
			break
		}

		if result.Address, err = getIpNet126(ip); err == nil {
			result.Type = "Net126"
			break
		}

		err = fmt.Errorf("get ip address fail")
		break
	}

	result.Timeout = time.Since(timeStart).Milliseconds()

	return result, err
}

// GetIpAddressAsync 获取ip的地址(异步)
func GetIpAddressAsync(ip string, cb func(result *JsonResult, err error)) error {
	go func() {
		result, err := GetIpAddress(ip)
		cb(result, err)
	}()
	return nil
}

type jsonBaiduOpenData struct {
	Status       string `json:"status"`
	T            string `json:"t"`
	SetCacheTime string `json:"set_cache_time"`
	Data         []struct {
		Location         string `json:"location"`
		TitleCont        string `json:"titlecont"`
		OrigIp           string `json:"origip"`
		OrigIpQuery      string `json:"origipquery"`
		ResourceId       string `json:"resourceid"`
		OriginQuery      string `json:"OriginQuery"`
		ExtendedLocation string `json:"ExtendedLocation"`
		ShareImage       int    `json:"shareImage"`
		ShowLikeShare    int    `json:"showLikeShare"`
		ShowLamp         string `json:"showlamp"`
		FetchKey         string `json:"fetchkey"`
		AppInfo          string `json:"appinfo"`
		RoleId           int    `json:"role_id"`
		DispType         int    `json:"disp_type"`
	} `json:"data"`
}

func getIpBaiduOpenData(ip string) (string, error) {
	url := fmt.Sprintf(
		"http://opendata.baidu.com/api.php?query=%v&co=&resource_id=6006&t=%v&ie=utf8&oe=utf8&format=json",
		ip, time.Now().Unix())

	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)

	var result jsonBaiduOpenData
	if err = json.Unmarshal(data, &result); err != nil {
		return "", err
	}

	//log.Println(string(data))

	if result.Status != "0" {
		return "", fmt.Errorf("query fail %v", result.Status)
	}

	if len(result.Data) == 0 {
		return "", fmt.Errorf("query fail. no data")
	}

	return result.Data[0].Location, nil
}

type jsonBaidu struct {
	Status  string `json:"status"`
	QueryId string `json:"QueryID"`
	SrcId   string `json:"Srcid"`
	Data    []struct {
		SrcId         string `json:"srcid"`
		ResourceId    string `json:"resourceid"`
		OriginQuery   string `json:"OriginQuery"`
		OrigIpQuery   string `json:"origipquery"`
		Query         string `json:"query"`
		OrigIp        string `json:"origip"`
		Location      string `json:"location"`
		UserIp        string `json:"userip"`
		ShowLamp      string `json:"showlamp"`
		TpLt          string `json:"tplt"`
		TitleCont     string `json:"titlecont"`
		RealUrl       string `json:"realurl"`
		ShowLikeShare string `json:"showLikeShare"`
		ShareImage    string `json:"shareImage"`
	} `json:"data"`
}

func getIpBaidu(ip string) (string, error) {
	url := fmt.Sprintf("https://sp0.baidu.com/8aQDcjqpAAV3otqbppnN2DJv/api.php?query=%v&resource_id=5809&t=%v&ie=utf8&oe=utf8&format=json",
		ip, time.Now().Unix())

	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)

	var result jsonBaidu
	if err = json.Unmarshal(data, &result); err != nil {
		return "", err
	}

	//log.Println(string(data))

	if result.Status != "0" {
		return "", fmt.Errorf("query fail %v", result.Status)
	}

	if len(result.Data) == 0 {
		return "", fmt.Errorf("query fail. no data")
	}

	return result.Data[0].Location, nil
}

func getIpPCOnline(ip string) (string, error) {
	url := fmt.Sprintf("http://whois.pconline.com.cn/ip.jsp?ip=%v", ip)

	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return mahonia.NewDecoder("gbk").ConvertString(string(data)), nil
}

type jsonTaoBao struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Area      string `json:"area"`
		Country   string `json:"country"`
		IspId     string `json:"isp_id"`
		QueryIp   string `json:"queryIp"`
		City      string `json:"city"`
		Ip        string `json:"ip"`
		Isp       string `json:"isp"`
		County    string `json:"county"`
		RegionId  string `json:"region_id"`
		Region    string `json:"region"`
		CountryId string `json:"country_id"`
		CityId    string `json:"city_id"`
	} `json:"data"`
}

func getIpTaoBao(ip string) (string, error) {
	url := fmt.Sprintf("http://ip.taobao.com/outGetIpInfo?ip=%v&accessKey=alibaba-inc", ip)

	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)

	//log.Println(string(data))

	var result jsonTaoBao
	if err = json.Unmarshal(data, &result); err != nil {
		return "", err
	}

	if result.Code != 0 {
		return "", fmt.Errorf("%v", result.Msg)
	}

	return fmt.Sprintf("%v %v%v %v", result.Data.Country, result.Data.Region, result.Data.City, result.Data.Isp), nil
}

func getIpNet126(ip string) (string, error) {
	url := fmt.Sprintf("http://ip.ws.126.net/ipquery?ip=%v", ip)

	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)

	var strData = mahonia.NewDecoder("gbk").ConvertString(string(data))

	arr := strings.SplitN(strData, "\"", 5)
	if len(arr) != 5 {
		return "", fmt.Errorf("fail")
	}

	return fmt.Sprintf("%v %v", arr[1], arr[3]), nil
}
