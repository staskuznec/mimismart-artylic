package arylic_api

import (
	"fmt"
	"strconv"
)

type PlayBackApi struct {
	api *ArylicAPI
}

func NewPlayBackApi(api *ArylicAPI) *PlayBackApi {
	return &PlayBackApi{api: api}
}

// Key	Value Description
//type	0: Main or standalone device
//1: Device is a Multiroom Guest
//
//ch	Active channel(s)
//0: Stereo
//1: Left
//2: Right
//
//mode	Playback mode
//0: Idling
//1: airplay streaming
//2: DLNA streaming
//10: Playing network content, e.g. vTuner, Home Media Share, Amazon Music, Deezer, etc.
//11: playing UDISK(Local USB disk on Arylic Device)
//20: playback start by HTTPAPI
//31: Spotify Connect streaming
//40: Line-In input mode
//41: Bluetooth input mode
//43: Optical input mode
//47: Line-In #2 input mode
//51: USBDAC input mode
//99: The Device is a Guest in a Multiroom Zone
//
//loop	Is a Combination of SHUFFLE and REPEAT modes
//0: SHUFFLE: disabled REPEAT: enabled - loop
//1: SHUFFLE: disabled REPEAT: enabled - loop once
//2: SHUFFLE: enabled REPEAT: enabled - loop
//3: SHUFFLE: enabled REPEAT: disabled
//4: SHUFFLE: disabled REPEAT: disabled
//5: SHUFFLE: enabled REPEAT: enabled - loop once
//
//eq	The current Equalizer setting
//
//status	Device status
//stop: no audio selected
//play: playing audio
//load: load ??
//pause: audio paused
//
//curpos	Current playing position (in ms)
//
//offset_pts	!! DOCUMENTATION IN PROGRESS !!
//
//totlen	Current track length (in ms)
//
//Title	[hexed string] of the track title
//
//Artist	[hexed string] of the artist
//
//Album	[hexed string] of the album
//
//alarmflag	!! DOCUMENTATION IN PROGRESS !!
//
//plicount	The total number of tracks in the playlist
//
//plicurr	Index of current track in playlist
//
//vol	Current volume
//Value range is from 0 - 100. So can be considered a linear percentage (0% to 100%)
//
//mute	The mute status
//0: Not muted
//1: Muted

type PlayerStatus struct {
	Type      string `json:"type"`
	Ch        string `json:"ch"`
	Mode      string `json:"mode"`
	Loop      string `json:"loop"`
	Eq        string `json:"eq"`
	Status    string `json:"status"`
	CurPos    string `json:"curpos"`
	OffsetPts string `json:"offset_pts"`
	TotLen    string `json:"totlen"`
	Title     string `json:"Title"`
	Artist    string `json:"Artist"`
	Album     string `json:"Album"`
	AlarmFlag string `json:"alarmflag"`
	PliCount  string `json:"plicount"`
	PliCurr   string `json:"plicurr"`
	Vol       string `json:"vol"`
	Mute      string `json:"mute"`
}

