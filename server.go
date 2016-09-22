package main

import (
	"./vse"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/robfig/config"
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

	cfg, err := config.ReadDefault("config.ini")
	exitIfError(err)

	port, err := cfg.Int(config.DEFAULT_SECTION, "port")
	exitIfError(err)

	tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf(":%d", port))
	exitIfError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	exitIfError(err)

	go onSignal(listener)

	number, err := cfg.Int(config.DEFAULT_SECTION, "number")
	exitIfError(err)
	printfln("There are %d players totally", number)

	master := int32(0)
	total := int32(number)
	players.Total = &total
	players.MasterId = &master
	players.List = make([]*vse.Player, 0, number)
	fmt.Println("Waiting for connection...")

	for i := 0; ; i++ {
		id := int32(i)
		if i == number {
			fmt.Println("Players meet up, starting game")
		}

		conn, err := listener.Accept()
		exitIfError(err)

		if i < number {
			printfln("netId %d connected", id)
			bufMap[conn] = make(chan []byte)
			go safelyHandle(conn)
			netId := &vse.NetId{NetId: &id}
			data, _ := proto.Marshal(netId)
			write(conn, compose(int32(0), data))
			write(conn, getPlayersMessage())
		} else {
			printfln("netId %d connected, but players have been full", id)
			continue
		}
	}
}

func onSignal(listener net.Listener) {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan)
	<-sigChan
	for conn := range bufMap {
		conn.Close()
	}
	fmt.Println("close all connections.")
	listener.Close()
	log.Fatal("close listener and exit.")
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
				msg := buf[beg : beg+n+int(x)]
				fmt.Println(msg)
				broadcast(msg)
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

func printfln(format string, a ...interface{}) (n int, err error) {
	return fmt.Println(fmt.Sprintf(format, a))
}

func exitIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
