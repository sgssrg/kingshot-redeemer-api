package model

type ScrapePlayerInfo struct {
	Pid      int    `json:"pid" example:"123456"`
	Kid      int    `json:"kid" example:"123"`
	Dname    string `json:"dName,omitempty" example:"MeowMeow"`
	Pfp      string `json:"pfp,omitempty" example:"<-pfp-url->"`
	Alliance string `json:"alliance,omitempty" example:"XYZ"`
}

type CustomScrapePlayerErrInfo struct {
	Pid     int    `json:"pid"`
	Type    int    `json:"type"`
	Message string `json:"message"`
}

type MightPulseScrapeAlliance struct {
	Ok       bool `json:"ok"`
	Aid      int  `json:"aid"`
	Alliance struct {
		Aid         int     `json:"aid"`
		Kid         int     `json:"kid"`
		Name        string  `json:"name"`
		Abbr        string  `json:"abbr"`
		LeaderUID   int     `json:"leader_uid"`
		LeaderName  string  `json:"leader_name"`
		Count       int     `json:"count"`
		MemberMax   int     `json:"member_max"`
		Power       int64   `json:"power"`
		PowerRank   int     `json:"power_rank"`
		Lv          int     `json:"lv"`
		Language    int     `json:"language"`
		Notice      string  `json:"notice"`
		Manifesto   string  `json:"manifesto"`
		Public      bool    `json:"public"`
		UpdatedAt   float64 `json:"updated_at"`
		FirstSeenAt float64 `json:"first_seen_at"`
		Flag        []int   `json:"flag"`
		Raw         struct {
			Aid       int    `json:"aid"`
			Name      string `json:"name"`
			Abbr      string `json:"abbr"`
			Flag      []int  `json:"flag"`
			FlagBg    []int  `json:"flag_bg"`
			Lv        int    `json:"lv"`
			Exp       int    `json:"exp"`
			Count     int    `json:"count"`
			Notice    string `json:"notice"`
			Manifesto string `json:"manifesto"`
			Condition struct {
				Level int `json:"level"`
				Power int `json:"power"`
			} `json:"condition"`
			Language  int `json:"language"`
			Leader    int `json:"leader"`
			RankNames []struct {
				Rank int    `json:"rank"`
				Name string `json:"name"`
			} `json:"rank_names"`
			LeaderName string `json:"leader_name"`
			Power      int64  `json:"power"`
			PowerRank  int    `json:"power_rank"`
			NoticeUID  int    `json:"notice_uid"`
			MemberMax  int    `json:"member_max"`
			NoticeTime int    `json:"notice_time"`
			HonorLevel int    `json:"honor_level"`
			DuelGrade  int    `json:"duel_grade"`
			Kid        int    `json:"kid"`
		} `json:"raw"`
		UpdatedAgo    string `json:"updated_ago"`
		FlagURL       string `json:"flag_url"`
		LanguageLabel string `json:"language_label"`
		PublicLabel   string `json:"public_label"`
		Condition     struct {
			Level int `json:"level"`
			Power int `json:"power"`
		} `json:"condition"`
		ConditionLevel     int    `json:"condition_level"`
		ConditionPower     int    `json:"condition_power"`
		TrackerPower       int64  `json:"tracker_power"`
		TrackerPowerRank   int    `json:"tracker_power_rank"`
		TrackerMemberCount int    `json:"tracker_member_count"`
		KingdomAlliances   int    `json:"kingdom_alliances"`
		Slug               string `json:"slug"`
		Permalink          string `json:"permalink"`
	} `json:"alliance"`
	Members []struct {
		UID                 int           `json:"uid"`
		Fid                 int           `json:"fid"`
		NickName            string        `json:"nick_name"`
		Power               int           `json:"power"`
		StoveLv             int           `json:"stove_lv"`
		Vip                 interface{}   `json:"vip"`
		IsVipActive         interface{}   `json:"is_vip_active"`
		Lv                  int           `json:"lv"`
		Kid                 int           `json:"kid"`
		Aid                 int           `json:"aid"`
		AllianceAbbr        string        `json:"alliance_abbr"`
		AllianceName        string        `json:"alliance_name"`
		X                   int           `json:"x"`
		Y                   int           `json:"y"`
		OfflineTime         interface{}   `json:"offline_time"`
		PowerUpdateTime     int           `json:"power_update_time"`
		LastActiveAt        float64       `json:"last_active_at"`
		Language            string        `json:"language"`
		Image               int           `json:"image"`
		Office              int           `json:"office"`
		UploadImage         interface{}   `json:"upload_image"`
		LordsEquipShow      interface{}   `json:"lords_equip_show"`
		LordsEquipmentsJSON string        `json:"lords_equipments_json"`
		BattleRecordJSON    string        `json:"battle_record_json"`
		EquipedskinJSON     string        `json:"equipedskin_json"`
		GloryJSON           interface{}   `json:"glory_json"`
		ArenaHeroesJSON     string        `json:"arena_heroes_json"`
		AllianceRank        int           `json:"alliance_rank"`
		ShieldEndtime       interface{}   `json:"shield_endtime"`
		BurnEndtime         interface{}   `json:"burn_endtime"`
		Kills               int           `json:"kills"`
		MysticTrial         int           `json:"mystic_trial"`
		PowerRank           int           `json:"power_rank"`
		KillsRank           int           `json:"kills_rank"`
		StoveRank           int           `json:"stove_rank"`
		MysticRank          int           `json:"mystic_rank"`
		TrailsJSON          string        `json:"trails_json"`
		RanksJSON           string        `json:"ranks_json"`
		RecordsJSON         string        `json:"records_json"`
		MapUpdatedAt        float64       `json:"map_updated_at"`
		StatsUpdatedAt      float64       `json:"stats_updated_at"`
		DetailUpdatedAt     float64       `json:"detail_updated_at"`
		ArenaUpdatedAt      float64       `json:"arena_updated_at"`
		FirstSeenAt         float64       `json:"first_seen_at"`
		LastSeenAt          float64       `json:"last_seen_at"`
		MigrantScore        int           `json:"migrant_score"`
		MigrantRank         int           `json:"migrant_rank"`
		MapKid              int           `json:"map_kid"`
		LordsEquipments     []interface{} `json:"lords_equipments"`
		BattleRecord        []interface{} `json:"battle_record"`
		Equipedskin         []struct {
			Skintype int `json:"skintype"`
			ID       int `json:"id"`
			ExpireTs int `json:"expire_ts"`
		} `json:"equipedskin"`
		Glory       interface{}   `json:"glory"`
		ArenaHeroes []interface{} `json:"arena_heroes"`
		Trails      struct {
			Num1 struct {
				TowerID int    `json:"tower_id"`
				Name    string `json:"name"`
				Score   int    `json:"score"`
				Rank    int    `json:"rank"`
				LevelID int    `json:"level_id"`
				Stage   int    `json:"stage"`
			} `json:"1"`
			Num2 struct {
				TowerID int    `json:"tower_id"`
				Name    string `json:"name"`
				Score   int    `json:"score"`
				Rank    int    `json:"rank"`
				LevelID int    `json:"level_id"`
				Stage   int    `json:"stage"`
			} `json:"2"`
			Num3 struct {
				TowerID int    `json:"tower_id"`
				Name    string `json:"name"`
				Score   int    `json:"score"`
				Rank    int    `json:"rank"`
				LevelID int    `json:"level_id"`
				Stage   int    `json:"stage"`
			} `json:"3"`
			Num4 struct {
				TowerID int    `json:"tower_id"`
				Name    string `json:"name"`
				Score   int    `json:"score"`
				Rank    int    `json:"rank"`
				LevelID int    `json:"level_id"`
				Stage   int    `json:"stage"`
			} `json:"4"`
			Num5 struct {
				TowerID int    `json:"tower_id"`
				Name    string `json:"name"`
				Score   int    `json:"score"`
				Rank    int    `json:"rank"`
				LevelID int    `json:"level_id"`
				Stage   int    `json:"stage"`
			} `json:"5"`
			Num6 struct {
				TowerID int    `json:"tower_id"`
				Name    string `json:"name"`
				Score   int    `json:"score"`
				Rank    int    `json:"rank"`
				LevelID int    `json:"level_id"`
				Stage   int    `json:"stage"`
			} `json:"6"`
		} `json:"trails"`
		Ranks struct {
			Num3 struct {
				Type   int   `json:"type"`
				Key    int   `json:"key"`
				Kid    int   `json:"kid"`
				Rank   int   `json:"rank"`
				Score  int   `json:"score"`
				Values []int `json:"values"`
				Raw    struct {
					Rank int `json:"rank"`
					Data struct {
						Key   int   `json:"key"`
						Value []int `json:"value"`
					} `json:"data"`
				} `json:"raw"`
			} `json:"3"`
			Num4 struct {
				Type   int   `json:"type"`
				Key    int   `json:"key"`
				Kid    int   `json:"kid"`
				Rank   int   `json:"rank"`
				Score  int   `json:"score"`
				Values []int `json:"values"`
				Raw    struct {
					Rank int `json:"rank"`
					Data struct {
						Key   int   `json:"key"`
						Value []int `json:"value"`
					} `json:"data"`
				} `json:"raw"`
			} `json:"4"`
			Num5 struct {
				Type   int   `json:"type"`
				Key    int   `json:"key"`
				Kid    int   `json:"kid"`
				Rank   int   `json:"rank"`
				Score  int   `json:"score"`
				Values []int `json:"values"`
				Raw    struct {
					Rank int `json:"rank"`
					Data struct {
						Key   int   `json:"key"`
						Value []int `json:"value"`
					} `json:"data"`
				} `json:"raw"`
			} `json:"5"`
			Num19 struct {
				Type   int   `json:"type"`
				Key    int   `json:"key"`
				Kid    int   `json:"kid"`
				Rank   int   `json:"rank"`
				Score  int   `json:"score"`
				Values []int `json:"values"`
				Raw    struct {
					Rank int `json:"rank"`
					Data struct {
						Key   int   `json:"key"`
						Value []int `json:"value"`
					} `json:"data"`
				} `json:"raw"`
			} `json:"19"`
			Num20 struct {
				Type   int   `json:"type"`
				Key    int   `json:"key"`
				Kid    int   `json:"kid"`
				Rank   int   `json:"rank"`
				Score  int   `json:"score"`
				Values []int `json:"values"`
				Raw    struct {
					Rank int `json:"rank"`
					Data struct {
						Key   int   `json:"key"`
						Value []int `json:"value"`
					} `json:"data"`
				} `json:"raw"`
			} `json:"20"`
			Num21 struct {
				Type   int   `json:"type"`
				Key    int   `json:"key"`
				Kid    int   `json:"kid"`
				Rank   int   `json:"rank"`
				Score  int   `json:"score"`
				Values []int `json:"values"`
				Raw    struct {
					Rank int `json:"rank"`
					Data struct {
						Key   int   `json:"key"`
						Value []int `json:"value"`
					} `json:"data"`
				} `json:"raw"`
			} `json:"21"`
			Num22 struct {
				Type   int   `json:"type"`
				Key    int   `json:"key"`
				Kid    int   `json:"kid"`
				Rank   int   `json:"rank"`
				Score  int   `json:"score"`
				Values []int `json:"values"`
				Raw    struct {
					Rank int `json:"rank"`
					Data struct {
						Key   int   `json:"key"`
						Value []int `json:"value"`
					} `json:"data"`
				} `json:"raw"`
			} `json:"22"`
			Num23 struct {
				Type   int   `json:"type"`
				Key    int   `json:"key"`
				Kid    int   `json:"kid"`
				Rank   int   `json:"rank"`
				Score  int   `json:"score"`
				Values []int `json:"values"`
				Raw    struct {
					Rank int `json:"rank"`
					Data struct {
						Key   int   `json:"key"`
						Value []int `json:"value"`
					} `json:"data"`
				} `json:"raw"`
			} `json:"23"`
			Num24 struct {
				Type   int   `json:"type"`
				Key    int   `json:"key"`
				Kid    int   `json:"kid"`
				Rank   int   `json:"rank"`
				Score  int   `json:"score"`
				Values []int `json:"values"`
				Raw    struct {
					Rank int `json:"rank"`
					Data struct {
						Key   int   `json:"key"`
						Value []int `json:"value"`
					} `json:"data"`
				} `json:"raw"`
			} `json:"24"`
			Num25 struct {
				Type   int   `json:"type"`
				Key    int   `json:"key"`
				Kid    int   `json:"kid"`
				Rank   int   `json:"rank"`
				Score  int   `json:"score"`
				Values []int `json:"values"`
				Raw    struct {
					Rank int `json:"rank"`
					Data struct {
						Key   int   `json:"key"`
						Value []int `json:"value"`
					} `json:"data"`
				} `json:"raw"`
			} `json:"25"`
			Num26 struct {
				Type  int `json:"type"`
				Key   int `json:"key"`
				Kid   int `json:"kid"`
				Rank  int `json:"rank"`
				Score int `json:"score"`
			} `json:"26"`
		} `json:"ranks"`
		Records []struct {
			ID    int `json:"id"`
			Value int `json:"value,omitempty"`
		} `json:"records"`
		RawBase           interface{}   `json:"raw_base"`
		RawDetail         interface{}   `json:"raw_detail"`
		TownCenterLevel   int           `json:"town_center_level"`
		VipLevel          interface{}   `json:"vip_level"`
		LastLogin         string        `json:"last_login"`
		Online            bool          `json:"online"`
		Inactive90D       bool          `json:"inactive_90d"`
		Gear              []interface{} `json:"gear"`
		GovernorID        int           `json:"governor_id"`
		LastChecked       string        `json:"last_checked"`
		ProfilePopulated  bool          `json:"profile_populated"`
		AvatarURL         string        `json:"avatar_url"`
		AvatarImageID     int           `json:"avatar_image_id"`
		AvatarSource      string        `json:"avatar_source"`
		AllianceRankLabel string        `json:"alliance_rank_label"`
		ShieldActive      bool          `json:"shield_active"`
		ShieldRemaining   interface{}   `json:"shield_remaining"`
		ShieldLabel       string        `json:"shield_label"`
		BurnActive        bool          `json:"burn_active"`
		BurnRemaining     interface{}   `json:"burn_remaining"`
		BurnLabel         string        `json:"burn_label"`
		TrailsList        []struct {
			TowerID int    `json:"tower_id"`
			Name    string `json:"name"`
			Score   int    `json:"score"`
			Rank    int    `json:"rank"`
			LevelID int    `json:"level_id"`
			Stage   int    `json:"stage"`
		} `json:"trails_list"`
		Leaderboards []struct {
			Name             string      `json:"name"`
			Value            int         `json:"value"`
			ValueLabel       interface{} `json:"value_label"`
			KingdomRank      int         `json:"kingdom_rank"`
			KingdomRankLabel string      `json:"kingdom_rank_label"`
			Trail            bool        `json:"trail"`
			RankType         int         `json:"rank_type"`
		} `json:"leaderboards"`
		Heroes    []interface{} `json:"heroes"`
		HeroCount int           `json:"hero_count"`
	} `json:"members"`
	MemberCount        int         `json:"member_count"`
	TrackerPower       int64       `json:"tracker_power"`
	TrackerPowerRank   int         `json:"tracker_power_rank"`
	TrackerMemberCount int         `json:"tracker_member_count"`
	KingdomAlliances   int         `json:"kingdom_alliances"`
	Permalink          string      `json:"permalink"`
	Fresh              bool        `json:"fresh"`
	CacheSeconds       int         `json:"cache_seconds"`
	Refreshed          bool        `json:"refreshed"`
	Live               interface{} `json:"live"`
}
