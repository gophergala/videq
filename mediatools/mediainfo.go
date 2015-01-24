package mediatools

import (
	alog "github.com/cenkalti/log"
	"github.com/codeskyblue/go-sh"
	_ "strings"
	"time"
)

type MediaInfo struct {
	log alog.Logger
}

func NewMediaInfo(log alog.Logger) *MediaInfo {
	m := new(MediaInfo)
	m.log = log
	return m
}

// fetched media info
func (m *MediaInfo) GetMediaInfo(fileName string) (output string, err error) {

	// timeout should be a session
	//	out, err := sh.Command("ping", "-t", "127.0.0.1").SetTimeout(time.Second * 60).Output()

	out, err := sh.Command("mediainfo", fileName).SetTimeout(time.Second * 60).Output()
	// fmt.Printf("output:(%s), err(%v)\n", string(out), err)
	if err == sh.ErrExecTimeout {
		m.log.Errorf("shell exec timeouteded.", err)
	}
	if err != nil {
		m.log.Errorf("sh.Command error:", err)
		return "", err
	}

	m.log.Infof("%s", string(out))

	//output = strings.Replace(string(out), "\n", "<br>", -1)

	//lines := strings.Split(string(out), "\n")
	//m.log.Infof("%s", len(lines))
	//m.log.Infof("%#v", lines)
	// m.log.Info("%#v", len(lines))

	// for _, line := range lines {
	// 	m.log.Infoln(line)
	// }

	return output, nil
}
