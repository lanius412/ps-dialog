package dialogs

import (
	"fmt"
	"strings"

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
	Prompt string
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

type Result string

func Message(msg string) *MessageBoxObj {
	obj := &MessageBoxObj{Msg: msg, Btn: Btn_OK, Icn: Icon_None}
	obj.Dlg.Title = "Message Box"
	return obj
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

func (obj *MessageBoxObj) Show() (Result, error) {
	cmd := fmt.Sprintf(`Add-Type -AssemblyName System.Windows.Forms;[System.Windows.Forms.MessageBox]::Show('%s', '%s', %d, %d)`, obj.Msg, obj.Dlg.Title, obj.Btn, obj.Icn)
	out, err := ps.Execute(cmd)
	if err != nil {
		return Result("Error"), err
	}
	return Result(out), nil
}

func InputBox() *InputBoxObj {
	obj := &InputBoxObj{Prompt: "Type in the box below"}
	obj.Dlg.Title = "Input Box"
	return obj
}

func (obj *InputBoxObj) Title(title string) *InputBoxObj {
	obj.Dlg.Title = title
	return obj
}

func (obj *InputBoxObj) Description(desc string) *InputBoxObj {
	obj.Prompt = desc
	return obj
}

func (obj *InputBoxObj) Show() (string, error) {
	cmd := fmt.Sprintf(`[void][Reflection.Assembly]::LoadWithPartialName('Microsoft.VisualBasic');[Microsoft.VisualBasic.Interaction]::InputBox('%s','%s')`, obj.Prompt, obj.Dlg.Title)
	out, err := ps.Execute(cmd)
	if err != nil {
		return "", err
	}
	return out, nil
}

func File() *FileDialogObj {
	obj := &FileDialogObj{InitialDir: "C:\\", Filter: "All files (*.*)|*.*", Multi: false}
	obj.Dlg.Title = "File Dialog"
	return obj
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

func (obj *FileDialogObj) Open() ([]string, Result, error) {
	cmd := fmt.Sprintf(`[void][System.Reflection.Assembly]::LoadWithPartialName('System.windows.forms');$fdlg = New-Object System.Windows.Forms.OpenFileDialog;$fdlg.Title = '%s';$fdlg.InitialDirectory = '%s';$fdlg.Filter = '%s';$fdlg.Multiselect = $%t;if($fdlg.ShowDialog() -eq [System.Windows.Forms.DialogResult]::OK){Write-Output $fdlg.FileNames} else {Write-Output 'Cancel'}`, obj.Dlg.Title, obj.InitialDir, obj.Filter, obj.Multi)
	out, err := ps.Execute(cmd)
	if err != nil {
		return []string{}, Result("Error"), err
	}
	if out == "Cancel" {
		return []string{}, Result(out), nil
	}
	return strings.Split(out, "\r\n"), Result("OK"), nil
}
