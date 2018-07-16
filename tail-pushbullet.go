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

func notify(p *pushbullet.Client, s string) error {
	log.Printf("Notifying: %s", s)
	devs, err := p.Devices()
	if err != nil {
		return err
	}
	for _, d := range devs {
		log.Printf("Pushing to %s: %s", d.Iden, s)
		p.PushNote(d.Iden, "tail match", s)
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

	t, err := tail.TailFile(flag.Args()[0], tail.Config{Follow: true})
	if err != nil {
		panic(fmt.Sprintf("TailFile: %v", err))
	}
	for line := range t.Lines {
		log.Printf("%s", line.Text)
		if re.MatchString(line.Text) {
			notify(p, line.Text)
		}
	}
}
