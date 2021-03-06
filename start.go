package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

// go build -ldflags="-H windowsgui"

var IsSpecialMode = walk.NewMutableCondition()

type MyMainWindow struct {
	*walk.MainWindow
}

func main() {

	mw := new(MyMainWindow)

	var toggleSpecialModePB *walk.PushButton
	var teDay, teDayCount, teName, teCount, teCountNoBlank *walk.TextEdit

	var openAction, showAboutBoxAction *walk.Action
	var recentMenu *walk.Menu

	if err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "마감 안내기",
		MenuItems: []MenuItem{
			Menu{
				Text: "파일 업로드",
				Items: []MenuItem{
					Action{
						AssignTo:    &openAction,
						Text:        "파일 추가",
						Enabled:     Bind("enabledCB.Checked"),
						Visible:     Bind("!openHiddenCB.Checked"),
						OnTriggered: mw.fileUploadAction_Triggered,
					},
					Menu{
						AssignTo: &recentMenu,
						Text:     "최근 파일",
						//OnTriggered: mw.recentFileAction_Triggered,
					},
					Separator{},
					Action{
						Text:        "종료",
						OnTriggered: func() { mw.Close() },
					},
				},
			},

			Menu{
				Text: "&Help",
				Items: []MenuItem{
					Action{
						AssignTo:    &showAboutBoxAction,
						Text:        "About",
						OnTriggered: mw.showAboutBoxAction_Triggered,
					},
				},
			},
		},

		ContextMenuItems: []MenuItem{
			ActionRef{&showAboutBoxAction},
		},

		MinSize: Size{270, 150},
		Layout:  VBox{},
		Children: []Widget{

			TextEdit{
				Text:     "D-DAY",
				AssignTo: &teDay, ReadOnly: true},

			TextEdit{
				Text:     "D-DAY까지 남은 날",
				AssignTo: &teDayCount, ReadOnly: true},
			TextEdit{
				Text:     "원고 이름",
				AssignTo: &teName, ReadOnly: true},
			TextEdit{
				Text:     "공백 포함 글자수",
				AssignTo: &teCount, ReadOnly: true},
			TextEdit{
				Text:     "공백 미포함 글자수",
				AssignTo: &teCountNoBlank, ReadOnly: true},

			PushButton{
				AssignTo: &toggleSpecialModePB,
				Text:     "항상 위 기능",
				OnClicked: func() {
					IsSpecialMode.SetSatisfied(!IsSpecialMode.Satisfied())

					if IsSpecialMode.Satisfied() {
						toggleSpecialModePB.SetText("항상 위 기능 켜짐")
					} else {
						toggleSpecialModePB.SetText("항상 위 기능 꺼짐")
					}
				},
			},

			PushButton{
				Text: "마감일 안내받기",
				OnClicked: func() {
					if teDay.Text() == "D-DAY" {
						return
					}
					day, name := teDay, teName
					Alarm(day.Text(), name.Text())
				},
			},
		},
	}.Create()); err != nil {
		walk.MsgBox(mw, "err", err.Error(), walk.MsgBoxIconInformation)
	}

	addRecentFileActions := func(conf Configuration) {

		a := walk.NewAction()
		a.SetText(conf.Filename)
		a.Triggered().Attach(func() {

			day := conf.Dday
			filename := conf.Filename

			count, countNoBlank := CountFile(filename)
			dayCount := GetDDay(day)
			teDay.SetText(day)
			teDayCount.SetText(strconv.Itoa(dayCount))
			teName.SetText(filename)
			teCount.SetText("공백 포함 " + count + " 자")
			teCountNoBlank.SetText("공백 미포함 " + countNoBlank + " 자")

		})
		recentMenu.Actions().Add(a)

	}

	jsonFile, err := os.Open(ConfFilePath)
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var conf []Configuration
	json.Unmarshal(byteValue, &conf)

	if err != nil {
		walk.MsgBox(mw, "err", err.Error(), walk.MsgBoxIconInformation)
	}

	// conf 의 json을 읽어왔는데 비어있다고 나오는중
	for _, v := range conf {

		walk.MsgBox(mw, "체크", v.Filename+v.Dday, walk.MsgBoxIconInformation)
		addRecentFileActions(v)

	}

	mw.Run()
}

func (mw *MyMainWindow) showAboutBoxAction_Triggered() {
	walk.MsgBox(mw, "About", `글 쓰시는 분들의 마감을 도와드립니다.
			20180703 히어로
				github @hero0926
		`, walk.MsgBoxIconInformation)
}

func (mw *MyMainWindow) fileUploadAction_Triggered() {
	Fileupload()
}
