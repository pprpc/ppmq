package topic

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pprpc/util/logs"
)

func Match(t1, t2 string) (code uint32) {
	if t1 == t2 {
		code = 3
		return
	}

	_t1 := removeEmpty(strings.Split(t1, "/"))
	_t2 := removeEmpty(strings.Split(t2, "/"))

	_lt1 := len(_t1)
	_lt2 := len(_t2)
	if _lt1 < _lt2 {
		for i, v := range _t1 {
			logs.Logger.Debugf("v: %s, t2[%d]: %s.", v, i, _t2[i])
			if v != _t2[i] {
				code = 0
				return
			}
		}
		code = 1
		return
	} else if _lt1 > _lt2 {
		for i, v := range _t2 {
			if v != _t1[i] {
				code = 0
				return
			}
		}
		code = 2
		return
	} else {
		for i, v := range _t2 {
			if v != _t1[i] {
				code = 0
				return
			}
		}
		code = 0
	}
	return
}

// Contains .
func Contains(t1 string, t2 []string) bool {
	var t uint32
	for _, row := range t2 {
		t = MatchV2(t1, row)
		if t > 1 {
			return true
		}
	}
	return false
}

// GetDiffTopic .
func GetDiffTopic(t1, t2 []string) (d []string) {
	for _, row := range t1 {
		if Contains(row, t2) == false {
			d = append(d, row)
		}
	}
	return d
}

func removeEmpty(a []string) (ret []string) {
	ret = a
	return
	for i, v := range a {
		if v == "" {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}

// MatchV2 Topic
// FIXME: not support "+"
func MatchV2(t1, t2 string) (code uint32) {
	if t1 == t2 {
		code = 3
		return
	}
	if t2 == "/#" || t2 == "/" || t2 == "#" {
		code = 2
		return
	}

	_t1 := strings.Split(t1, "/")
	_t2 := strings.Split(t2, "/")
	_lt1 := len(_t1)
	_lt2 := len(_t2)
	if _lt2 > _lt1 {
		// t1 = /a/b
		// t2 = /a/b/, /a/b/#
		for i, v := range _t1 {
			if v != _t2[i] {
				code = 0
				return
			}
		}
		if _t2[_lt2-1] == "" || _t2[_lt2-1] == "#" {
			code = 2
		}
		return
	} else if _lt2 == _lt1 {
		for i, v := range _t2 {
			if v != _t1[i] && i != _lt1-1 {
				code = 0
				return
			} else {
				if (v == "" || v == "#") && i == _lt1-1 {
					code = 2
				}
			}
		}
		return
	} else {
		for i, v := range _t2 {
			if v != _t1[i] && i != _lt2-1 {
				code = 0
				return
			} else if i == _lt2-1 {
				if v == "" || v == "#" {
					code = 2
				}
			}
		}
		return
	}
}

// CheckSubTopic check subscribe topic
func CheckSubTopic(topic string) (err error) {
	if topic == "" {
		err = fmt.Errorf("subscribe:  topic is null")
	}
	m, _ := regexp.MatchString("^[-._a-zA-Z0-9/]+[#]{0,1}$", topic)
	if m != true {
		err = fmt.Errorf("subscribe topic: %s, illegal", topic)
	}
	return
}

// CheckPubTopic check publish topic
func CheckPubTopic(topic string) (err error) {
	if topic == "" {
		err = fmt.Errorf("publish:  topic is null")
	}
	m, _ := regexp.MatchString("^[-._a-zA-Z0-9/]+[a-zA-Z0-9]{1}$", topic)
	if m != true {
		err = fmt.Errorf("publish topic: %s, illegal", topic)
	}
	return
}
