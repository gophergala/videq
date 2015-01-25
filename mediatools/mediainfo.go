package mediatools

import (
	"errors"
	"fmt"
	alog "github.com/cenkalti/log"
	"github.com/codeskyblue/go-sh"
	"github.com/gophergala/videq/config"
	"os"
	"path/filepath"
	"strconv"
	//	"strconv"
	"strings"
	"time"
)

type MediaInfo struct {
	log         alog.Logger
	config      config.Config
	resolutions map[string]VideoResolution
}

func NewMediaInfo(log alog.Logger) *MediaInfo {
	m := new(MediaInfo)
	m.log = log
	m.resolutions = resolutions
	return m
}

// http://en.wikipedia.org/wiki/H.264/MPEG-4_AVC
// Commonly supported resolutions and aspect ratios include:
// 854 x 480 (16:9 480p)
// 1280 x 720 (16:9 720p)
// 1920 x 1080 (16:9 1080p)
// 640 x 480 (4:3 480p)
// 1280 x 1024 (5:4)
// 1920 x 1440 (4:3)

var resolutions = map[string]VideoResolution{
	"854 x 480 (16:9 480p)":    {Width: 854, Height: 480, AspectRatio: "16:9", AspectRatioInt: float64(16) / float64(9), Short: "480p"},
	"1280 x 720 (16:9 720p)":   {Width: 1280, Height: 720, AspectRatio: "16:9", AspectRatioInt: float64(16) / float64(9), Short: "720p"},
	"1920 x 1080 (16:9 1080p)": {Width: 1920, Height: 1080, AspectRatio: "16:9", AspectRatioInt: float64(16) / float64(9), Short: "1080p"},
	"640 x 480 (4:3 480p)":     {Width: 640, Height: 480, AspectRatio: "4:3", AspectRatioInt: float64(4) / float64(3), Short: "480p"},
	"1280 x 1024 (5:4)":        {Width: 1280, Height: 1024, AspectRatio: "5:4", AspectRatioInt: float64(4) / float64(3), Short: ""},
	"1920 x 1440 (4:3)":        {Width: 1920, Height: 1440, AspectRatio: "4:3", AspectRatioInt: float64(4) / float64(3), Short: ""},
}

type VideoResolution struct {
	Width          int     `json:"width"`
	Height         int     `json:"height"`
	AspectRatio    string  `json:"aspectratio"`
	AspectRatioInt float64 `json:"aspectratiofloat"`
	Short          string  `json:"short"` // shorthand name for a family of video display resolutions
}

type MediaFileInfo struct {
	FileName        string
	FileSize_bytes  string
	VideoCount      int
	AudioCount      int
	Duration_ms     string
	Duration        time.Duration
	Duration_string string
	Format          string
	CodecID         string
	Resolution      string
	Width           string
	Height          string
	Standard        string
	Codec           string
	Bitrate_bps     string
	Framerate       string
	AspectRatio     string
	Audio           string
}

/*
FileName: r2w_1080p.mov
FileSize_bytes: 104937987
VideoCount: 1
AudioCount: 1
Duration_ms: 90125
Format: MPEG-4
CodecID: qt
Resolution: 1920x816
Width: 1920
Height: 816
Standard:
Codec: AVC Main@L4.0
Bitrate_bps: 9185470
Framerate: 24.000 fps
AspectRatio: 2.35:1
Audio: English 128 Kbps CBR 2 chnls AAC LC
*/

