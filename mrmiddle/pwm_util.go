package mrmiddle

import (
	"os"

	"gobot.io/x/gobot/sysfs"
)

// pwmPath returns pwm base path
func pwmPath() string {
	return "/sys/class/pwm/pwmchip0"
}

// pwmExportPath returns export path
func pwmExportPath() string {
	return pwmPath() + "/export"
}

// pwmUnExportPath returns unexport path
func pwmUnExportPath() string {
	return pwmPath() + "/unexport"
}

// pwmDutyCyclePath returns duty_cycle path for specified pin
func pwmDutyCyclePath(pin string) string {
	return pwmPath() + "/pwm" + pin + "/duty_cycle"
}

// pwmPeriodPath returns period path for specified pin
func pwmPeriodPath(pin string) string {
	return pwmPath() + "/pwm" + pin + "/period"
}

// pwmEnablePath returns enable path for specified pin
func pwmEnablePath(pin string) string {
	return pwmPath() + "/pwm" + pin + "/enable"
}
func writeSysfsFile(path string, data []byte) (i int, err error) {
	file, err := sysfs.OpenFile(path, os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		return
	}

	return file.Write(data)
}

func pwmEnable(pin string, val string) (err error) {
	_, err = writeSysfsFile(pwmEnablePath(pin2pwmpin(pin)), []byte(val))
	return
}

// writeDuty writes value to pwm duty cycle path
func writeDuty(pin string, duty string) (err error) {
	_, err = writeSysfsFile(pwmDutyCyclePath(pin2pwmpin(pin)), []byte(duty))
	return
}

// export writes pin to pwm export path
func export(pin string) (err error) {
	_, err = writeSysfsFile(pwmExportPath(), []byte(pin2pwmpin(pin)))
	return
}

// export writes pin to pwm unexport path
func unexport(pin string) (err error) {
	_, err = writeSysfsFile(pwmUnExportPath(), []byte(pin2pwmpin(pin)))
	return
}

func pin2pwmpin(pin string) string {
	return map[string]string{"3": "0", "5": "1", "6": "2", "9": "3"}[pin]
}
