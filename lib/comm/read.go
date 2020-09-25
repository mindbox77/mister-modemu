package comm

import (
	"time"

	"github.com/jorgefuertes/mister-modemu/lib/console"
)

// ReadLoop - Read loop
func ReadLoop() {
	console.Info("CONN/RX", "Listening…")
	cBuf := make([]byte, 1, 1024)
	cBuf = []byte{}
	rBuf := make([]byte, 1, 1)

	// read loop
	for {
		n, err := s.Read(rBuf)
		if err != nil {
			if err.Error() != "EOF" {
				console.Error("CONN/RX", err.Error())
			}
			time.Sleep(100 * time.Millisecond)
			continue
		}
		if rBuf[0] == bs || rBuf[0] == del {
			if len(cBuf) > 0 {
				cBuf = cBuf[:len(cBuf)-1]
			}
			if status.echo {
				writeByte(&[]byte{bs})
			}
			continue
		}
		if status.echo {
			writeByte(&rBuf)
			if rBuf[0] == cr {
				writeByte(&[]byte{lf})
			}
		}
		for i := 0; i < n; i++ {
			cBuf = append(cBuf, rBuf[i])
			if len(cBuf) > 2048 {
				console.Error("CONN/RX", "Command buffer limit reached")
				write(er)
				cBuf = []byte{}
				break
			}
		}
		if cBuf[len(cBuf)-1] == lf || cBuf[len(cBuf)-1] == cr {
			if string(cBuf[0:2]) == "AT" {
				cmd := bufToStr(&cBuf)
				res := parseCmd(cmd)
				console.Debug("CONN/REPLY", res)
				write(res + "\r\n")
			}
			cBuf = []byte{}
		}
	}
}
