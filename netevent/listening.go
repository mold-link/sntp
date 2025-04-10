// package netevent
// author: btfak.com
// create: 2013-7-20
package netevent

import (
	"net"
	"strconv"
)

type _reactor struct {
	udp_listeners  map[int]UdpClient
	tcp_clients    map[int]TcpClient
	unix_listeners map[string]UnixHandler
	udp_conn       map[int]*net.UDPConn
	tcp_listeners  map[int]*net.TCPListener
	unix_conn      map[string]*net.UnixListener
	timer          []*LaterCalling
}

func (p *_reactor) ListenUnix(addr string, unix UnixHandler) {
	if p.unix_listeners == nil {
		p.unix_listeners = make(map[string]UnixHandler)
	}
	if p.unix_conn == nil {
		p.unix_conn = make(map[string]*net.UnixListener)
	}
	p.unix_listeners[addr] = unix
	laddr, err := net.ResolveUnixAddr("unix", addr)
	if err != nil {
		return
	}
	c, erl := net.ListenUnix("unix", laddr)
	if erl != nil {
	} else {
		p.unix_conn[addr] = c
	}
}

func (p *_reactor) ListenUdp(port int, udp UdpClient) {
	laddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(port))
	if err == nil {
		p.listenUdp(laddr, udp)
	} else {
		return
	}
}

func (p *_reactor) ListenTcp(port int, tcp TcpClient) {
	p.initReactor()
	p.tcp_clients[port] = tcp
	laddr, err := net.ResolveTCPAddr("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return
	}
	c, erl := net.ListenTCP("tcp", laddr)
	if erl != nil {
	} else {
		p.tcp_listeners[port] = c
	}
	transport := new(tcpTransport)
	transport.setListener(c)
	tcp.SetTcpTransport(transport)
}

// function for inner-file usage
// helper function for the udp listening of ipv4 and ipv6
func (p *_reactor) listenUdp(addr *net.UDPAddr, udp UdpClient) {
	p.initReactor()
	p.udp_listeners[addr.Port] = udp
	c, erl := net.ListenUDP("udp", addr)
	if erl != nil {
	} else {
		p.udp_conn[addr.Port] = c
	}
	transport := new(udpTransport)
	transport.setConn(c)
	udp.SetUdpTransport(transport)
}

// init the reactor
func (p *_reactor) initReactor() {
	if p.udp_listeners == nil {
		p.udp_listeners = make(map[int]UdpClient)
	}
	if p.udp_conn == nil {
		p.udp_conn = make(map[int]*net.UDPConn)
	}
	if p.tcp_clients == nil {
		p.tcp_clients = make(map[int]TcpClient)
	}
	if p.tcp_listeners == nil {
		p.tcp_listeners = make(map[int]*net.TCPListener)
	}
}
