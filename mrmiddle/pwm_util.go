package mrmiddle

import (
	"os"

	"gobot.io/x/gobot/sysfs"
)

// pwmPath returns pwm base path
func pwmPath() string {
	return "/sys/class/pwm/pwmchip0"
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

func pin2pwmpin(pin string) string {
	return map[string]string{"3": "0", "5": "1", "6": "2", "9": "3"}[pin]
}
