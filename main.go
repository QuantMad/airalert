package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type AlertData struct {
	States map[string]State
}

type State struct {
	Enabled bool
}

const (
	IsStarted  int = 0
	IsOnAlert  int = 1
	IsNotAlert int = 2

	UpdatingInterval time.Duration = 10

	UndefinedState string = "undefined"
	ConfigFileName string = "airalert.conf"

	AirAlertAPI string = "https://emapa.fra1.cdn.digitaloceanspaces.com/statuses.json"
)

func main() {
	var currentState string = tryLoadCurrentState()

	for currentState == UndefinedState {
		handleClearRun()
		currentState = tryLoadCurrentState()
	}

	var currentStatus int = IsStarted

	for {
		var alertData = loadAlertData()

		if alertData.States[currentState].Enabled && currentStatus == IsNotAlert {
			println("Air alert started")
			currentStatus = IsOnAlert
			exec.Command("mpg123", "./res/male_air_on.mp3").Run()
		} else if !alertData.States[currentState].Enabled && currentStatus == IsOnAlert {
			println("Air alert expired")
			currentStatus = IsNotAlert
			exec.Command("mpg123", "./res/male_air_off.mp3").Run()
		}
		println("\t", currentState, ":", alertData.States[currentState].Enabled)

		time.Sleep(UpdatingInterval * time.Second)
	}
}

func tryLoadCurrentState() string {
	bytes, err := ioutil.ReadFile(getConfigPath())

	if err != nil {
		println("Cannot read config file")

		return UndefinedState
	}

	loadedStateName := string(bytes)
	println("Loaded state name:", loadedStateName)

	return loadedStateName
}

func handleClearRun() {
	println("Handling clear run")
	for {
		currentStateName := getStateNameFromUser()
		trySaveCurrentState(currentStateName)

		break
	}
}

func loadAlertData() AlertData {
	for {
		time.Sleep(1 * time.Second)
		rawJSON, err := http.Get(AirAlertAPI)

		if err != nil {
			println(err)
			continue
		}

		var alertData AlertData
		err = json.NewDecoder(rawJSON.Body).Decode(&alertData)

		if err != nil {
			println(err)
			continue
		}

		return alertData
	}
}

func getStateNameFromUser() string {
	alertData := loadAlertData()
	println("Select your state: ")
	keys := make([]string, 0, len(alertData.States))
	var i int = 0
	for k := range alertData.States {
		keys = append(keys, k)
		println(i, ":", keys[i])
		i++
	}

	var reader = bufio.NewReader(os.Stdin)
	for {
		print(">: ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		index, err := strconv.Atoi(text)

		if err != nil {
			println("Invalid input:", err)
			continue
		}

		keysLen := len(keys)

		if index < 0 || index >= keysLen {
			println("Invalid input:", "Index must be in range between 0 and", keysLen)
			continue
		}

		return keys[index]
	}
}

func trySaveCurrentState(stateName string) {
	configPath := getConfigPath()

	configFile, err := os.Create(configPath)

	if err != nil {
		panic("Cannot create config file\n\t" + err.Error())
	}

	defer configFile.Close()

	_, err2 := configFile.WriteString(stateName)

	if err2 != nil {
		panic("Cannot write to config file\n\t" + err2.Error())
	}

	println("config file", configPath, "saved with value:", stateName)
}

func getConfigPath() string {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		panic("Cannot open user directory:\n\t" + err.Error())
	}

	return homeDir + "/." + ConfigFileName
}

func boolToStatus(boolStatus bool) int {
	if boolStatus {
		return IsOnAlert
	} else {
		return IsNotAlert
	}
}
