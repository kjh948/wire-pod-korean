package sdkapp

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/digital-dream-labs/api/go/jdocspb"
	"github.com/digital-dream-labs/hugh/grpc/client"
	"github.com/fforchino/vector-go-sdk/pkg/vector"
	"github.com/fforchino/vector-go-sdk/pkg/vectorpb"
	"github.com/kercre123/chipper/pkg/logger"
	jdocsserver "github.com/kercre123/chipper/pkg/servers/jdocs"
	"github.com/kercre123/chipper/pkg/vars"

)

// the big workaround

// the escape pod CA cert only gets appended to the cert store when a jdocs connection is created
// this doesn't happen at every boot
// this utilizes Vector's connCheck to see if a bot has disconnected from the server for more than 10 seconds
// if it has, it will pull jdocs from the bot which will cause the CA cert to get appended to the store

// setting JDOCS_PINGER_ENABLED=false will disable jdocs pinger
var PingerEnabled bool = true

func pingJdocs(target string) {
	ctx := context.Background()
	target = strings.Split(target, ":")[0]
	var serial string
	matched := false
	for _, robot := range vars.BotInfo.Robots {
		if strings.TrimSpace(strings.ToLower(robot.IPAddress)) == strings.TrimSpace(strings.ToLower(target)) {
			matched = true
			serial = robot.Esn
		}
	}

	if !matched {
		logger.Println("jdocs pinger error: serial did not match any bot in bot json")
		logger.Println("Error pinging jdocs")
		return
	}
	robotTmp, err := NewWP(serial, false)
	if err != nil {
		logger.Println(err)
		return
	}
	_, err = robotTmp.Conn.BatteryState(ctx, &vectorpb.BatteryStateRequest{})
	if err != nil {
		robotTmp, err = NewWP(serial, true)
		if err != nil {
			logger.Println(err)
			logger.Println("Error pinging jdocs")
			return
		}
		_, err = robotTmp.Conn.BatteryState(ctx, &vectorpb.BatteryStateRequest{})
		if err != nil {
			logger.Println("Error pinging jdocs, likely unauthenticated")
			return
		}
	}
	resp, err := robotTmp.Conn.PullJdocs(ctx, &vectorpb.PullJdocsRequest{
		JdocTypes: []vectorpb.JdocType{vectorpb.JdocType_ROBOT_SETTINGS},
	})
	if err != nil {
		logger.Println(err)
		logger.Println("Failed to pull jdocs")
		return
	}
	logger.Println("Successfully got jdocs from " + serial)
	// write to file
	var jdoc jdocspb.Jdoc
	jdoc.DocVersion = resp.NamedJdocs[0].Doc.DocVersion
	jdoc.FmtVersion = resp.NamedJdocs[0].Doc.FmtVersion
	jdoc.ClientMetadata = resp.NamedJdocs[0].Doc.ClientMetadata
	jdoc.JsonDoc = resp.NamedJdocs[0].Doc.JsonDoc
	vars.AddJdoc("vic:"+serial, "vic.RobotSettings", jdoc)

}

var jdocsTargets []string
var jdocsTimers []int
var jdocsShouldPing []bool
var jdocsTimerStarted []bool
var jdocsTimerReset []bool

func startJdocsTimer(target string) {
	var jdocsBotNum int
	for num, ip := range jdocsTargets {
		if ip == target {
			jdocsBotNum = num
		}
	}
	if !jdocsTimerStarted[jdocsBotNum] {
		jdocsTimerStarted[jdocsBotNum] = true
		jdocsShouldPing[jdocsBotNum] = false
		logger.Println("Starting jdocs pinger timer for " + target)
		go func() {
			// wait 10 seconds
			for {
				time.Sleep(time.Second * 1)
				jdocsTimers[jdocsBotNum] = jdocsTimers[jdocsBotNum] + 1
				if jdocsTimers[jdocsBotNum] == 10 {
					logger.Println("No connCheck from " + target + " in more than 10 seconds, will ping jdocs on next check")
					jdocsShouldPing[jdocsBotNum] = true
					jdocsTimerStarted[jdocsBotNum] = false
					return
				}
				if jdocsTimerReset[jdocsBotNum] {
					jdocsTimers[jdocsBotNum] = 0
					//logger.Println("Resetting timer to 0 for bot " + target)
					jdocsTimerReset[jdocsBotNum] = false
				}
			}
		}()
	}
}

func NewWP(serial string, useGlobal bool) (*vector.Vector, error) {
	target := ""
	guid := ""
	if serial == "" {
		log.Fatal("please use the -serial argument and set it to your robots serial number")
		return nil, fmt.Errorf("Configuration options missing")
	}
	wirepodPath := os.Getenv("WIREPOD_HOME")
	if len(wirepodPath) == 0 {
		wirepodPath = "."
	}
	matched := false
	for _, robot := range vars.BotInfo.Robots {
		if strings.EqualFold(serial, robot.Esn) {
			matched = true
			target = robot.IPAddress + ":443"
			guid = robot.GUID
			break
		}
	}
	if !matched {
		log.Println("vector-go-sdk error: serial did not match any bot in bot json")
		return nil, errors.New("vector-go-sdk error: serial did not match any bot in bot json")
	}
	c, err := client.New(
		client.WithTarget(target),
		client.WithInsecureSkipVerify(),
	)
	if err != nil {
		return nil, err
	}
	if err := c.Connect(); err != nil {
		return nil, err
	}
	return vector.New(
		vector.WithTarget(target),
		vector.WithSerialNo(serial),
		vector.WithToken(guid),
	)
}

func jdocsPingTimer(target string) bool {
	for num, ip := range jdocsTargets {
		if ip == target {
			var returnValue bool = jdocsShouldPing[num]
			startJdocsTimer(target)
			jdocsTimerReset[num] = true
			if returnValue {
				jdocsShouldPing[num] = false
			}
			return returnValue
		}
	}
	jdocsTargets = append(jdocsTargets, target)
	jdocsTimers = append(jdocsTimers, 0)
	jdocsShouldPing = append(jdocsShouldPing, false)
	jdocsTimerStarted = append(jdocsTimerStarted, false)
	jdocsTimerReset = append(jdocsTimerReset, false)
	startJdocsTimer(target)
	return true
}

func connCheck(w http.ResponseWriter, r *http.Request) {
	switch {
	default:
		http.Error(w, "not found", http.StatusNotFound)
		return
	case strings.Contains(r.URL.Path, "/ok"):
		if PingerEnabled {
			//	logger.Println("connCheck request from " + r.RemoteAddr)
			robotTarget := strings.Split(r.RemoteAddr, ":")[0] + ":443"
			robotTargetCheck := strings.Split(r.RemoteAddr, ":")[0]
			jsonB, _ := os.ReadFile(jdocsserver.InfoPath)
			json := string(jsonB)
			if strings.Contains(json, strings.TrimSpace(robotTargetCheck)) {
				ping := jdocsPingTimer(robotTarget)
				if ping {
					pingJdocs(robotTarget)
				}
			}
		}
		fmt.Fprintf(w, "ok")
		return
	}
}
