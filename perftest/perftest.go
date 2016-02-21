package main

import (
	"bytes"
	"bufio"
	"os/exec"
	"fmt"
	"errors"
	"net/http"
	"tad-demo/common"
)

var cfg = common.NewConfig()
var restcommApi = common.NewRestcommApi(cfg.Service.Restcomm, cfg.Auth.User, cfg.Auth.Pass)

/*
type Subscription struct {
	acceptedQueue chan string
}

func (s Subscription) Receive() string {
	return <-s.acceptedQueue
}

*/
type State struct{
	started bool
	cmd *exec.Cmd
	buffer *bytes.Buffer
}

type PerfTest struct {
	state *State
	//subscription Subscription
}

func (t* PerfTest) RegisterNumber() {
	common.Info.Println("\tRegister number:", "6666", "7777")
	
	common.NewIncomingPhoneNumber("", "6666").CreateOrUpdate(restcommApi, fmt.Sprintf("http://%s/start", cfg.GetExternalAddress(30666)))
	common.NewIncomingPhoneNumber("", "7777").CreateOrUpdate(restcommApi, fmt.Sprintf("http://%s/stop", cfg.GetExternalAddress(30666)))
}

func (t* PerfTest) Subscribe() {
	common.Info.Println("Start perftest")

	//t.subscription = Subscription{acceptedQueue: make(chan string, 100)}

	go func() {
		http.HandleFunc("/start", t.handlerStart)
		http.HandleFunc("/stop", t.handlerStop)
		http.HandleFunc("/status", t.handlerStatus)
		err := http.ListenAndServe(fmt.Sprintf(":%d", 30666), nil)
		if err != nil {
			panic(err)
		}
	}()
}

func (t* PerfTest) handlerStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if !t.state.started {
		fmt.Fprintf(w, "<b>State:</b> Test is not run")
		return
	}
	if t.state.cmd == nil {
		fmt.Fprintf(w, "<b>State:</b> Test is run. But cmd is null. Seems internal error")
		return
	}
	var text string
	if t.state.buffer != nil {
		text = t.state.buffer.String()
	}
	fmt.Fprintf(w, "<b>State:</b> Test in progress<br>%s", text)
}

func (t* PerfTest) handlerStart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/xml")
	err := t.start()
	if err == nil {
		fmt.Fprintf(w, "<Response><Say>Test will be run</Say><Hangup/></Response>")
	}else{
		fmt.Fprintf(w, "<Response><Say>%s</Say><Hangup/></Response>", err.Error())
	}
}

func (t PerfTest) handlerStop(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/xml")
	err := t.stop()
	if err == nil {
		fmt.Fprintf(w, "<Response><Say>Test will be stopped</Say><Hangup/></Response>")
	}else{
		fmt.Fprintf(w, "<Response><Say>%s</Say><Hangup/></Response>", err.Error())
	}
}

func (t* PerfTest) start() error {
	common.Info.Print("PerfTest.start:", t.state.started)
	if t.state.started {
		common.Info.Print("Test alredy started")
		return errors.New("Test alredy started")
	}
	common.Trace.Print("Set ++")
	t.state.started = true
	err := t.startAb()
	if err != nil {
		common.Trace.Print("Reset started !!!")
		t.state.started = false
		common.Error.Print("Can not run test", err)
		return errors.New("Internal Error. Can not run test")
	}
	return nil
}

func (t* PerfTest) stop() error {
	common.Info.Print("PerfTest.stop: ", t.state.started)
	if !t.state.started {
		common.Info.Print("Test is not run")
		return errors.New("Test is not run")
	}
	if t.state.cmd != nil {
		common.Trace.Print("Killprocess ---")
		t.state.cmd.Process.Kill()
		t.state.cmd = nil
	}
	common.Trace.Print("Reset started !!!")
	t.state.started = false;
	return nil
}

func (t* PerfTest) startAb() error{
	cmdName := "ab"
	cmdArgs := []string{"-n", "2000000", "-c", "1000", "-w", "-r", "http://54.92.251.25:30790/test.xml"}

	cmd := exec.Command(cmdName, cmdArgs...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(cmdReader)
	var buffer bytes.Buffer
	go func() {
		for scanner.Scan() {
			buffer.WriteString(scanner.Text() + "\n")
		}
	}()
	err = cmd.Start()
	common.Info.Print("!! after start = ", err)
	if err != nil {
		return err
	}
	t.state.cmd = cmd
	t.state.buffer = &buffer
	go func() {
		common.Info.Print("cmd.wait")
	    t.state.cmd.Wait()
	    common.Info.Print("after cmd.wait")
	    t.state.cmd = nil
	    common.Trace.Print("Reset started !!!")
	    t.state.started = false
	    common.Info.Print("after reset")
	}()
	return nil
}

func (t* PerfTest) makeHttpCall(){
	if !t.state.started {
		common.Trace.Print("Stop do it internal")
		return
	}
	resp, err := http.Get("http://54.92.251.25:30790/test.xml")
	if err != nil {
		common.Error.Println("Can't send request:", err)
		return
	}
	resp.Body.Close()
}

func (sms* PerfTest) Await() {
	common.Info.Println("wait for ctrl+c")
	common.WaitCtrlC()
}