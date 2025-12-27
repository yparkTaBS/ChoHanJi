package handlers

type RouteToken string

var (
	POSTRoom               RouteToken = "/api/room"
	POSTCharacter          RouteToken = "/api/character"
	GETPlayerEvent         RouteToken = "/api/player/event"
	GETRoomAdmin           RouteToken = "/api/waiting/room/admin"
	POSTGameStart          RouteToken = "/api/game/start"
	GETPlayerGameStatus    RouteToken = "/api/game/player"
	GETAdminGameStatus     RouteToken = "/api/game/admin"
	POSTSubmitMoves        RouteToken = "/api/game/move"
	POSTSubmitAttacks      RouteToken = "/api/game/attack"
	POSTSubmitAttackResult RouteToken = "/api/game/attack/result"
	POSTSubmitBonusAttacks RouteToken = "/api/game/bonusAttack"
	POSTSubmitSkip         RouteToken = "/api/game/skip"
	POSTProceed            RouteToken = "/api/game/proceed"
)
