package main

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os/exec"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func initGit(c *gin.Context) {

	url := c.DefaultQuery("url", "")
	if url == "" {
		c.String(500, "url is null")
		return
	}
	config := Config{}
	fmt.Println("url:", url)
	getConfig(&config)
	fmt.Println(config)
	url = "http://" + config.Username + ":" + config.Password + "@" + url
	fmt.Println(url)
	cmd := exec.Command("/bin/bash", "-c", " cd /cicd/git && git clone "+url)

	// 创建获取命令输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		c.String(500, "Error:can not obtain stdout pipe for command:%s\n", err)
		return
	}
	// 执行命令
	if err := cmd.Start(); err != nil {
		c.String(500, "Error:The command is err,", err)
		return
	}

	// 读取所有输出
	bytes1, err := ioutil.ReadAll(stdout)
	if err != nil {
		c.String(500, "ReadAll Stdout:", err.Error())
		return
	}
	if err := cmd.Wait(); err != nil {
		c.String(500, "wait:", err.Error()+"\n", string(bytes1[:]))
		return
	}
	c.String(200, "stdout:\n\n %s", bytes1)
}

func initGit2(c *gin.Context) {

	// 升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	mt := websocket.TextMessage
	defer ws.Close()

	url := c.DefaultQuery("url", "")
	if url == "" {
		ws.WriteMessage(mt, []byte("Error:name is null\n"))
		return
	}
	config := Config{}
	fmt.Println("url:", url)
	getConfig(&config)
	fmt.Println(config)
	url = "http://" + config.Username + ":" + config.Password + "@" + url
	fmt.Println(url)
	cmd := exec.Command("/bin/bash", "-c", " cd /cicd/git && git clone "+url)

	// 创建获取命令输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		err1 := ws.WriteMessage(mt, []byte("Error:can not obtain stdout pipe for command:"+err.Error()))
		if err1 != nil {
			return
		}
		return
	}
	// 执行命令
	if err := cmd.Start(); err != nil {
		err1 := ws.WriteMessage(mt, []byte("Error:The command is err:"+err.Error()))
		if err1 != nil {
			return
		}
		return
	}

	// 读取所有输出
	for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		tmp = bytes.Trim(tmp, "\u0000")
		err1 := ws.WriteMessage(mt, tmp)
		if err1 != nil {
			return
		}
		if err != nil {
			break
		}
	}

	if err := cmd.Wait(); err != nil {
		err1 := ws.WriteMessage(mt, []byte("Error:wait:"+err.Error()))
		if err1 != nil {
			return
		}
		return
	}
	err1 := ws.WriteMessage(mt, []byte("OK"))
	if err1 != nil {
		return
	}
}

func cicd(c *gin.Context) {

	name := c.DefaultQuery("name", "")
	branch := c.DefaultQuery("branch", "master")
	project := c.DefaultQuery("project", "")
	jarName := c.DefaultQuery("jarName", "4")
	js := c.DefaultQuery("js", "5")

	if name == "" {
		c.String(500, "name is null")
		return
	}
	if project == "" {
		c.String(500, "project is null")
		return
	}

	cmdStr := " /usr/www/cicd.sh " + name + " " + branch + " " + project + " " + jarName + " " + js
	// c.String(200, cmdStr)
	cmd := exec.Command("/bin/bash", "-c", " nsenter -m -u -i -n -p  -t 1 sh -c "+cmdStr)

	// 创建获取命令输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		c.String(500, "Error:can not obtain stdout pipe for command:%s\n", err)
		return
	}
	// 执行命令
	if err := cmd.Start(); err != nil {
		c.String(500, "Error:The command is err,", err)
		return
	}

	// 读取所有输出
	bytes1, err := ioutil.ReadAll(stdout)
	if err != nil {
		c.String(500, "ReadAll Stdout:", err.Error())
		return
	}
	if err := cmd.Wait(); err != nil {
		c.String(500, "wait:", err.Error()+"\n", string(bytes1[:]))
		return
	}
	c.String(200, "stdout:\n\n %s", bytes1)
}

func cicd2(c *gin.Context) {

	name := c.DefaultQuery("name", "")
	branch := c.DefaultQuery("branch", "pos_dev#3.2")

	// 升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	mt := websocket.TextMessage
	defer ws.Close()

	if name == "" {
		err1 := ws.WriteMessage(mt, []byte("Error:name is null\n"))
		if err1 != nil {
			return
		}
		return
	}

	cmd := exec.Command("/bin/bash", "-c", "/usr/www/cicd.sh "+name+" "+branch)

	// 创建获取命令输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		err1 := ws.WriteMessage(mt, []byte("Error:can not obtain stdout pipe for command:"+err.Error()))
		if err1 != nil {
			return
		}
		return
	}
	// 执行命令
	if err := cmd.Start(); err != nil {
		err1 := ws.WriteMessage(mt, []byte("Error:The command is err:"+err.Error()))
		if err1 != nil {
			return
		}
		return
	}

	// 读取所有输出
	for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		tmp = bytes.Trim(tmp, "\u0000")
		err1 := ws.WriteMessage(mt, tmp)
		if err1 != nil {
			return
		}
		if err != nil {
			break
		}
	}

	if err := cmd.Wait(); err != nil {
		err1 := ws.WriteMessage(mt, []byte("Error:wait:"+err.Error()))
		if err1 != nil {
			return
		}
		return
	}
	err1 := ws.WriteMessage(mt, []byte("OK"))
	if err1 != nil {
		return
	}
}

func getConfig(config *Config) {
	file, err := ioutil.ReadFile("/cicd/config.yml")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = yaml.Unmarshal(file, config)
	if err != nil {
		fmt.Println(err)
		return
	}
}

type Config struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}
