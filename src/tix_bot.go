package main

import (
	"context"
	"fmt"
	"log"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func main() {
	log.Println("開始執行腳本")

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
		Value:  "qav8f6f75gbq5ivk9eapognjjb",
		Domain: "tixcraft.com",
	}

	log.Println("正在訪問網站")
	//var ticketPageURL string
	err := chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			return network.SetCookies([]*network.CookieParam{sidCookie}).Do(ctx)
		}),
		chromedp.Navigate("https://tixcraft.com/activity/detail/24_lesmis"),

		chromedp.WaitVisible(`.buy`, chromedp.ByQuery),
		//chromedp.Click(`#onetrust-accept-btn-handler`, chromedp.NodeVisible),
		//chromedp.AttributeValue(`.buy a`, "href", &ticketPageURL, nil),
		chromedp.Click(`.buy`, chromedp.NodeVisible),

		chromedp.WaitVisible(`#gameList`, chromedp.ByQuery),
		chromedp.Click(`//tr[contains(., '2024/01/07 (日)  15:00')]//button[contains(@class, 'btn-primary')]`, chromedp.BySearch),

		chromedp.WaitVisible(`li.select_form_b`, chromedp.ByQuery),
		chromedp.Click(`li.select_form_b:first-of-type`, chromedp.ByQuery),

		chromedp.WaitVisible(`#TicketForm_ticketPrice_03`, chromedp.ByQuery),
		chromedp.SetValue(`#TicketForm_ticketPrice_03`, "2", chromedp.ByQuery),
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
