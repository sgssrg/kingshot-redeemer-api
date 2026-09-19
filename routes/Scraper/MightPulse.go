package scraper

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"gitlab.com/ribonin/apis/kingshot-redeem/model"
	"gitlab.com/ribonin/apis/kingshot-redeem/routes/Scraper/lib"
)

func MPAlliancePlayerScaperRouter() echo.HandlerFunc {
	return func(c *echo.Context) error {
		kid := c.QueryParam("kid")
		alliance := c.QueryParam("alliance")
		if kid == "" {
			return c.JSON(http.StatusBadRequest, model.UpdatePlayerSSEValidator{WrongField: "kid", Message: "Please enter a non-empty string."})
		}
		if alliance == "" {
			return c.JSON(http.StatusBadRequest, model.UpdatePlayerSSEValidator{WrongField: "alliance", Message: "Please enter a non-empty string."})
		}

		aD, RespCode, err := lib.MPAlliancePlayerScaper(&kid, &alliance)
		if RespCode == 404 {
			c.JSON(http.StatusBadRequest, model.MPScrapeAllianceResp{ErrorState: true, Message: "Alliance Doesn't Exist in MightPulse"})
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.MPScrapeAllianceResp{ErrorState: true, Message: "Unable to that Alliance"})
		}
		// Finding alliance Leader's info
		pI := buildIndexByID(&aD.Members)
		players := transformMembers(*aD)

		AllianceLeader := model.PlayerInfo{Pid: uint(pI[aD.Alliance.LeaderUID].Fid), Kid: uint(pI[aD.Alliance.LeaderUID].Kid), Dname: pI[aD.Alliance.LeaderUID].NickName, Pfp: pI[aD.Alliance.LeaderUID].AvatarURL, Alliance: pI[aD.Alliance.LeaderUID].AllianceAbbr}

		AllianceInfo := model.AllianceInfo{AllianceName: &aD.Alliance.Name, AllianceShortCode: &aD.Alliance.Abbr, AllianceMemberCount: &aD.MemberCount, LastUpdatedAgo: aD.Alliance.UpdatedAgo, AlliancePlayers: &players, AllianceLeader: &AllianceLeader}

		return c.JSON(http.StatusOK, AllianceInfo)
	}
}

func buildIndexByID(players *[]model.MPMemberScrape) map[int]model.MPMemberScrape {
	index := make(map[int]model.MPMemberScrape, len(*players))
	for _, p := range *players {
		index[p.UID] = p
	}
	return index
}

func transformMembers(src model.MightPulseScrapeAlliance) []model.PlayerInfo {
	out := make([]model.PlayerInfo, 0, len(src.Members))
	for _, m := range src.Members {
		out = append(out, model.PlayerInfo{
			Pid:      uint(m.UID),
			Kid:      uint(m.Kid),
			Dname:    m.NickName,
			Pfp:      m.AvatarURL,
			Alliance: m.AllianceName,
		})
	}
	return out
}
