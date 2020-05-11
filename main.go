package main
import (
	"encoding/json"
	"fmt"
	"net/http"
	"io/ioutil"
//	"os"
	"bytes"
	"flag"
)

// {"errcode":0,
// "errmsg":"ok",
// "access_token":"PWqiNwIkF70uw8Z-Ug3OC3RLroXAGa65P4uziFRA3BGuT-g7lQx7dQxrc9IWWqPuGJ9AJATvDedbNVAGJobbcQmZssxDLkeuDBSNASOkBHiJz2LdJf4ex0HPEpGw4U-rMtwf3c5hZf2JWqbBxC6TsJwirivZL8DznfuFXmbSJRzn47tKTUlqoHMWAjU5c1H8WNI5KVXXdywChPC5qD9FjQ",
// "expires_in":7200}
// {"errcode":40013,
// "errmsg":"invalid corpid, hint: [1589168158_51_1be5fdc8eca30f4f84082c4ce408d934], from ip: 139.227.136.148, more info at https://open.work.weixin.qq.com/devtool/query?e=40013"}
type AccessToken struct {
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type Text struct{
	Content string `json:"content"`
}
type Message struct {
	Touser  string `json:"touser"`
	Toparty string `json:"toparty"`
	Totag   string `json:"totag"`
	Msgtype string `json:"msgtype"`
	Agentid int    `json:"agentid"`
	Text    Text `json:"text"`
	Safe                   int `json:"safe"`
	EnableIDTrans          int `json:"enable_id_trans"`
	EnableDuplicateCheck   int `json:"enable_duplicate_check"`
	DuplicateCheckInterval int `json:"duplicate_check_interval"`
}


type Response struct {
	Errcode      int    `json:"errcode"`
	Errmsg       string `json:"errmsg"`
	Invaliduser  string `json:"invaliduser"`
	Invalidparty string `json:"invalidparty"`
	Invalidtag   string `json:"invalidtag"`
}

type Config struct {
	ID           string  `json:"id"`
	Secret       string  `json:"secret"`
	MessageConfig      Message  `json:"message_config"`
}

var id = "ww5c3f7cf7bdda96f1"
var secret = "9zNWJgUp5sPF9GTXr_5sMMooizpWjUna-lqd4vUlhbU"

var config_file = flag.String("c", "./wechat.json", "输入配置文件")
var text = flag.String("t", "hello world", "输入发送内容")

func main() {
	flag.Parse()
    fmt.Println("-c:", *config_file)
    fmt.Println("-t:", *text)
	fmt.Println("其他参数：", flag.Args())
	
	contents,err := ioutil.ReadFile(*config_file)
	if err != nil {
		panic(err)
	}
	var config Config
	err = json.Unmarshal(contents, &config)
	if err != nil {
		panic(err)
	}

	token := fetch_access_token(config.ID, config.Secret)
	fmt.Println(token.Errcode, token.Errmsg, token.AccessToken)

	text := Text{Content: *text}
	config.MessageConfig.Text = text
	send(token.AccessToken, config.MessageConfig)
}

func fetch_access_token(id, secret string) AccessToken {
	var url = fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s",id, secret)
	ret, err := http.Get(url)
	if err != nil {
        panic(err)
	}
	defer ret.Body.Close()
	body, err := ioutil.ReadAll(ret.Body)
	if err != nil {
        panic(err)
	}
	var msg AccessToken
	err = json.Unmarshal(body, &msg)
	if err != nil {
		panic(err)
	}
	return msg
}

func send(token string, message Message) int {
	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s",token)
	msg, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(msg))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
	defer resp.Body.Close()
	// statuscode := resp.StatusCode
    // hea := resp.Header
	body, _ := ioutil.ReadAll(resp.Body)
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		panic(err)
	}
	return response.Errcode
}

