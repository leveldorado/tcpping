package tcpping

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestParseHosts(t *testing.T) {
	Convey("Single ip", t, func() {
		ip := "145.34.34.33"
		hosts := ParseHosts(ip)
		So(len(hosts), ShouldEqual, 1)
		So(hosts[0], ShouldEqual, ip)
	})

	Convey("Single domain", t, func() {
		domain := "google.com"
		hosts := ParseHosts(domain)
		So(len(hosts), ShouldEqual, 1)
		So(hosts[0], ShouldEqual, domain)
	})

	Convey("Single CIDR", t, func() {
		cidr := "145.34.34.33/30"
		hosts := ParseHosts(cidr)
		So(len(hosts), ShouldEqual, 4)
		So(hosts[0], ShouldEqual, "145.34.34.32")
		So(hosts[1], ShouldEqual, "145.34.34.33")
		So(hosts[2], ShouldEqual, "145.34.34.34")
		So(hosts[3], ShouldEqual, "145.34.34.35")
	})

	Convey("ip, domain, cidr", t, func() {
		in := "164.34.34.5,google.com,145.34.34.33/30"
		hosts := ParseHosts(in)
		So(len(hosts), ShouldEqual, 6)
		So(hosts, ShouldContain, "145.34.34.32")
		So(hosts, ShouldContain, "145.34.34.33")
		So(hosts, ShouldContain, "145.34.34.34")
		So(hosts, ShouldContain, "145.34.34.35")
		So(hosts, ShouldContain, "164.34.34.5")
		So(hosts, ShouldContain, "google.com")
	})
}

func TestParsePorts(t *testing.T) {
	Convey("Single port", t, func() {
		ports, err := ParsePorts("34")
		So(err, ShouldBeNil)
		So(len(ports), ShouldEqual, 1)
		So(ports[0], ShouldEqual, 34)
	})

	Convey("Single range", t, func() {
		ports, err := ParsePorts("34-36")
		So(err, ShouldBeNil)
		So(len(ports), ShouldEqual, 3)
		So(ports, ShouldContain, uint16(34))
		So(ports, ShouldContain, uint16(35))
		So(ports, ShouldContain, uint16(36))
	})

	Convey("port and range", t, func() {
		ports, err := ParsePorts("28,34-36")
		So(err, ShouldBeNil)
		So(len(ports), ShouldEqual, 4)
		So(ports, ShouldContain, uint16(34))
		So(ports, ShouldContain, uint16(35))
		So(ports, ShouldContain, uint16(36))
		So(ports, ShouldContain, uint16(28))
	})
}
