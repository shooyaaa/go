package cron

import (
	"fmt"

	"github.com/shooyaaa/config"
	. "github.com/shooyaaa/core/library"
	. "github.com/shooyaaa/core/network"
	"github.com/shooyaaa/utils"
)

var GlobalTicker = NewTicker()

type Cron struct {
}

func (c Cron) Run() Module {
	telegram := utils.Telegram{}

	cron := GlobalTicker.Cron(1800, telegram.Cron)
	cron.AddLimitation(HourUnit, Op_Bigger, 9)
	cron.AddLimitation(HourUnit, Op_Smaller, 16)
	cron.AddLimitation(DayUnit, Op_Bigger, 0)
	cron.AddLimitation(DayUnit, Op_Smaller, 6)

	smzdm := utils.Telegram{}
	smzdm.AddJobs(&utils.Smzdm{})
	GlobalTicker.Cron(1800, smzdm.Do)

	//GlobalTicker.Cron(30, Monitor)

	return GlobalTicker
}

func Monitor() {
	up, down := MonitorHosts()
	telegram := utils.Telegram{}
	if len(up) > 0 {
		for _, u := range up {
			telegram.SendText(fmt.Sprintf("host %s is up with ip %s", u.Name, u.Ip), config.TelegramUserID)
		}
	}

	if len(down) > 0 {
		for _, d := range down {
			telegram.SendText(fmt.Sprintf("host %s is down with ip %s", d.Name, d.Ip), config.TelegramUserID)
		}
	}
}
