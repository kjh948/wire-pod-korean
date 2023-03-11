package wirepod_vosk
import (
	// "encoding/json"
	"log"
	"strconv"
	"os"
	"fmt"
	"time"
	"bytes"
	// "strings"
	"os/exec"

	vosk "github.com/alphacep/vosk-api/go"
	"github.com/kercre123/chipper/pkg/logger"
	"github.com/kercre123/chipper/pkg/vars"
	sr "github.com/kercre123/chipper/pkg/wirepod/speechrequest"
	resty "github.com/go-resty/resty/v2"

	"encoding/json"
)

var Name string = "vosk"

var model *vosk.VoskModel
var modelLoaded bool
var fp *os.File
var fp_asr *os.File

// var client = resty.New()
var clie = resty.New()

func Init() error {
	if vars.APIConfig.PastInitialSetup {
		vosk.SetLogLevel(-1)
		if modelLoaded {
			logger.Println("A model was already loaded, freeing")
			model.Free()
		}
		sttLanguage := vars.APIConfig.STT.Language
		if len(sttLanguage) == 0 {
			sttLanguage = "en-US"
		}
		// Open model
		modelPath := "../vosk/models/" + sttLanguage + "/model"
		logger.Println("Opening VOSK model (" + modelPath + ")")
		aModel, err := vosk.NewModel(modelPath)
		if err != nil {
			log.Fatal(err)
			return err
		}
		model = aModel
		modelLoaded = true
		logger.Println("VOSK initiated successfully")

	}

	// client = resty.New()
	return nil
}
const ShellToUse = "bash"

func Shellout(command string) (string, string, error) {
    var stdout bytes.Buffer
    var stderr bytes.Buffer
    cmd := exec.Command(ShellToUse, "-c", command)
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    err := cmd.Run()
    return stdout.String(), stderr.String(), err
}

func STT(req sr.SpeechRequest) (string, error) {
	logger.Println("(Bot " + strconv.Itoa(req.BotNum) + ", Google ASR) Processing...")

	curWD, _ := os.Getwd()
	curTime := time.Now()	
	fname := fmt.Sprintf("%s/dump/raw_%d%d%d-%d-%d-%d.raw",curWD, curTime.Year(),curTime.Month(),curTime.Day(),curTime.Hour(),curTime.Minute(),curTime.Second())
	fp,_ = os.OpenFile(fname, os.O_CREATE|os.O_WRONLY, 0644)

	speechIsDone := false
	sampleRate := 16000.0

	rec, err := vosk.NewRecognizer(model, sampleRate)
	if err != nil {
		log.Fatal(err)
	}
	rec.SetWords(1)
	rec.AcceptWaveform(req.FirstReq)
	for {
		var chunk []byte
		req, chunk, err = sr.GetNextStreamChunk(req)
		if err != nil {
			return "", err
		}
		// has to be split into 320 []byte chunks for VAD
		req, speechIsDone = sr.DetectEndOfSpeech(req)

		fp.Write(chunk)

		if speechIsDone {
			break
		}
		
	}

	fp.Close()
	// fp_asr.Close()

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
