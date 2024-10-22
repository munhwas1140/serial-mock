package serial_mock

import (
	"fmt"
	"os"
	"syscall"
	"time"
	"unsafe"
)

type MockSerialPorts map[string]*MockSerialPort

func New() MockSerialPorts {
	return MockSerialPorts{}
}

func (ms MockSerialPorts) Add(name string, count int, handler func(*MockSerialPort) error) error {
	for i := 1; i <= count; i++ {
		mock := &MockSerialPort{
			name:    name,
			handler: handler,
		}

		ptm, err := os.OpenFile("/dev/ptmx", os.O_RDWR, 0)
		if err != nil {
			return err
		}

		if err := ptsUnlock(ptm); err != nil {
			return err
		}

		sname, err := ptsName(ptm)
		if err != nil {
			return err
		}

		time.Sleep(50 * time.Millisecond)
		mock.pts, err = os.OpenFile(sname, os.O_RDWR|syscall.O_NOCTTY, 0)
		if err != nil {
			return err
		}
		mock.ptmFd = os.NewFile(ptm.Fd(), fmt.Sprintf("name_%d", i))

		ms[sname] = mock
	}
	return nil
}

func ptsName(f *os.File) (string, error) {
	var n uintptr
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), syscall.TIOCGPTN, uintptr(unsafe.Pointer(&n)))
	if err != 0 {
		return "", err
	}
	return fmt.Sprintf("/dev/pts/%d", n), nil
}

func ptsUnlock(f *os.File) error {
	var u uintptr
	// use TIOCSPTLCK with a zero valued arg to clear the slave pty lock
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), syscall.TIOCSPTLCK, uintptr(unsafe.Pointer(&u)))
	if err != 0 {
		return err
	}
	return nil
}

type MockSerialPort struct {
	ptmFd   *os.File
	pts     *os.File
	name    string
	handler func(*MockSerialPort) error
}
