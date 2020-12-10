package notifier

import (
	"strings"
	"fmt"
	"github.com/stefanprodan/mgob/pkg/config"
)

func SendNotification( planID string, subject string, body string, warn bool, config config.Plan) error {

	var errSMTP error
	var errSlack error
	var errCommands []error

	if config.SMTP != nil {
		n := SmtpNotificator{ config.SMTP }
		errSMTP = n.sendNotification(planID, subject, body , warn)
	}
	if config.Slack != nil {
		n := SlackNotificator{ config.Slack }
		errSlack = n.sendNotification(planID, subject, body, warn)		
	}
	cmdError := false
	if len := len(config.CmdPipe) ; len>0  {
		errCommands = make([]error, len)
		for i , cmdPipe := range(config.CmdPipe) {
			n := CmdPipeNotificator{cmdPipe}
			errCommands[i] = n.sendNotification(planID, subject, body, warn)		
			cmdError = cmdError || (errCommands[i] != nil)
		}
	}
	if nil != errSMTP || nil != errSlack || cmdError {
		return wrappErorrs(errSMTP, errSlack, errCommands)
	}
	return nil
}


// wrappErorrs assumes we hava at least one not nil error passed as parameter
func wrappErorrs( errSmtp error, errSlack error, errCommands []error ) error {


	var errstrings []string 

	if nil != errSmtp {
		errstrings = append(errstrings, errSmtp.Error())
	}

	if nil != errSlack {
		errstrings = append(errstrings, errSlack.Error())
	}

	for _ , e := range(errCommands) {
		errstrings = append(errstrings, e.Error())
	}

    return fmt.Errorf(strings.Join(errstrings, "\n"))
}