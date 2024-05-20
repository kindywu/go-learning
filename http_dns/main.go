package main

import (
	"context"
	"fmt"
	"net/url"

	"github.com/miekg/dns"
	"github.com/mosajjal/dnsclient"
)

func main() {
	// question := "example.com."
	// question := "google.com."
	// question := "qq.com."
	question := "roblox.com."
	msg := dns.Msg{}
	msg.RecursionDesired = true
	msg.SetQuestion(question, dns.TypeA)

	u := "https://cloudflare-dns.com/dns-query"
	uri, _ := url.Parse(u)
	c, err := dnsclient.NewDoHClient(*uri, true, "")

	if err != nil {
		println(err)
	}
	defer c.Close()
	reply, _, err := c.Query(context.Background(), &msg)
	if err != nil {
		println(err)
	}
	for _, r := range reply {
		fmt.Printf("%v\n", r)
	}
}
