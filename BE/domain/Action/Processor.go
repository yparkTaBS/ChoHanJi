package Action

import (
	"ChoHanJi/domain/Death"
	"ChoHanJi/domain/Fight"
	"ChoHanJi/domain/Game"
	"ChoHanJi/domain/Player"
	"ChoHanJi/domain/PlayerBlocker"
	"ChoHanJi/domain/Room"
	"ChoHanJi/domain/Team"
	"ChoHanJi/domain/TileFlag"
	"ChoHanJi/domain/UpdateMessage"
	"ChoHanJi/driven/sse/SSEHub"
	crand "crypto/rand"
	"encoding/json"
	"errors"
	"math/big"
)

type CurrentFights interface {
	Create(roomId Room.Id, gameType Game.Type, attId, defId Player.Id) (*Fight.Struct, error)
}

var _ CurrentFights = (*Fight.CurrentFights)(nil)

type DeathList interface {
	CheckIfDead(roomId Room.Id, playerId Player.Id) bool
	GetListOfDead(roomId Room.Id) []Player.Id
	Reset(roomId Room.Id)
}

var _ DeathList = (*Death.List)(nil)

type IPlayerBlocker interface {
	Block(roomId Room.Id, playerId Player.Id) error
	BlockPair(roomId Room.Id, p1, p2 Player.Id) error
	Unblock(roomId Room.Id, playerId Player.Id) error
	WaitUntilUnblocked(roomId Room.Id, playerId Player.Id) error
	WaitUntilAllAreUnblocked(roomId Room.Id) error
}

var _ IPlayerBlocker = (*PlayerBlocker.Struct)(nil)

type IHub interface {
	Publish(roomId, adminId, messageType, messageBody string) error
	PublishToAll(roomId, messageType, messageBody string) error
}

var _ IHub = (*SSEHub.Struct)(nil)

type Processor struct {
	r   Room.Rooms
	cf  CurrentFights
	dl  DeathList
	pb  IPlayerBlocker
	hub IHub
}

func NewProcessor(r Room.Rooms, cf CurrentFights, dl DeathList, pb IPlayerBlocker, hub IHub) *Processor {
	return &Processor{r, cf, dl, pb, hub}
}

