package handlers

type RouteToken string

var (
	POSTRoom       RouteToken = "/api/room"
	POSTCharacter  RouteToken = "/api/character"
	GETPlayerEvent RouteToken = "/api/player/event"
	GETRoomAdmin   RouteToken = "/api/waiting/room/admin"
)
