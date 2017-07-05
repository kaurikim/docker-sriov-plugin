package driver

import (
	"fmt"
	"path/filepath"
	"os"
)

const (
	netSysDir = "/sys/class/net"
	netDevPrefix = "device"
	netdevDriverDir = "device/driver"
	netdevUnbindFile = "unbind"
	netdevBindFile = "bind"

	netDevMaxVFCountFile = "sriov_totalvfs"
	netDevCurrentVFCountFile = "sriov_numvfs"
	netDevVFDevicePrefix = "virtfn"
)

func netDevDeviceDir(netDevName string) (string) {
	devDirName := netSysDir + "/" + netDevName + "/" + netDevPrefix
	return devDirName
}

func netdevGetMaxVFCount(name string) (int, error) {
	devDirName := netDevDeviceDir(name)

	maxDevFile := fileObject {
				Path: devDirName + "/" + netDevMaxVFCountFile,
		       }

	maxVfs, err := maxDevFile.ReadInt()
	if err != nil {
		return 0, err
	} else {
		fmt.Println("max_vfs = ", maxVfs)
		return maxVfs, nil
	}
}

func netdevSetMaxVFCount(name string, maxVFs int) (error) {
	devDirName := netDevDeviceDir(name)

	maxDevFile := fileObject {
				Path: devDirName + "/" + netDevCurrentVFCountFile,
		       }

	return maxDevFile.WriteInt(maxVFs)
}

func netdevGetEnabledVFCount(name string) (int, error) {
	devDirName := netDevDeviceDir(name)

	maxDevFile := fileObject {
				Path: devDirName + "/" + netDevCurrentVFCountFile,
		       }

	curVfs, err := maxDevFile.ReadInt()
	if err != nil {
		return 0, err
	} else {
		fmt.Println("cur_vfs = ", curVfs)
		return curVfs, nil
	}
}

func netdevEnableSRIOV(name string) (error) {
	var maxVFCount int
	var err error

	devDirName := netDevDeviceDir(name)

	devExist := dirExists(devDirName)
	if !devExist {
		return fmt.Errorf("device not found")
	}

	maxVFCount, err = netdevGetMaxVFCount(name)
	if err != nil {
		fmt.Println("netdevice found", name, maxVFCount)
		return err
	}

	if maxVFCount != 0 {
		return netdevSetMaxVFCount(name, maxVFCount)	
	} else {
		return fmt.Errorf("sriov unsupported")
		return nil
	}
}

func netdevDisableSRIOV(name string) (error) {
	devDirName := netDevDeviceDir(name)

	devExist := dirExists(devDirName)
	if !devExist {
		return fmt.Errorf("device not found")
	}

	return netdevSetMaxVFCount(name, 0)	
}

func vfNetdevNameFromParent(parentNetdev string, vfDir string) (string) {

	devDirName := netDevDeviceDir(parentNetdev)

	vfNetdev, _ := lsFilesWithPrefix(devDirName + "/" + vfDir + "/" + "net", "", false)
	if len(vfNetdev) <= 0 {
		return ""
	} else {
		return vfNetdev[0]
	}
}

func vfPCIDevNameFromVfDir(parentNetdev string, vfDir string) (string) {
	link := filepath.Join(netSysDir, parentNetdev, netDevPrefix, vfDir) 
	pciDevDir, err := os.Readlink(link)
	if err != nil {
		return ""
	}
	if (len(pciDevDir) <=3) {
		return ""
	}

	return pciDevDir[3:len(pciDevDir)]
}

func unbindVF(parentNetdev string, vfPCIDevName string) (error) {
	cmdFile := filepath.Join(netSysDir, parentNetdev, netdevDriverDir, netdevUnbindFile) 

	cmdFileObj := fileObject {
				Path: cmdFile,
		       }

	return cmdFileObj.Write(vfPCIDevName)
}

func bindVF(parentNetdev string, vfPCIDevName string) (error) {
	cmdFile := filepath.Join(netSysDir, parentNetdev, netdevDriverDir, netdevBindFile) 

	cmdFileObj := fileObject {
				Path: cmdFile,
		       }

	return cmdFileObj.Write(vfPCIDevName)
}

func vfDevList(name string) ([]string, error) {
	var vfDirList []string
	var i int
	devDirName := netDevDeviceDir(name)

	virtFnDirs, err := lsFilesWithPrefix(devDirName, netDevVFDevicePrefix, true)

	if (err != nil) {
		return nil, err
	}

	i = 0
	for _, vfDir := range virtFnDirs {
		vfDirList = append(vfDirList, vfDir)
		fmt.Println("virtual device name = ", vfDirList[i])
		i++
	}
	return vfDirList, nil
}
