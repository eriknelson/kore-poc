package comm

import (
	log "github.com/sirupsen/logrus"
)

type CmdDelegate struct {
	IngressMessage IngressMessage
	Submatches     []string

	response string
}

func NewCmdDelegate(im IngressMessage, subm []string) CmdDelegate {
	return CmdDelegate{
		IngressMessage: im,
		Submatches:     subm,
		response:       "",
	}
}

func (d *CmdDelegate) SendResponse(response string) {
	log.Debugf("CmdDelegate::SendResponse: %s", response)
	d.response = response
}
