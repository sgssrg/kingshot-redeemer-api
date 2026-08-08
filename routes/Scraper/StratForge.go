package scraper

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"

	"gitlab.com/ribonin/apis/kingshot-redeem/model"
	"gitlab.com/ribonin/apis/kingshot-redeem/routes/Scraper/lib"
)

// StratForgePlayerScraperAPI godoc
// @Summary      Scrape player profile from StratForge
// @Description  Fetches player profile data from stratforge.tools by player FID (Fleet ID). Returns player display name, kingdom ID, alliance, and profile picture.
// @Tags         Scraper
// @Accept       json
// @Produce      json
// @Param        fid   path      string  true  "Player Fleet ID"
// @Success      202   {object}  model.ScrapePlayerInfo
// @Failure      409   {object}  model.CustomScrapePlayerErrInfo "Player ID doesn't exist"
// @Failure      500   {object}  model.CustomScrapePlayerErrInfo "Scraping failed or invalid Player ID"
// @Failure      502   {object}  map[string]string "Failed to fetch from stratforge.tools"
// @Router       /scraper/stratforge/{fid} [get]

func StratForgePlayerScraperAPI(c *echo.Context) error {
	fid := c.Param("fid")
	slog.Info("Fetch Started for fid -" + fid)
	pidInt, _ := strconv.Atoi(fid)

	pInfo, errInfo, err := lib.StratForgePlayerScraper(fid)
	if err != nil {
		slog.Error("Error scraping player", "fid", fid, "err", err)
		return c.JSON(http.StatusInternalServerError, model.CustomScrapePlayerErrInfo{
			Pid:     pidInt,
			Type:    1,
			Message: "Player unable to be scraped / PlayerID is wrong",
		})
	}

	if errInfo != nil {
		return c.JSON(http.StatusConflict, errInfo)
	}
	slog.Info("Player ", "Fetched", pInfo)
	return c.JSON(http.StatusAccepted, pInfo)
}
