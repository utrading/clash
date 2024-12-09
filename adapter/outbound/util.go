package outbound

import (
	"net"
	"strconv"
	"time"

	"github.com/utrading/clash/component/resolver"
	C "github.com/utrading/clash/constant"
	"github.com/utrading/clash/transport/socks5"

	"github.com/Dreamacro/protobytes"
)

func tcpKeepAlive(c net.Conn) {
	if tcp, ok := c.(*net.TCPConn); ok {
		tcp.SetKeepAlive(true)
		tcp.SetKeepAlivePeriod(30 * time.Second)
	}
}

func serializesSocksAddr(metadata *C.Metadata) []byte {
	buf := protobytes.BytesWriter{}

	addrType := metadata.AddrType()
	buf.PutUint8(uint8(addrType))

	p, _ := strconv.ParseUint(metadata.DstPort, 10, 16)
	switch addrType {
	case socks5.AtypDomainName:
		buf.PutUint8(uint8(len(metadata.Host)))
		buf.PutString(metadata.Host)
	case socks5.AtypIPv4:
		buf.PutSlice(metadata.DstIP.To4())
	case socks5.AtypIPv6:
		buf.PutSlice(metadata.DstIP.To16())
	}

	buf.PutUint16be(uint16(p))
	return buf.Bytes()
}

func resolveUDPAddr(network, address string) (*net.UDPAddr, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}

	ip, err := resolver.ResolveIP(host)
	if err != nil {
		return nil, err
	}
	return net.ResolveUDPAddr(network, net.JoinHostPort(ip.String(), port))
}

func safeConnClose(c net.Conn, err error) {
	if err != nil {
		c.Close()
	}
}
