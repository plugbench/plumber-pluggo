package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/plugbench/nats_cli"

	"github.com/plugbench/plumber-pluggo/plumb"
)

var (
	attr      = flag.String("a", "", "set message attributes")
	src       = flag.String("s", "plumb", "set message source (default is plumb)")
	dst       = flag.String("d", "", "set message destination (default is cmd.show.data.plumb or cmd.show.data.text if -i)")
	mediaType = flag.String("t", "text/plain", "set the media type (default is text/plain)")
	wdir      = flag.String("w", "", "set message working directory (default is current directory)")
	showdata  = flag.Bool("i", false, "read data from stdin and send to cmd.show.data.text")
)

func workingDirectory() (string, error) {
	if strings.HasPrefix(*wdir, "file://") {
		return *wdir, nil
	}
	var dir string
	if *wdir == "" {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			return "", err
		}
	} else if strings.HasPrefix(*wdir, "/") {
		dir = *wdir
	} else {
		base, err := os.Getwd()
		if err != nil {
			return "", err
		}
		dir = filepath.ToSlash(filepath.Join(base, *wdir))
	}
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}
	hostname, _ := os.Hostname()
	return fmt.Sprintf("file://%s%s", hostname, dir), nil
}

func main() {
	flag.Parse()

	subject := *dst
	if subject == "" {
		if *showdata {
			subject = "cmd.show.data.text"
		} else {
			subject = "cmd.show.data.plumb"
		}
	}

	msg := nats.NewMsg(subject)
	msg.Header.Add(plumb.Source, *src)
	msg.Header.Add(plumb.ContentType, *mediaType)
	wdir, err := workingDirectory()
	if err != nil {
		log.Fatal(err)
	}
	msg.Header.Add(plumb.Base, wdir)

	attributes, err := plumb.ParseAttributes(*attr)
	if err != nil {
		log.Fatalf("parsing attributes: %v", err)
	}
	for k, v := range attributes {
		msg.Header.Add(k, v)
	}

	if *showdata {
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalf("reading stdin: %v", err)
		}
		msg.Data = bytes
	} else {
		msg.Data = []byte(strings.Join(flag.Args(), " "))
	}

	natsCfg, err := nats_cli.LoadConfigFromEnvironment()
	if err != nil {
		log.Fatalf("loading NATS config: %v", err)
	}

	nc, err := natsCfg.Connect()
	if err != nil {
		log.Fatalf("connecting to NATS: %v", err)
	}
	defer nc.Close()

	reply, err := nc.RequestMsg(msg, time.Second*10)
	if err != nil {
		log.Fatalf("sending NATS message: %v", err)
	}
	if string(reply.Data) != "ok" {
		log.Fatal(string(reply.Data))
	}
}
