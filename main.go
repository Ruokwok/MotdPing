//MotdBE协议封装
package main

import (
	"encoding/hex"
	"net"
	"strconv"
	"strings"
	"time"
	"fmt"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("使用方法:")
		fmt.Println("\tmotd version\t查看版本")
		fmt.Println("\tmotd <address>\t请求服务器MOTD")
		fmt.Println()
		fmt.Println("Motd Ping 项目采用Mozilla Public License Version 2.0协议开源")
		fmt.Println("https://github.com/Ruokwok/MotdPing");
		return
	}
	var address string
	address = args[1]
	if address == "version" {
		fmt.Println("Motd Ping v1.0")
		return
	}
	if find := strings.Contains(address, ":"); !find {
		address = address + ":19132"
	}
	fmt.Println("正在请求 " + address + " 的MOTD数据\n")
	data, err := MotdBE(address)
	echo(data)
	if err != nil {
		fmt.Println(err)
	}
}

func echo(data *MotdBEInfo) {
	//fmt.Println(data)
	fmt.Println("状态=> " + data.Status)
	fmt.Println("Motd=> " + data.Motd)
	fmt.Println("协议版本=> " + strconv.Itoa(data.Agreement))
	fmt.Println("游戏版本=> " + data.Version)
	fmt.Println("在线人数=> " + strconv.Itoa(data.Online) + "/" + strconv.Itoa(data.Max))
	fmt.Println("存档名称=> " + data.LevelName)
	fmt.Println("游戏模式=> " + data.GameMode)
	fmt.Println("延迟=> " + strconv.FormatInt(data.Delay, 10) + "ms")
	fmt.Println()
}

//MotdBE信息
type MotdBEInfo struct {
	Status    string `json:"status"`     //服务器状态
	Host      string `json:"host"`       //服务器Host
	Motd      string `json:"motd"`       //Motd信息
	Agreement int    `json:"agreement"`  //协议版本
	Version   string `json:"version"`    //支持的游戏版本
	Online    int    `json:"online"`     //在线人数
	Max       int    `json:"max"`        //最大在线人数
	LevelName string `json:"level_name"` //存档名字
	GameMode  string `json:"gamemode"`   //游戏模式
	Delay     int64  `json:"delay"`      //连接延迟
}

//@description: 通过UDP请求获取MotdBE信息
//@param {string} Host 服务器地址，nyan.xyz:19132
//@return {*MotdBEInfo}
//@return {error}
func MotdBE(Host string) (*MotdBEInfo, error) {
	if Host == "" {
		MotdInfo := &MotdBEInfo{
			Status: "offline",
		}
		return MotdInfo, nil
	}

	// 创建连接
	socket, err := net.Dial("udp", Host)
	if err != nil {
		MotdInfo := &MotdBEInfo{
			Status: "offline",
		}
		return MotdInfo, err
	}
	defer socket.Close()
	// 发送数据
	time1 := time.Now().UnixNano() / 1e6 //记录发送时间
	senddata, _ := hex.DecodeString("0100000000240D12D300FFFF00FEFEFEFEFDFDFDFD12345678")
	_, err = socket.Write(senddata)
	if err != nil {
		MotdInfo := &MotdBEInfo{
			Status: "offline",
		}
		return MotdInfo, err
	}
	// 接收数据
	UDPdata := make([]byte, 4096)
	socket.SetReadDeadline(time.Now().Add(5 * time.Second)) //设置读取五秒超时
	_, err = socket.Read(UDPdata)
	if err != nil {
		MotdInfo := &MotdBEInfo{
			Status: "offline",
		}
		return MotdInfo, err
	}
	time2 := time.Now().UnixNano() / 1e6 //记录接收时间
	//解析数据
	if err == nil {
		MotdData := strings.Split(string(UDPdata), ";")
		Agreement, _ := strconv.Atoi(MotdData[2])
		Online, _ := strconv.Atoi(MotdData[4])
		Max, _ := strconv.Atoi(MotdData[5])
		MotdInfo := &MotdBEInfo{
			Status:    "online",
			Host:      Host,
			Motd:      MotdData[1],
			Agreement: Agreement,
			Version:   MotdData[3],
			Online:    Online,
			Max:       Max,
			LevelName: MotdData[7],
			GameMode:  MotdData[8],
			Delay:     time2 - time1,
		}
		return MotdInfo, nil
	}

	MotdInfo := &MotdBEInfo{
		Status: "offline",
	}
	return MotdInfo, err
}
