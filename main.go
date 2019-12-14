package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/anaskhan96/soup"
)

var nodesURL = "https://www.youneed.win/free-ss"
var proxy = "http://127.0.0.1:1087"

type ssnode struct {
	Addr string
	Port int
	Pass string
	Cryp string
	Time int
	Tag  string
}

func getNodes(nodesURL string) []ssnode {
	proxyURL, err := url.Parse(proxy)
	proxyClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
	resp, err := soup.GetWithClient(nodesURL, proxyClient)
	check(err)
	doc := soup.HTMLParse(resp)
	links := doc.Find("section", "class", "context").FindAll("td")
	items := []string{}
	for _, link := range links {
		if link.Text() != "" {
			items = append(items, link.Text())
		}
	}
	nodeNum := len(items) / 4
	nodes := make([]ssnode, nodeNum)
	for i := range nodes {
		nodes[i].Addr = items[0]
		nodes[i].Port, err = strconv.Atoi(items[1])
		nodes[i].Pass = items[2]
		nodes[i].Cryp = items[3]
		items = items[4:]
	}
	return nodes
}

func testNodeSpeed(nodes []ssnode) []ssnode {
	for i := range nodes {
		cmdStr := "ping -c 2 " + nodes[i].Addr
		args := strings.Split(cmdStr, " ")
		cmd := exec.Command(args[0], args[1:]...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			nodes[i].Time = 10000
		} else {
			timeReg := regexp.MustCompile(`stddev = .*?\/(\d+).\d+\/.*`)
			nodes[i].Time, err = strconv.Atoi(string(timeReg.FindSubmatch(out)[1]))
		}
	}
	return nodes
}

func sortNodes(nodes []ssnode) []ssnode {
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Time < nodes[j].Time
	})
	return nodes
}

func generateURLs(nodes []ssnode) []string {
	URLs := []string{}
	for i := range nodes {
		nodes[i].Tag = "Node " + strconv.Itoa(i)
		userinfo := nodes[i].Cryp + ":" + nodes[i].Pass
		encoded := base64.StdEncoding.EncodeToString([]byte(userinfo))
		tagLen := len(nodes[i].Tag)
		utf8Tag := ""
		for j := 0; j < tagLen; j++ {
			utf8Tag += fmt.Sprintf("%%%2X", nodes[i].Tag[j])
		}
		ssURI := "ss://" + encoded + "@" + nodes[i].Addr + ":" + strconv.Itoa(nodes[i].Port) + "/?#" + utf8Tag
		URLs = append(URLs, ssURI)
	}
	return URLs
}

func write(data []string) {
	home, err := os.UserHomeDir()
	check(err)
	dir := filepath.Join(home, "Desktop", "nodes.txt")
	f, err := os.Create(dir)
	check(err)
	defer f.Close()
	for i := range data {
		_, err := f.WriteString(data[i] + "\n")
		check(err)
	}
	return
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	nodes := getNodes(nodesURL)
	testedNodes := testNodeSpeed(nodes)
	sortedNodes := sortNodes(testedNodes)
	URLs := generateURLs(sortedNodes)
	write(URLs)
}
