package main

import "github.com/bit4bit/glivo"
import "github.com/bit4bit/gfsocket"

import (
	"fmt"
	"os"
)


func HandleCall(call *glivo.Call) {
	defer call.Close()

	call.Answer()
	digits := glivo.NewChainDigits(call)
	digits.SetNumDigits(5)
	digits.SetInvalidDigitsSound("/invalid.wav")
	digits.SetRetries(4)
	digits.SetValidDigits("123456789")
	
	
	rst, _ := digits.Speak("sum one more one").Question("2")
	if rst {
		fmt.Println("HOO OK")
	}else{
		fmt.Println("HOO WRONG")
	}

	cl, _ := digits.Speak("please digit something").CollectInput()

	digits.Speak("Your input:").Speak(cl).Do()
	digits.Speak("Thanks for using glivo :)").Do()
	call.Hangup()
}

func main() {
	fscmd, err := gfsocket.NewFS("172.168.1.120:8021", "ClueCon")

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fscmd.SetDebug(false)
	fsout, err := glivo.NewFS("", "8084")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fsout.Start(HandleCall)

	fmt.Println(fscmd.Api("originate {hangup_after_bridge=false,originate_early_media=false}user/peter '&socket(172.168.1.115:8084 async full)'"))



	fsout.Stop()
}
