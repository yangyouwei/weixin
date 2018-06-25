//目的适用于zabbix微信报警
//程序需要三个参数，有顺序的。不要搞错。
//fmt.Println("args1 is usercount" )
//fmt.Println("args2 is the mesages's title" )
//fmt.Println("args3 is messages's content" )
//在zabbix中创建用户报警媒介时，要使用企业微信中的成员详细信息中的账号
//看注释,按需修改三处
//modify by yangyouwei 2017/11/16 下午10:54:43
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	//发送消息使用导的url
	sendurl = `https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=`
	//获取token使用导的url
	get_token = `https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=`
)

var requestError = errors.New("request error,check url or network")

type access_token struct {
	Access_token string `json:"access_token"`
	Expires_in   int    `json:"expires_in"`
}

//定义一个简单的文本消息格式
type send_msg struct {
	Touser  string            `json:"touser"`
	Toparty string            `json:"toparty"`
	Totag   string            `json:"totag"`
	Msgtype string            `json:"msgtype"`
	Agentid int               `json:"agentid"`
	Text    map[string]string `json:"text"`
	Safe    int               `json:"safe"`
}

type send_msg_error struct {
	Errcode int    `json:"errcode`
	Errmsg  string `json:"errmsg"`
}

var Usage = func() {
	fmt.Println("Usage: COMMAND args1 args2 args3")
	fmt.Println("args1 is usercount")
	fmt.Println("args2 is the mesages's title")
	fmt.Println("args3 is messages's content")
}

func main() {
	args := os.Args

	if args == nil || len(args) < 2 {
		Usage() //如果用户没有输入,或参数个数不够,则调用该函数提示用户
		return
	}
	touser := &args[1]
	agid := 1000002 //按需修改为agentid,删除了部分我的信息
	agentid := &agid
	h := args[2]
	head := &h
	txt := args[3]
	content := &txt
	c := "ww5ac145aab9167a00" //按需修改为corpid,删除了部分我的信息
	corpid := &c
	s := "uhTcjUD8DfohJ2igU9jTtvNW-8w1ldq7Q0RzqfrbODU" //按需修改为secret，删除了部分我的信息
	corpsecret := &s

	var m send_msg = send_msg{Touser: *touser, Msgtype: "text", Agentid: *agentid, Text: map[string]string{"content": *head + "\n" + *content}}

	///-p "wx2468f5838693e123" -s "JbjkM1jYq8g3GaHjOTgj27y4n4_7Dsv4FV94I5BMRSrBsm_aTsMUVJMhGu_DFGDSF"
	token, err := Get_token(*corpid, *corpsecret)
	if err != nil {
		println(err.Error())
		return
	}
	buf, err := json.Marshal(m)
	if err != nil {
		return
	}
	err = Send_msg(token.Access_token, buf)
	if err != nil {
		println(err.Error())
	}
}

//发送消息.msgbody 必须是 API支持的类型
func Send_msg(Access_token string, msgbody []byte) error {
	body := bytes.NewBuffer(msgbody)
	resp, err := http.Post(sendurl+Access_token, "application/json", body)
	if resp.StatusCode != 200 {
		return requestError
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	var e send_msg_error
	err = json.Unmarshal(buf, &e)
	if err != nil {
		return err
	}
	if e.Errcode != 0 && e.Errmsg != "ok" {
		return errors.New(string(buf))
	}
	return nil
}

//通过corpid 和 corpsecret 获取token
func Get_token(corpid, corpsecret string) (at access_token, err error) {
	resp, err := http.Get(get_token + corpid + "&corpsecret=" + corpsecret)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = requestError
		return
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(buf, &at)
	if at.Access_token == "" {
		err = errors.New("corpid or corpsecret error.")
	}
	return
}

func Parse(jsonpath string) ([]byte, error) {
	var zs = []byte("//")
	File, err := os.Open(jsonpath)
	if err != nil {
		return nil, err
	}
	defer File.Close()
	var buf []byte
	b := bufio.NewReader(File)
	for {
		line, _, err := b.ReadLine()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}
		line = bytes.TrimSpace(line)
		if len(line) <= 0 {
			continue
		}
		index := bytes.Index(line, zs)
		if index == 0 {
			continue
		}
		if index > 0 {
			line = line[:index]
		}
		buf = append(buf, line...)
	}
	return buf, nil
}
