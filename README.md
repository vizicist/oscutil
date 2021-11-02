# osc2midi

A simple OSC server to send MIDI messages.  It listens on an OSC port for messages of the form:

    /midi {status} {data1} {data2}
    
The {data1} and/or {data2} can be omitted if the status byte is for MIDI messages of only 1 or 2 bytes.

Options of osc2midi:
<pre>
  -list
        list MIDI I/O
  -output string
        MIDI Output Name
  -port int
        OSC port
  -verbose
        verbose mode
</pre>
        
Example use:

    osc2midi -verbose -port=2222 "-output=01. Internal MIDI"
