/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
*	 ____   ___   ____ _  _______ ___   ___  _
*	|  _ \ / _ \ / ___| |/ /_   _/ _ \ / _ \| |
*	| | | | | | | |   | ' /  | || | | | | | | |
*	| |_| | |_| | |___| . \  | || |_| | |_| | |___
*	|____/ \___/ \____|_|\_\ |_| \___/ \___/|_____|
*
*	https://github.com/yingzhuo/docktool
* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */
package cmd

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/subchen/go-cli"
	"github.com/yingzhuo/docktool/cnf"
	"github.com/yingzhuo/docktool/value"
	jtcp "github.com/yingzhuo/jing/tcp"
)

// 结果集合
var collection = newWaitingResultCollection()
var queue chan int

func ActionWait(c *cli.Context) {

	defer close(queue)

	cnf.InitGlobalConfig()

	logrus.Debugf("command: \"%v\"", c.Name())
	logrus.Debugf("pwd: \"%v\"", cnf.GlobalPWD)
	logrus.Debugf("binary dir: \"%v\"", cnf.GlobalBinaryDir)

	if c.NArg() != 0 {
		panic("too many parameters for sub-command wait")
	}

	list := getList()
	count := len(list)
	timeoutFlag := false
	queue = make(chan int, count)

	for threadId, addr := range list {
		logrus.Debugf("waiting: \"%v\"", addr)
		go doWait(threadId, addr, &timeoutFlag)
	}

	go doTimeout(cnf.WaitTimeout.Get(), &timeoutFlag)

	s1 := time.Now().UnixNano()
	for {
		select {
		case threadId := <-queue:
			result := collection.get(threadId)
			Println(result.String())

			if cnf.WaitLogic == "ANY" {
				goto t1
			}

			if collection.size() == count {
				goto t1
			}
		}
	}

t1:

	if cnf.GlobalDebugMode {
		s2 := time.Now().UnixNano()
		cost := time.Duration(s2 - s1)
		logrus.Debugf("cost: %v", cost)
	}
}

const (
	ok      = "ok"
	timeout = "timeout"
)

func getList() value.WaitList {
	list := cnf.WaitList

	// 读取环境变量
	for _, kv := range os.Environ() {
		if strings.HasPrefix(kv, cnf.WaitEnvPrefix) {
			if parts := strings.Split(kv, "="); len(parts) == 2 {
				list.Add(parts[1])
			}
		}
	}

	ret := value.WaitList{}
	dict := make(map[string]bool)

	// 去除重复
	for _, it := range list {
		if !dict[it] {
			ret.Add(it)
			dict[it] = true
		}
	}
	return ret
}

func doWait(threadId int, addr string, timeoutFlag *bool) {
	result := newWaitingResult(threadId, addr)

	defer func() {
		collection.add(result)
		queue <- threadId
	}()

	for {
		if jtcp.IsReachable(addr) {
			result.status = ok
			return
		} else {
			if *timeoutFlag {
				result.status = timeout
				return
			} else {
				nap()
			}
		}
	}
}

func nap() {
	time.Sleep(time.Millisecond * 200)
}

func doTimeout(timeout time.Duration, quitVar *bool) {

	if timeout <= 0 || *quitVar {
		return
	}

	select {
	case <-time.After(timeout):
		*quitVar = true
	}
}

// ---------------------------

type waitingResult struct {
	threadId string
	addr     string
	status   string
}

func newWaitingResult(threadId interface{}, addr string) *waitingResult {
	return &waitingResult{
		threadId: fmt.Sprintf("%v", threadId),
		addr:     addr,
		status:   "",
	}
}

func (e *waitingResult) String() string {
	return fmt.Sprintf("%-8s: %s", e.status, e.addr)
}

// ---------------------------

type waitingResultCollection struct {
	mutex *sync.Mutex
	dict  map[string]*waitingResult
}

func (e *waitingResultCollection) add(ele *waitingResult) {
	defer e.mutex.Unlock()
	e.mutex.Lock()
	e.dict[ele.threadId] = ele
}

func (e *waitingResultCollection) size() int {
	defer e.mutex.Unlock()
	e.mutex.Lock()
	return len(e.dict)
}

func (e *waitingResultCollection) get(threadId interface{}) *waitingResult {
	defer e.mutex.Unlock()
	e.mutex.Lock()
	return e.dict[fmt.Sprintf("%v", threadId)]
}

func newWaitingResultCollection() *waitingResultCollection {
	return &waitingResultCollection{
		mutex: &sync.Mutex{},
		dict:  map[string]*waitingResult{},
	}
}
