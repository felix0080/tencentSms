package main

import (
	"fmt"
	"time"
	"math/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
	"log"
	"io/ioutil"
	"encoding/json"
	"errors"
)
const RND  = 100000*100000
const FORMAT  ="%v"
var t *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
type Request struct{
	Tel Tel `json:"tel,omitempty"`
	Sign string `json:"sign,omitempty"`
	TplId string `json:"tpl_id,omitempty"`
	Params []string `json:"params,omitempty"`
	Sig string `json:"sig,omitempty"`
	Time int64 `json:"time,omitempty"`
	Extend string
	Ext string
}
type Tel struct {
	Nationcode string `json:"nationcode,omitempty"`
	Mobile string `json:"mobile,omitempty"`
}
type Result struct{
	Result int `json:"result,omitempty"`//0表示成功(计费依据)，非0表示失败
	Errmsg string `json:"errmsg,omitempty"`//result非0时的具体错误信息
	Ext string `json:"ext,omitempty"`//用户的session内容，腾讯server回包中会原样返回
	Sid string `json:"sid,omitempty"`//标识本次发送id，标识一次短信下发记录
	Fee int `json:"fee,omitempty"`//短信计费的条数
}
type Tphone struct {
	StrAppKey string
	AppId string
	TempId string
}

func (tp *Tphone) Send(phone string,code string) error {
	now:=time.Now().Unix()
	client := &http.Client{}
	random := tp.rand()
	var r Request
	var t Tel
	t.Mobile=phone
	t.Nationcode="86"
	r.Tel=t
	r.Time=now
	r.TplId=tp.TempId
	r.Sig=tp.Sig(random,now,phone)
	r.Params=[]string{code}
	path:=fmt.Sprintf("https://yun.tim.qq.com/v5/tlssmssvr/sendsms?sdkappid=%s&random=%s",tp.AppId,random)
	b,err:=json.Marshal(r)
	if err!=nil {
		log.Println(err)
		return err
	}
	req, err := http.NewRequest("POST",path , strings.NewReader(string(b)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "name=anny")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var result Result
	json.Unmarshal(body,&result)
	if  result.Result != 0{
		return errors.New(result.Errmsg)
	}
	fmt.Println(result)
	return nil
}
func (tp *Tphone) rand() string {
	vcode := fmt.Sprintf(FORMAT, t.Int63n(RND))
	return (vcode)
}
func (tp *Tphone) Sig(random string,time int64,mobile string) string {
	hash := sha256.New()
	s:=fmt.Sprintf(`appkey=%s&random=%v&time=%v&mobile=%v`,tp.StrAppKey,random,time,mobile)
	hash.Write([]byte(s))
	md := hash.Sum(nil)
	mdStr := hex.EncodeToString(md)
	return mdStr
}
func main() {
	t:=Tphone{
		"",
		"",
		"",
	}
	t.Send("","746234")
}
