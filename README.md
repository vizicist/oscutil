# osc2midi

A simple OSC server to send MIDI messages.  It listens on an OSC port for messages of the form:

    /midi {status} {data1} {data2}

and sends a MIDI message to a MIDI port.
    
The {data1} and/or {data2} can be omitted if the status byte
is for a MIDI message of only 1 or 2 bytes.

Options of osc2midi:
<pre>
  -list
        list MIDI I/O
  -midiport string
        MIDI port Name
  -oscport int
        OSC port
  -verbose
        verbose mode
</pre>
        
Example use:

    osc2midi -verbose -oscport=2222 "-midiport=01. Internal MIDI"
