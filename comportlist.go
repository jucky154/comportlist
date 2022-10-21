package main

import (
	_ "embed"
	"zylo/reiwa"
	"github.com/tadvi/winc"
	"github.com/mastrolinux/go-serial-native"
	"strings"
	"strconv"
	"golang.org/x/text/encoding/japanese"
        "golang.org/x/text/transform"
        "io/ioutil"
)

const winsize = "comportlistwindow"

type ComportView struct {
	list *winc.ListView
}

var comportview ComportView

type ComportItem struct {
	Name string
	Description string
	Manufacturer string
}

func (item  ComportItem) Text() (text []string) {
	text = append(text, item.Name)
	text = append(text, item.Description)
	text = append(text, item.Manufacturer)
	return
}

func (item  ComportItem) ImageIndex() int {
	return 0
}

func convutf8(str string) string{
	if len(str)==0{
		return "-"
	}

        ret, err := ioutil.ReadAll(transform.NewReader(strings.NewReader(str), japanese.ShiftJIS.NewDecoder()))
        if err != nil {
                return "表示できない文字コードです"
        }

        return string(ret)
}

func comportUpdate() {
	//listを消す
	comportview.list.DeleteAllItems()

	//情報を取得
	ports, err := serial.ListPorts()
	if err != nil {
		comportview.list.AddItem(ComportItem{
		Name : "ポートに接続できませんでした・何もありません",
		Description : "-",
		Manufacturer : "-",
		})
		return
	}

	if len(ports) == 0 {
		comportview.list.AddItem(ComportItem{
		Name : "シリアルポートに何も接続がありません",
		Description : "-",
		Manufacturer : "-",
		})
		return
	}

	for _, info := range ports {
		comportview.list.AddItem(ComportItem{
		Name : convutf8(info.Name()),
		Description : convutf8(info.Description()),
		Manufacturer : convutf8(info.USBManufacturer()),
		})
	}
	return
}

var mainWindow *winc.Form


func makewindow() {
	// --- Make Window
	mainWindow = newForm(nil)
	x, _ := strconv.Atoi(reiwa.GetINI(winsize, "x"))
	y, _ := strconv.Atoi(reiwa.GetINI(winsize, "y"))
	w, _ := strconv.Atoi(reiwa.GetINI(winsize, "w"))
	h, _ := strconv.Atoi(reiwa.GetINI(winsize, "h"))
	if w <= 0 || h <= 0 {
		w = 720
		h = 140
	}

	mainWindow.SetSize(w, h)
	if x <= 0 || y <= 0 {
		mainWindow.Center()
	} else {
		mainWindow.SetPos(x, y)
	}
	mainWindow.SetText("ComPort一覧")

	comportview.list = winc.NewListView(mainWindow)
	comportview.list.EnableEditLabels(false)
	comportview.list.AddColumn("COM番号", 100)
	comportview.list.AddColumn("詳細", 300)
	comportview.list.AddColumn("製造企業", 300)

	btn := winc.NewPushButton(mainWindow)
	btn.SetText("接続情報の更新")
	btn.OnClick().Bind(func(e *winc.Event) {
		comportUpdate()
	})

	dock := winc.NewSimpleDock(mainWindow)
	dock.Dock(btn, winc.Top)
	dock.Dock(comportview.list, winc.Fill)

	mainWindow.Show()
	comportUpdate()
	return
}

func init() {
	reiwa.OnLaunchEvent = onLaunchEvent
	reiwa.PluginName = "infoEs"
}

func onLaunchEvent() {
	reiwa.RunDelphi(`PluginMenu.Add(op.Put(MainMenu.CreateMenuItem(), "Name", "PluginComportlistWindow"))`)
	reiwa.RunDelphi(`op.Put(MainMenu.FindComponent("PluginComportlistWindow"), "Caption", "ComPort一覧")`)

	reiwa.HandleButton("MainForm.MainMenu.PluginComportlistWindow", func(num int){
		makewindow()
	})	
}