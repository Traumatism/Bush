package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func inc(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		if ip[i]++; ip[i] > 0 {
			break
		}
	}
}

func parseTarget(target string, ports []int) {
	if _, err := os.Stat(target); err == nil {
		file, err := os.Open(target)

		if err != nil {
			log.Fatal(err)
		}

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			parseTarget(scanner.Text(), ports)
		}

	}

	if strings.Contains(target, "/") {
		ip, ipnet, err := net.ParseCIDR(target)
		if err != nil {
			log.Fatal(err)
		}

		for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
			for _, port := range ports {
				fmt.Printf("%s:%d\n", ip, port)
			}
		}
	}
}

func parsePorts(port string) []int {
	ports := []int{}

	for _, port_range := range strings.Split(port, ",") {
		if strings.Contains(port_range, "-") {
			port_range_split := strings.Split(port_range, "-")

			start, err := strconv.Atoi(port_range_split[0])

			if err != nil {
				log.Fatal(err)
			}

			end, err := strconv.Atoi(port_range_split[1])

			if err != nil {
				log.Fatal(err)
			}

			for i := start; i <= end; i++ {
				if 0 < i && i < 65536 {
					ports = append(ports, i)
				}
			}

		} else {
			port_int, err := strconv.Atoi(port_range)

			if err != nil {
				log.Fatal(err)
			}

			if 0 < port_int && port_int < 65536 && err == nil {
				ports = append(ports, port_int)
			}
		}
	}

	return ports
}

func main() {

	ports := flag.String("ports", "25565", "Port range to use (nmap fmt)")
	cidr := flag.String("cidr", "149.202.86.0/24", "CIDR to use (can be a file too)")

	flag.Parse()

	parseTarget(*cidr, parsePorts(*ports))

}
