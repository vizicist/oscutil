package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/hypebeast/go-osc/osc"
	"github.com/vizicist/portmidi"
)

func GetOutputStream(outputname string) (stream *portmidi.Stream, err error) {

	ndevices := portmidi.CountDevices()
	for n := 0; n < ndevices; n++ {
		devid := portmidi.DeviceID(n)
		dev := portmidi.Info(devid)
		if dev.IsOutputAvailable && dev.Name == outputname {
			log.Printf("MIDI Output %d is %s\n", devid, dev.Name)
			stream, err = portmidi.NewOutputStream(devid, 1, 0)
			if err != nil {
				return nil, fmt.Errorf("portmidi.NewOutputStream: err=%s", err)
			}
			return stream, nil
		}
	}
	return nil, fmt.Errorf("no MIDI output named %s", outputname)
}
func ListOutputs() {

	ndevices := portmidi.CountDevices()
	for n := 0; n < ndevices; n++ {
		devid := portmidi.DeviceID(n)
		dev := portmidi.Info(devid)
		if dev.IsOutputAvailable {
			log.Printf("MIDI Output %d is %s\n", devid, dev.Name)
		}
		if dev.IsInputAvailable {
			log.Printf("MIDI Input %d is %s\n", devid, dev.Name)
		}
	}
}

/*
	status := 0xb0 | (s.channel - 1)
	e := portmidi.Event{
		Timestamp: portmidi.Time(),
		Status:    int64(status),
		Data1:     int64(0x7b),
		Data2:     int64(0x00),
	}
	oscsendEvent(s, []portmidi.Event{e})
*/

// SendNote sends MIDI output for a Note
/*
	s, err := m.getSoundOutput(n.Sound)
	if err != nil {
		log.Printf("OscmidiDevice.SendNote error: %s\n", err)
		return
	}
	var status uint8
	switch n.TypeOf {
	case NOTEON:
		status = 0x90
	case NOTEOFF:
		status = 0x80
	default:
		log.Printf("SendNote can't YET handle Note TypeOf=%v\n", n.TypeOf)
		return
	}
	// NOTE: s.channel is 1-based, but MIDI output is 0-based
	status |= (s.channel - 1)
	e := portmidi.Event{
		Timestamp: portmidi.Time(),
		Status:    int64(status),
		Data1:     int64(n.Pitch),
		Data2:     int64(n.Velocity),
	}
	if debug {
		log.Printf("MIDI.SendNote status=0x%0x pitch=%d velocity=%d\n", status, n.Pitch, n.Velocity)
	}
}
*/

func usage() {
	log.Printf("Usage: osc2midi [-l] [-p oscport] [-o midioutput]")
}

func main() {
	log.Printf("Main start\n")

	list := flag.Bool("list", false, "List MIDI I/O")
	verbose := flag.Bool("verbose", false, "Verbose mode")
	port := flag.Int("port", 0, "OSC port")
	output := flag.String("output", "", "MIDI Output Name")

	flag.Parse()

	portmidi.Initialize()

	if *list {
		ListOutputs()
		return
	}
	if *output == "" {
		usage()
		return
	}
	stream, err := GetOutputStream(*output)
	if err != nil {
		log.Printf("GetOutputStream: err=%s", err)
		return
	}
	if *verbose {
		log.Printf("stream=%+v\n", stream)
	}
	client := osc.NewClient("127.0.0.1", *port)
	if *verbose {
		log.Printf("client=%+v\n", client)
	}

	d := osc.NewStandardDispatcher()

	err = d.AddMsgHandler("*", func(msg *osc.Message) {
		log.Printf("Got OSC! msg=%+v\n", msg)
	})
	if err != nil {
		log.Printf("ERROR! %s\n", err.Error())
		return
	}

	addr := fmt.Sprintf("127.0.0.1:%d", port)
	server := &osc.Server{
		Addr:       addr,
		Dispatcher: d,
	}
	if *verbose {
		log.Printf("Now listening for OSC on port %d\n", port)
	}
	go server.ListenAndServe()
	log.Printf("Blocking forever\n")
	select {} // block forever
}
