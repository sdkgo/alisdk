package alisdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const smsApiUrl string = "http://dysmsapi.aliyuncs.com/"

type smsClient struct {
	accessKeyId     string
	accessKeySecret string
}

type SmsOptions struct {
	TemplateParam   string
	SmsUpExtendCode string
	OutId           string
}

type SmsResponse struct {
	RequestId string
	Code      string
	Message   string
	BizId     string
}

func NewSmsClient(keyId string, keySecret string) *smsClient {
	return &smsClient{accessKeyId: keyId, accessKeySecret: keySecret}
}

func (c *smsClient) SendSms(phoneNumbers string, signName string,
	templateCode string, options *SmsOptions) (*SmsResponse, error) {
	var params = make(map[string]string)
	// 公共参数
	addCommonParams(params, c.accessKeyId)
	// 接口参数
	params["PhoneNumbers"] = phoneNumbers
	params["SignName"] = signName
	params["TemplateCode"] = templateCode
	if options != nil {
		if options.TemplateParam != "" {
			params["TemplateParam"] = options.TemplateParam
		}
		if options.SmsUpExtendCode != "" {
			params["SmsUpExtendCode"] = options.SmsUpExtendCode
		}
		if options.OutId != "" {
			params["OutId"] = options.OutId
		}
	}
	// 签名
	signedStr, paramStr := PopSignature(params, c.accessKeySecret)
	url := smsApiUrl + "?Signature=" + specialUrlEncode(signedStr) + "&" + paramStr
	// 发送请求
	res, err := request(url)
	if err != nil {
		return nil, err
	}
	result := new(SmsResponse)
	if err := json.Unmarshal(res, result); err != nil {
		return nil, errors.New(fmt.Sprintf("Send sms failed, parse json result failed：%s", err.Error()))
	}
	return result, nil
}

type Time time.Time

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"2006-01-02 15:04:05"`, string(data), time.Local)
	if err != nil {
		return err
	}
	*t = Time(now)
	return
}
func (t Time) MarshalJSON() ([]byte, error) {
	b := []byte(time.Time(t).Format(`"2006-01-02 15:04:05"`))
	return b, nil
}

type smsDetail struct {
	PhoneNum     string
	SendStatus   int
	ErrCode      string
	TemplateCode string
	Content      string
	SendDate     Time
	ReceiveDate  string
	OutId        string
}
type dtoTemp struct {
	SmsSendDetailDTO []smsDetail
}
type queryResponse struct {
	RequestId         string
	Code              string
	Message           string
	TotalCount        int
	SmsSendDetailDTOs *dtoTemp
}

func (c *smsClient) QuerySendDetails(phoneNumber string, sendDate time.Time, bizId string) ([]smsDetail, error) {
	return c.QuerySendDetailsPaged(phoneNumber, sendDate, bizId, 50, 1)
}
func (c *smsClient) QuerySendDetailsPaged(phoneNumber string, sendDate time.Time, bizId string,
	pageSize int, currentPage int) ([]smsDetail, error) {
	var params = make(map[string]string)
	// 公共参数
	addCommonParams(params, c.accessKeyId)
	params["Action"] = "QuerySendDetails"
	// 接口参数
	params["PhoneNumber"] = phoneNumber
	params["SendDate"] = sendDate.Format("20060102")
	if bizId != "" {
		params["BizId"] = bizId
	}
	params["PageSize"] = strconv.Itoa(pageSize)
	params["CurrentPage"] = strconv.Itoa(currentPage)
	// 签名
	signedStr, paramStr := PopSignature(params, c.accessKeySecret)
	url := smsApiUrl + "?Signature=" + specialUrlEncode(signedStr) + "&" + paramStr
	// 发送请求
	res, err := request(url)
	if err != nil {
		return nil, err
	}
	//fmt.Println(string(res))
	result := new(queryResponse)
	if err := json.Unmarshal(res, result); err != nil {
		return nil, errors.New(fmt.Sprintf("Query sms detail failed, parse json result failed：%s", err.Error()))
	}
	return result.SmsSendDetailDTOs.SmsSendDetailDTO, nil
}

func request(reqUrl string) ([]byte, error) {
	//fmt.Println(reqUrl)
	res, err := http.Get(reqUrl)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		body, _ := ioutil.ReadAll(res.Body)
		res.Body.Close()
		return nil, errors.New(fmt.Sprintf("Send sms failed:\n %s", body))
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}
	return body, nil
}

func addCommonParams(params map[string]string, accessKeyId string) {
	params["AccessKeyId"] = accessKeyId
	params["Timestamp"] = time.Now().UTC().Format("2006-01-02T15:04:05Z")
	params["Format"] = "JSON"
	params["SignatureMethod"] = "HMAC-SHA1"
	params["SignatureVersion"] = "1.0"
	params["SignatureNonce"] = randomUuid()
	params["Action"] = "SendSms"
	params["Version"] = "2017-05-25"
	params["RegionId"] = "cn-hangzhou"
}
