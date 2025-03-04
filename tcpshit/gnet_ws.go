package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync/atomic"
	"syscall"
	"time"

	_ "net/http/pprof"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
)

type wsServer struct {
	gnet.BuiltinEventEngine

	addr      string
	multicore bool
	eng       gnet.Engine
	connected int64
}

func (wss *wsServer) OnBoot(eng gnet.Engine) gnet.Action {
	wss.eng = eng
	logging.Infof("echo server with multi-core=%t is listening on %s", wss.multicore, wss.addr)
	return gnet.None
}

func (wss *wsServer) OnOpen(c gnet.Conn) ([]byte, gnet.Action) {
	c.SetContext(new(wsCodec))
	atomic.AddInt64(&wss.connected, 1)
	return nil, gnet.None
}

func (wss *wsServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	if err != nil {
		logging.Warnf("error occurred on connection=%s, %v\n", c.RemoteAddr().String(), err)
	}
	atomic.AddInt64(&wss.connected, -1)
	logging.Infof("conn[%v] disconnected", c.RemoteAddr().String())
	return gnet.None
}

func (wss *wsServer) OnTraffic(c gnet.Conn) (action gnet.Action) {
	wsc := c.Context().(*wsCodec)
	if wsc.readBufferBytes(c) == gnet.Close {
		return gnet.Close
	}
	ok, action := wsc.upgrade(c)
	if !ok {
		return
	}

	if wsc.buf.Len() <= 0 {
		return gnet.None
	}
	messages, err := wsc.Decode(c)
	if err != nil {
		log.Println(err)
		return gnet.Close
	}
	if messages == nil {
		return
	}
	// go func(conn gnet.Conn) {
	// 	ticker := time.NewTicker(5 * time.Second) // 每30秒发送一次
	// 	defer ticker.Stop()

	// 	logging.Infof("install heartbeat ticker")
	// 	for {
	// 		select {
	// 		case <-ticker.C:
	// 			// 发送心跳包
	// 			err := wsutil.WriteServerMessage(conn, ws.OpPing, []byte("PING"))
	// 			if err != nil {
	// 				fmt.Println("Error sending heartbeat:", err)
	// 				return
	// 			}
	// 			logging.Infof("Sent heartbeat")

	// 		// 处理接收到的消息
	// 		default:
	// 			resp, op, err := wsutil.ReadClientData(conn)
	// 			if err != nil {
	// 				fmt.Println("Error reading data:", err)
	// 				return
	// 			}
	// 			logging.Infof("recv after ping: %s, %d", string(resp), op)
	// 		}
	// 	}
	// }(c)

	for _, message := range messages {
		logging.Infof("recv from client[%s]  %v", c.LocalAddr().String(), string(message.Payload))
		err = wsutil.WriteServerMessage(c, message.OpCode, message.Payload)
		if err != nil {
			log.Println(err)
			return gnet.Close
		}
	}
	return gnet.None
}

func (wss *wsServer) OnTick() (delay time.Duration, action gnet.Action) {
	logging.Infof("[connected-count=%v]", atomic.LoadInt64(&wss.connected))
	return 10 * time.Second, gnet.None
}

type wsCodec struct {
	upgraded bool
	buf      bytes.Buffer
	wsMsgBuf wsMessageBuf
}

type wsMessageBuf struct {
	firstHeader *ws.Header
	curHeader   *ws.Header
	cachedBuf   bytes.Buffer
}

type readWrite struct {
	io.Reader
	io.Writer
}

