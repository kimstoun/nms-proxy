package linkDb

import pb "gaoyl/pb"
import "fmt"
import "errors"

var testAddNum int32
var testCheckUn int32
var testCheckLn int32
var testWaitNum int32

func AddToUnLinkedPort(pp pb.PortParameter) error {
	muDb.Lock()
	defer muDb.Unlock()
	testAddNum = testAddNum + 1
	fmt.Printf("testAddNum = %v\r\n", testAddNum)
	pk := PortKey{pp.AppName, pp.PortName}
	//先查询在全连接库里面是否已经存在
	if _, state := getPortFromLinkInfo(pk); state != PORTNOTEXIST {
		return errors.New("在全连接库里面已经有该连接\r\n")
	}

	err := addToUnLinkedPort(pp)
	if err == nil {
		WaitSignal.L.Lock()
		WaitSignal.Broadcast()
		WaitSignal.L.Unlock()
	}
	return err
}

func WaitPortBeConfiged(pk PortKey) uint32 {
	state := GetPortState(pk)
	testWaitNum = testWaitNum + 1
	fmt.Printf("testWaitNum =  %d \r\n", testWaitNum)
	fmt.Printf("pk %v ;state : %d\r\n", pk, state)
	if state == PORTNOTEXIST || state == LINKCONFIGOK || state == LINKCONFIGERROR {
		return state
	}
	WaitSignal.L.Lock()
	for {
		WaitSignal.Wait()
		state := GetPortState(pk)
		fmt.Printf("pk %v ;state : %d\r\n", pk, state)
		if state == PORTNOTEXIST || state == LINKCONFIGOK || state == LINKCONFIGERROR {
			WaitSignal.L.Unlock()
			return state
		}
		testWaitNum = testWaitNum + 1
		fmt.Printf("testWaitNum =  %d \r\n", testWaitNum)
	}
}

func checkUnLinkedPortRouting() {
	WaitSignal.L.Lock()
	for {
		WaitSignal.Wait()
		testCheckLn = testCheckLn + 1
		fmt.Printf("testCheckUn = %d\r\n", testCheckLn)
		if state := ScanUnLinkedPort(); state == STATECHANGED {
			WaitSignal.Broadcast()
		}

	}
	WaitSignal.L.Unlock()

}

func checkLinkInfoRouting() {
	WaitSignal.L.Lock()
	for {
		WaitSignal.Wait()
		testCheckUn = testCheckUn + 1
		fmt.Printf("testCheckLn =  %d \r\n", testCheckUn)
		if state := ScanLinkInfo(); state == STATECHANGED {
			WaitSignal.Broadcast()
		}
	}
	WaitSignal.L.Unlock()
}
