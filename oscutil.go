package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/hypebeast/go-osc/osc"
	"github.com/vizicist/portmidi"
)

var Verbose bool
var Midiout *portmidi.Stream

// oscutil listen {oscport}
// oscutil send {oscport} {/addr} {args} ...
// oscutil midi {oscport} {midiport}

func main() {

	if len(os.Args) < 2 {
		fmt.Printf("Usage:\n" +
			"  oscutil listen {oscport}\n" +
			"  oscutil send {oscport} {/addr} {args...}\n" +
			"  listmidi\n" +
			"  servemidi {oscport} {midiport}\n")
		return
	}
	switch os.Args[1] {
	case "listen":
		doListen(os.Args[2:])
	case "send":
		doSend(os.Args[2:])
	case "servemidi":
		doServeMidi(os.Args[2:])
	case "listmidi":
		doListMidi()
	}
}

func doSend(args []string) {
	if len(args) < 2 {
		fmt.Printf("Missing port and/or message on send\n")
		return
	}
	oscport, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("Bad osc port value: %s\n", err)
		return
	}
	client := osc.NewClient("127.0.0.1", oscport)
	if args[1][0] != '/' {
		fmt.Printf("OSC message must start with /\n")
		return
	}
	msg := osc.NewMessage(args[1])
	for n := 2; n < len(args); n++ {
		// we try to deduce the type of the argument
		num, err := strconv.Atoi(args[n])
		if err != nil {
			// Conversion to an integer didn't work, try float
			flt, err := strconv.ParseFloat(args[n], 32)
			if err != nil {
				msg.Append(args[n])
			} else {
				msg.Append(float32(flt))
			}
		} else {
			msg.Append(int32(num))
		}
	}
	client.Send(msg)
}

func doListMidi() {
	portmidi.Initialize()
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

func doServeMidi(args []string) {

	portmidi.Initialize()

	if len(args) != 2 {
		fmt.Printf("Missing OSC port and/or MIDI port\n")
		return
	}
	oscport, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("Bad osc port value: %s\n", err)
		return
	}

	midiport := args[1]

	Midiout, err = GetOutputStream(midiport)
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

	addr := fmt.Sprintf("127.0.0.1:%d", oscport)
	server := &osc.Server{
		Addr:       addr,
		Dispatcher: d,
	}
	if Verbose {
		fmt.Printf("Now listening for OSC on port %d\n", oscport)
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

func handleOSC(msg *osc.Message) {
	switch msg.Address {
	default:
		fmt.Printf("Unrecognized OSC message: %s\n", msg)
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

func doListen(args []string) {

	if len(args) == 0 {
		fmt.Printf("Missing osc port number\n")
		return
	}
	oscport, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("Bad osc port value: %s\n", err)
		return
	}

	d := osc.NewStandardDispatcher()
	err = d.AddMsgHandler("*", func(msg *osc.Message) {
		fmt.Printf("time: %d  msg: %s", time.Now().UnixMilli(), msg.Address)
		for n := 0; n < len(msg.Arguments); n++ {
			fmt.Printf(" %v", msg.Arguments[n])
		}
		fmt.Printf("\n")
	})
	if err != nil {
		fmt.Printf("AddMsgHandler: err=%s\n", err)
		return
	}

	addr := fmt.Sprintf("127.0.0.1:%d", oscport)
	server := &osc.Server{
		Addr:       addr,
		Dispatcher: d,
	}
	if Verbose {
		fmt.Printf("Now listening for OSC on port %d\n", oscport)
	}
	startOSC(server) // never returns
}
