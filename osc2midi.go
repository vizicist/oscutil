package main

import (
	"flag"
	"fmt"

	"github.com/hypebeast/go-osc/osc"
	"github.com/vizicist/portmidi"
)

var Verbose bool
var Midiout *portmidi.Stream

func main() {

	plist := flag.Bool("list", false, "list MIDI I/O")
	pverbose := flag.Bool("verbose", false, "verbose mode")
	pport := flag.Int("port", 0, "OSC port")
	poutput := flag.String("output", "", "MIDI Output Name")

	flag.Parse()

	Verbose = *pverbose
	portmidi.Initialize()

	if *plist {
		ListOutputs()
		return
	}
	if *poutput == "" {
		flag.Usage()
		return
	}
	var err error
	Midiout, err = GetOutputStream(*poutput)
	if err != nil {
		fmt.Printf("GetOutputStream: err=%s", err)
		return
	}

	d := osc.NewStandardDispatcher()
	err = d.AddMsgHandler("*", handleOSC)
	if err != nil {
		fmt.Printf("AddMsgHandler: err=%s\n", err)
		return
	}

	addr := fmt.Sprintf("127.0.0.1:%d", *pport)
	server := &osc.Server{
		Addr:       addr,
		Dispatcher: d,
	}
	if Verbose {
		fmt.Printf("Now listening for OSC on port %d\n", *pport)
	}
	startOSC(server) // never returns
}

func GetOutputStream(outputname string) (stream *portmidi.Stream, err error) {

	ndevices := portmidi.CountDevices()
	for n := 0; n < ndevices; n++ {
		devid := portmidi.DeviceID(n)
		dev := portmidi.Info(devid)
		if dev.IsOutputAvailable && dev.Name == outputname {
			stream, err = portmidi.NewOutputStream(devid, 1, 0)
			if err != nil {
				return nil, fmt.Errorf("portmidi.NewOutputStream: err=%s", err)
			}
			if Verbose {
				fmt.Printf("Opened MIDI output: %s\n", dev.Name)
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
			fmt.Printf("MIDI Output %d is %s\n", devid, dev.Name)
		}
		if dev.IsInputAvailable {
			fmt.Printf("MIDI Input %d is %s\n", devid, dev.Name)
		}
	}
}

func handleOSC(msg *osc.Message) {
	if Verbose {
		fmt.Printf("handleOSC: OSC message = %s\n", msg)
	}
	switch msg.Address {
	case "/midi":
		tags, _ := msg.TypeTags()
		_ = tags
		nargs := msg.CountArguments()

		switch {
		case nargs == 0:
			fmt.Printf("OSC /midi message: no arguments?\n")
			return
		case nargs > 3:
			fmt.Printf("OSC /midi message: too many arguments?\n")
			return
		}

		var b []int = []int{0, 0, 0}
		var err error
		for n := 0; n < nargs; n++ {
			b[n], err = argAsInt(msg, n)
			if err != nil {
				fmt.Printf("OSC /midi message: err=%s\n", err)
				return
			}
		}
		if Verbose {
			fmt.Printf("handleOSC: sending MIDI bytes %d %d %d\n", b[0], b[1], b[2])
		}
		Midiout.WriteShort(int64(b[0]), int64(b[1]), int64(b[2]))
	}
}

func startOSC(server *osc.Server) {
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("ListenAndServer: err=%s\n", err)
		return
	}
}

func argAsInt(msg *osc.Message, index int) (i int, err error) {
	arg := msg.Arguments[index]
	switch v := arg.(type) {
	case int32:
		i = int(v)
	case int64:
		i = int(v)
	default:
		err = fmt.Errorf("expected an int in OSC argument index=%d", index)
	}
	return i, err
}
