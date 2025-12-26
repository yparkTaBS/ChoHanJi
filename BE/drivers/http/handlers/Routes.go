package handlers

type RouteToken string

var (
	POSTRoom            RouteToken = "/api/room"
	POSTCharacter       RouteToken = "/api/character"
	GETPlayerEvent      RouteToken = "/api/player/event"
	GETRoomAdmin        RouteToken = "/api/waiting/room/admin"
	POSTGameStart       RouteToken = "/api/game/start"
	GETPlayerGameStatus RouteToken = "/api/game/player"
	GETAdminGameStatus  RouteToken = "/api/game/admin"

	POSTSubmitMoves RouteToken = "/api/game/Move"
)
