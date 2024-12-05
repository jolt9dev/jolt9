package ps

import (
	"bufio"
	"fmt"
	o "os"
	"runtime"

	"github.com/jolt9dev/jolt9/pkg/os"
)

const (
	ARCH     = runtime.GOARCH
	PLATFORM = runtime.GOOS
)

var (
	Args     = o.Args[1:]
	Stderr   = o.Stderr
	Stdout   = o.Stdout
	Stdin    = o.Stdin
	ExecPath = o.Args[0]
	history  = []string{}
	reader   *bufio.Reader
	writer   *bufio.Writer
	eol      = []byte(os.EOL)
)

func init() {
	current, _ := Cwd()
	history = append(history, current)
}

func Cwd() (string, error) {
	return o.Getwd()
}

func Exit(code int) {
	o.Exit(code)
}

func Pid() int {
	return o.Getpid()
}

func Ppid() int {
	ppid := o.Getppid()
	if ppid == 0 {
		return -1
	}

	return ppid
}

func Uid() int {
	return o.Getuid()
}

func Gid() int {
	return o.Getgid()
}

func Euid() int {
	return o.Geteuid()
}

func Egid() int {
	return o.Getegid()
}

func Pushd(path string) error {
	history = append(history, path)
	return o.Chdir(path)
}

func Popd() error {
	if len(history) == 1 {
		return nil
	}

	last := history[len(history)-1]
	history = history[:len(history)-1]
	return o.Chdir(last)
}

func Read(b []byte) (int, error) {
	if reader == nil {
		reader = bufio.NewReader(o.Stdin)
	}

	return reader.Read(b)
}

func ReadLine() (string, error) {
	if reader == nil {
		reader = bufio.NewReader(o.Stdin)
	}

	b, _, e := reader.ReadLine()
	return string(b), e
}

func WriteBytes(b []byte) (int, error) {
	return Stdout.Write(b)
}

func WriteRune(r rune) (int, error) {
	if writer == nil {
		writer = bufio.NewWriter(o.Stdout)
	}

	return writer.WriteRune(r)
}

func WriteString(s string) (int, error) {
	if writer == nil {
		writer = bufio.NewWriter(Stdout)
	}

	b, err := writer.WriteString(s)
	if err != nil {
		return b, err
	}
	writer.Flush()
	return b, err
}

func Writef(format string, a ...interface{}) (int, error) {
	msg := fmt.Sprintf(format, a...)
	return WriteString(msg)
}

func Writeln(s string) (int, error) {

	n, err := WriteString(s)
	if err != nil {
		return n, err
	}

	n2, err := writer.Write(eol)
	return n + n2, err
}
