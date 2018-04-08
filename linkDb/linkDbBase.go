package linkDb

import "fmt"
import "sync"
import pb "github.com/kimstoun/nms-proxy/pb"
import "errors"

type LinkKey struct {
	SendPortKey PortKey
	RecvPortKey PortKey
}
type PortKey struct {
	AppName  string
	PortName string
}

const (
	SENDPORT = 1 + iota
	RECVPORT
	DEPIPORT
)
const (
	LINKWAITBECONFIG = 1 + iota
	LINKCONFIGOK
	LINKCONFIGERROR
	PORTWAITBECONFIG
	PORTNOTEXIST
)
const (
	STATECHANGED = 1 + iota
	STATEUNCHANGED
)

var unLinkedPort map[PortKey]pb.PortParameter
var linkInfo map[LinkKey]pb.LinkParameter
var muDb sync.Mutex
var muSignal sync.Mutex
var WaitSignal *sync.Cond

func init() {
	unLinkedPort = make(map[PortKey]pb.PortParameter)
	linkInfo = make(map[LinkKey]pb.LinkParameter)
	WaitSignal = sync.NewCond(&muSignal)
	go checkUnLinkedPortRouting()
	go checkLinkInfoRouting()
}

//检测是否有符合条件的两个端口，把这两个端口移到linkinfo里面去
func ScanUnLinkedPort() uint32 {
	muDb.Lock()
	defer muDb.Unlock()
	for k, v := range unLinkedPort {
		pk := PortKey{v.RemoteAppName, v.RemotePortName}
		if pv, ok := unLinkedPort[pk]; ok {
			addPortsToLinkInfo(v, pv)
			deletFromUnLinkedPort(k)
			deletFromUnLinkedPort(pk)
			return STATECHANGED
		}
	}
	return STATEUNCHANGED
}

func ScanLinkInfo() uint32 {
	muDb.Lock()
	defer muDb.Unlock()
	state := STATEUNCHANGED
	for k, v := range linkInfo {
		if v.LinkState == LINKWAITBECONFIG {
			//TODO 配置链路
			//修改链路状态
			fmt.Println("链路配置成功\r\n")
			changeLinkState(k, LINKCONFIGOK)
			state = STATECHANGED
		} else if v.LinkState == LINKCONFIGERROR {
			deletLinkInfo(k)
			state = STATECHANGED
		} else {
			//正常状态
		}
	}
	return uint32(state)
}

func GetPortState(pk PortKey) uint32 {
	muDb.Lock()
	defer muDb.Unlock()
	state := uint32(0)
	if _, state = getPortFromUnLinkedPort(pk); state == PORTNOTEXIST {
		_, state = getPortFromLinkInfo(pk)
	}
	return state
}

func GetInfoByRioId(rioId int32) (map[PortKey]pb.PortParameter, map[LinkKey]pb.LinkParameter) {
	tUnLinked := make(map[PortKey]pb.PortParameter)
	tLinkInfo := make(map[LinkKey]pb.LinkParameter)
	for k, v := range unLinkedPort {
		if v.RioId == rioId {
			tUnLinked[k] = v
		}
	}

	for k, v := range linkInfo {
		if v.SendPort.RioId == rioId || v.RecvPort.RioId == rioId {
			tLinkInfo[k] = v
		}
	}
	return tUnLinked, tLinkInfo
}

func GetInfoByPortName(portName string) (map[PortKey]pb.PortParameter, map[LinkKey]pb.LinkParameter) {
	tUnLinked := make(map[PortKey]pb.PortParameter)
	tLinkInfo := make(map[LinkKey]pb.LinkParameter)
	for k, v := range unLinkedPort {
		if v.PortName == portName {
			tUnLinked[k] = v
		}
	}

	for k, v := range linkInfo {
		if v.SendPort.PortName == portName || v.RecvPort.PortName == portName {
			tLinkInfo[k] = v
		}
	}
	return tUnLinked, tLinkInfo

}

func GetInfoByAppName(appName string) (map[PortKey]pb.PortParameter, map[LinkKey]pb.LinkParameter) {
	tUnLinked := make(map[PortKey]pb.PortParameter)
	tLinkInfo := make(map[LinkKey]pb.LinkParameter)
	for k, v := range unLinkedPort {
		if v.AppName == appName {
			tUnLinked[k] = v
		}
	}

	for k, v := range linkInfo {
		if v.SendPort.AppName == appName || v.RecvPort.AppName == appName {
			tLinkInfo[k] = v
		}
	}
	return tUnLinked, tLinkInfo

}

func GetAllInfo() (map[PortKey]pb.PortParameter, map[LinkKey]pb.LinkParameter) {
	return unLinkedPort, linkInfo
}

/*****半连接库的基础操作*************/
//往未配对端口数据库中添加端口
func addToUnLinkedPort(pp pb.PortParameter) error {
	pk := PortKey{pp.AppName, pp.PortName}
	if v, ok := unLinkedPort[pk]; ok {
		fmt.Println(v)
		return errors.New("该端口已经存在")
	}
	unLinkedPort[pk] = pp
	return nil

}

func deletFromUnLinkedPort(pk PortKey) error {
	if _, ok := unLinkedPort[pk]; ok {
		delete(unLinkedPort, pk)
		return nil
	}
	return errors.New("该端口不存不存在")
}

func getPortFromUnLinkedPort(pk PortKey) (pb.PortParameter, uint32) {
	if v, ok := unLinkedPort[pk]; ok {
		return v, PORTWAITBECONFIG
	}
	return pb.PortParameter{}, PORTNOTEXIST

}

/********************全连接库的接触操作***************/
func addPortsToLinkInfo(pp1 pb.PortParameter, pp2 pb.PortParameter) error {
	lp := new(pb.LinkParameter)
	lp.SendPort = &pp1
	lp.RecvPort = &pp2
	lp.LinkState = LINKWAITBECONFIG
	lk := new(LinkKey)
	lk.SendPortKey.AppName = pp1.AppName
	lk.SendPortKey.PortName = pp1.PortName
	lk.RecvPortKey.AppName = pp2.AppName
	lk.RecvPortKey.PortName = pp2.PortName
	if _, ok := linkInfo[*lk]; ok {
		return errors.New("mvport 链路已经存在")
	}
	linkInfo[*lk] = *lp
	return nil

}

func deletLinkInfo(lk LinkKey) error {
	if _, ok := linkInfo[lk]; ok {
		delete(linkInfo, lk)
		return nil
	}
	return errors.New("deletLinkInfo 链路不存在")

}
func changeLinkState(lk LinkKey, linkState uint32) error {
	if v, ok := linkInfo[lk]; ok {
		v.LinkState = linkState
		linkInfo[lk] = v
		return nil
	}
	return errors.New("ChangeLinkState 没有找到链路")

}

func getPortFromLinkInfo(pk PortKey) (pb.LinkParameter, uint32) {
	for k, v := range linkInfo {
		if (k.SendPortKey == pk) || (k.RecvPortKey == pk) {
			return v, v.LinkState
		}

	}
	return pb.LinkParameter{}, PORTNOTEXIST
}

//打印未配对的端口
func printUnLinkedPort() {
	for _, v := range unLinkedPort {
		fmt.Println(v)
	}
}

//打印已配置好的链路信息
func printLinfInfo() error {
	fmt.Println(linkInfo)
	return nil
}
