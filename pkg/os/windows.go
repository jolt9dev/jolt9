//go:build windows

package os

const (
	FAMILY   = "windows"
	POSIX    = false
	EOL      = "\r\n"
	PATH_SEP = ";"
	DIR_SEP  = "\\"
	DEV_NULL = "NUL"
)

func IsWsl() bool {
	return false
}

func IsWindows() bool {
	return true
}

func IsDarwin() bool {
	return false
}
