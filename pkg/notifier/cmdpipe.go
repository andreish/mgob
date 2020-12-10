package notifier

import (
	"fmt"
	"io"
	"os/exec"
	log "github.com/sirupsen/logrus"
	"github.com/stefanprodan/mgob/pkg/config"
)

type CmdPipeNotificator struct {
	*config.CmdPipe
}
 

func (config *CmdPipeNotificator) sendNotification(planID string, subject string, body string, warn bool) error {

	cmd := exec.Command(config.Command, config.Args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.WithField("plan", planID).Error(err)
	}

	go func() {
		defer stdin.Close()
		pipedata := fmt.Sprintf("{\"text\": \"%s: %s\"}", subject, body)
		io.WriteString(stdin, pipedata)
	}()


	result, err := cmd.CombinedOutput() 
	log.WithField("plan", planID).Errorf("On demand backup cmdpipe <%s> notification output:\n%s", config.Command, result)

	if( err != nil )	{
		log.WithField("plan", planID).Errorf("On demand backup cmdpipe notification error:\n%s", err)
	}
	return err 
}