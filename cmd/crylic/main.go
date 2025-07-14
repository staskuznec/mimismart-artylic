package main

import (
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"arylic-multiroom/internal/config"
	arylic_api "arylic-multiroom/internal/transport/api/arylic-api"

	"github.com/lmittmann/tint"
)

func main() {

	// Получаем имя исполняемого файла и формируем имя конфигурационного файла
	execName := filepath.Base(os.Args[0])
	configFileName := strings.TrimSuffix(execName, filepath.Ext(execName)) + ".yaml"
	cfg := config.LoadConfig(configFileName)

	var level slog.Level
	switch cfg.LogLevel {
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelDebug
	}
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level: level, // уровень из конфига
	}))
	slog.SetDefault(logger)
	if len(os.Args) < 3 {
		errorStr := fmt.Sprintf("Использование: %s 3 параметра 1,2 параметр игнорируется 3 должен иметь формат:  \n<ip>|<section>|<command>|<value> в формате hex\n", os.Args[0])
		logger.Error(errorStr)
		os.Exit(1)
	}
	hexArg := os.Args[3]
	decoded, err := hex.DecodeString(hexArg)
	if err != nil {
		logger.Error("Ошибка декодирования hex:", err)
		os.Exit(1)
	}
	cleaned := strings.TrimRight(string(decoded), "\x00")
	parts := strings.SplitN(cleaned, "|", 4)
	if len(parts) != 4 {
		logger.Error("Неверный формат строки, ожидалось ip|section|command|value")
		os.Exit(1)
	}
	//logger.Debug("Полученные параметры: ", decoded)
	ip := parts[0]
	section := parts[1]
	command := parts[2]
	value := parts[3]
	paramStr := fmt.Sprintf("ip: %s, section: %s, command: %s, value: %s", ip, section, command, value)
	logger.Debug(paramStr)

	api := arylic_api.NewAPI(ip)
	//fmt.Printf("Отправка команды '%s' на устройство с IP %s\n", command, ip)
	switch {
	case section == "player":
		err = handlePlayer(api, command, value)
	case section == "network":

	case section == "deviceInfo":

	default:

	}
	if err != nil {
		logger.Error("Ошибка: ", err)
		os.Exit(1)
	}
}

func handlePlayer(api *arylic_api.ArylicAPI, command, value string) error {
	switch command {
	case "onepause":
		return api.PlayBack.OnePause()
	case "resume":
		return api.PlayBack.Resume()
	case "stop":
		return api.PlayBack.Stop()
	case "next":
		return api.PlayBack.Next()
	case "prev":
		return api.PlayBack.Prev()
	case "volumeUP":
		return api.PlayBack.VolumeUp()
	case "volumeDown":
		return api.PlayBack.VolumeDown()
	case "setVolume":
		if value == "" {
			return fmt.Errorf("значение для setVolume не указано")
		}
		return api.PlayBack.SetVolume(value)
	case "setMute":
		return api.PlayBack.Mute()
	case "nextSource":
		return api.PlayBack.SetNextInputSource()
	case "seekBack":
		return api.PlayBack.SeekBack(5000) // 5 секунд назад
	case "seekForward":
		return api.PlayBack.SeekForward(5000) // 5 секунд вперед
	case "loopmode":
		return api.PlayBack.SetShuffleAndRepeat(value)
	case "playM3U":
		if value == "" {
			return fmt.Errorf("значение для playM3U не указано")
		}
		return api.PlayBack.PlayM3U(value)
	case "playURL":
		if value == "" {
			return fmt.Errorf("значение для playURL не указано")
		}
		return api.PlayBack.PlayUrl(value)
	default:
		return fmt.Errorf("неизвестная команда: %s", command)
	}
}
