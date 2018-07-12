package main // import "github.com/juliens/tcpproxy"

import (
	"encoding/json"
	"flag"
	"net"
	"os"

	"crypto/tls"
	"fmt"

	"github.com/sirupsen/logrus"
)

var port, certFile, keyFile string

func init() {
	flag.StringVar(&port, "port", "8080", "give me a port number")
	flag.StringVar(&certFile, "cert", "", "cert file")
	flag.StringVar(&keyFile, "key", "", "key file")
}

func main() {
	flag.Parse()

	var err error
	var l net.Listener
	fmt.Println(certFile, keyFile)
	if certFile != "" && keyFile != "" {
		config := &tls.Config{}
		config.Certificates = make([]tls.Certificate, 1)
		config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			logrus.Fatalf("Invalid certificate: %v", err)
		}
		l, err = tls.Listen("tcp", ":"+port, config)
		fmt.Printf("Start TCP WhoamI on port %s with cert %s and key %s\n", port, certFile, keyFile)
	} else {
		l, err = net.Listen("tcp", ":"+port)
		fmt.Printf("Start TCP WhoamI on port %s\n", port)
	}

	if err != nil {
		logrus.Fatalf("Error while creating listener: %v", err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			logrus.Errorf("Error while Accept: %v", err)
		}

		_, err = conn.Write(getData())
		if err != nil {
			logrus.Errorf("Error while Write: %v", err)
		}
		conn.Close()
	}
}

func getData() []byte {
	hostname, _ := os.Hostname()
	data := struct {
		Hostname string   `json:"hostname,omitempty"`
		IP       []string `json:"ip,omitempty"`
	}{
		hostname,
		[]string{},
	}

	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			data.IP = append(data.IP, ip.String())
		}
	}
	bytes, _ := json.Marshal(data)
	return bytes
}