func (w *wsCodec) upgrade(c gnet.Conn) (ok bool, action gnet.Action) {
	if w.upgraded {
		ok = true
		return
	}
	buf := &w.buf
	tmpReader := bytes.NewReader(buf.Bytes())
	oldLen := tmpReader.Len()

	hs, err := ws.Upgrade(readWrite{tmpReader, c})
	skipN := oldLen - tmpReader.Len()
	if err != nil {
		log.Println(err)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return
		}
		buf.Next(skipN)
		logging.Infof("conn[%v] [err=%v]", c.RemoteAddr().String(), err.Error())
		action = gnet.Close
		return
	}
	buf.Next(skipN)
	logging.Infof("conn[%v] upgrade websocket protocol! Handshake: %v", c.RemoteAddr().String(), hs)
	if err != nil {
		logging.Infof("conn[%v] [err=%v]", c.RemoteAddr().String(), err.Error())
		action = gnet.Close
		return
	}
	ok = true
	w.upgraded = true
	return
}
func (w *wsCodec) readBufferBytes(c gnet.Conn) gnet.Action {
	size := c.InboundBuffered()
	buf := make([]byte, size, size)
	read, err := c.Read(buf)
	if err != nil {
		logging.Infof("read err! %w", err)
		return gnet.Close
	}
	if read < size {
		logging.Infof("read bytes len err! size: %d read: %d", size, read)
		return gnet.Close
	}
	// logging.Infof("recv from client[%s]  %v", c.LocalAddr().String(), string(buf))
	w.buf.Write(buf)
	return gnet.None
}

func (w *wsCodec) Decode(c gnet.Conn) (outs []wsutil.Message, err error) {
	messages, err := w.readWsMessages()
	if err != nil {
		logging.Infof("Error reading message! %v", err)
		return nil, err
	}
	if messages == nil || len(messages) <= 0 {
		return
	}
	for _, message := range messages {
		if message.OpCode.IsControl() {
			err = wsutil.HandleClientControlMessage(c, message)
			if err != nil {
				log.Println(err)
				return
			}
			continue
		}
		if message.OpCode == ws.OpText || message.OpCode == ws.OpBinary {
			outs = append(outs, message)
		}
	}
	return
}

func (w *wsCodec) readWsMessages() (messages []wsutil.Message, err error) {
	msgBuf := &w.wsMsgBuf
	in := &w.buf
	for {
		if msgBuf.curHeader == nil {
			if in.Len() < ws.MinHeaderSize {
				return
			}
			var head ws.Header
			if in.Len() >= ws.MaxHeaderSize {
				head, err = ws.ReadHeader(in)
				if err != nil {
					log.Println(err)
					return messages, err
				}
			} else {
				tmpReader := bytes.NewReader(in.Bytes())
				oldLen := tmpReader.Len()
				head, err = ws.ReadHeader(tmpReader)
				skipN := oldLen - tmpReader.Len()
				if err != nil {
					log.Println(err)
					if err == io.EOF || err == io.ErrUnexpectedEOF {
						return messages, nil
					}
					in.Next(skipN)
					return nil, err
				}
				in.Next(skipN)
			}

			msgBuf.curHeader = &head
			err = ws.WriteHeader(&msgBuf.cachedBuf, head)
			if err != nil {
				log.Println(err)
				return nil, err
			}
		}
		dataLen := (int)(msgBuf.curHeader.Length)
		if dataLen > 0 {
			if in.Len() >= dataLen {
				_, err = io.CopyN(&msgBuf.cachedBuf, in, int64(dataLen))
				if err != nil {
					log.Println(err)
					return
				}
			} else {
				fmt.Println(in.Len(), dataLen)
				logging.Infof("incomplete data")
				return
			}
		}
		if msgBuf.curHeader.Fin {
			messages, err = wsutil.ReadClientMessage(&msgBuf.cachedBuf, messages)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			msgBuf.cachedBuf.Reset()
		} else {
			logging.Infof("The data is split into multiple frames")
		}
		msgBuf.curHeader = nil
	}
}

// https://community.sap.com/t5/sap-codejam-blog-posts/scaling-websockets-challenges/ba-p/13579624
func main() {
	//Increase resources limitations
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}

	// Enable pprof hooks
	go func() {
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Fatalf("pprof failed: %v", err)
		}
	}()

	var port int
	var multicore bool

	// Example command: go run main.go --port 8080 --multicore=true
	flag.IntVar(&port, "port", 7777, "server port")
	flag.BoolVar(&multicore, "multicore", true, "multicore")
	flag.Parse()

	wss := &wsServer{addr: fmt.Sprintf("tcp://localhost:%d", port), multicore: multicore}

	fmt.Println("start serving")

	// Start serving!
	log.Println("server exits:", gnet.Run(wss, wss.addr, gnet.WithMulticore(multicore), gnet.WithReusePort(true), gnet.WithTicker(true)))
}
