package address

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const broadcastAddr = "255.255.255.255:3333"

// Registrar service's structure
type Registrar struct {
	name            string
	token           string
	address         map[string]string
	listenerAddress string
	NewAddress      chan string
}

// NewRegistrar creates a new instance of registrar service
// and starts address listener and speaker
func NewRegistrar(nodeName string, token string) *Registrar {
	address := Registrar{}
	address.name = nodeName
	address.token = token
	address.address = make(map[string]string)
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
		go address.handleAddress(buffer[0:n])
	}
}

func (address *Registrar) handleAddress(addressData []byte) {
	log.Debug("Handling new address")
	addressParts, err := decrypt(address.token, addressData)
	if err != nil {
		return
	}
	parts := strings.Split(addressParts, "|")
	if len(parts) != 2 {
		return
	}
	name := parts[0]
	remoteAddr := parts[1]
	if name == address.name {
		return
	}
	if parts[1] == getLocalIP() {
		return
	}
	val, ok := address.address[name]

	notify := false
	if ok {
		if val != parts[1] {
			address.address[name] = remoteAddr
			notify = true
		}
	} else {
		address.address[name] = remoteAddr
		notify = true
	}
	if notify {
		log.Debugf("Added new address %v", remoteAddr)
		address.NewAddress <- remoteAddr
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
	ser := fmt.Sprintf("%s|%s", address.name, getLocalIP())
	return encrypt(address.token, ser)
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
