package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	sheetFlag string
)
var t_x86 string

type PhysicalMachines struct {
	PhysicalMachines map[string]PhysicalMachineSingle `yaml:"physical_machines"`
}

func RunCmd() *cobra.Command {
	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "Create a kvm VM",
		Long:  "vm run --sheet Sheet1",
		Run: func(cmd *cobra.Command, args []string) {
			switch sheetFlag {
			case "Sheet1":
				FromYamlReadConfig("./Sheet1Config.yaml")
			}
		},
	}
	runCmd.Flags().StringVarP(&sheetFlag, "sheet", "s", "", "Assignment sheet page (require)")
	//必选
	runCmd.MarkFlagRequired("sheet")

	return runCmd
}

// 创建配置
func createFile(filePath, content string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

func handleErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
func FromYamlReadConfig(file string) {
	t_x86 = `<domain type='kvm' xmlns:qemu='http://libvirt.org/schemas/domain/qemu/1.0'>
    <name>%s</name>
    <vcpu>%d</vcpu>
    <memory>%d</memory>
    <os>
        <type arch='x86_64'>hvm</type>
        <boot dev='hd'/>
    </os>
    <devices>
        <emulator>/usr/libexec/qemu-kvm</emulator>
        <disk type='file' device='disk'>
            <source file='/%s'/>
            <target dev='vda'/>
            <driver name='qemu' type='qcow2' cache='none'/>
        </disk>
        <disk type='file' device='cdrom'>
            <source file='/%s'/>
            <target dev='hdc'/>
            <readonly/>
            <driver name='qemu' type='raw'/>
        </disk>
        <disk type='file' device='disk'>
            <source file='/%s'/>
            <target dev='vdb'/>
            <driver name='qemu' type='qcow2' cache='none'/>
        </disk>
        <controller type='scsi' index='0' model='virtio-scsi'>
            <address type='pci' domain='0x0000' bus='0x00' slot='0x09' function='0x0'/>
        </controller>
        <interface type='bridge'>
            <source bridge='BRIDGE1'/>
            <mac address='%s'/>
            <model type='virtio'/>
        </interface>
        <graphics type='vnc' listen='0.0.0.0' port='-1' passwd='Gl_#9uzh'/>
        <input type='tablet' bus='usb'/>
        <serial type='pty'>
            <target port='0'/>
        </serial>
        <console type='pty'>
            <target type='serial' port='0'/>
        </console>
    </devices>
    <features>
        <acpi/>
    </features>
    <clock offset='localtime'/>
</domain>`
	yamlFile, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("error reading YAML file: %v", err)
	}

	var data PhysicalMachines
	err = yaml.Unmarshal(yamlFile, &data)
	if err != nil {
		log.Fatalf("error unmarshalling YAML: %v", err)
	}
	// 遍历所有物理机
	for _, physicalMachine := range data.PhysicalMachines {
		// mySSH := &ssh.SSH{
		// 	IP:       physicalMachine.IP,
		// 	Username: physicalMachine.User,
		// 	Password: physicalMachine.Password,
		// }
		// si := ssh.NewSSHClient(mySSH)
		// content, err := si.RunCommand("ls")
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// fmt.Println(string(content))
		// sfi := ssh.NewSSHClient(mySSH)
		// sfi.Copy("./Xshell7_lenovo.exe", "/apps/")
		err := os.Mkdir(physicalMachine.IP, 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}
		for _, vm := range physicalMachine.ConnectedVMs {
			IP := vm.Eth0
			Netmask := vm.Eth0Netmask
			Gateway := vm.Eth0Gateway
			Hostname := vm.Hostname
			Cpu := vm.CPU
			Memory := vm.Memory
			Mac := vm.Eth0Mac
			Osdisk := "apps/vmimages/" + Hostname + "/sys_disk"
			Datadisk := "apps/vmimages/" + Hostname + "/data_disk"
			Initiso := "apps/vmimages/" + Hostname + "/init_iso"

			// 创建文本文件
			err := createFile(physicalMachine.IP+"/"+Hostname+".txt", fmt.Sprintf("ETH0_IP=%s\nETH0_NETMASK=%s\nETH0_GATEWAY=%s\nHOSTNAME=%s\n", IP, Netmask, Gateway, Hostname))
			handleErr(err)

			// 创建 YAML 文件
			err = createFile(physicalMachine.IP+"/"+Hostname+".yaml", fmt.Sprintf(t_x86, Hostname, Cpu, Memory*1024, Osdisk, Datadisk, Initiso, Mac))
			handleErr(err)

		}

	}
}
