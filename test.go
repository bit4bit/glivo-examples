package main

import "github.com/bit4bit/glivo"
import "github.com/bit4bit/glivo/chain"
import "github.com/bit4bit/gfsocket"

import (
	"fmt"
	"os"
	"flag"
	"syscall"
	"os/signal"
	"log"
)
var logger = log.New(os.Stdout, "glivo-test", log.LstdFlags)

var clientIP = flag.String("listener-address", "127.0.0.1:8084", "where listener Freeswitch outbound")
var serverIP = flag.String("freeswitch-address", "", "where is the Freeswitch")
var userTest = flag.String("user", "peter", "name of user to call for tests")
var whatTest = flag.String("test","digits", "what test to run: digits")
var onlyCall = flag.Bool("only-call",false, "only do originate")
var onlyServer = flag.Bool("only-server", false, "only start listener")
func usage() {
	fmt.Printf("Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}


func HandleCall(call *glivo.Call) {
	defer call.Close()

	call.Answer()
	digits := chain.NewChainDigits(call)
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
	flag.Parse()

	if(flag.NFlag() == 0) {
		usage()
	}



	if false == *onlyCall {


		fsout, err := glivo.NewFS(*clientIP, logger)
		if err != nil {
			logger.Fatal(err.Error())
			os.Exit(1)
		}
		defer func(){
			fsout.Stop()
		}()
		fsout.Start(HandleCall)
		if *onlyServer {
			logger.Printf("Starting listener server %s\n", *clientIP)
			sc := make(chan os.Signal)
			signal.Notify(sc, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
			<-sc
			os.Exit(0)
		}
	}



	fscmd, err := gfsocket.NewFS(*serverIP, "ClueCon")

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
		
	logger.Println(fscmd.Api(fmt.Sprintf("originate {hangup_after_bridge=false,originate_early_media=false}user/%s '&socket(%s async full)'", *userTest, *clientIP)))

	



}
