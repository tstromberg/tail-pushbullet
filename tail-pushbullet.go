// tail-pushbullet lets you tail a file and make push-bullet notifications on special events.
package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"

	"github.com/hpcloud/tail"
	"github.com/xconstruct/go-pushbullet"
)

var APIKey = flag.String("api_key", "", "PushBullet API key")
var Regex = flag.String("regex", "", "Regex to notify on")
var NotifyOnStart = flag.Bool("notify_on_start", true, "Notify at startup")

func notify(p *pushbullet.Client, title string, body string) error {
	log.Printf("Notifying: title=%s body=%s", title, body)
	devs, err := p.Devices()
	if err != nil {
		return err
	}
	for _, d := range devs {
		p.PushNote(d.Iden, title, body)
		if err != nil {
			return fmt.Errorf("Could not push to %+v: %v", d, err)
		}
	}
	return nil
}

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		panic("No file specified")
	}
	if *Regex == "" {
		panic("No --regex specified")
	}
	if *APIKey == "" {
		panic("No --api_key specified")
	}
	re := regexp.MustCompile(*Regex)
	p := pushbullet.New(*APIKey)
	target := flag.Args()[0]

	if *NotifyOnStart {
		notify(p, fmt.Sprintf("tail-pushbullet now watching %s", target), fmt.Sprintf("Regex: %s", *Regex))
	}

	t, err := tail.TailFile(target, tail.Config{Follow: true})
	if err != nil {
		panic(fmt.Sprintf("TailFile: %v", err))
	}
	for line := range t.Lines {
		log.Printf("%s", line.Text)
		if re.MatchString(line.Text) {
			notify(p, fmt.Sprintf("%s matched %s", target, *Regex), line.Text)
		}
	}
}
