package main

import "github.com/bit4bit/glivo"
import "github.com/bit4bit/glivo/dptools"
import "github.com/bit4bit/gfsocket"

import (
	"fmt"
	"os"
	"flag"
	"syscall"
	"os/signal"
	"time"
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


func HandleCall(call *glivo.Call, userData interface{}) {
	//defer call.Close()
	call.WaitAnswer()
	call.Answer()
	dptools := dptools.NewDPTools(call)
	aleg, bleg := dptools.Bridge("user/1001")

	if aleg.Content["Variable_originate_disposition"] == "SUCCESS" {
		fmt.Print(aleg.Content)
		fmt.Printf("Aleg: %s", aleg.Content["Variable_hangup_cause"])
	} else {
		fmt.Printf("Aleg: %s", aleg.Content["Variable_originate_disposition"])
	}
	if bleg != nil {
		fmt.Printf("Bleg: %s", aleg.Content["Variable_bridge_hangup_cause"])
	}
	fmt.Print("\n")


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
		fsout.Start(HandleCall, nil)
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
		
	logger.Println(fscmd.Api(fmt.Sprintf("originate {hangup_after_bridge=false,originate_early_media=false}user/%s '&socket(%s sync full)'", *userTest, *clientIP)))

	time.Sleep(20 * time.Second)



}

