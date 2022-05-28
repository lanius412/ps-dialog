# ps-dialog
Windows Dialog using powershell for Go

# Features
* Message Box  
 |- Button Type - https://docs.microsoft.com/ja-jp/dotnet/api/system.windows.forms.messageboxbuttons?view=windowsdesktop-6.0  
 |_ Icon Type - https://docs.microsoft.com/ja-jp/dotnet/api/system.windows.forms.messageboxicon?view=windowsdesktop-6.0
* Input Box
* File Dialog

# Usage
```
import dialog "github.com/lanius412/ps-dialog"

func main() {
  /* Message Box */
  msgBox := dialog.Message("Do you want to try again?").Title("Message")
  result, err := msgBox.Button(dialog.Btn_AbortRetryIgnore).Icon(dialog.Icon_Exclamation).Show()
  if err != nil {
    panic(err)
  }
  fmt.Println(result) // Abort, Retry or Ignore
  
  /* Input Box */
  psword, _ := dialog.InputBox().Title("Input").Description("Type password").Show()
  fmt.Println(psword)
  
  /* File Dialog */
  fileDlg := dialog.File().Title("File")
  filepaths, result, _ := fileDlg.SetStartDir("C:\\").SetFileter("text file(*.txt)", "txt").Multiple().Open()
  if result != "Cancel" {
    fmt.Println(filepaths) // [C:\Users\[username]\Downloads\sample1.txt, C:\Users\dev_win\Downloads\sample2.txt]
  }
}
```
