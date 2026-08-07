package player

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v5"
	"gitlab.com/ribonin/apis/kingshot-redeem/db"
)

// DeletePlayer removes a player by PiD.
// @Summary Delete a player
// @Description Removes a player from the database by player ID.
// @Tags Player
// @Produce json
// @Param pid path int true "Player ID"
// @Success 202 {object} idealResponse
// @Failure 500 {object} map[string]interface{} "Delete failed"
// @Router /player/del/{pid} [delete]
func DeletePlayer(c *echo.Context) error {
	pid := c.Param("pid")
	dbApp := db.InitDB()
	fmt.Println(os.Getenv("DB_URI"))
	defer dbApp.Pool.Close()

	slog.Info(pid)
	ctx := context.Background()
	pidInt, _ := strconv.Atoi(pid)

	err := dbApp.Queries.DeletePlayer(ctx, int32(pidInt))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"id":      pid,
			"deleted": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusAccepted, idealResponse{Pid: pid, Deleted: true})
}

type idealResponse struct {
	Pid     string `json:"pid"`
	Deleted bool   `json:"deleted"`
}
