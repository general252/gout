package ucontainer

import (
	"io"
	"sync"
)

type PacketIOBuffer struct {
	mutex sync.Mutex
	p     *Pool
	q     *Queue

	notify chan struct{}
	subs   bool
	closed bool
}

func NewPacketBuffer() *PacketIOBuffer {
	return &PacketIOBuffer{
		notify: make(chan struct{}),
		p:      NewPoolWithSize(5 * 1024 * 1024), // 5MB
		q:      NewQueue(),
	}
}

func (c *PacketIOBuffer) Write(packet []byte) (int, error) {
	if packet == nil || len(packet) == 0 {
		return 0, io.ErrShortWrite
	}

	c.mutex.Lock()

	if c.closed {
		c.mutex.Unlock()
		return 0, io.ErrClosedPipe
	}

	var notify chan struct{}

	if c.subs {
		// readers are waiting.  Prepare to notify, but only
		// actually do it after we release the lock.
		notify = c.notify
		c.notify = make(chan struct{})
		c.subs = false
	}

	buf, err := c.p.Get(len(packet))
	if err != nil {
		return 0, err
	}

	n := copy(buf, packet)

	c.q.Push(buf)

	c.mutex.Unlock()

	if notify != nil {
		close(notify)
	}

	return n, nil
}

func (c *PacketIOBuffer) Read() ([]byte, error) {
	for {
		c.mutex.Lock()

		if c.q.Len() > 0 {
			v, _ := c.q.Pop()

			c.mutex.Unlock()

			pkt, ok := v.([]byte)
			if !ok || pkt == nil {
				return nil, io.ErrUnexpectedEOF
			}

			return pkt, nil
		}

		if c.closed {
			c.mutex.Unlock()
			return nil, io.EOF
		}

		notify := c.notify
		c.subs = true
		c.mutex.Unlock()

		select {
		case <-notify:
		}
	}
}

func (c *PacketIOBuffer) Close() error {
	c.mutex.Lock()

	if c.closed {
		c.mutex.Unlock()
		return nil
	}

	notify := c.notify
	c.closed = true

	c.mutex.Unlock()

	close(notify)

	return nil
}

// Count returns the number of packets in the buffer.
func (c *PacketIOBuffer) Count() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.q.Len()
}
