package wirepod_vosk
import (
	"strconv"
	"os"
	"fmt"
	"time"
	"github.com/kercre123/chipper/pkg/logger"
	sr "github.com/kercre123/chipper/pkg/wirepod/speechrequest"
	resty "github.com/go-resty/resty/v2"
	"encoding/json"
)

var Name string = "vosk"

var fp *os.File
var clie = resty.New()

func Init() error {

	return nil
}

func STT(req sr.SpeechRequest) (string, error) {
	logger.Println("(Bot " + strconv.Itoa(req.BotNum) + ", Google ASR) Processing...")

	curWD, _ := os.Getwd()
	curTime := time.Now()	
	fname := fmt.Sprintf("%s/dump/raw_%d%d%d-%d-%d-%d.raw",curWD, curTime.Year(),curTime.Month(),curTime.Day(),curTime.Hour(),curTime.Minute(),curTime.Second())
	fp,_ = os.OpenFile(fname, os.O_CREATE|os.O_WRONLY, 0644)

	speechIsDone := false

	for {
		var chunk []byte
		req, chunk, _ = sr.GetNextStreamChunk(req)
		req, speechIsDone = sr.DetectEndOfSpeech(req)
		fp.Write(chunk)

		if speechIsDone {
			break
		}
		
	}

	fp.Close()

	type Body struct {
		Name string `json:"command"`
		Age  string    `json:"file"`
	}
	var body = Body { "asr_raw", fname}
	bodyData, _ := json.Marshal(body)	
	fmt.Println(string(bodyData))

	resp, _ := clie.R().
		EnableTrace().
		SetHeader("Content-Type", "application/json").		
		SetBody(string(bodyData)).
		Post("http://127.0.0.1:8888/chat")

	fmt.Println("  Body       :\t", resp.String())

	asrOut := resp.String()

	logger.Println("Google ASR Bot " + strconv.Itoa(req.BotNum) + " Transcribed text: " + asrOut)


	return asrOut, nil
}
