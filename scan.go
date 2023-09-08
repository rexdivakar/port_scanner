// main.go

package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

var serviceMapping = map[int]string{
	1:    "TCP Port Service Multiplexer (TCPMUX)",
	2:    "CompressNET Management Utility",
	3:    "CompressNET Compression Process",
	5:    "Remote Job Entry (RJE)",
	7:    "Echo Protocol",
	9:    "Discard Protocol",
	11:   "Active Users (systat service)",
	13:   "Daytime Protocol",
	17:   "Quote of the Day (QOTD)",
	18:   "Message Send Protocol (MSP)",
	19:   "Character Generator Protocol (CHARGEN)",
	20:   "FTP Data (File Transfer Protocol)",
	21:   "FTP Control (File Transfer Protocol)",
	22:   "SSH (Secure Shell)",
	23:   "Telnet",
	25:   "SMTP (Simple Mail Transfer Protocol)",
	53:   "DNS (Domain Name System)",
	80:   "HTTP (Hypertext Transfer Protocol)",
	110:  "POP3 (Post Office Protocol version 3)",
	123:  "Network Time Protocol (NTP)",
	137:  "NetBIOS Name Service",
	138:  "NetBIOS Datagram Service",
	139:  "NetBIOS Session Service",
	143:  "IMAP (Internet Message Access Protocol)",
	161:  "Simple Network Management Protocol (SNMP)",
	443:  "HTTPS (HTTP Secure)",
	465:  "Simple Mail Transfer Protocol Secure",
	587:  "SMTP (Alternative Port)",
	993:  "Internet Message Access Protocol Secure",
	995:  "POP3S (Secure POP3)",
	3306: "MySQL Database",
	5432: "PostgresSQL Database",
	8080: "HTTP (Alternative Port)",
	8081: "HTTP (Alternative Port)",

	// Add more port details as needed
}

func scanPort(target string, port int, wg *sync.WaitGroup, serviceMapping map[int]string) {
	defer wg.Done()

	address := fmt.Sprintf("%s:%d", target, port)
	conn, err := net.DialTimeout("tcp", address, 1*time.Second)
	if err != nil {
		return // Port closed
	}
	defer conn.Close()

	fmt.Printf("Port %d is open - %s\n", port, getServiceName(port, serviceMapping))
}

func scanIP(ip string, ports []int, wg *sync.WaitGroup, serviceMapping map[int]string) {
	defer wg.Done()

	for _, port := range ports {
		address := fmt.Sprintf("%s:%d", ip, port)
		conn, err := net.DialTimeout("tcp", address, 1*time.Second)
		if err != nil {
			continue // Port closed
		}
		defer conn.Close()

		fmt.Printf("Port %d is open on IP %s - %s\n", port, ip, getServiceName(port, serviceMapping))
	}
}

func getServiceName(port int, serviceMapping map[int]string) string {
	if service, ok := serviceMapping[port]; ok {
		return service
	}
	return "Unknown"
}

func main() {
	var target string
	flag.StringVar(&target, "target", "example.com", "Target host or IP to scan")
	var ipRange string
	flag.StringVar(&ipRange, "iprange", "", "IP range to scan (e.g., 192.168.1.1-10)")
	flag.Parse()

	if ipRange != "" {
		// IP scanning mode
		ports := []int{1, 2, 3, 4, 5} // Add the ports you want to scan here
		ipStart, ipEnd, err := parseIPRange(ipRange)
		if err != nil {
			fmt.Println("Error parsing IP range:", err)
			return
		}

		var wg sync.WaitGroup

		for ip := ipStart; ip <= ipEnd; ip++ {
			ipStr := fmt.Sprintf("%s.%d", target, ip)
			wg.Add(1)
			go scanIP(ipStr, ports, &wg, serviceMapping)
		}

		wg.Wait()
	} else {
		// Port scanning mode
		// Define the range of ports you want to scan
		startPort := 1
		endPort := 1024

		var wg sync.WaitGroup

		for port := startPort; port <= endPort; port++ {
			wg.Add(1)
			go scanPort(target, port, &wg, serviceMapping)
		}

		wg.Wait()
	}

	// Print IP address if the target is a website URL
	if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
		ips, err := net.LookupHost(target)
		if err != nil {
			fmt.Println("Error looking up IP address:", err)
		} else {
			fmt.Println("IP addresses for", target+":")
			for _, ip := range ips {
				fmt.Println(ip)
			}
		}
	}
}

func parseIPRange(ipRange string) (int, int, error) {
	parts := strings.Split(ipRange, "-")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid IP range format")
	}

	start, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}

	end, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, err
	}

	return start, end, nil
}
