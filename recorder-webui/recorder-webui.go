package main

import (
	"os/exec"
	"io"
	"fmt"
	"errors"
	"net/http"
	"tad-demo/common"
	"html/template"
)

var cfg = common.NewConfig()

type State struct{
	started bool
	cmd *exec.Cmd
	writer *io.WriteCloser
}

type RecorderState struct {
	PlayerState string
	StartVisibility string
	RecordingVisibility string
	LatestRecord string
	LatestRecordVisibility string
}

type PerfTest struct {
	state *State
	//subscription Subscription
}

func (t* PerfTest) Start() {
	common.Info.Println("Start webui")

	//t.subscription = Subscription{acceptedQueue: make(chan string, 100)}

	go func() {
		http.HandleFunc("/start", t.handlerStart)
		http.HandleFunc("/stop", t.handlerStop)
		http.HandleFunc("/", t.handlerStatus)
		err := http.ListenAndServe(fmt.Sprintf(":%d", 8090), nil)
		if err != nil {
			panic(err)
		}
	}()
}

func (self* PerfTest) handlerStatus(w http.ResponseWriter, r *http.Request) {
	var state *RecorderState

	_, body, err := common.Get("http://" + cfg.Service.Host + ":8080/last")

	var latestRecord string
	if err == nil {
		latestRecord = string(body)
		common.Trace.Println("latest record:", latestRecord)
	}
	
	latestRecordVisibility := "none"
	if latestRecord != "" {
		latestRecordVisibility = "block"
	}

	latestRecord = "http://" + cfg.Service.Host + ":8080/" + latestRecord

	if !self.state.started {
		state = &RecorderState{
					PlayerState: "Recorder is not started: Press 'Start' button",
					StartVisibility: "block",
					RecordingVisibility: "none",
					LatestRecord: latestRecord,
					LatestRecordVisibility: latestRecordVisibility}
	} else if self.state.started && self.state.cmd == nil {
		state = &RecorderState{
					PlayerState: "Internal error: Invalid recorder state",
					StartVisibility: "block",
					RecordingVisibility: "none",
					LatestRecordVisibility: "none"}
	}else if self.state.started && self.state.cmd != nil {
		state = &RecorderState{
					PlayerState: "Recorder in progress: Drop conference to get link",
					StartVisibility: "none",
					RecordingVisibility: "block",
					LatestRecordVisibility: "none" }
	}

	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("page.html")
	t.Execute(w, state)
}

func (self* PerfTest) handlerStart(w http.ResponseWriter, r *http.Request) {
	self.start()
	http.Redirect(w, r, "/", 301)
}

func (self PerfTest) handlerStop(w http.ResponseWriter, r *http.Request) {
	self.stop()
	http.Redirect(w, r, "/", 301)
}

func (t* PerfTest) start() error {
	common.Info.Print("PerfTest.start:", t.state.started)
	if t.state.started {
		common.Info.Print("Test alredy started")
		return errors.New("Test alredy started")
	}
	common.Trace.Print("Set ++")
	t.state.started = true
	err := t.startCommand()
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
	if t.state.writer != nil {
		common.Trace.Print("Killprocess ---")
		writer := *t.state.writer;
		writer.Write([]byte("terminate\n"))
		writer.Close()
		//t.state.cmd.Process.Kill()
		//t.state.cmd = nil
	}
	common.Trace.Print("Reset started !!!")
	t.state.started = false;
	return nil
}

func (t* PerfTest) startCommand() error{
	cmdName := "./start_recorder.sh"

	cmd := exec.Command(cmdName)
	writer, err := cmd.StdinPipe()
	if err != nil {
		common.Error.Print("can not get input pipe", err)
		return err
	}
	err = cmd.Start()
	common.Info.Print("!! after start = ", err)
	if err != nil {
		return err
	}
	t.state.cmd = cmd
	t.state.writer = &writer
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

func (sms* PerfTest) Await() {
	common.Info.Println("wait for ctrl+c")
	common.WaitCtrlC()
}