//go:build windows

package dialog

import (
	"fmt"

	ps "github.com/Tobotobo/powershell"
)

type Dlg struct {
	Title string
}

type MessageBoxObj struct {
	Dlg
	Msg      string
	Btn, Icn int
}

type InputBoxObj struct {
	Dlg
	Promtpt string
}

type FileDialogObj struct {
	Dlg
	InitialDir, Filter string
	Multi              bool
}

const (
	Btn_OK = iota
	Btn_OKCancel
	Btn_AbortRetryIgnore
	Btn_YesNoCancel
	Btn_YesNo
	Btn_RetryCancel
	Btn_CancelTryContinue

	Icon_None  = 0
	Icon_Error = 16
	Icon_Hand
	Icon_Stop
	Icon_Question    = 32
	Icon_Exclamation = 48
	Icon_Warning
	Icon_Asterisk = 64
	Icon_Information
)

func Message(msg string) *MessageBoxObj {
	return &MessageBoxObj{Msg: msg, Btn: Btn_OK, Icn: Icon_None}
}

func (obj *MessageBoxObj) Title(title string) *MessageBoxObj {
	obj.Dlg.Title = title
	return obj
}

func (obj *MessageBoxObj) Button(ButtonType int) *MessageBoxObj {
	obj.Btn = ButtonType
	return obj
}

func (obj *MessageBoxObj) Icon(IconType int) *MessageBoxObj {
	obj.Icn = IconType
	return obj
}

func (obj *MessageBoxObj) Show() string {
	out, err := ps.Execute(fmt.Sprintf(`[System.Windows.Forms.MessageBox]::Show(%s, %s, %d, %d)`, obj.Dlg.Title, obj.Msg, obj.Btn, obj.Icn))
	if err != nil {
		panic(err)
	}
	return out
}

func InputBox() *InputBoxObj {
	return &InputBoxObj{Promtpt: ""}
}

func (obj *InputBoxObj) Title(title string) *InputBoxObj {
	obj.Dlg.Title = title
	return obj
}

func (obj *InputBoxObj) Description(desc string) *InputBoxObj {
	obj.Promtpt = desc
	return obj
}

func (obj *InputBoxObj) Show() string {
	out, err := ps.Execute(fmt.Sprintf(`[Microsoft.VisualBasic.Interaction]::InputBox(%s, %s)`, obj.Promtpt, obj.Title))
	if err != nil {
		panic(err)
	}
	return out
}

func File() *FileDialogObj {
	return &FileDialogObj{InitialDir: "C:\\", Filter: "All files (*.*)|*.*", Multi: false}
}

func (obj *FileDialogObj) Title(title string) *FileDialogObj {
	obj.Dlg.Title = title
	return obj
}

func (obj *FileDialogObj) SetStartDir(dir string) *FileDialogObj {
	obj.InitialDir = dir
	return obj
}

func (obj *FileDialogObj) SetFilter(desc, ext string) *FileDialogObj {
	obj.Filter = fmt.Sprintf(`%s (*.%s)|*.%s`, desc, ext, ext)
	return obj
}

func (obj *FileDialogObj) Multiple() *FileDialogObj {
	obj.Multi = true
	return obj
}

func (obj *FileDialogObj) Open() string {
	out, err := ps.Execute(fmt.Sprintf(`
		[void][System.Reflection.Assembly]::LoadWithPartialName("System.windows.forms")
		$fdlg = New-Object System.Windows.Forms.OpenFileDialog
		$fdlg.Title = %s
		$fdlg.InitialDirectory = %s
		$fdlg.Filter = %s
		$fdlg.Multiselect = $true
	`, obj.Dlg.Title, obj.InitialDir, obj.Filter))
	if err != nil {
		panic(err)
	}
	return out
}
