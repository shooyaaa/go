package cron

import (
	. "github.com/shooyaaa/core/library"
	"github.com/shooyaaa/utils"
)

var GlobalTicker = NewTicker()

type Cron struct {
}

func (c Cron) Run() Module {
	telegram := utils.Telegram{}

	cron := GlobalTicker.Cron(10, telegram.Cron)
	cron.AddLimitation(HourUnit, Op_Bigger, 9)
	//cron.AddLimitation(HourUnit, Op_Smaller, 16)
	cron.AddLimitation(DayUnit, Op_Bigger, 0)
	cron.AddLimitation(DayUnit, Op_Smaller, 6)

	smzdm := utils.Telegram{}
	smzdm.AddJobs(&utils.Smzdm{})
	GlobalTicker.Cron(1800, smzdm.Do)

	return GlobalTicker
}
