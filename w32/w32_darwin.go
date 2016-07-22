package w32

func ShellExecuteAndWait(hwnd HWND, lpOperation, lpFile, lpParameters, lpDirectory string, nShowCmd int) error {
	panic("Unsupported OS")
}

func ShellExecuteEx(pExecInfo *SHELLEXECUTEINFO) error {
	panic("Unsupported OS")
}
