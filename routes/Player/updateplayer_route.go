package player

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v5"
	"gitlab.com/ribonin/apis/kingshot-redeem/db"
	"gitlab.com/ribonin/apis/kingshot-redeem/model"
	"resty.dev/v3"
)

// UpdateAllPlayer updates real-time gamedata at the momement of all players in the database and makes a local copy to your database with the current data available ingame.
// @Summary Updates the Db with latest Player
// @Description Loads all players from the database and update every single of them, streaming per-player results as Server-Sent Events (event: update-res) followed by a summary (event: update-fin) about the rows it replaced with recent information. Rate limited.
// @Tags Player
// @Produce text/event-stream
// @Success 200 {string} string "update-res and update-fin event" example(event: update-res\ndata: {"updated":true,"player":{"pid":123,"kid":789,"dNname":"updated-player-name","pfp":"<updated-player-pfp-url>,"alliance":"DES"}}\n\n event: update-fin\ndata: {"count":123 })
// @Router /player/update [patch]

func UpdateAllPlayer(c *echo.Context) error {

	reqContext := c.Request().Context()
	res := c.Response()

	// 1. Set required SSE headers
	res.Header().Set("Content-Type", "text/event-stream")
	res.Header().Set("Cache-Control", "no-cache")
	res.Header().Set("Connection", "keep-alive")
	res.Header().Set("X-Accel-Buffering", "no") // Prevents Nginx/reverse proxies from buffering the stream

	// Ensure the underlying ResponseWriter supports HTTP flushing
	flusher, ok := res.(http.Flusher)
	if !ok {
		slog.Error("Streaming unsupported: response writer does not implement http.Flusher")
		return echo.NewHTTPError(http.StatusInternalServerError, "Streaming unsupported")
	}

	dbApp := db.InitDB()
	defer dbApp.Pool.Close()

	client := resty.New()
	defer client.Close()

	metaJsonData := model.UpdateMetaSSE{
		Count: 0,
	}

	ctx := context.Background()
	p, _ := dbApp.Queries.GetAllPlayers(ctx)
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-reqContext.Done():
				return
			case <-ticker.C:
				fmt.Fprintf(res, ": keep-alive\n\n")
				flusher.Flush()
			}
		}
	}()

	for _, pi := range p {
		select {
		case <-reqContext.Done():
			slog.Info("Client disconnected, stopping process")
			return nil
		default:
		}

		// Fetching player using Scraper
		// /scraper/player/

		var pD model.PlayerInfo
		pid := strconv.Itoa(int(pi.Pid))
		slog.Info("Updating Started for - " + pid)
		url := os.Getenv("host") + "/scraper/player/" + pid
		pF, _ := client.R().SetResult(&pD).SetRetryCount(3).SetHeader("Content-Type", "application/json").Get(url)
		if pF.StatusCode() == 202 {
			p, _ := dbApp.Queries.UpdatePlayer(ctx, db.UpdatePlayerParams{
				Pid:      pi.Pid,
				Kid:      pD.Kid,
				Dname:    pgtype.Text{String: pD.Dname, Valid: true},
				Pfp:      pgtype.Text{String: pD.Pfp, Valid: true},
				Alliance: pgtype.Text{String: pD.Alliance, Valid: true},
			})
			fmt.Println(p)
			slog.Info("Player Updated with PiD - " + pid)
			updateEvent := model.UpdatePlayerSSE{
				Updated: true,
				Player:  p,
			}

			// 2. Marshal payload to JSON
			jsonData, err := json.Marshal(updateEvent)
			if err != nil {
				slog.Error("Failed to marshal event", "error", err)
				continue
			}

			// 3. Format as standard SSE payload: "data: <json>\n\n"
			if _, err := fmt.Fprintf(res, "event: update-res\ndata: %s\n\n", jsonData); err != nil {
				return err
			}

			// 4. Immediately flush buffer to the client
			flusher.Flush()
			// sleep(1)
			metaJsonData.Count++
		} else {
			continue
		}

	}
	metaJsonBytes, err := json.Marshal(metaJsonData)
	if err != nil {
		slog.Error("Failed to marshal event", "error", err)
		return nil
	}

	// FIX 2: Fixed duplicate "data: data:" -> "data: %s\n\n"
	fmt.Fprintf(res, "event: update-fin\ndata: %s\n\n", metaJsonBytes)
	flusher.Flush()

	return nil
}
func sleep(ms int) {
	time.Sleep(time.Duration(ms) * time.Second)
}