/*
fetches media info
usage:
mt := mediatools.NewMediaInfo(log)
minfo, err := mt.GetMediaInfo("_test/master_1080.mp4")
if err != nil {
	log.Fatal(err)
}
log.Infof("%#v", minfo)
*/
func (m *MediaInfo) GetMediaInfo(fileName string) (fileInfo MediaFileInfo, err error) {

	// timeout should be a session
	//	out, err := sh.Command("ping", "-t", "127.0.0.1").SetTimeout(time.Second * 60).Output()

	mediaInfoParams := `--Inform=General;FileName:: %FileName%.%FileExtension%\r\nFileSize_bytes:: %FileSize%\r\nVideoCount:: %VideoCount%\r\nAudioCount:: %AudioCount%\r\nDuration_ms:: %Duration%\r\nFormat:: %Format%\r\nCodecID:: %CodecID%\r\n
Video;Resolution:: %Width%x%Height%\r\nWidth:: %Width%\r\nHeight:: %Height%\r\nStandard:: %Standard%\r\nCodec:: %Codec/String% %Format_Profile%\r\nBitrate_bps:: %BitRate%\r\nFramerate:: %FrameRate% fps\r\nAspectRatio:: %DisplayAspectRatio/String%\r\n
Audio;Audio:: %Language/String% %BitRate/String% %BitRate_Mode% %Channel(s)% chnls %Codec/String%\r\n
Text;%Language/String%
Text_Begin;Subs:
Text_Middle;, 
Text_End;.\r\n
`
	out, err := sh.Command("mediainfo", mediaInfoParams, fileName).SetTimeout(time.Second * 60).Output()
	// fmt.Printf("output:(%s), err(%v)\n", string(out), err)
	if err == sh.ErrExecTimeout {
		m.log.Errorf("shell exec timeouteded.", err)
	}
	if err != nil {
		m.log.Errorf("sh.Command error:", err)
		return fileInfo, err
	}

	//output = strings.Replace(string(out), "\n", "<br>", -1)

	lines := strings.Split(string(out), "\n")
	//fileInfo := new(MediaFileInfo)

	for _, line := range lines {
		paramArr := strings.Split(line, "::")
		if len(paramArr) != 2 {
			continue
		}

		paramName := strings.Trim(paramArr[0], " ")
		paramValue := strings.Trim(paramArr[1], " ")

		switch paramName {
		case `FileName`:
			fileInfo.FileName = paramValue
		case `FileSize_bytes`:
			fileInfo.FileSize_bytes = paramValue
		case `VideoCount`:
			val, err := strconv.ParseInt(paramValue, 10, 0) // int
			if err == nil {
				fileInfo.VideoCount = int(val)
			} else {
				fileInfo.VideoCount = 0
			}
		case `AudioCount`:
			val, err := strconv.ParseInt(paramValue, 10, 0) // int
			if err == nil {
				fileInfo.AudioCount = int(val)
			} else {
				fileInfo.AudioCount = 0
			}

		case `Duration_ms`:
			fileInfo.Duration_ms = paramValue
			// durationDurationInt64, err := strconv.ParseInt(fileInfo.Duration_ms, 10, 64)
			// fileInfo.Duration = time.Millisecond * durationDurationInt64 / 100
			dur, err := time.ParseDuration(paramValue + "ms")
			if err == nil {
				fileInfo.Duration = dur
				fileInfo.Duration_string = fmt.Sprintf("%s", dur)
			} else {
				m.log.Error(dur)
			}
		case `Format`:
			fileInfo.Format = paramValue
		case `CodecID`:
			fileInfo.CodecID = paramValue
		case `Resolution`:
			fileInfo.Resolution = paramValue
		case `Width`:
			fileInfo.Width = paramValue
		case `Height`:
			fileInfo.Height = paramValue
		case `Standard`:
			fileInfo.Standard = paramValue
		case `Codec`:
			fileInfo.Codec = paramValue
		case `Bitrate_bps`:
			fileInfo.Bitrate_bps = paramValue
		case `Framerate`:
			fileInfo.Framerate = paramValue
		case `AspectRatio`:
			fileInfo.AspectRatio = paramValue
		case `Audio`:
			fileInfo.Audio = paramValue
		}
	}

	return fileInfo, nil
}

// CheckMedia checks if video file is ok for encoding
// TODO: more checks
// 1. limit output resolutions ONLY to the same of original video or LOWER ones
// 2. input format check

func (m *MediaInfo) CheckMedia(fileName string) (ok bool, fileInfo MediaFileInfo, res map[string]VideoResolution, err error) {
	res = m.resolutions

	exists, err := m.checkIfFileExists(fileName)
	if err != nil {
		return false, fileInfo, nil, err
	}
	if exists == false {
		return false, fileInfo, nil, errors.New(fmt.Sprintf("File '%s' does not exists.", fileName))
	}

	fileInfo, err = m.GetMediaInfo(fileName)
	if err != nil {
		m.log.Error(err)
		return false, fileInfo, nil, err
	}

	//maxDuration := time.Second * 60
	maxDuration := time.Minute * 5
	if fileInfo.Duration > maxDuration {
		m.log.Infoln(fileInfo.Duration)
		return false, fileInfo, nil, errors.New(fmt.Sprintf("File '%s' is too long. Max duration: %s, File duration: %s", fileName, maxDuration, fileInfo.Duration_string))
	}

	if fileInfo.VideoCount == 0 {
		return false, fileInfo, nil, errors.New(fmt.Sprintf("File '%s' is no video.", fileName))
	}

	// TODO - more meaningfull checks

	return true, fileInfo, m.resolutions, nil
}

func (m *MediaInfo) checkIfFileExists(fileName string) (bool, error) {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (m *MediaInfo) returnBaseFilename(fileName string) string {
	// var extension = filepath.Ext(fileName)
	// var name = fileName[0 : len(fileName)-len(extension)]
	// return name

	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func (m *MediaInfo) deleteFile(file string) (err error) {
	// check if file/folder exists
	if _, err := os.Stat(file); !os.IsNotExist(err) {
		// vrati errror samo ako je neki zesci error, ne na file not found
		return err
	}
	return nil
}

// "bytes"
// "encoding/json"

// m.log.Debug(m.resolutions)
// b, err := json.Marshal(m.resolutions)
// if err != nil {
// 	m.log.Fatal(err)
// }
// var out bytes.Buffer
// json.Indent(&out, b, "=", "\t")
// out.WriteTo(os.Stdout)
