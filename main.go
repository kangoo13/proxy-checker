package main

import (
	"bufio"
	"github.com/alexflint/go-arg"
	"github.com/kangoo13/cloudflare-proxy-ban-checker/cfprxchecker"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	var args cfprxchecker.Args
	var inputArgs struct {
		Website         string `arg:"positional,required" help:"website to test proxy against cloudflare"`
		ProxyList       string `arg:"positional,required" help:"path to the proxyList"`
		GoodProxiesPath string `default:"good.txt" help:"path to the good proxies identified"`
		BadProxiesPath  string `default:"bad.txt" help:"path to the bad proxies identified"`
		TimeoutProxy    int64  `default:"5" help:"timeout proxy duration"`
	}

	arg.MustParse(&inputArgs)

	if inputArgs.GoodProxiesPath != "" {
		args.GoodProxiesOutputFile = OpenFile(inputArgs.GoodProxiesPath)
		defer args.GoodProxiesOutputFile.Close()
	}

	if inputArgs.BadProxiesPath != "" {
		args.BadProxiesOutputFile = OpenFile(inputArgs.BadProxiesPath)
		defer args.BadProxiesOutputFile.Close()
	}

	args.WebsiteToCrawl = inputArgs.Website
	args.ProxyList = FileToStringSlice(inputArgs.ProxyList)
	args.TimeoutProxy = time.Duration(inputArgs.TimeoutProxy) * time.Second

	cfprxchecker.CheckProxiesAgainstCloudFlare(&args)
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
