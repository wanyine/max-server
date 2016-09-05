package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	// "strconv"
	// "strings"
	// "./max"
	"./vse"
	"github.com/golang/protobuf/proto"
)

var bufMap = make(map[net.Conn](chan []byte))
var players = vse.Players{}
var roles = make(map[int32]int32)

func main() {
	service := ":7777"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	// checkError(err)
	fmt.Println("Please input the players number(1-5):")
	rd := bufio.NewReader(os.Stdin)
	var num = setPlayersNumber(rd)
	var master int32 = 0
	fmt.Println("Waiting for connection...")

	var id int32 = 0
	players.Total = &num
	players.MasterId = &master
	players.List = make([]*vse.Player, 0, num)

	for ; id < num; id++ {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("connection error", err)
			continue
		}

		fmt.Printf("netId %d connected\n", id)
		bufMap[conn] = make(chan []byte)
		go handleRequest(conn)

		netId := &vse.NetId{NetId: &id}
		data, _ := proto.Marshal(netId)

		write(conn, compose(int32(0), data))
		write(conn, getPlayersMessage())
	}
	fmt.Println("Players meet up, game started")

	err = listener.Close()
	if err != nil {
		log.Fatal("close listener failed", err)
	}

	for {
	}
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
	log.Printf("%T,%v\n", players, players)
	data, _ := proto.Marshal(&players)
	return compose(int32(2), data)
}

func setPlayersNumber(rd *bufio.Reader) int32 {
	in, err := rd.ReadBytes('\n')
	if err != nil || len(in) > 2 || in[0] > '4' || in[0] < '1' {
		fmt.Println("wrong number, please input again:")
		return setPlayersNumber(rd)
	} else {
		return bytesToInt32(append([]byte{0, 0, 0}, in[0]-'0'))
	}
}

func handleRequest(conn net.Conn) {

	go handleWirte(conn, bufMap[conn])

	for {
		buf := make([]byte, 1024)
		l, e := conn.Read(buf)
		if e != nil {
			continue
		}

		id := bytesToInt32(buf[:4])
		switch id {
		case 1:
			msg := buf[4:l]
			player := vse.Player{}
			proto.Unmarshal(msg, &player)
			if _, ok := roles[*(player.NetId)]; ok {
				continue
			} else {
				roles[*(player.NetId)] = *(player.ClientId)
			}
			broadcast(getPlayersMessage())
		default:
			fmt.Println(buf[:l])
			broadcast(buf[:l])
		}
	}

	defer func() { // will it be executed when exception occured?
		conn.Close()
		fmt.Println("close")
		delete(bufMap, conn)
	}()
}

func compose(id int32, msg []byte) []byte {
	msg = append(int32ToBytes(id), msg...)
	return append(proto.EncodeVarint(uint64(len(msg))), msg...)
}

func int32ToBytes(num int32) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, &num)
	// buf := bytes.NewBuffer([]byte{})
	// binary.Write(&buf, binary.BigEndian, &num)
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

// func checkError(err string) {
// }
