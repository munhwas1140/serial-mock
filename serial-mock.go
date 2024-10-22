package serial_mock

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

type MockSerialPortMap map[string]*MockSerialPort

type MockSerialPort struct {
	ptys    []*os.File
	name    string
	handler func() error
}

func New() *MockSerialPortMap {
	return &MockSerialPortMap{}
}

func Add(name string, count int, fn func() error) {

}

func (mm *MockSerialPortMap) Run() {

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
