package todo

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	pfxAuthorName = "author "
	pfxAuthorMail = "author-mail "
	pfxTime       = "author-time "
)

func Blame(path string, lineNo int) (BlameInfo, error) {
	var blameInfo BlameInfo
	lineArg := fmt.Sprintf("-L %d,%d", lineNo, lineNo)
	args := []string{
		"blame",
		"--porcelain",
		"--incremental",
		lineArg,
		path,
	}

	cmd := exec.Command("git", args...)

	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		return blameInfo, err
	}
	if err := cmd.Start(); err != nil {
		return blameInfo, err
	}

	scanner := bufio.NewScanner(bufio.NewReader(stdOut))
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, pfxAuthorMail):
			value := strings.TrimPrefix(line, pfxAuthorMail)
			blameInfo.Mail = strings.Trim(value, "<>")
		case strings.HasPrefix(line, pfxAuthorName):
			blameInfo.Name = strings.TrimPrefix(line, pfxAuthorName)
		case strings.HasPrefix(line, pfxTime):
			value := strings.TrimPrefix(line, pfxTime)
			unixTime, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return blameInfo, err
			}
			blameInfo.Time = time.Unix(unixTime, 0)
		}
	}
	return blameInfo, nil
}
