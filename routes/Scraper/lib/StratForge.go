package lib

import (
	"log/slog"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"gitlab.com/ribonin/apis/kingshot-redeem/model"
)

func StratForgePlayerScraper(fid string) (*model.ScrapePlayerInfo, *model.CustomScrapePlayerErrInfo, error) {
	slog.Info("Fetch Started for fid -" + fid)
	link := "https://stratforge.tools/p/" + fid
	var pInfo model.ScrapePlayerInfo
	pid, _ := strconv.Atoi(fid)
	pInfo.Pid = pid

	var errInfo *model.CustomScrapePlayerErrInfo // pointer, starts nil

	cly := colly.NewCollector(
		colly.AllowedDomains("stratforge.tools"),
	)

	// Player name
	cly.OnHTML("h1.font-extrabold.text-2xl.text-bone.tracking-tight", func(e *colly.HTMLElement) {
		text := strings.TrimSpace(e.Text)
		pInfo.Dname = text

		if text == "We don't have that one" {
			errInfo = &model.CustomScrapePlayerErrInfo{
				Pid:     pid,
				Type:    0,
				Message: "PlayerID Doesn't Exist",
			}
		}
	})

	// Kingdom ID
	cly.OnHTML("a.text-forge-text", func(e *colly.HTMLElement) {
		href := e.Attr("href") // "/k/1406"
		slog.Info("Anchor href", "href", href)

		// simplest way: strip prefix
		if strings.HasPrefix(href, "/k/") {
			kidStr := strings.TrimPrefix(href, "/k/")
			if kid, err := strconv.Atoi(kidStr); err == nil {
				pInfo.Kid = kid
				slog.Info("Kid found", "value", kid)
			}
		}
	})

	// Alliance tag
	cly.OnHTML("span.ml-1\\.5.font-mono.text-dim.text-xs", func(e *colly.HTMLElement) {
		re := regexp.MustCompile(`[^a-zA-Z0-9]`)
		pInfo.Alliance = re.ReplaceAllString(e.Text, "")
	})

	// Profile picture
	cly.OnHTML("img.shrink-0.rounded-full.object-cover", func(e *colly.HTMLElement) {
		pInfo.Pfp = e.Attr("src")
	})

	cly.OnError(func(r *colly.Response, err error) {
		slog.Error("Request failed", "url", r.Request.URL, "err", err)
	})

	if err := cly.Visit(link); err != nil {
		return &pInfo, errInfo, err
	}

	return &pInfo, errInfo, nil
}
