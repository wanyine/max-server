package main

import (
	"./vse"
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	"net"
	"os"
	"os/signal"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

var bufMap = make(map[net.Conn](chan []byte))
var players = vse.Players{}
var roles = make(map[int32]int32)

func main() {

	service := ":7777"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	exitIfError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	exitIfError(err)

	go onSignal(listener)

	fmt.Println("Please input the players number(1-5):")
	rd := bufio.NewReader(os.Stdin)
	num := setPlayersNumber(rd)
	master := int32(0)
	players.Total = &num
	players.MasterId = &master
	players.List = make([]*vse.Player, 0, num)

	fmt.Println("Waiting for connection...")

	for id := int32(0); ; id++ {
		if id == num {
			fmt.Println("Players meet up, starting game")
		}
		conn, err := listener.Accept()
		exitIfError(err)
		if id < num {
			fmt.Printf("netId %d connected\n", id)
			bufMap[conn] = make(chan []byte)
			go safelyHandle(conn)
			netId := &vse.NetId{NetId: &id}
			data, _ := proto.Marshal(netId)
			write(conn, compose(int32(0), data))
			write(conn, getPlayersMessage())
		} else {
			fmt.Printf("netId %d connected, but players have been full\n", id)
			continue
		}
	}
}

func onSignal(listener net.Listener) {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan)
	<-sigChan
	for conn, _ := range bufMap {
		conn.Close()
	}
	fmt.Println("close all connections.")
	listener.Close()
	fmt.Println("close listener.")
}

func getPlayersMessage() []byte {

	players.List = make([]*vse.Player, 0, len(roles))
	for key, value := range roles {
		k, v := key, value
		player := vse.Player{
			NetId:    &k,
			ClientId: &v,
		}
		players.List = append(players.List, &player)
	}
	data, _ := proto.Marshal(&players)
	return compose(int32(2), data)
}

func setPlayersNumber(rd *bufio.Reader) int32 {
	in, err := rd.ReadBytes('\n')
	if err != nil || in[0] > '4' || in[0] < '1' {
		fmt.Println("wrong number, please input again:")
		return setPlayersNumber(rd)
	} else {
		return bytesToInt32(append([]byte{0, 0, 0}, in[0]-'0'))
	}
}

func safelyHandle(conn net.Conn) {

	defer func() {
		if err := recover(); err != nil {
			_ = err.(Error)
			delete(bufMap, conn)
			log.Println(err)
		}
	}()

	go handleWirte(conn, bufMap[conn])

	for {
		buf := make([]byte, 1024)
		num, err := conn.Read(buf)
		if err != nil {
			panic(Error(err.Error()))
		}

		beg := 0
		for beg < num {
			x, n := proto.DecodeVarint(buf[beg:num])
			id := bytesToInt32(buf[beg+n : beg+n+4])

			switch id {
			case 1:
				msg := buf[beg+n+4 : beg+n+int(x)]
				player := vse.Player{}
				if e := proto.Unmarshal(msg, &player); e != nil {
					log.Println(e)
					continue
				}
				if _, ok := roles[*(player.NetId)]; ok {
					continue
				} else {
					roles[*(player.NetId)] = *(player.ClientId)
				}
				broadcast(getPlayersMessage())
			default:
				fmt.Println(buf[:num])
				broadcast(buf[:num])
			}

			beg += int(x) + n
		}
	}
}

func compose(id int32, msg []byte) []byte {
	msg = append(int32ToBytes(id), msg...)
	return append(proto.EncodeVarint(uint64(len(msg))), msg...)
}

func int32ToBytes(num int32) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, &num)
	return buf.Bytes()
}

func bytesToInt32(slice []byte) int32 {
	var p int32
	binary.Read(bytes.NewBuffer(slice), binary.BigEndian, &p)
	return p
}

func write(conn net.Conn, msg []byte) {
	bufMap[conn] <- msg
}

func broadcast(msg []byte) {

	for _, ch := range bufMap {
		ch <- msg
	}
}

func handleWirte(conn net.Conn, sendChan <-chan []byte) {
	for data := range sendChan {
		_, err := conn.Write(data)
		if err != nil {
			continue
		}
	}
}

func exitIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
