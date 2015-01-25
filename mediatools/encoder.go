package mediatools

import (
	"errors"
	"fmt"
	"github.com/codeskyblue/go-sh"
)

func (m *MediaInfo) EncodeVideoFile(fileLoc string, fileName string) (err error) {
	file := fileLoc + fileName
	m.log.Infoln(file)

	outFileMP4, err := m.encodeMP4(fileLoc, fileName)
	if err != nil {
		// ako je mp4 puko kod encoding, ne mozemo nastaviti jer svi ostali zele taj file kao source
		return err
	}
	outFileOGG, err := m.encodeOGG(fileLoc, outFileMP4)
	outFileWEBM, err := m.encodeWEBM(fileLoc, outFileMP4)

	outFileJPG, err := m.generateThumbnailJPG(fileLoc, outFileMP4)

	m.log.Debug("Created: ", outFileMP4, outFileOGG, outFileWEBM, outFileJPG)
	return nil
}

//
// HandBrakeCLI -i _test/master_1080.mp4 -o _test/out/master_1080_2.mp4 -e x264 -q 22 -r 15 -B 64 -X 480 -O -x level=4.0:ref=9:bframes=16:b-adapt=2:direct=auto:analyse=all:8x8dct=0:me=tesa:merange=24:subme=11:trellis=2:fast-pskip=0:vbv-bufsize=25000:vbv-maxrate=20000:rc-lookahead=60
// http://thanhsiang.org/faqing/node/196
/*
alternativno UMEJESTO handbrakea
These are the options i used to convert to a H.264/AAC .mp4 format for html5 playback (i think it may help out for other guys with this this problem in some way):
ffmpeg -i input.flv -vcodec mpeg4 -acodec aac output.mp4
*/
func (m *MediaInfo) encodeMP4(fileLoc string, fileName string) (fileNameOut string, err error) {
	fileSource := fileLoc + fileName
	//fileNameOut = m.returnBaseFilename(fileName) + ".mp4"
	fileNameOut = "encoded.mp4"
	fileDestination := fileLoc + fileNameOut

	maxWidth := "1280"
	extraParams := `level=4.0:ref=9:bframes=16:b-adapt=2:direct=auto:analyse=all:8x8dct=0:me=tesa:merange=24:subme=11:trellis=2:fast-pskip=0:vbv-bufsize=25000:vbv-maxrate=20000:rc-lookahead=60`

	// test fast preset
	//	out, err := sh.Command("HandBrakeCLI", "-i", fileSource, "-o", fileDestination, "-e", "x264", "-q", "22", "-r", "15", "-B", "64", "-X", maxWidth, "-O", "-x", extraParams).Output()

	// alen's web master preset
	out, err := sh.Command("HandBrakeCLI", "-i", fileSource, "-o", fileDestination,
		"-e", "x264", "-b", "500", "-B", "64", "-X", maxWidth, "--keep-display-aspect", "--two-pass", "-O", "-x", extraParams).Output()
	if err == sh.ErrExecTimeout {
		m.log.Errorf("shell exec timeouteded.", err)
	}
	if err != nil {
		m.log.Errorf("sh.Command error:", err)
		return "", err
	}

	m.log.Debugf("output:(%s), err(%v)\n", string(out), err)

	ok, err := m.checkIfFileExists(fileDestination)
	if err != nil {
		return "", err
	}
	if ok == false {
		return "", errors.New(fmt.Sprintf("File '%s' does not exists. Encoding failed?", fileDestination))
	}

	return fileNameOut, nil

}

//
// ffmpeg2theora Master_1080.mp4 --two pass --videobitrate 900 -x 1280 -y 720
//
func (m *MediaInfo) encodeOGG(fileLoc string, fileName string) (fileNameOut string, err error) {
	fileSource := fileLoc + fileName
	// fileNameOut = m.returnBaseFilename(fileName) + ".ogg"
	fileNameOut = "encoded.ogg"
	fileDestination := fileLoc + fileNameOut

	maxWidth := "1280"
	maxHeight := "720"

	out, err := sh.Command("ffmpeg2theora", fileSource, "-o", fileDestination, "--two pass", "--videobitrate", "900", "-x", maxWidth, "-y", maxHeight).Output()
	if err == sh.ErrExecTimeout {
		m.log.Errorf("shell exec timeouteded.", err)
	}
	if err != nil {
		m.log.Errorf("sh.Command error:", err)
		return "", err
	}

	//m.log.Debug(out)
	m.log.Debugf("output:(%s), err(%v)\n", string(out), err)
	ok, err := m.checkIfFileExists(fileDestination)
	if err != nil {
		return "", err
	}
	if ok == false {
		return "", errors.New(fmt.Sprintf("File '%s' does not exists. Encoding failed?", fileDestination))
	}

	return fileNameOut, nil
}

