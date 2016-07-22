package w32

type (
	DWORD     uint32
	HANDLE    uintptr
	HINSTANCE HANDLE
	HKEY      HANDLE
	HWND      HANDLE
	ULONG     uint32
	LPCTSTR   uintptr
	LPVOID    uintptr
)

// https://msdn.microsoft.com/en-us/library/windows/desktop/ms681382(v=vs.85).aspx
const (
	ERROR_BAD_FORMAT = 11
)

// https://msdn.microsoft.com/en-us/library/windows/desktop/bb759784(v=vs.85).aspx
const (
	SE_ERR_FNF             = 2
	SE_ERR_PNF             = 3
	SE_ERR_ACCESSDENIED    = 5
	SE_ERR_OOM             = 8
	SE_ERR_DLLNOTFOUND     = 32
	SE_ERR_SHARE           = 26
	SE_ERR_ASSOCINCOMPLETE = 27
	SE_ERR_DDETIMEOUT      = 28
	SE_ERR_DDEFAIL         = 29
	SE_ERR_DDEBUSY         = 30
	SE_ERR_NOASSOC         = 31
)

// https://msdn.microsoft.com/en-us/library/windows/desktop/bb759784(v=vs.85).aspx
const (
	SEE_MASK_NOCLOSEPROCESS = 0x00000040
)

type SHELLEXECUTEINFO struct {
	cbSize         DWORD
	fMask          ULONG
	hwnd           HWND
	lpVerb         LPCTSTR
	lpFile         LPCTSTR
	lpParameters   LPCTSTR
	lpDirectory    LPCTSTR
	nShow          int
	hInstApp       HINSTANCE
	lpIDList       LPVOID
	lpClass        LPCTSTR
	hkeyClass      HKEY
	dwHotKey       DWORD
	hIconOrMonitor HANDLE
	hProcess       HANDLE
}
