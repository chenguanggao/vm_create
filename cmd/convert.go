/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type PhysicalMachineSingle struct {
	IP           string                 `yaml:"ip"`
	User         string                 `yaml:"user"`
	Password     string                 `yaml:"password"`
	ConnectedVMs []VirtualMachineSingle `yaml:"connected_vms"`
}

type VirtualMachineSingle struct {
	Hostname    string `yaml:"hostname"`
	OS          string `yaml:"os"`
	CPU         int    `yaml:"cpu"`
	Memory      int    `yaml:"memory"`
	OsDisk      int    `yaml:"osdisk"`
	DataDisk    int    `yaml:"datadisk"`
	Eth0        string `yaml:"eth0"`
	Eth0Vlan    int    `yaml:"eth0vlan"`
	Eth0Gateway string `yaml:"eth0gateway"`
	Eth0Mac     string `yaml:"eth0mac"`
	Eth0Net     string `yaml:"eth0net"`
	Eth0Netmask string `yaml:"eth0netmask"`
}

type PhysicalMachineDouble struct {
	IP           string                 `yaml:"ip"`
	User         string                 `yaml:"user"`
	Password     string                 `yaml:"password"`
	ConnectedVMs []VirtualMachineDouble `yaml:"connected_vms"`
}

type VirtualMachineDouble struct {
	Hostname    string `yaml:"hostname"`
	OS          string `yaml:"os"`
	CPU         int    `yaml:"cpu"`
	Memory      int    `yaml:"memory"`
	OsDisk      int    `yaml:"osdisk"`
	DataDisk    int    `yaml:"datadisk"`
	Eth0        string `yaml:"eth0"`
	Eth0Vlan    int    `yaml:"eth0vlan"`
	Eth0Gateway string `yaml:"eth0gateway"`
	Eth0Mac     string `yaml:"eth0mac"`
	Eth0Net     string `yaml:"eth0net"`
	Eth0Netmask string `yaml:"eth0netmask"`
	Eth1        string `yaml:"eth1"`
	Eth1Vlan    int    `yaml:"eth1vlan"`
	Eth1Gateway string `yaml:"eth1gateway"`
	Eth1Mac     string `yaml:"eth1mac"`
	Eth1Net     string `yaml:"eth1net"`
	Eth1Netmask string `yaml:"eth1netmask"`
}

type PhysicalMachineThree struct {
	IP           string                `yaml:"ip"`
	User         string                `yaml:"user"`
	Password     string                `yaml:"password"`
	ConnectedVMs []VirtualMachineThree `yaml:"connected_vms"`
}

type VirtualMachineThree struct {
	Hostname    string `yaml:"hostname"`
	OS          string `yaml:"os"`
	CPU         int    `yaml:"cpu"`
	Memory      int    `yaml:"memory"`
	OsDisk      int    `yaml:"osdisk"`
	DataDisk    int    `yaml:"datadisk"`
	Eth0        string `yaml:"eth0"`
	Eth0Vlan    int    `yaml:"eth0vlan"`
	Eth0Gateway string `yaml:"eth0gateway"`
	Eth0Mac     string `yaml:"eth0mac"`
	Eth0Net     string `yaml:"eth0net"`
	Eth0Netmask string `yaml:"eth0netmask"`
	Eth1        string `yaml:"eth1"`
	Eth1Vlan    int    `yaml:"eth1vlan"`
	Eth1Gateway string `yaml:"eth1gateway"`
	Eth1Mac     string `yaml:"eth1mac"`
	Eth1Net     string `yaml:"eth1net"`
	Eth1Netmask string `yaml:"eth1netmask"`
	Eth2        string `yaml:"eth2"`
	Eth2Vlan    int    `yaml:"eth2vlan"`
	Eth2Gateway string `yaml:"eth2gateway"`
	Eth2Mac     string `yaml:"eth2mac"`
	Eth2Net     string `yaml:"eth2net"`
	Eth2Netmask string `yaml:"eth2netmask"`
}

func ConvertCmd() *cobra.Command {
	var convertCmd = &cobra.Command{
		Use:   "convert",
		Short: "Parse excel into a usable yaml file",
		Long:  `vm  --config parameter.xlsx analyze `,
		Run: func(cmd *cobra.Command, args []string) {
			switch sheetFlag {
			case "Sheet1":
				convertExcelSheet1(rootOpt.cfgFile, sheetFlag)
			case "Sheet2":
				convertExcelSheet2(rootOpt.cfgFile, sheetFlag)
			case "Sheet3":
				convertExcelSheet3(rootOpt.cfgFile, sheetFlag)
			}
		},
	}
	convertCmd.Flags().StringVarP(&sheetFlag, "sheet", "s", "", "Assignment sheet page (require)")
	//必选
	convertCmd.MarkFlagRequired("sheet")
	return convertCmd
}