/*
ffmpeg -i _test/master_1080.mp4 -pass 1 -passlogfile hattrick.webm -keyint_min 0 -g 250 -skip_threshold 0 -vcodec libvpx -b 600k -s 1280x720 -aspect 16:9 -an -y hattrick.webm
Output file is empty, nothing was encoded (check -ss / -t / -frames parameters if used)

ffmpeg -i _test/master_1080.mp4 -pass 2 -passlogfile hattrick.webm -keyint_min 0 -g 250 -skip_threshold 0 -vcodec libvpx -b 600k -s 1280x720 -aspect 16:9 -acodec libvorbis -y hattrick.webm
*/

func (m *MediaInfo) encodeWEBM(fileLoc string, fileName string) (fileNameOut string, err error) {
	fileSource := fileLoc + fileName
	// fileNameOut = m.returnBaseFilename(fileName) + ".webm"
	fileNameOut = "encoded.webm"
	fileDestination := fileLoc + fileNameOut

	// out, err := sh.
	// 	Command("ffmpeg", "-i", fileSource, "-pass", "1", "-passlogfile", fileDestination, "-keyint_min", "0", "-g", "250", "-skip_threshold", "0", "-vcodec", "libvpx", "-b", "600k", "-s", "1280x720", "-aspect", "16:9", "-an", "-y", fileDestination).
	// 	Command("ffmpeg", "-i", fileSource, "-pass", "2", "-passlogfile", fileDestination, "-keyint_min", "0", "-g", "250", "-skip_threshold", "0", "-vcodec", "libvpx", "-b", "600k", "-s", "1280x720", "-aspect", "16:9", "-acodec", "libvorbis", "-y", fileDestination).
	// 	Run()
	out, err := sh.
		Command("ffmpeg", "-i", fileSource, "-pass", "1", "-passlogfile", fileDestination, "-keyint_min", "0", "-g", "250", "-skip_threshold", "0", "-vcodec", "libvpx", "-b", "600k", "-s", "1280x720", "-aspect", "16:9", "-an", "-y", fileDestination).
		Output()

	if err == sh.ErrExecTimeout {
		m.log.Errorf("shell exec timeouteded.", err)
	}
	if err != nil {
		m.log.Errorf("sh.Command error:", err)
		return "", err
	}

	out, err = sh.
		Command("ffmpeg", "-i", fileSource, "-pass", "2", "-passlogfile", fileDestination, "-keyint_min", "0", "-g", "250", "-skip_threshold", "0", "-vcodec", "libvpx", "-b", "600k", "-s", "1280x720", "-aspect", "16:9", "-acodec", "libvorbis", "-y", fileDestination).
		Output()

	if err == sh.ErrExecTimeout {
		m.log.Errorf("shell exec timeouteded.", err)
	}
	if err != nil {
		m.log.Errorf("sh.Command error:", err)
		return "", err
	}

	//m.log.Debug(out)
	m.log.Debugf("output:(%s), err(%v)\n", string(out), err)
	ok, err := m.checkIfFileExists(fileDestination)
	if err != nil {
		return "", err
	}
	if ok == false {
		return "", errors.New(fmt.Sprintf("File '%s' does not exists. Encoding failed?", fileDestination))
	}

	return fileNameOut, nil
}

// https://trac.ffmpeg.org/wiki/Create%20a%20thumbnail%20image%20every%20X%20seconds%20of%20the%20video
// ffmpeg -i master_1080.mp4 -t 0.001 -ss 7 -vframes 1 -y -f mjpeg master_test.jpg
// ffmpeg -i input.flv -ss 00:00:14.435 -f image2 -vframes 1 out.png

func (m *MediaInfo) generateThumbnailJPG(fileLoc string, fileName string) (fileNameOut string, err error) {
	fileSource := fileLoc + fileName
	// fileNameOut = m.returnBaseFilename(fileName) + ".webm"
	fileNameOut = "encoded.jpg"
	fileDestination := fileLoc + fileNameOut

	out, err := sh.Command("ffmpeg", "-i", fileSource, "-t", "0.001", "-ss", "7", "-vframes", "1", "-y", "-f", "mjpeg", fileDestination).Output()
	if err == sh.ErrExecTimeout {
		m.log.Errorf("shell exec timeouteded.", err)
	}
	if err != nil {
		m.log.Errorf("sh.Command error:", err)
		return "", err
	}

	//m.log.Debug(out)
	m.log.Debugf("output:(%s), err(%v)\n", string(out), err)
	ok, err := m.checkIfFileExists(fileDestination)
	if err != nil {
		return "", err
	}
	if ok == false {
		return "", errors.New(fmt.Sprintf("File '%s' does not exists. Encoding failed?", fileDestination))
	}

	return fileNameOut, nil
}
