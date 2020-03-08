package main

import (
	"bufio"
	"github.com/alexflint/go-arg"
	"github.com/kangoo13/proxy-checker/core"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	var inputArgs = InputArgs{}

	p := arg.MustParse(&inputArgs)

	if inputArgs.CloudFlareBypass && inputArgs.WebsiteCloudFlare == "" {
		p.Fail("If you wish to test bypass cloudflare, please provide website to test against protected by cloudflare with -w")
	}

	outputFile := OpenFile(inputArgs.GoodProxiesPath)
	defer outputFile.Close()
	workingProxies := core.CheckProxies(&core.ProxyChecker{
		ProxyList:    FileToStringSlice(inputArgs.ProxyList),
		TimeoutProxy: time.Duration(inputArgs.TimeoutProxy) * time.Second,
		AvailableCheckers: core.AvailableCheckers{
			CloudFlareBypass: inputArgs.CloudFlareBypass,
			Basic:            true,
		},
		WebsiteCloudFlare: inputArgs.WebsiteCloudFlare,
	})

	datawriter := bufio.NewWriter(outputFile)

	defer datawriter.Flush()
	for _, data := range workingProxies {
		_, _ = datawriter.WriteString(data + "\n")
	}
}

func OpenFile(filePath string) *os.File {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend|0660)
	if err != nil {
		log.Fatalf("[OpenFile] error while opening file %s", err)
	}

	return f
}

func FileToStringSlice(filePath string) []string {
	var (
		fileTextLines []string
	)
	readFile, err := os.Open(filePath)

	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}

	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		newLine := fileScanner.Text()
		if strings.Trim(newLine, "\n\t\r") != "" {
			fileTextLines = append(fileTextLines, fileScanner.Text())
		}
	}

	return fileTextLines
}

type InputArgs struct {
	ProxyList       string `arg:"positional,required" help:"path to the proxyList"`
	GoodProxiesPath string `arg:"positional" default:"good_proxies.txt" help:"path to the good proxies identified"`
	TimeoutProxy    int64  `arg:"-t" default:"4" help:"timeout proxy duration in seconds"`
	availableCheckers
	WebsiteCloudFlare string `arg:"-w" help:"website to test proxy against cloudflare, needed if -cf"`
}

type availableCheckers struct {
	CloudFlareBypass bool `arg:"-c" help:"if activated, check if proxy bypass cloudflare"`
	Basic            bool `arg:"-b" default:"true" help:"if activated, simply checks if proxy is working"`
}