// getPlayerStatus получает статус плеера по HTTP API.
func (p *PlayBackApi) getPlayerStatus() (*PlayerStatus, error) {
	var status PlayerStatus
	err := p.api.DoAPIRequest("GET", "getPlayerStatus", &status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

//SelectInputSource выбирает источник входного сигнала для плеера.
//The audio source that has to be switched
//wifi: wifi mode
//line-in: line analogue input
//bluetooth: bluetooth mode
//optical: optical digital input
//co-axial: co-axial digital input
//line-in2: line analogue input #2
//udisk: UDisk mode
//PCUSB: USBDAC mode

var InputSources = []string{
	"wifi",
	"line-in",
	"bluetooth",
	"optical",
	"co-axial",
	"line-in2",
	"udisk",
	"PCUSB",
}

// isValidSource проверяет, является ли указанный источник входного сигнала допустимым.
func isValidSource(source string) bool {
	for _, s := range InputSources {
		if s == source {
			return true
		}
	}
	return false
}

func (p *PlayBackApi) SelectInputSource(source string) error {
	if !isValidSource(source) {
		return fmt.Errorf("неизвестный источник входного сигнала: %s", source)
	}
	return p.api.DoAPIRequest("GET", "setPlayerCmd:source:"+source, nil)
}

// SetNextInputSource переключает плеер на следующий источник входного сигнала.
func (p *PlayBackApi) SetNextInputSource() error {
	status, err := p.getPlayerStatus()
	if err != nil {
		return err
	}
	current := status.Mode
	idx := -1
	for i, s := range InputSources {
		if s == current {
			idx = i
			break
		}
	}
	nextIdx := (idx + 1) % len(InputSources)
	return p.SelectInputSource(InputSources[nextIdx])
}

// PlayUrl отправляет команду на воспроизведение указанного файла URL в плеере.
func (p *PlayBackApi) PlayUrl(url string) error {
	if url == "" {
		return fmt.Errorf("пустой URL для воспроизведения")
	}
	// Отправляем команду на воспроизведение URL
	err := p.api.DoAPIRequest("GET", "setPlayerCmd:play:"+url, nil)
	if err != nil {
		return err
	}
	return nil
}

// PlayM3U отправляет команду на воспроизведение указанного M3U-файла/плейлиста в плеере.
func (p *PlayBackApi) PlayM3U(m3u string) error {
	if m3u == "" {
		return fmt.Errorf("пустой M3U-файл для воспроизведения")
	}
	// Отправляем команду на воспроизведение M3U-файла
	err := p.api.DoAPIRequest("GET", "setPlayerCmd:m3u:play:"+m3u, nil)
	if err != nil {
		return err
	}
	return nil
}

// PlaySelectedTrack отправляет команду на воспроизведение выбранного трека по индексу в плеере.
func (p *PlayBackApi) PlaySelectedTrack(track string) error {
	if track == "" {
		return fmt.Errorf("пустой трек для воспроизведения")
	}
	// Отправляем команду на воспроизведение выбранного трека
	err := p.api.DoAPIRequest("GET", "setPlayerCmd:playindex:"+track, nil)
	if err != nil {
		return err
	}
	return nil
}

// SetShuffleAndRepeat устанавливает режимы перемешивания и повторения воспроизведения плеера.
func (p *PlayBackApi) SetShuffleAndRepeat(mode string) error {
	return p.api.DoAPIRequest("GET", "setPlayerCmd:loopmode:"+mode, nil)
}

// Pause отправляет команду на паузу плеера.
func (p *PlayBackApi) Pause() error {
	return p.api.DoAPIRequest("GET", "setPlayerCmd:pause", nil)
}

// Resume отправляет команду на возобновление воспроизведения плеера.
func (p *PlayBackApi) Resume() error {
	err := p.api.DoAPIRequest("GET", "setPlayerCmd:resume", nil)
	if err != nil {
		return err
	}
	return nil
}

// OnePause отправляет команду на переключение состояния паузы плеера.
func (p *PlayBackApi) OnePause() error {
	err := p.api.DoAPIRequest("GET", "setPlayerCmd:onepause", nil)
	if err != nil {
		return err
	}
	return nil
}

// Stop отправляет команду на остановку плеера.
func (p *PlayBackApi) Stop() error {
	err := p.api.DoAPIRequest("GET", "setPlayerCmd:stop", nil)
	if err != nil {
		return err
	}
	return nil
}

// Prev отправляет команду на переход к предыдущему треку в плеере.
func (p *PlayBackApi) Prev() error {
	err := p.api.DoAPIRequest("GET", "setPlayerCmd:prev", nil)
	if err != nil {
		return err
	}
	return nil
}

// Next отправляет команду на переход к следующему треку в плеере.
func (p *PlayBackApi) Next() error {
	err := p.api.DoAPIRequest("GET", "setPlayerCmd:next", nil)
	if err != nil {
		return err
	}
	return nil
}

// setMute устанавливает состояние звука (включен/выключен) плеера.
func (p *PlayBackApi) setMute(mute bool) error {
	val := "1"
	if !mute {
		val = "0"
	}
	err := p.api.DoAPIRequest("GET", "setPlayerCmd:mute:"+val, nil)
	if err != nil {
		return err
	}
	return err
}

// Seeking отправляет команду на перемотку плеера на указанную позицию в миллисекундах.
func (p *PlayBackApi) Seeking(position int) error {
	if position < 0 {
		return fmt.Errorf("недопустимая позиция: %d", position)
	}
	// Отправляем команду на перемотку плеера
	err := p.api.DoAPIRequest("GET", fmt.Sprintf("setPlayerCmd:seek:%d", position), nil)
	if err != nil {
		return err
	}
	return nil
}

// SeekBack отправляет команду на перемотку плеера назад на указанное количество миллисекунд.
func (p *PlayBackApi) SeekBack(milliseconds int) error {
	// получаем текущий статус плеера
	if milliseconds < 0 {
		return fmt.Errorf("недопустимое количество миллисекунд: %d", milliseconds)
	}
	status, err := p.getPlayerStatus()
	if err != nil {
		return err
	}
	// вычисляем новую позицию
	currentPos := status.CurPos
	currentPosInt, err := strconv.Atoi(currentPos)
	if err != nil {
		return err
	}
	newPos := currentPosInt - milliseconds
	if newPos < 0 {
		newPos = 0 // не допускаем отрицательную позицию
	}
	// отправляем команду на перемотку плеера назад
	err = p.api.DoAPIRequest("GET", fmt.Sprintf("setPlayerCmd:seek:%d", newPos), nil)
	if err != nil {
		return err
	}
	return nil
}

// SeekForward отправляет команду на перемотку плеера вперед на указанное количество миллисекунд.
func (p *PlayBackApi) SeekForward(milliseconds int) error {
	// получаем текущий статус плеера
	if milliseconds < 0 {
		return fmt.Errorf("недопустимое количество миллисекунд: %d", milliseconds)
	}
	status, err := p.getPlayerStatus()
	if err != nil {
		return err
	}
	// вычисляем новую позицию
	currentPos := status.CurPos
	// преобразуем текущую позицию в целое число
	currentPosInt, err := strconv.Atoi(currentPos)
	if err != nil {
		return err
	}
	newPos := currentPosInt - milliseconds
	totalLen, err := strconv.Atoi(status.TotLen)
	if newPos > totalLen {
		newPos = totalLen // не допускаем превышение общей длины трека
	}
	// отправляем команду на перемотку плеера вперед
	err = p.api.DoAPIRequest("GET", fmt.Sprintf("setPlayerCmd:seek:%d", newPos), nil)
	if err != nil {
		return err
	}
	return nil
}

// SetVolume устанавливает громкость плеера в процентах (от 0 до 100).
func (p *PlayBackApi) SetVolume(volume string) error {
	intVolume, err := strconv.Atoi(volume)
	if err != nil {
		return fmt.Errorf("недопустимое значение громкости: %s, должно быть числом", volume)
	}
	if intVolume < 0 || intVolume > 100 {
		return fmt.Errorf("недопустимая громкость: %d, должна быть в диапазоне от 0 до 100", volume)
	}
	// Отправляем команду на установку громкости плеера
	err = p.api.DoAPIRequest("GET", fmt.Sprintf("setPlayerCmd:vol:%s", volume), nil)
	if err != nil {
		return err
	}
	return nil
}

// VolumeUp увеличивает громкость плеера на какойто % в плеере зашито значение.
func (p *PlayBackApi) VolumeUp() error {
	err := p.api.DoAPIRequest("GET", "setPlayerCmd:vol%2b%2b", nil)
	if err != nil {
		return err
	}
	return nil
}

// VolumeDown уменьшает громкость плеера на какойто % в плеере зашито значение.
func (p *PlayBackApi) VolumeDown() error {
	err := p.api.DoAPIRequest("GET", "setPlayerCmd:vol--", nil)
	if err != nil {
		return err
	}
	return nil
}

// Mute переключает состояние звука плеера между включенным и выключенным.
func (p *PlayBackApi) Mute() error {
	status, err := p.getPlayerStatus()
	if err != nil {
		return err
	}
	return p.setMute(status.Mute == "0")
}

// GetPlaylistTrackCount возвращает количество треков в текущем плейлисте.
func (p *PlayBackApi) GetPlaylistTrackCount() (int, error) {
	var count int
	err := p.api.DoAPIRequest("GET", "GetTrackNumber", &count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
