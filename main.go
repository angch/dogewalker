package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// We're using `doge` as the command to run.
// Fastest way to do it is to install rust toolchain using rustup https://rustup.rs/
// and run a `cargo install dns-doge`

var blocked_ips = map[string]string{
	"175.139.142.25": "tmnet",
}

func main() {
	app := "doge"

	dnsmethods := []string{"-U", "-T", "-S", "-H"}
	dnsservers := []string{"1.1.1.1", "8.8.8.8", "1.9.1.9", "9.9.9.9"}
	dohservers := map[string]string{
		"1.1.1.1": "https://cloudflare-dns.com/dns-query",
		"8.8.8.8": "https://dns.google/dns-query",
		"9.9.9.9": "https://dns9.quad9.net/dns-query",
	}
	domains := []string{"google.com", "facebook.com", "twitter.com", "www.artstation.com"}

	for _, domain := range domains {
		for _, method := range dnsmethods {
			for _, server := range dnsservers {
				args := []string{"--time"}

				if method == "-H" {
					s, ok := dohservers[server]
					if !ok {
						// No DoH for this DNS server,
						continue
					}
					args = append(args, method, "@"+s)
					// args = append(args, "-n="+server, method, "@"+s)
				} else {
					args = append(args, method, "@"+server)
				}
				args = append(args, domain)

				log.Println("Running command: ", app, args)
				cmd := exec.Command(app, args...)
				stdout, err := cmd.Output()

				if err != nil {
					fmt.Println(err.Error())
					// return
				}

				// Lazy, FIXME
				isblocked := ""
				for ip, who := range blocked_ips {
					if strings.Contains(string(stdout), ip) {
						isblocked = who
						break
					}
				}

				// Print the output
				stdout = []byte(strings.TrimSuffix(string(stdout), "\n"))
				fmt.Println(string(stdout))
				if isblocked != "" {
					fmt.Println("Blocked by", isblocked)
				}
				fmt.Println()
			}

		}
	}

}
