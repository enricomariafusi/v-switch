package tap

import (
	"fmt"
	"log"
	"syscall"
	"unsafe"
)

// Here we use syscalls, which are written in c. So that, this is go c-style

type TapConn struct {
	fd     int
	ifname string
}

func (tap_conn *TapConn) Open(mtu uint, name string) (err error) {

	tap_conn.ifname = name
	log.Printf("[TAP][SYS] Interface name is <%s>", tap_conn.ifname)

	// Open the tap/tun device

	tap_conn.fd, err = syscall.Open("/dev/net/tun", syscall.O_RDWR, syscall.S_IRUSR|syscall.S_IWUSR|syscall.S_IRGRP|syscall.S_IROTH)
	if err != nil {
		return fmt.Errorf("Error opening device /dev/net/tun: %s", err)
	}

	// Prepare a struct ifreq structure for TUNSETIFF with tap settings
	// IFF_TAP: tap device, IFF_NO_PI: no extra packet information
	ifr_flags := uint16(syscall.IFF_TAP | syscall.IFF_NO_PI)
	// FIXME: Assumes little endian
	ifr_struct := make([]byte, 32)
	ifr_struct[16] = byte(ifr_flags)
	ifr_struct[17] = byte(ifr_flags >> 8)
	copy(ifr_struct[0:15], tap_conn.ifname)
	r0, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(tap_conn.fd), syscall.TUNSETIFF, uintptr(unsafe.Pointer(&ifr_struct[0])))
	if r0 != 0 {
		tap_conn.Close()
		return fmt.Errorf("Error setting tun/tap type: %s", err)
	} else {
		log.Printf("[TAP][SYS] TAP device %s created", tap_conn.ifname)
	}

	// Create a raw socket for our tap interface, so we can set the MTU
	tap_sockfd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, syscall.ETH_P_ALL)
	if err != nil {
		tap_conn.Close()
		return fmt.Errorf("Error creating packet socket: %s", err)
	}
	// We won't need the socket after we've set the MTU and brought the interface up
	defer syscall.Close(tap_sockfd)

	// Bind the raw socket to our tap interface
	err = syscall.BindToDevice(tap_sockfd, tap_conn.ifname)
	if err != nil {
		tap_conn.Close()
		return fmt.Errorf("Error binding packet socket to tap interface: %s", err)
	}

	// Prepare a ifreq structure for SIOCSIFMTU with MTU setting
	ifr_mtu := mtu
	// FIXME: Assumes little endian
	ifr_struct[16] = byte(ifr_mtu)
	ifr_struct[17] = byte(ifr_mtu >> 8)
	ifr_struct[18] = byte(ifr_mtu >> 16)
	ifr_struct[19] = byte(ifr_mtu >> 24)
	r0, _, err = syscall.Syscall(syscall.SYS_IOCTL, uintptr(tap_sockfd), syscall.SIOCSIFMTU, uintptr(unsafe.Pointer(&ifr_struct[0])))
	if r0 != 0 {
		tap_conn.Close()
		return fmt.Errorf("Error setting MTU on tap interface: %s", err)
	}

	// Get the current interface flags in ifr_struct
	r0, _, err = syscall.Syscall(syscall.SYS_IOCTL, uintptr(tap_sockfd), syscall.SIOCGIFFLAGS, uintptr(unsafe.Pointer(&ifr_struct[0])))
	if r0 != 0 {
		tap_conn.Close()
		return fmt.Errorf("Error getting tap interface flags: %s", err)
	}
	// Update the interface flags to bring the interface up
	// FIXME: Assumes little endian
	ifr_flags = uint16(ifr_struct[16]) | (uint16(ifr_struct[17]) << 8)
	ifr_flags |= syscall.IFF_UP | syscall.IFF_RUNNING
	ifr_struct[16] = byte(ifr_flags)
	ifr_struct[17] = byte(ifr_flags >> 8)
	r0, _, err = syscall.Syscall(syscall.SYS_IOCTL, uintptr(tap_sockfd), syscall.SIOCSIFFLAGS, uintptr(unsafe.Pointer(&ifr_struct[0])))
	if r0 != 0 {
		tap_conn.Close()
		return fmt.Errorf("Error bringing up tap interface: %s", err)
	}

	return nil
}

func (tap_conn *TapConn) Close() {
	syscall.Close(tap_conn.fd)
}

func (tap_conn *TapConn) Read(b []byte) (n int, err error) {
	return syscall.Read(tap_conn.fd, b)
}

func (tap_conn *TapConn) Write(b []byte) (n int, err error) {
	return syscall.Write(tap_conn.fd, b)
}
