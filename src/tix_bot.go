package main

import (
	"context"
	"fmt"
	"log"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type InputData struct {
	url  string
	sid  string
	date string
	num  string
}

var app = tview.NewApplication()
var text = tview.NewTextView().
	SetTextColor(tcell.ColorGreen).
	SetText("(a) to add a new contact \n(q) to quit")
var form = tview.NewForm()
var pages = tview.NewPages()
var flex = tview.NewFlex()
var inputdata InputData

func setInputData() *tview.Form {

	form.AddInputField("Url", "", 20, nil, func(url string) {
		inputdata.url = url
	})

	form.AddInputField("SID", "", 20, nil, func(sid string) {
		inputdata.sid = sid
	})

	form.AddInputField("Date", "", 20, nil, func(date string) {
		inputdata.date = date
	})

	form.AddInputField("Num", "", 20, nil, func(num string) {
		inputdata.num = num
	})

	form.AddButton("Save", func() {
		pages.SwitchToPage("Menu")
	})

	return form
}

func main() {
	log.Println("開始執行腳本")

	// get input data

	flex.SetDirection(tview.FlexRow).
		AddItem(text, 0, 1, false)

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 113 {
			app.Stop()
		} else if event.Rune() == 97 {
			form.Clear(true)
			setInputData()
			pages.SwitchToPage("Add Input Data")
		}
		return event
	})

	pages.AddPage("Menu", flex, true, true)
	pages.AddPage("Add Input Data", form, true, false)

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
	//

	allocatorCtx, cancel := chromedp.NewExecAllocator(context.Background(), append(
		chromedp.DefaultExecAllocatorOptions[:],
		// 確保不包括 chromedp.Headless 選項
		chromedp.Flag("headless", false), // 這行可以明確禁用無頭模式（通常不是必需的，除非在默認選項中包含了無頭模式）
	)...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocatorCtx)
	defer cancel()

	sidCookie := &network.CookieParam{
		Name:   "SID",
		Value:  inputdata.sid,
		Domain: "tixcraft.com",
	}

	log.Println("正在訪問網站")
	//var ticketPageURL string
	err := chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			return network.SetCookies([]*network.CookieParam{sidCookie}).Do(ctx)
		}),
		chromedp.Navigate(inputdata.url),

		chromedp.WaitVisible(`.buy`, chromedp.ByQuery),
		//chromedp.Click(`#onetrust-accept-btn-handler`, chromedp.NodeVisible),
		//chromedp.AttributeValue(`.buy a`, "href", &ticketPageURL, nil),
		chromedp.Click(`.buy`, chromedp.NodeVisible),

		chromedp.WaitVisible(`#gameList`, chromedp.ByQuery),
		chromedp.Click(`//tr[contains(., '2024/01/07 (日)  15:00')]//button[contains(@class, 'btn-primary')]`, chromedp.BySearch),

		chromedp.WaitVisible(`li.select_form_b`, chromedp.ByQuery),
		chromedp.Click(`li.select_form_b:first-of-type`, chromedp.ByQuery),

		chromedp.WaitVisible(`#TicketForm_ticketPrice_03`, chromedp.ByQuery),
		chromedp.SetValue(`#TicketForm_ticketPrice_03`, inputdata.num, chromedp.ByQuery),
	)

	if err != nil {
		log.Printf("遇到錯誤：%v", err)
	} else {
		log.Println("腳本執行完成")
	}

	var input string
	fmt.Println("按下 Enter 鍵結束...")
	fmt.Scanln(&input)
}