func (p *Processor) Process(roomId Room.Id, attacks []AttackStruct, moves []MoveStruct, bonusAttacks []BonusAttackStruct) error {
	defer func() {
		p.dl.Reset(roomId)
	}()

	changes := UpdateMessage.New()

	fm, found := p.r[roomId]
	if !found {
		return errors.New("room not found")
	}

	// Track which dead players we've already processed so we don't double-drop / double-respawn.
	processedDead := make(map[Player.Id]struct{})

	// Resolve newly-dead players:
	// - drop item where they died
	// - move to spawn immediately (removes them from old tile)
	resolveNewDeaths := func() error {
		for _, deadId := range p.dl.GetListOfDead(roomId) {
			if _, done := processedDead[deadId]; done {
				continue
			}

			player, ok := fm.Players[deadId]
			if !ok {
				return errors.New("player not found")
			}

			deathX, deathY := player.X, player.Y

			// 1) Drop item at death tile BEFORE moving to spawn
			if player.Bag != nil {
				deathTile, err := fm.Map.GetTile(deathX, deathY)
				if err != nil {
					return err
				}

				item := player.Bag
				prevX, prevY := item.X, item.Y

				item.X, item.Y = deathX, deathY
				deathTile.AddItem(item)

				player.Bag = nil
				changes.UpsertPlayer(deadId, deathX, deathY, deathX, deathY, nil)
				changes.UpsertItem(item.Id, deathX, deathY, prevX, prevY)
			}

			// 2) Move dead player to spawn immediately
			spawnX, spawnY := fm.Map.GetSpawn(Team.Enum(player.TeamNumber))

			changes.UpsertPlayer(deadId, spawnX, spawnY, deathX, deathY, nil)

			if err := p.movePlayerOnMap(fm, player, spawnX, spawnY, deathX, deathY); err != nil {
				return err
			}

			processedDead[deadId] = struct{}{}
		}
		return nil
	}

	// -------------------
	// Phase: Attack
	// -------------------
	if err := p.hub.PublishToAll(string(roomId), "Phase", "Attack"); err != nil {
		return err
	}

	for _, attack := range attacks {
		attackerId := attack.AttackerId
		defenderId := attack.DefenderId

		// Dead cannot act
		if p.dl.CheckIfDead(roomId, attackerId) {
			continue
		}

		if defenderId == "treasure" {
			attacker, found := fm.Players[attackerId]
			if !found {
				return errors.New("attacker does not exist")
			}
			if err := fm.Map.DisperseItems(Team.Enum(attacker.TeamNumber)); err != nil {
				return err
			}
			continue
		}

		// Dead cannot receive attacks
		if p.dl.CheckIfDead(roomId, defenderId) {
			continue
		}

		if err := p.startFight(roomId, attackerId, defenderId); err != nil {
			return err
		}
	}

	if err := p.pb.WaitUntilAllAreUnblocked(roomId); err != nil {
		return err
	}

	// ✅ Immediately remove dead from tiles and put them at spawn (drop item where they died).
	if err := resolveNewDeaths(); err != nil {
		return err
	}

	// -------------------
	// Phase: Move
	// -------------------
	if err := p.hub.PublishToAll(string(roomId), "Phase", "Move"); err != nil {
		return err
	}

	for _, move := range moves {
		playerId := move.Id

		// Dead should not move
		if p.dl.CheckIfDead(roomId, playerId) {
			continue
		}

		pl, ok := fm.Players[playerId]
		if !ok {
			return errors.New("player not found")
		}

		changes.UpsertPlayer(playerId, move.X, move.Y, move.PrevX, move.PrevY, nil)

		if err := p.movePlayerOnMap(fm, pl, move.X, move.Y, move.PrevX, move.PrevY); err != nil {
			return err
		}
	}

	// -------------------
	// Phase: BonusAttack
	// -------------------
	if err := p.hub.PublishToAll(string(roomId), "Phase", "BonusAttack"); err != nil {
		return err
	}

	for _, bonusAttack := range bonusAttacks {
		attackerId := bonusAttack.Id

		// Dead cannot bonus attack
		if p.dl.CheckIfDead(roomId, attackerId) {
			continue
		}

		attacker, found := fm.Players[attackerId]
		if !found {
			return errors.New("player not found")
		}

		tile, err := fm.Map.GetTile(bonusAttack.X, bonusAttack.Y)
		if err != nil {
			return err
		}

		// Find alive enemy defender
		var defenderId Player.Id
		for _, pl := range tile.Player {
			if pl.TeamNumber == attacker.TeamNumber {
				continue
			}
			if p.dl.CheckIfDead(roomId, pl.Id) {
				continue
			}
			defenderId = pl.Id
			break
		}
		if defenderId == "" {
			continue
		}

		if err := p.startFight(roomId, attackerId, defenderId); err != nil {
			return err
		}
	}

	if err := p.pb.WaitUntilAllAreUnblocked(roomId); err != nil {
		return err
	}

	// ✅ Immediately resolve newly-dead from BonusAttack.
	if err := resolveNewDeaths(); err != nil {
		return err
	}

	// -------------------
	// Phase: CollisionResolution
	// -------------------
	if err := p.hub.PublishToAll(string(roomId), "Phase", "CollisionResolution"); err != nil {
		return err
	}

	tiles := fm.Map.GetNonEmptyTiles()

	for _, tile := range tiles {
		// skip special
		if tile.Flag == TileFlag.SPAWN || tile.Flag == TileFlag.TREASURE_CHEST {
			continue
		}

		if len(tile.Player) < 2 {
			continue
		}

		var t0 []Player.Id
		var t1 []Player.Id

		for _, pl := range tile.Player {
			if p.dl.CheckIfDead(roomId, pl.Id) {
				continue
			}
			if pl.TeamNumber == 0 {
				t0 = append(t0, pl.Id)
			} else {
				t1 = append(t1, pl.Id)
			}
		}

		if len(t0) == 0 || len(t1) == 0 {
			continue
		}

		removeAtSwap := func(list []Player.Id, idx int) []Player.Id {
			last := len(list) - 1
			list[idx] = list[last]
			return list[:last]
		}

		side, err := randIndex(2)
		if err != nil {
			return err
		}

		var champ Player.Id
		champTeam := side

		if champTeam == 0 {
			idx, err := randIndex(len(t0))
			if err != nil {
				return err
			}
			champ = t0[idx]
			t0 = removeAtSwap(t0, idx)
		} else {
			idx, err := randIndex(len(t1))
			if err != nil {
				return err
			}
			champ = t1[idx]
			t1 = removeAtSwap(t1, idx)
		}

		for len(t0) > 0 && len(t1) > 0 {
			if p.dl.CheckIfDead(roomId, champ) {
				// champ already resolved + moved to spawn, replace
				if champTeam == 0 {
					if len(t0) == 0 {
						break
					}
					idx, err := randIndex(len(t0))
					if err != nil {
						return err
					}
					champ = t0[idx]
					t0 = removeAtSwap(t0, idx)
				} else {
					if len(t1) == 0 {
						break
					}
					idx, err := randIndex(len(t1))
					if err != nil {
						return err
					}
					champ = t1[idx]
					t1 = removeAtSwap(t1, idx)
				}
				continue
			}

			var challenger Player.Id
			var challengerIdx int

			if champTeam == 0 {
				challengerIdx, err = randIndex(len(t1))
				if err != nil {
					return err
				}
				challenger = t1[challengerIdx]
			} else {
				challengerIdx, err = randIndex(len(t0))
				if err != nil {
					return err
				}
				challenger = t0[challengerIdx]
			}

			if p.dl.CheckIfDead(roomId, challenger) {
				// should not happen often if resolveNewDeaths ran, but be safe
				if champTeam == 0 {
					t1 = removeAtSwap(t1, challengerIdx)
				} else {
					t0 = removeAtSwap(t0, challengerIdx)
				}
				continue
			}

			if err := p.startFight(roomId, champ, challenger); err != nil {
				return err
			}

			if err := p.pb.WaitUntilUnblocked(roomId, champ); err != nil {
				return err
			}
			if err := p.pb.WaitUntilUnblocked(roomId, challenger); err != nil {
				return err
			}

			// ✅ Immediately process deaths from this collision fight
			if err := resolveNewDeaths(); err != nil {
				return err
			}

			champDead := p.dl.CheckIfDead(roomId, champ)
			challDead := p.dl.CheckIfDead(roomId, challenger)

			if !champDead && !challDead {
				challDead = true
			}

			if challDead {
				if champTeam == 0 {
					t1 = removeAtSwap(t1, challengerIdx)
				} else {
					t0 = removeAtSwap(t0, challengerIdx)
				}
				continue
			}

			if champDead {
				// challenger becomes champ
				champ = challenger
				if champTeam == 0 {
					t1 = removeAtSwap(t1, challengerIdx)
					champTeam = 1
				} else {
					t0 = removeAtSwap(t0, challengerIdx)
					champTeam = 0
				}
				continue
			}
		}
	}

	if err := p.pb.WaitUntilAllAreUnblocked(roomId); err != nil {
		return err
	}

	// ✅ Resolve any remaining deaths (safety)
	if err := resolveNewDeaths(); err != nil {
		return err
	}

	// -------------------
	// Item pickup (alive only)
	// -------------------
	tiles = fm.Map.GetNonEmptyTiles()
	for _, tile := range tiles {
		if len(tile.Items) == 0 || len(tile.Player) == 0 {
			continue
		}

		if err := shufflePlayers(tile.Player); err != nil {
			return err
		}

		for _, pl := range tile.Player {
			if p.dl.CheckIfDead(roomId, pl.Id) {
				continue
			}
			if pl.Bag != nil {
				continue
			}
			if len(tile.Items) == 0 {
				break
			}

			it := tile.Items[0]
			tile.Items = tile.Items[1:]

			prevX, prevY := it.X, it.Y
			it.X, it.Y = -1, -1
			changes.UpsertItem(it.Id, -1, -1, prevX, prevY)

			pl.Bag = it

			itemId := it.Id
			changes.UpsertPlayer(pl.Id, pl.X, pl.Y, pl.X, pl.Y, &itemId)
		}
	}

	env := UpdateEnvelope{
		MessageType: "Update",
		Message:     changes,
	}

	payload, err := json.Marshal(env)
	if err != nil {
		return err
	}

	if err := p.hub.PublishToAll(string(roomId), "Update", string(payload)); err != nil {
		return err
	}

	return nil
}

