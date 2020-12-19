// +build !darwin,!dragonfly,!freebsd,!netbsd,!openbsd,!linux,!windows

package gateway

func (inte *Interface) discoverGatewayOS() error {
	return errNotImplemented
}

func (inte *Interface) discoverGatewayInterfaceOS() error {
	return errNotImplemented
}
