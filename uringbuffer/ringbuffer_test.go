package uringbuffer

import (
	"bytes"
	"log"
	"testing"
	"time"
)

func TestCreateError(t *testing.T) {
	_, err := New[int](1000)
	if err.Error() != "size must be a power of two" {
		t.FailNow()
	}
}

func TestPushBeforePull(t *testing.T) {
	r, err := New[[]byte](1024)
	if err != nil {
		t.FailNow()
	}
	defer r.Close()

	data := bytes.Repeat([]byte{0x01, 0x02, 0x03, 0x04}, 1024/4)

	r.Push(data)
	ret, ok := r.Pull()
	if !ok {
		t.FailNow()
	}
	if !bytes.Equal(data, ret) {
		t.FailNow()
	}
}

func TestPullBeforePush(t *testing.T) {
	r, err := New[[]byte](1024)
	if err != nil {
		t.FailNow()
	}
	defer r.Close()

	data := bytes.Repeat([]byte{0x01, 0x02, 0x03, 0x04}, 1024/4)

	done := make(chan struct{})
	go func() {
		defer close(done)
		ret, ok := r.Pull()
		if !ok {
			t.FailNow()
		}
		if !ok {
			t.FailNow()
		}
		if !bytes.Equal(ret, data) {
			t.FailNow()
		}
	}()

	time.Sleep(100 * time.Millisecond)

	r.Push(data)
	<-done
}

func TestClose(t *testing.T) {
	r, err := New[[]byte](1024)
	if err != nil {
		t.FailNow()
	}

	done := make(chan struct{})
	go func() {
		defer close(done)

		_, ok := r.Pull()
		if !ok {
			t.FailNow()
		}

		_, ok = r.Pull()
		if ok {
			t.FailNow()
		}
	}()

	r.Push([]byte{0x01, 0x02, 0x03, 0x04})

	r.Close()
	<-done

	r.Reset()

	r.Push([]byte{0x05, 0x06, 0x07, 0x08})

	_, ok := r.Pull()
	if !ok {
		t.FailNow()
	}
}

func BenchmarkPushPullContinuous(b *testing.B) {
	r, _ := New[[]byte](1024 * 8)
	defer r.Close()

	data := make([]byte, 1024)

	for n := 0; n < b.N; n++ {
		done := make(chan struct{})
		go func() {
			defer close(done)
			for i := 0; i < 1024*8; i++ {
				r.Push(data)
			}
		}()

		for i := 0; i < 1024*8; i++ {
			r.Pull()
		}

		<-done
	}
}

func BenchmarkPushPullPaused5(b *testing.B) {
	r, _ := New[[]byte](128)
	defer r.Close()

	data := make([]byte, 1024)

	for n := 0; n < b.N; n++ {
		done := make(chan struct{})
		go func() {
			defer close(done)
			for i := 0; i < 128; i++ {
				r.Push(data)
				time.Sleep(5 * time.Millisecond)
			}
		}()

		for i := 0; i < 128; i++ {
			r.Pull()
		}

		<-done
	}
}

func BenchmarkPushPullPaused10(b *testing.B) {
	r, _ := New[[]byte](1024 * 8)
	defer r.Close()

	data := make([]byte, 1024)

	for n := 0; n < b.N; n++ {
		done := make(chan struct{})
		go func() {
			defer close(done)
			for i := 0; i < 128; i++ {
				r.Push(data)
				time.Sleep(10 * time.Millisecond)
			}
		}()

		for i := 0; i < 128; i++ {
			r.Pull()
		}

		<-done
	}
}

func ExampleNew() {
	ringBuffer, err := New[int](1024)
	if err != nil {
		return
	}

	go func() {
		for i := 0; i < 10; i++ {
			ringBuffer.Push(i)
		}

		time.Sleep(time.Second * 5)
		ringBuffer.Close()
	}()

	for {
		v, ok := ringBuffer.Pull()
		log.Println(v, ok)
		if !ok {
			break
		}
	}

	// output:
	//
}
