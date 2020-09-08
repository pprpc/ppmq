package msg

import (
	"fmt"

	"github.com/pprpc/util/common"
	"xcthings.com/ppmq/ppmqcli"
	"github.com/pprpc/core"
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
