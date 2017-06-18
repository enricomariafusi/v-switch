package tap

import (
	"errors"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

const (
	cIFF_TUN         = 0x0001 // Not to be used
	cIFF_TAP         = 0x0002 // This is to make the device to behave as a TAP
	cIFF_NO_PI       = 0x1000
	cIFF_MULTI_QUEUE = 0x0100 // being able to write and read at the same moment
)

type device struct {
	nr, nw string
	r, w   *os.File
}

var (
	ErrBusy        = errors.New("device is already in use")
	ErrNotReady    = errors.New("device is not ready")
	ErrExhausted   = errors.New("no devices are available")
	ErrUnsupported = errors.New("device is unsupported on this platform")
)

// Interface represents a TUN/TAP network interface
type Interface interface {
	// return name of TUN/TAP interface
	Name() string

	// implement io.Reader interface, read bytes into p from TUN/TAP interface
	Read(p []byte) (n int, err error)

	// implement io.Writer interface, write bytes from p to TUN/TAP interface
	Write(p []byte) (n int, err error)

	// implement io.Closer interface, must be called when done with TUN/TAP interface
	Close() error

	// return string representation of TUN/TAP interface
	String() string
}

func (d *device) Name() string                { return d.nr }
func (d *device) String() string              { return d.nr }
func (d *device) Close() error                { return errors.New(d.w.Close().Error() + d.r.Close().Error()) }
func (d *device) Read(p []byte) (int, error)  { return d.r.Read(p) }
func (d *device) Write(p []byte) (int, error) { return d.w.Write(p) }

func newTAP(name string) (Interface, error) {

	// creating parallel  kernel pipe for reading

	file_r, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	iface_r, err := createTuntapInterface(file_r.Fd(), name, cIFF_TAP|cIFF_NO_PI|cIFF_MULTI_QUEUE)
	if err != nil {
		return nil, err
	}

	// now open the parallel pipeline for writing

	file_w, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	iface_w, err := createTuntapInterface(file_w.Fd(), name, cIFF_TAP|cIFF_NO_PI|cIFF_MULTI_QUEUE)
	if err != nil {
		return nil, err
	}

	return &device{nr: iface_r, nw: iface_w, r: file_r, w: file_w}, nil
}

type tuntapInterface struct {
	Name  [0x10]byte
	Flags uint16
	pad   [0x28 - 0x10 - 2]byte
}

func createTuntapInterface(fd uintptr, name string, flags uint16) (string, error) {
	var req tuntapInterface
	req.Flags = flags
	copy(req.Name[:], name)

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TUNSETIFF), uintptr(unsafe.Pointer(&req)))
	if errno != 0 {
		return "", errno
	}

	return strings.Trim(string(req.Name[:]), "\x00"), nil
}
