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
	"gitlab.com/ribonin/apis/kingshot-redeem/model"
)

// DeletePlayer removes a player by PiD.
// @Summary Delete a player
// @Description Removes a player from the database by player ID.
// @Tags Player
// @Produce json
// @Param pid path int true "Player ID"
// @Success 202 {object} model.DeletePlayerResponse "Player deleted" example({"pid":123,"deleted":true})
// @Failure 500 {object} model.DeletePlayerErrorResponse "Delete failed" example({"pid":123,"deleted":false,"error":"database error"})
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
		return c.JSON(http.StatusInternalServerError, model.DeletePlayerErrorResponse{
			Pid:     pid,
			Deleted: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusAccepted, model.DeletePlayerResponse{Pid: pid, Deleted: true})
}
