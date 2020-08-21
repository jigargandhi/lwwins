package address

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const broadcastAddr = "255.255.255.255:3333"

// Registrar service's structure
type Registrar struct {
	id              int
	token           string
	address         map[int]string
	listenerAddress string
	NewAddress      chan string
}

// NewRegistrar creates a new instance of registrar service
// and starts address listener and speaker
func NewRegistrar(id int, token string) *Registrar {
	address := Registrar{}
	address.id = id
	address.token = token
	address.address = make(map[int]string)
	address.NewAddress = make(chan string)
	address.listenerAddress = "0.0.0.0:3333"
	go address.addressListener()
	go address.addressSpeaker()
	return &address
}

func (address *Registrar) addressListener() {
	connection, err := net.ListenPacket("udp", address.listenerAddress)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer connection.Close()
	for {
		buffer := make([]byte, 1024)
		n, _, err := connection.ReadFrom(buffer)
		if err != nil {
			log.Printf("Error reading from buffer %v", err.Error())
		}
		// addressData is of the format id|address|secretToken
		addressData := string(buffer[0:n])
		go address.handleAddress(addressData)
	}
}

func (address *Registrar) handleAddress(addressData string) {
	log.Debug("Handling new address")
	parts := strings.Split(addressData, "|")
	if len(parts) != 3 {
		return
	}
	id, error := strconv.Atoi(parts[0])
	if error != nil {
		log.Printf("Error occured while doing parsing Id from %v: %v ", addressData, error.Error())
		return
	}
	if id == address.id {
		return
	}
	if parts[1] != address.token {
		return
	}
	if parts[2] == getLocalIP() {
		return
	}
	val, ok := address.address[id]

	notify := false
	if ok {
		if val != parts[2] {
			address.address[id] = parts[2]
			notify = true
		}
	} else {
		address.address[id] = parts[2]
		notify = true
	}
	if notify {
		log.Debugf("Added new address %v", parts[2])
		address.NewAddress <- parts[2]
	}
}

func (address *Registrar) addressSpeaker() {
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

func (address *Registrar) serialize() []byte {
	ser := fmt.Sprintf("%d|%s|%s", address.id, address.token, getLocalIP())
	return []byte(ser)
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

// ForAddress calls the handler for each address
func (address *Registrar) ForAddress(handler func(value string)) {
	for _, value := range address.address {
		if value != getLocalIP() {
			go handler(value)
		}
	}
}
