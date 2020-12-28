package controller

import (
	"github.com/3hajk/vending_machine/controller/board"
	"github.com/rivo/users"
	"html/template"
)

type Option struct {
	Value    int8
	Id, Text string
	Selected bool
}

type BoardConfig struct {
	Title            template.HTML
	Data             *board.PourConfig
	Tare             *board.TareConfig
	Outs             *board.SensorConfig
	Sanitize         *board.SanitizeConfig
	BoardStatus      *board.VendingStatus
	Success          template.HTML
	Mode             bool
	SanitizeTemplate *[]Option
	Errors           map[string]string
	User             users.User
	UserEmail        string
	UserRole         int
}