func convertToInt(value string, field string) int {
	if value == "" {
		return 0
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("error converting %s to int: %v", field, err)
	}
	return intValue
}
func convertExcelSheet1(cfgfile string, sheet string) {
	xlsx, err := excelize.OpenFile(cfgfile)
	if err != nil {
		log.Fatal(err)
	}

	rows := xlsx.GetRows(sheet)
	physicalMachines := make(map[string]PhysicalMachineSingle)
	for _, row := range rows[2:] { // Skip the header row
		physicalIP := row[0]
		user := row[1]
		password := row[2]
		vm := VirtualMachineSingle{
			Hostname:    row[3],
			OS:          row[4],
			CPU:         convertToInt(row[5], "CPU"),
			Memory:      convertToInt(row[6], "Memory"),
			OsDisk:      convertToInt(row[7], "OSDisk"),
			DataDisk:    convertToInt(row[8], "DataDisk"),
			Eth0:        row[9],
			Eth0Vlan:    convertToInt(row[10], "VLAN1"),
			Eth0Gateway: row[11],
			Eth0Mac:     row[12],
			Eth0Net:     row[13],
			Eth0Netmask: row[14],
		}
		if pm, ok := physicalMachines[physicalIP]; ok {
			pm.ConnectedVMs = append(pm.ConnectedVMs, vm)
			physicalMachines[physicalIP] = pm
		} else {
			physicalMachines[physicalIP] = PhysicalMachineSingle{
				IP:           physicalIP,
				User:         user,
				Password:     password,
				ConnectedVMs: []VirtualMachineSingle{vm},
			}
		}
		// 序列化物理机信息
		data := map[string]interface{}{
			"physical_machines": physicalMachines,
		}
		yamlData, err := yaml.Marshal(data)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		// 将 YAML 数据写入文件
		file, err := os.Create(sheetFlag + "Config.yaml")
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		defer file.Close()

		_, err = file.Write(yamlData)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

	}
	fmt.Println("YAML configuration file for physical machines created successfully.")
}
func convertExcelSheet2(cfgfile string, sheet string) {
	fmt.Println(sheet)
	xlsx, err := excelize.OpenFile(cfgfile)
	if err != nil {
		log.Fatal(err)
	}

	rows := xlsx.GetRows(sheet)
	physicalMachines := make(map[string]PhysicalMachineDouble)
	for _, row := range rows[2:] { // Skip the header row
		physicalIP := row[0]
		user := row[1]
		password := row[2]
		vm := VirtualMachineDouble{
			Hostname:    row[3],
			OS:          row[4],
			CPU:         convertToInt(row[5], "CPU"),
			Memory:      convertToInt(row[6], "Memory"),
			OsDisk:      convertToInt(row[7], "OSDisk"),
			DataDisk:    convertToInt(row[8], "DataDisk"),
			Eth0:        row[9],
			Eth0Vlan:    convertToInt(row[10], "VLAN1"),
			Eth0Gateway: row[11],
			Eth0Mac:     row[12],
			Eth0Net:     row[13],
			Eth0Netmask: row[14],
			Eth1:        row[15],
			Eth1Vlan:    convertToInt(row[16], "VLAN2"),
			Eth1Gateway: row[17],
			Eth1Mac:     row[18],
			Eth1Net:     row[19],
			Eth1Netmask: row[20],
		}
		if pm, ok := physicalMachines[physicalIP]; ok {
			pm.ConnectedVMs = append(pm.ConnectedVMs, vm)
			physicalMachines[physicalIP] = pm
		} else {
			physicalMachines[physicalIP] = PhysicalMachineDouble{
				IP:           physicalIP,
				User:         user,
				Password:     password,
				ConnectedVMs: []VirtualMachineDouble{vm},
			}
		}
		// 序列化物理机信息
		data := map[string]interface{}{
			"physical_machines": physicalMachines,
		}
		yamlData, err := yaml.Marshal(data)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		// 将 YAML 数据写入文件
		file, err := os.Create(sheetFlag + "Config.yaml")
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		defer file.Close()

		_, err = file.Write(yamlData)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

	}
	fmt.Println("YAML configuration file for physical machines created successfully.")
}
func convertExcelSheet3(cfgfile string, sheet string) {
	xlsx, err := excelize.OpenFile(cfgfile)
	if err != nil {
		log.Fatal(err)
	}

	rows := xlsx.GetRows(sheet)
	physicalMachines := make(map[string]PhysicalMachineThree)
	for _, row := range rows[2:] { // Skip the header row
		physicalIP := row[0]
		user := row[1]
		password := row[2]
		vm := VirtualMachineThree{
			Hostname:    row[3],
			OS:          row[4],
			CPU:         convertToInt(row[5], "CPU"),
			Memory:      convertToInt(row[6], "Memory"),
			OsDisk:      convertToInt(row[7], "OSDisk"),
			DataDisk:    convertToInt(row[8], "DataDisk"),
			Eth0:        row[9],
			Eth0Vlan:    convertToInt(row[10], "VLAN1"),
			Eth0Gateway: row[11],
			Eth0Mac:     row[12],
			Eth0Net:     row[13],
			Eth0Netmask: row[14],
			Eth1:        row[15],
			Eth1Vlan:    convertToInt(row[16], "VLAN2"),
			Eth1Gateway: row[17],
			Eth1Mac:     row[18],
			Eth1Net:     row[19],
			Eth1Netmask: row[20],
			Eth2:        row[21],
			Eth2Vlan:    convertToInt(row[22], "VLAN3"),
			Eth2Gateway: row[23],
			Eth2Mac:     row[24],
			Eth2Net:     row[25],
			Eth2Netmask: row[26],
		}
		if pm, ok := physicalMachines[physicalIP]; ok {
			pm.ConnectedVMs = append(pm.ConnectedVMs, vm)
			physicalMachines[physicalIP] = pm
		} else {
			physicalMachines[physicalIP] = PhysicalMachineThree{
				IP:           physicalIP,
				User:         user,
				Password:     password,
				ConnectedVMs: []VirtualMachineThree{vm},
			}
		}
		// 序列化物理机信息
		data := map[string]interface{}{
			"physical_machines": physicalMachines,
		}
		yamlData, err := yaml.Marshal(data)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		// 将 YAML 数据写入文件
		file, err := os.Create(sheetFlag + "Config.yaml")
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		defer file.Close()

		_, err = file.Write(yamlData)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

	}
	fmt.Println("YAML configuration file for physical machines created successfully.")
}
