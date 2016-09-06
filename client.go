package main

import (
	"./vse"
	"bytes"
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		log.Println("Please tell the address of servie, keen:7777 e.g.")
		os.Exit(0)
	}
	conn, err := net.DialTimeout("tcp", os.Args[1], 3*time.Second)
	if err != nil {
		log.Fatal("connect failed", err)
		os.Exit(-1)
	}
	log.Println("connected")
	for {
		data := make([]byte, 1024)
		num, err := conn.Read(data)
		if err != nil {
			continue
		}

		log.Println(data[:num])

		beg := 0
		for beg < num {
			x, n := proto.DecodeVarint(data[beg:num])
			id := readInt32(data[beg+n : beg+n+4])
			msg := data[beg+n+4 : beg+n+int(x)]
			switch id {
			case 0:
				netId := &vse.NetId{}
				proto.Unmarshal(msg, netId)
				log.Println(netId)

				random := rand.New(rand.NewSource(99))
				clientId := random.Int31() % 6
				player := vse.Player{
					NetId:    netId.NetId,
					ClientId: &clientId,
				}
				buf, _ := proto.Marshal(&player)
				msg := append([]byte{0, 0, 0, 1}, buf...)
				conn.Write(append(proto.EncodeVarint(uint64(len(msg))), msg...))

			case 2:
				players := &vse.Players{}
				proto.Unmarshal(msg, players)
				log.Println(players)
			default:
				log.Println(data[:num])
			}
			beg += n + int(x)
		}
	}
}

func readInt32(slice []byte) int32 {
	var p int32
	binary.Read(bytes.NewBuffer(slice), binary.BigEndian, &p)
	return p
}
