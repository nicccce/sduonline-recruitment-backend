package util

import (
	"fmt"
	"sduonline-recruitment/pkg/conf"
)

type Code2SessionResp struct {
	SessionKey string `json:"session_key"`
	Openid     string `json:"openid"`
	Unionid    string `json:"unionid"`
	WxError
}
type WxError struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

func GetCode2Session(code string) (*Code2SessionResp, error) {
	var resp Code2SessionResp
	_, err := client.R().SetResult(&resp).
		ForceContentType("application/json").
		SetQueryParams(map[string]string{
			"appid":      conf.Conf.WxAppID,
			"secret":     conf.Conf.WxAppSecret,
			"js_code":    code,
			"grant_type": "authorization_code",
		}).Get("https://api.weixin.qq.com/sns/jscode2session")
	if err != nil {
		return nil, err
	}
	if resp.Errcode != 0 {
		return nil, fmt.Errorf(resp.Errmsg)
	}
	return &resp, nil
}
