package main

import (
	"quizapi/common"
	"quizapi/modules/eventsmodule"
	"quizapi/modules/matchmodule"
	"quizapi/modules/questionsmodule"
	"quizapi/modules/timemodule"
	"quizapi/modules/usersmodule"
	"quizapi/modules/wsmodule"
)

var packages []common.Package = []common.Package{
	MainPackage{},
	timemodule.Package{},
	eventsmodule.Package{},
	usersmodule.Package{},
	questionsmodule.Package{},
	wsmodule.Package{},
	matchmodule.Package{},
}
