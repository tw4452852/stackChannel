package stackChannel

import (
	"bytes"
	"io"
	"log"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

type stackChannel struct {
	channel    chan []byte   //just one slice element inchannel
	readBuffer *bytes.Buffer //for save extra read elements
	lastError  error
}

func newStackChannel() *stackChannel {
	c := &stackChannel{
		channel:    make(chan []byte, 1),
		readBuffer: new(bytes.Buffer),
	}
	go c.readLoop()
	return c
}

//satisfy io.Writer
func (c *stackChannel) Write(p []byte) (int, error) {
	toChannel := p
	for {
		select {
		case c.channel <- toChannel:
			return len(p), nil
		case old := <-c.channel:
			toChannel = append(old, toChannel...)
		}
	}
	panic("not reachable!")
}

//satisfy io.Reader
func (c *stackChannel) Read(p []byte) (int, error) {
	//read until enough!
	have, need := 0, len(p)
	for {
		c, e := c.readBuffer.Read(p)
		//some error except EOF happened
		if e != nil && e != io.EOF {
			return have, e
		}
		have += c
		p = p[c:]
		//finished?
		if have >= need {
			return have, nil
		}
		//need continue
		runtime.Gosched()
	}
	panic("not reachable")
}

//satisfy io.Closer
func (c *stackChannel) Close() error {
	close(c.channel)
	return c.lastError
}

//forward the stuff in channel to readBuffer
func (c *stackChannel) readLoop() error {
	//log the recent error
	for {
		select {
		case p, open := <-c.channel:
			if open != true {
				return c.lastError
			}
			if _, err := c.readBuffer.Write(p); err != nil {
				log.Println("write to buffer failed:", err)
				c.lastError = err
			}
		default:
			//TODO:if buffer is too large, flush sth to disk?
			//nothing to do
		}
	}
	panic("not reachable!")
}