func getRandomGame() (Game.Type, error) {
	nBig, err := crand.Int(crand.Reader, big.NewInt(int64(len(Game.List))))
	if err != nil {
		return "", err
	}
	return Game.List[int(nBig.Int64())], nil
}

func (p *Processor) startFight(roomId Room.Id, attackerId, defenderId Player.Id) error {
	if err := p.pb.WaitUntilUnblocked(roomId, attackerId); err != nil {
		return err
	}
	if err := p.pb.WaitUntilUnblocked(roomId, defenderId); err != nil {
		return err
	}

	if err := p.pb.BlockPair(roomId, attackerId, defenderId); err != nil {
		return err
	}

	game, err := getRandomGame()
	if err != nil {
		_ = p.pb.Unblock(roomId, attackerId)
		_ = p.pb.Unblock(roomId, defenderId)
		return err
	}

	fight, err := p.cf.Create(roomId, game, attackerId, defenderId)
	if err != nil {
		_ = p.pb.Unblock(roomId, attackerId)
		_ = p.pb.Unblock(roomId, defenderId)
		return err
	}

	msg, err := json.Marshal(fight)
	if err != nil {
		_ = p.pb.Unblock(roomId, attackerId)
		_ = p.pb.Unblock(roomId, defenderId)
		return err
	}

	if err := p.hub.Publish(string(roomId), string(attackerId), "Fight", string(msg)); err != nil {
		_ = p.pb.Unblock(roomId, attackerId)
		_ = p.pb.Unblock(roomId, defenderId)
		return err
	}
	if err := p.hub.Publish(string(roomId), string(defenderId), "Fight", string(msg)); err != nil {
		_ = p.pb.Unblock(roomId, attackerId)
		_ = p.pb.Unblock(roomId, defenderId)
		return err
	}

	return nil
}

