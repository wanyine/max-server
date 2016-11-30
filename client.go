package main

import (
	"./vse"
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

func main() {
	var (
		fHost = flag.String("host", "", "Host to dial")
		fPort = flag.Int("port", 0, "Port to listen")
	)
	flag.Parse()
	if *fHost == "" || *fPort == 0 {
		fmt.Println("Example: go run client.go -host localhost -port 7777")
		flag.PrintDefaults()
		os.Exit(1)
	}
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", *fHost, *fPort), 3*time.Second)
	if err != nil {
		log.Fatal("connect failed", err)
		os.Exit(1)
	}
	log.Println("connected")
	for {
		data := make([]byte, 1024)
		num, err := conn.Read(data)
		if err != nil {
			continue
		}

		beg := 0
		for beg < num {
			x, n := proto.DecodeVarint(data[beg:num])
			id := readInt32(data[beg+n : beg+n+4])
			msg := data[beg+n+4 : beg+n+int(x)]
			switch id {
			case 0:
				netID := vse.NetId{}
				proto.Unmarshal(msg, &netID)
				log.Println(netID)

				random := rand.New(rand.NewSource(99))
				clientID := random.Int31() % 6
				player := vse.Player{
					NetId:    netID.NetId,
					ClientId: &clientID,
				}
				buf, _ := proto.Marshal(&player)
				msg := append([]byte{0, 0, 0, 1}, buf...)
				conn.Write(append(proto.EncodeVarint(uint64(len(msg))), msg...))
				msg = []byte{5, 0, 0, 0, 7, 1, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2,
					5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2,
					5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2,
					5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2,
					5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2,
					5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2,
					5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2,
					5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2,
					5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2,
					5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2,
					5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9, 2, 5, 0, 0, 0, 9}
				conn.Write(msg)

			case 2:
				players := &vse.Players{}
				proto.Unmarshal(msg, players)
				log.Println(players)
			default:
				log.Println(data[beg : beg+n+int(x)])
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
