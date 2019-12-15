package address

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const broadcastAddr = "255.255.255.255:3333"

// Registrar service's structure
type Registrar struct {
	id      int
	token   string
	address map[int]string
}

// Make creates a new instance of registrar service
func Make(id int, token string) *Registrar {
	address := Registrar{}
	address.id = id
	address.token = token
	address.address = make(map[int]string)
	return &address
}

// Start
func (address *Registrar) Start() {
	go addressListener(address)
	go addressSpeaker(address)

}

func addressListener(address *Registrar) {
	//localAddress, _ := net.ResolveUDPAddr("udp", "3333")
	//connection, err := net.ListenUDP("udp", localAddress)
	connection, err := net.ListenPacket("udp", "0.0.0.0:3333")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer connection.Close()
	for {
		buffer := make([]byte, 1024)
		_, _, err := connection.ReadFrom(buffer)
		if err != nil {
			fmt.Println(err.Error())
		}
		stringdata := string(buffer)
		parts := strings.Split(stringdata, "|")
		//fmt.Println("Received data ", stringdata)
		if len(parts) != 3 {
			continue
		}
		id, error := strconv.Atoi(parts[0])
		if error != nil {
			fmt.Println("Error occured while doing str ", error.Error())
			continue
		}
		if id == address.id {
			continue
		}
		if parts[1] != address.token {
			continue
		}
		//fmt.Println(fmt.Sprintf("Found node with Id %d and address %s", id, parts[2]))
		address.address[id] = parts[2]
	}
}

func addressSpeaker(address *Registrar) {
	connection, _ := net.Dial("udp", broadcastAddr)
	defer connection.Close()
	message := address.serialize()
	for {
		time.Sleep(2 * time.Second)
		_, err := connection.Write(message)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func (registrar *Registrar) serialize() []byte {
	ser := fmt.Sprintf("%d|%s|%s", registrar.id, registrar.token, getLocalIP())
	return []byte(ser)
}

func getLocalAddress() string {
	return "127.0.0.1"
}

func getLocalIP() string {
	var localIP string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		os.Exit(1)
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				localIP = ipnet.IP.String()
			}
		}
	}
	return localIP
}
