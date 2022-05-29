/*
Package ps-dialog provides common dialog boxes, such as MessageBox, InputBox, and Open/Save FileDialogs.

Usage
	import (
		"fmt"
		"os"

		dialog "github.com/lanius412/ps-dialog"
	)

	func main() {
		result, err := dialog.MessageBox("Message Content").Title("MessageBox").Button(dialog.Btn_YesNo).Icon(dialog.Icon_Question).Show()
		if err != nil {
			panic(err)
		}
		fmt.Println(result) // Yes

		in, _ := dialog.InputBox().Title("InputBox").Description("Type in the below box").Show()
		fmt.Println(in) // Hello

		homeDir, _ := os.UserHomeDir()
		fileDlg := dialog.File().Title("FileDialog").StartDir(homeDir).ExtFilter("TEXT File (*.txt)|*.txt")
		filepaths, _, _ := fileDlg.Open().Multiple().Load()
		fmt.Println(filepaths) // [C:\[username]\Desktop\sample.txt, C:\[username]/Desktop\sample2.txt]

		filepath, _, _ := fileDlg.Save().Load()
		fmt.Println(filepath) // C:\[username]\Desktop\save.txt
	}


*/

package dialogs

import (
	"fmt"
	"strings"

	ps "github.com/Tobotobo/powershell"
)

type Dlg struct {
	Title string
}

// Result is the value of the button was clicked
type Result string

type MessageBoxObj struct {
	Dlg
	Msg      string
	Btn, Icn int
}

//This function returns default MessageBoxObject (Title: "Message Box", ButtonType: OK, IconType: None)
func Message(msg string) *MessageBoxObj {
	obj := &MessageBoxObj{Msg: msg, Btn: Btn_OK, Icn: Icon_None}
	obj.Dlg.Title = "Message Box"
	return obj
}

func (obj *MessageBoxObj) Title(title string) *MessageBoxObj {
	obj.Dlg.Title = title
	return obj
}

// Button Pattern on MessageBox
const (
	Btn_OK = iota
	Btn_OKCancel
	Btn_AbortRetryIgnore
	Btn_YesNoCancel
	Btn_YesNo
	Btn_RetryCancel
	Btn_CancelTryContinue
)

func (obj *MessageBoxObj) Button(ButtonType int) *MessageBoxObj {
	obj.Btn = ButtonType
	return obj
}

// Icon Pattern on MessageBox
const (
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

type InputBoxObj struct {
	Dlg
	Prompt string
}

// This function returns default InputBoxObject (Title: "Input Box", Description: "Type in the below Box")
func InputBox() *InputBoxObj {
	obj := &InputBoxObj{Prompt: "Type in the below Box"}
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

type FileDialogObj struct {
	Dlg
	InitialDir, Filter string
}

// This function returns default FileDialogObject (Title: "File Dialog", InitialDir: "C:\", Filter: "All files (*.*)|*.*")
func File() *FileDialogObj {
	obj := &FileDialogObj{InitialDir: "C:\\", Filter: "All files (*.*)|*.*"}
	obj.Dlg.Title = "File Dialog"
	return obj
}

func (obj *FileDialogObj) Title(title string) *FileDialogObj {
	obj.Dlg.Title = title
	return obj
}

func (obj *FileDialogObj) StartDir(dir string) *FileDialogObj {
	obj.InitialDir = dir
	return obj
}

func (obj *FileDialogObj) ExtFilter(desc, ext string) *FileDialogObj {
	obj.Filter = fmt.Sprintf(`%s (*.%s)|*.%s`, desc, ext, ext)
	return obj
}

type OpenFileDialogObj struct {
	FileDialogObj
	Multi bool
}

func (obj *FileDialogObj) Open() *OpenFileDialogObj {
	return &OpenFileDialogObj{FileDialogObj: *obj}
}

func (obj *OpenFileDialogObj) Multiple() *OpenFileDialogObj {
	obj.Multi = true
	return obj
}

func (obj *OpenFileDialogObj) Load() ([]string, Result, error) {
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

type SaveFileDialogObj struct {
	FileDialogObj
	OverwriteWarning bool
	OverwriteForce   string
}

func (obj *FileDialogObj) Save() *SaveFileDialogObj {
	return &SaveFileDialogObj{FileDialogObj: *obj, OverwriteWarning: true, OverwriteForce: ""}
}

func (obj *SaveFileDialogObj) OverwriteWarningDisable() *SaveFileDialogObj {
	obj.OverwriteWarning = false
	return obj
}

func (obj *SaveFileDialogObj) OverwriteForceEnable() *SaveFileDialogObj {
	obj.OverwriteForce = "-Force"
	return obj
}

func (obj *SaveFileDialogObj) Load() (string, Result, error) {
	cmd := fmt.Sprintf(`[void][System.Reflection.Assembly]::LoadWithPartialName('System.windows.forms');$fdlg = New-Object System.Windows.Forms.SaveFileDialog;$fdlg.Title = '%s';$fdlg.InitialDirectory = '%s';$fdlg.Filter = '%s';$fdlg.OverwritePrompt = $%t;if($fdlg.ShowDialog() -eq [System.Windows.Forms.DialogResult]::OK){New-Item $fdlg.FileName %s} else {Write-Output 'Cancel'}`, obj.Dlg.Title, obj.InitialDir, obj.Filter, obj.OverwriteWarning, obj.OverwriteForce)
	out, err := ps.Execute(cmd)
	if err != nil {
		return "", Result("Error"), err
	}
	if out == "Cancel" {
		return "", Result(out), nil
	}
	return out, Result("OK"), nil
}
