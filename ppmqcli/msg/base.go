package msg

import (
	"fmt"

	"xcthings.com/hjyz/common"
	"xcthings.com/ppmq/ppmqcli"
	"xcthings.com/pprpc"
)

type PPMQMsg struct {
	Cli *ppmqcli.PpmqCli
}

func New(cli *ppmqcli.PpmqCli) (msg *PPMQMsg) {
	msg = new(PPMQMsg)
	msg.Cli = cli
	return
}

//
func getMsgid(pre string) string {
	return fmt.Sprintf("%s-%d-%d", pre, pprpc.GetSeqID(), common.GetTimeMs())
}
