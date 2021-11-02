# oscutil

A utility for sending and listening for OSC messages.  Usage:

    oscutil listen {oscport}
    oscutil send {oscport} {/addr} {args...}
    
This utility also includes a server that you can use to send MIDI messages.
For that, use:
 
     oscutil servemidi {oscport} {midiport}

For exampe:

    oscutil servemidi 2222 "01. Internal MIDI"

will listen on OSC port 2222 for messages of the form:

    /midi {status} {data1} {data2}

and sends a MIDI message to a MIDI port.
    
The {data1} and/or {data2} can be omitted if the status byte
is for a MIDI message of only 1 or 2 bytes.

