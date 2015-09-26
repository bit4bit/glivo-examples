package main

import "github.com/bit4bit/glivo"
import "github.com/bit4bit/glivo/chain"
import "github.com/bit4bit/gfsocket"

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var logger = log.New(os.Stdout, "glivo-test", log.LstdFlags)

var clientIP = flag.String("listener-address", "127.0.0.1:8084", "where listener Freeswitch outbound")
var serverIP = flag.String("freeswitch-address", "", "where is the Freeswitch")
var userTest = flag.String("user", "peter", "name of user to call for tests")
var whatTest = flag.String("test", "digits", "what test to run: digits")
var onlyCall = flag.Bool("only-call", false, "only do originate")
var onlyServer = flag.Bool("only-server", false, "only start listener")
var gateway = flag.String("gateway", "user", "gateway")

var waitHangup = make(chan bool)

func usage() {
	fmt.Printf("Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func HandleCall(call *glivo.Call, userData interface{}) {
	//defer call.Close()
	call.WaitAnswer()
	call.Answer()
	digits := chain.NewChainDigits(call)
	digits.SetNumDigits(5)
	digits.SetInvalidDigitsSound("/invalid.wav")
	digits.SetRetries(4)
	digits.SetValidDigits("123456789")

	rst, _ := digits.Speak("sum one more one").Question("2")
	if rst {
		fmt.Println("HOO OK")
	} else {
		fmt.Println("HOO WRONG")
	}

	cl, _ := digits.Speak("please digit something").CollectInput()
	digits.Speak("Your input:").Speak(cl).Do()
	digits.Speak("Thanks for using glivo :)").Do()
	call.Hangup()
	waitHangup <- true
}

func main() {
	flag.Parse()

	if flag.NFlag() == 0 {
		usage()
	}

	if false == *onlyCall {

		fsout, err := glivo.Listen(*clientIP, logger)
		if err != nil {
			logger.Fatal(err.Error())
			os.Exit(1)
		}
		defer func() {
			fsout.Stop()
		}()
		go fsout.Serve(HandleCall, nil)
		if *onlyServer {
			logger.Printf("Starting listener server %s\n", *clientIP)
			sc := make(chan os.Signal)
			signal.Notify(sc, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
			<-sc
			os.Exit(0)
		}
	}

	fscmd, err := gfsocket.Dial(*serverIP, "ClueCon")

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	ret := fscmd.Api(fmt.Sprintf("originate {hangup_after_bridge=true,ignore_early_media=true}%s/%s '&socket(%s sync full)'", *gateway, *userTest, *clientIP))
	fmt.Print(ret)
	if ret.Status == "+OK" {
		<-waitHangup
	} else {
		logger.Fatal(ret.Content)
	}

}
