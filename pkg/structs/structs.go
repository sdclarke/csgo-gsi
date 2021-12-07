package structs

type State struct {
	Provider   *provider
	Map        *csmap
	Round      *round
	Player     *player
	AllPlayers map[string]*player
	Previously *State
	Added      *State
	Auth       *auth
}

type provider struct {
	Appid     int64
	Name      string
	Steamid   string
	Timestamp float64
	Version   int64
}

type csmap struct {
	Name   string
	Phase  string
	Round  int64
	TeamCT *team `json:"team_ct"`
	TeamT  *team `json:"team_t"`
}

type round struct {
	Phase   string
	WinTeam string `json:"win_team"`
	Bomb    string
}

type player struct {
	Activity   string
	MatchStats *matchStats `json:"match_stats"`
	Name       string
	Team       string
	State      *playerState
	Steamid    string
	Weapons    map[string]*weapon
}

type team struct {
	Score                int64
	ConsecutiveRoundLoss int64 `json:"consecutive_round_loss"`
}

type matchStats struct {
	Assists int64
	Deaths  int64
	Kills   int64
	Mvps    int64
	Score   int64
}

type playerState struct {
	Armor       int64
	Burning     int64
	EquipValue  int64 `json:"equip_value"`
	Flashed     int64
	Health      int64
	Helmet      bool
	Money       int64
	RoundKillhs int64 `json:"round_killhs"`
	RoundKills  int64 `json:"round_kills"`
	Smoked      int64
}

type weapon struct {
	Name        string
	PaintKit    string `json:"paint_kit"`
	State       string
	Type        string
	AmmoClip    int64 `json:"ammo_clip"`
	AmmoClipMax int64 `json:"ammo_clip_max"`
	AmmoReserve int64 `json:"ammo_reserve"`
}
type auth struct {
	Token string
}
