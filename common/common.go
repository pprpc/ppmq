package common

import "strconv"

const (
	SESSIONKEEP  int32 = 0
	SESSIONCLEAR int32 = 1

	PINGNORMAL int32 = 0
	PINGSLEEP  int32 = 1

	ConsumeNew       int32 = 1
	ConsumeOld       int32 = 2
	ConsumeNewAndOld int32 = 3

	ONLINE  int32 = 1
	OFFLINE int32 = 2

	FIRSTSEND  int32 = 0
	REPEATSEND int32 = 1

	//MSGSTATUSClearSession  ClearSession = 1 时，相关消息的状态
	MSGSTATUSClearSession int32 = 88
)

// GetUserID .
func GetUserID(account string) int64 {
	t, err := strconv.Atoi(account)
	if err != nil {
		return 0
	} else {
		return int64(t)
	}
}
