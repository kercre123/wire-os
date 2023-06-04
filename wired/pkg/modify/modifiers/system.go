package modifiers

import "os/exec"

func RunBashCommand(command string) {
	exec.Command("/bin/bash", "-c", command).Run()
}

func SetCPUFreq(performance bool) {
	var cpufreq string
	var memfreq string
	var cpugov string
	if performance {
		cpufreq = "1267200"
		memfreq = "800000"
		cpugov = "performance"
	} else {
		cpufreq = "833333"
		memfreq = "550000"
		cpugov = "interactive"
	}
	RunBashCommand("echo " + cpufreq + " > /sys/devices/system/cpu/cpu0/cpufreq/scaling_max_freq")
	RunBashCommand("echo disabled > /sys/kernel/debug/msm_otg/bus_voting")
	RunBashCommand("echo 0 > /sys/kernel/debug/msm-bus-dbg/shell-client/update_request")
	RunBashCommand("echo 1 > /sys/kernel/debug/msm-bus-dbg/shell-client/mas")
	RunBashCommand("echo 512 > /sys/kernel/debug/msm-bus-dbg/shell-client/slv")
	RunBashCommand("echo 0 > /sys/kernel/debug/msm-bus-dbg/shell-client/ab")
	RunBashCommand("echo active clk2 0 1 max " + memfreq + " > /sys/kernel/debug/rpm_send_msg/message")
	RunBashCommand("echo 1 > /sys/kernel/debug/msm-bus-dbg/shell-client/update_request")
	RunBashCommand("echo " + cpugov + " > /sys/devices/system/cpu/cpu0/cpufreq/scaling_governor")
}

func HigherPerformance_Apply() error {
	SetCPUFreq(true)
	return nil
}

func HigherPerformance_Remove() error {
	SetCPUFreq(false)
	return nil
}

func HigherPerformance_Init(applied bool) error {
	if applied {
		SetCPUFreq(true)
	} else {
		SetCPUFreq(false)
	}
	return nil
}
