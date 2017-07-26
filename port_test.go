package tcpping

import (
	. "github.com/smartystreets/goconvey/convey"
	"io"
	"net"
	"testing"
)

func TestPortCheck(t *testing.T) {
	localhost := "127.0.0.1"
	openL, _ := net.Listen("tcp", ":6000")
	Convey("Should be open", t, func() {
		go shouldBeOpenListenerF(openL)
		p := Port{Port: 6000}
		p.check(localhost)
		openL.Close()
		So(p.Status, ShouldEqual, PortStatusOpen)
	})

	Convey("Should be closed", t, func() {
		p := Port{Port: 6002}
		p.check(localhost)
		openL.Close()
		So(p.Status, ShouldEqual, PortStatusClosed)
	})

}

func shouldBeOpenListenerF(l net.Listener) {
	conn, _ := l.Accept()
	io.Copy(conn, conn)
	conn.Close()
}