func randIndex(n int) (int, error) {
	if n <= 0 {
		return 0, errors.New("randIndex: n must be > 0")
	}
	nBig, err := crand.Int(crand.Reader, big.NewInt(int64(n)))
	if err != nil {
		return 0, err
	}
	return int(nBig.Int64()), nil
}

func shufflePlayers(players []*Player.Struct) error {
	// Fisher–Yates shuffle using randIndex (crypto/rand)
	for i := len(players) - 1; i > 0; i-- {
		j, err := randIndex(i + 1)
		if err != nil {
			return err
		}
		players[i], players[j] = players[j], players[i]
	}
	return nil
}

func (p *Processor) movePlayerOnMap(fm *Room.Room, player *Player.Struct, newX, newY, prevX, prevY int) error {
	// No movement: still ensure authoritative coords are correct
	if prevX == newX && prevY == newY {
		player.X = newX
		player.Y = newY
		return nil
	}

	oldTile, err := fm.Map.GetTile(prevX, prevY)
	if err != nil {
		return err
	}

	newTile, err := fm.Map.GetTile(newX, newY)
	if err != nil {
		return err
	}

	oldTile.RemovePlayer(player.Id)
	newTile.AddPlayer(player)

	player.X = newX
	player.Y = newY
	return nil
}

type UpdateEnvelope struct {
	MessageType string      `json:"MessageType"`
	Message     interface{} `json:"Message"`
}
