package processreqs

import (
	"strconv"

	// pb "github.com/digital-dream-labs/api/go/chipperpb"
	"github.com/kercre123/chipper/pkg/logger"
	"github.com/kercre123/chipper/pkg/vars"
	"github.com/kercre123/chipper/pkg/vtt"
	sr "github.com/kercre123/chipper/pkg/wirepod/speechrequest"
	ttr "github.com/kercre123/chipper/pkg/wirepod/ttr"

	"fmt"
	resty "github.com/go-resty/resty/v2"	
	"encoding/json"
)

var clie = resty.New()

func (s *Server) ProcessIntentGraph(req *vtt.IntentGraphRequest) (*vtt.IntentGraphResponse, error) {
	sr.BotNum = sr.BotNum + 1
	var successMatched bool
	speechReq := sr.ReqToSpeechRequest(req)
	var transcribedText string
	if !isSti {
		var err error
		transcribedText, err = sttHandler(speechReq)
		if err != nil {
			sr.BotNum = sr.BotNum - 1
			ttr.IntentPass(req, "intent_system_noaudio", "voice processing error", map[string]string{"error": err.Error()}, true, speechReq.BotNum)
			return nil, nil
		}
		successMatched = ttr.ProcessTextAll(req, transcribedText, matchListList, intentsList, speechReq.IsOpus, speechReq.BotNum)
	} else {
		intent, slots, err := stiHandler(speechReq)
		if err != nil {
			if err.Error() == "inference not understood" {
				logger.Println("No intent was matched")
				sr.BotNum = sr.BotNum - 1
				ttr.IntentPass(req, "intent_system_noaudio", "voice processing error", map[string]string{"error": err.Error()}, true, speechReq.BotNum)
				return nil, nil
			}
			logger.Println(err)
			sr.BotNum = sr.BotNum - 1
			ttr.IntentPass(req, "intent_system_noaudio", "voice processing error", map[string]string{"error": err.Error()}, true, speechReq.BotNum)
			return nil, nil
		}
		ttr.ParamCheckerSlotsEnUS(req, intent, slots, speechReq.IsOpus, speechReq.BotNum, speechReq.Device)
		sr.BotNum = sr.BotNum - 1
		return nil, nil
	}
	if !successMatched {
		logger.Println("No intent was matched.")
		if vars.APIConfig.Knowledge.Enable && vars.APIConfig.Knowledge.Provider == "openai" && len([]rune(transcribedText)) >= 6 {
			
			type Body struct {
				Name string `json:"command"`
				Age  string    `json:"text"`
			 }
			var body = Body { "chatgpt_tts_wav", transcribedText}
			bodyData, _ := json.Marshal(body)	
			fmt.Println(string(bodyData))
		
			resp, _ := clie.R().
				EnableTrace().
				SetHeader("Content-Type", "application/json").		
				SetBody(string(bodyData)).
				Post("http://127.0.0.1:8888/chat")
		
			fmt.Println("  Body       :\t", resp.String())

			return nil, nil
		}
		sr.BotNum = sr.BotNum - 1
		ttr.IntentPass(req, "intent_system_noaudio", transcribedText, map[string]string{"": ""}, false, speechReq.BotNum)
		return nil, nil
	}
	sr.BotNum = sr.BotNum - 1
	logger.Println("Bot " + strconv.Itoa(speechReq.BotNum) + " request served.")
	return nil, nil
}
