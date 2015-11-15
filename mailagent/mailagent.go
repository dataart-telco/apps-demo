package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"gopkg.in/gomail.v1"
	"gopkg.in/redis.v3"
	"io"
	"os"
	"tad-demo/common"
	"time"
)

const INCOMING_PHONE_PREFIX = "conf:"
const KEY_GATHER = "gather"

var prevMd5Hash string

func main() {

	fmt.Println("Start mailagent")

	user := flag.String("user", "", "Smtp user")
	pass := flag.String("pass", "", "Smtp password")
	timer := flag.Int("timer", 5, "timer in minutes")

	cfg := common.NewConfig()
	delay := time.Duration(*timer) * time.Minute

	db := common.NewDbClient(cfg.Service.Redis)

	if *user == "" || *pass == "" {
		fmt.Println("No user or password")
		return
	}

	execMe := func() {
		incomingFile := dumpIncoming(db)
		//acceptedFile := dumpAccepted(db)
		defer os.Remove(incomingFile)
		//		defer os.Remove(acceptedFile)

		md5Hash, err := ComputeMd5(incomingFile)
		if err == nil {
			if md5Hash != "" && prevMd5Hash == md5Hash {
				common.Trace.Println("Ignore sendign file. hashes the same: " + md5Hash)
				return
			}
			prevMd5Hash = md5Hash
		}
		err = sendEmail(*user, *pass, incomingFile /*, acceptedFile*/)
		if err != nil {
			fmt.Println("Send message error", err)
			return
		}
	}

	execMe()

	ticker := time.NewTicker(delay)
	go func() {
		for {
			select {
			case <-ticker.C:
				execMe()
			}
		}
	}()

	fmt.Println("press Ctrl+C")
	common.WaitCtrlC()

	ticker.Stop()
}

func ComputeMd5(filePath string) (string, error) {
	var result []byte
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(result)), nil
}

func dumpIncoming(db *redis.Client) string {
	fmt.Println(time.Now(), "dumpIncoming->")

	now := time.Now()
	fileName := "dump_incomming" + now.Format("2006_01_02_15_04_05") + ".txt"
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for _, k := range db.MGet(db.Keys(INCOMING_PHONE_PREFIX + "*").Val()...).Val() {
		f.WriteString(fmt.Sprintln(k))
	}
	return fileName
}

func dumpAccepted(db *redis.Client) string {
	fmt.Println(time.Now(), "dumpPromted->")

	now := time.Now()
	fileName := "dump_accepted" + now.Format("2006_01_02_15_04_05") + ".txt"
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for _, e := range db.LRange(KEY_GATHER, 0, 1000).Val() {
		f.WriteString(fmt.Sprintln(e))
	}
	return fileName
}

func sendEmail(user string, pass string, incomingFile string /*, acceptedFile string*/) error {
	fmt.Println("sendEmail->")
	_, month, day := time.Now().Date()

	msg := gomail.NewMessage()
	msg.SetHeader("From", user)
	msg.SetHeader("To", user)
	msg.SetHeader("Subject", fmt.Sprintf("Redis dump - %d/%d", day, month))
	msg.SetBody("text/html", "Redis dump is attached")

	f1, err := gomail.OpenFile(incomingFile)
	if err == nil {
		msg.Attach(f1)
	} else {
		fmt.Println("ERROR: can't open file", incomingFile)
	}

	/*	f2, err := gomail.OpenFile(acceptedFile)
		if err == nil {
			msg.Attach(f2)
		}else{
			fmt.Println("ERROR: can't open file", acceptedFile)
		}
	*/

	mailer := gomail.NewMailer("smtp.gmail.com", user, pass, 465)
	if err := mailer.Send(msg); err != nil {
		return err
	}
	return nil
}
