package portforward

import (
	"errors"
	"net"
	"strings"

	"kubevirt.io/client-go/kubecli"
	"kubevirt.io/client-go/log"
)

type portForwarder struct {
	kind, namespace, name string
	resource              portforwardableResource
}

type portforwardableResource interface {
	PortForward(name string, port int, protocol string) (kubecli.StreamInterface, error)
}

func (p *portForwarder) startForwarding(address *net.IPAddr, port forwardedPort) error {
	log.Log.Infof("forwarding %s %s:%d to %d", port.protocol, address, port.local, port.remote)
	if port.protocol == protocolUDP {
		return p.startForwardingUDP(address, port)
	}

	if port.protocol == protocolTCP {
		return p.startForwardingTCP(address, port)
	}

	return errors.New("unknown protocol: " + port.protocol)
}

func handleConnectionError(err error, port forwardedPort) {
	if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
		log.Log.Errorf("error handling connection for %d: %v", port.local, err)
	}
}
