"use client";

import { use, useCallback, useEffect, useMemo, useRef, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import TeamTileClass from "@/components/ui/TeamTileClass";
import TeamTextClass from "@/components/ui/TeamTextClass";
import DirectionalControls, { Direction } from "@/components/DirectionalControls";
import Engine from "@/controller/Engine";
import Item from "@/model/Item";
import Player, { PlayerClass, PlayerInstance } from "@/model/Player";
import { Message } from "@/model/SSEMessage";
import { GameConnected } from "@/model/SSEMessages/GameConnected";
import { Flag, Teams } from "@/model/Tile";
import { Flag as FlagIcon, Package, Sword } from "lucide-react";
import Change from "@/model/Change";
import { Button } from "@/components/ui/button";

type RenderedGrid = ReturnType<Engine["RenderAll"]>;
const DELTAS: Record<Direction, [number, number]> = {
  up: [0, -1],
  down: [0, 1],
  left: [-1, 0],
  right: [1, 0],
};

export default function Page({
  params,
}: {
  params: Promise<{ roomId: string; playerId: string }>;
}) {
  const esRef = useRef<EventSource | null>(null);
  const [renderedGrid, setRenderedGrid] = useState<RenderedGrid | null>(null);
  const [players, setPlayers] = useState<Player[]>([]);
  const [me, setMe] = useState<PlayerInstance | null>(null);
  const [items, setItems] = useState<Item[]>([]);
  const [connected, setConnected] = useState(false);
  const engineRef = useRef<Engine | null>(null);
  const [mapSize, setMapSize] = useState<{ width: number; height: number }>({ width: 0, height: 0 });
  const [sightRadius, setSightRadius] = useState(2);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [actionError, setActionError] = useState<string | null>(null);
  const [actionMessage, setActionMessage] = useState<string | null>(null);
  const [movementCapacity, setMovementCapacity] = useState(0);
  const [remainingMovement, setRemainingMovement] = useState(0);
  const [hasMoved, setHasMoved] = useState(false);
  const [hasAttacked, setHasAttacked] = useState(false);
  const [hasSkipped, setHasSkipped] = useState(false);

  const { roomId, playerId } = use(params);

  useEffect(() => {
    const es = new EventSource(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}api/game/player?roomId=${roomId}&playerId=${playerId}`
    );
    esRef.current = es;

    es.onopen = () => console.log("SSE connected");
    es.onmessage = (event) => {
      console.log("SSE message:", event.data);
      try {
        const data = JSON.parse(event.data) as { MessageType?: string; Message?: unknown };
        if (!data.MessageType) return;
        const baseMessage = new Message(data.MessageType);

        if (baseMessage.MessageType === "ping") {
          return;
        }

        if (baseMessage.MessageType === "Connection") {
          const msgBody = data.Message as GameConnected;

          const engine = new Engine(msgBody.MapWidth, msgBody.MapHeight);
          engine.Initialize(msgBody.Tiles ?? [], msgBody.Players ?? [], msgBody.Items ?? []);
          engineRef.current = engine;
          setMapSize({ width: msgBody.MapWidth, height: msgBody.MapHeight });

          const myPlayer = (msgBody.Players ?? []).find(p => p.Id === playerId);
          if (!myPlayer) {
            throw new Error("I'm not in the game?");
          }

          const myInstance = new PlayerInstance(myPlayer);
          setMe(myInstance);

          const otherPlayers = (msgBody.Players ?? []).filter(p => p.Id !== playerId);
          setPlayers(otherPlayers);
          setItems(msgBody.Items ?? []);
          setConnected(true);

          const classMovement = myInstance.ClassInfo.MovementSpeed;
          setMovementCapacity(classMovement);
          setRemainingMovement(classMovement);
          setHasMoved(false);
          setHasAttacked(false);
          setHasSkipped(false);
          const sight = myInstance.Class === PlayerClass.Thief ? 3 : 2;
          setSightRadius(sight);

          setRenderedGrid(
            engine.RenderParts(
              myInstance.CurrentX - sight,
              myInstance.CurrentY - sight,
              2 * sight + 1,
              2 * sight + 1,
              false
            )
          );
        }
      } catch (error) {
        console.log("Failed to parse SSE message", error);
      }
    };
    es.onerror = (err) => console.log("SSE error", err);

    return () => {
      es.close();
      esRef.current = null;
    };
  }, [playerId, roomId]);

  const grid = renderedGrid ?? [];
  const cols = grid[0]?.length ?? 0;
  const currentPlayer = useMemo(
    () => me ?? players.find((player) => player.Id === playerId),
    [me, players, playerId]
  );

  const centerPosition = useMemo(() => {
    if (!grid.length) return { r: 0, c: 0 };
    return {
      r: Math.floor(grid.length / 2),
      c: Math.floor((grid[0]?.length ?? 1) / 2),
    };
  }, [grid]);

  const renderAroundPlayer = useCallback(
    (player: PlayerInstance) => {
      if (!engineRef.current) return;
      const sight = player.Class === PlayerClass.Thief ? 3 : 2;
      setSightRadius(sight);
      setRenderedGrid(
        engineRef.current.RenderParts(
          player.CurrentX - sight,
          player.CurrentY - sight,
          2 * sight + 1,
          2 * sight + 1,
          false
        )
      );
    },
    []
  );

  const showDirection = useCallback(
    (direction: Direction) => {
      if (!me || !grid.length || hasAttacked || remainingMovement <= 0) return false;
      if (hasSkipped) return false;

      const [dx, dy] = DELTAS[direction];
      const targetX = me.CurrentX + dx;
      const targetY = me.CurrentY + dy;

      if (targetX < 0 || targetY < 0 || targetX >= mapSize.width || targetY >= mapSize.height) {
        return false;
      }

      const targetRow = centerPosition.r + dy;
      const targetCol = centerPosition.c + dx;

      const targetTile = grid[targetRow]?.[targetCol];
      if (!targetTile) return false;

      const [playerNames, , targetFlag, targetTeam] = targetTile;
      const hasBlockingEnemy = !!playerNames && targetTeam !== me.Team;

      return targetFlag !== Flag.INACCESSIBLE && !hasBlockingEnemy;
    },
    [centerPosition.c, centerPosition.r, grid, hasAttacked, hasSkipped, mapSize.height, mapSize.width, me, remainingMovement]
  );

  const handleMove = useCallback(
    async (direction: Direction) => {
      if (!me || !engineRef.current) return;
      if (hasSkipped) {
        setActionError("You have already skipped this turn.");
        return;
      }
      if (hasAttacked) {
        setActionError("You cannot move after attacking.");
        return;
      }
      if (remainingMovement <= 0) {
        setActionError("No movement remaining for this turn.");
        return;
      }
      if (!showDirection(direction)) return;

      const [dx, dy] = DELTAS[direction];
      const targetX = me.CurrentX + dx;
      const targetY = me.CurrentY + dy;

      setIsSubmitting(true);
      setActionError(null);
      setActionMessage(null);

      try {
        const targetTile = grid[centerPosition.r + dy]?.[centerPosition.c + dx];
        const [playerNames, , , targetTeam] = targetTile ?? [];

        if (playerNames && targetTeam !== me.Team) {
          setActionError("An opposing player is blocking that tile.");
          return;
        }

        const response = await fetch(
          `${process.env.NEXT_PUBLIC_API_BASE_URL}api/game/move?roomId=${roomId}`,
          {
            method: "POST",
            body: JSON.stringify({
              X: targetX,
              Y: targetY,
              PrevX: me.CurrentX,
              PrevY: me.CurrentY,
              Id: me.Id,
            }),
          }
        );

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || `Failed to submit move (${response.status})`);
      }

        engineRef.current.Update([
          new Change(targetX, targetY, me.CurrentX, me.CurrentY, "Player", me.Id),
        ]);

        setMe((current) =>
          current
            ? {
              ...current,
              CurrentX: targetX,
              CurrentY: targetY,
              X: targetX,
              Y: targetY,
            }
            : current
        );

        setHasMoved(true);
        setRemainingMovement((prev) => Math.max(prev - 1, 0));
        renderAroundPlayer({
          ...me,
          CurrentX: targetX,
          CurrentY: targetY,
          X: targetX,
          Y: targetY,
        });

        setActionMessage(`Moved to (${targetX}, ${targetY}). Movement left: ${Math.max(remainingMovement - 1, 0)}`);
      } catch (error) {
        setActionError(error instanceof Error ? error.message : "Failed to process move");
      } finally {
        setIsSubmitting(false);
      }
    },
    [hasAttacked, hasSkipped, me, remainingMovement, renderAroundPlayer, roomId, showDirection]
  );

  const handleSkip = useCallback(async () => {
    if (!me) return;
    setIsSubmitting(true);
    setActionError(null);
    setActionMessage(null);

    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}api/game/skip?roomId=${roomId}`,
        {
          method: "POST",
          body: JSON.stringify({ Id: me.Id }),
        }
      );

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || `Failed to skip turn (${response.status})`);
      }

      setHasMoved(false);
      setHasAttacked(true);
      setRemainingMovement(0);
      setHasSkipped(true);
      setActionMessage("Turn skipped");
    } catch (error) {
      setActionError(error instanceof Error ? error.message : "Failed to skip turn");
    } finally {
      setIsSubmitting(false);
    }
  }, [me, roomId]);

  const getAttackTarget = useCallback(
    (direction: Direction) => {
      if (!me) return null;

      const [dx, dy] = DELTAS[direction];
      const targetX = me.CurrentX + dx;
      const targetY = me.CurrentY + dy;

      if (targetX < 0 || targetY < 0 || targetX >= mapSize.width || targetY >= mapSize.height) {
        return null;
      }

      const targetTile = grid[centerPosition.r + dy]?.[centerPosition.c + dx];
      if (!targetTile) return null;

      const [, , targetFlag, targetTeam] = targetTile;
      const enemyPlayer = players.find(
        (player) => player.X === targetX && player.Y === targetY && player.Team !== me.Team
      );
      const enemyChest = targetFlag === Flag.TREASURE_CHEST && targetTeam !== Teams.Neutral && targetTeam !== me.Team;

      if (enemyPlayer) {
        return { type: "player" as const, targetId: enemyPlayer.Id, coords: { x: targetX, y: targetY } };
      }

      if (enemyChest) {
        return { type: "chest" as const, coords: { x: targetX, y: targetY } };
      }

      return null;
    },
    [centerPosition.c, centerPosition.r, grid, mapSize.height, mapSize.width, me, players]
  );

  const canAttackDirection = useCallback(
    (direction: Direction) => {
      if (!me || hasAttacked || hasSkipped) return false;
      if (hasMoved) return false;

      return !!getAttackTarget(direction);
    },
    [getAttackTarget, hasAttacked, hasMoved, hasSkipped, me]
  );

  const handleAttack = useCallback(
    async (direction: Direction) => {
      if (!me) return;

      if (hasSkipped) {
        setActionError("You have already skipped this turn.");
        return;
      }

      if (hasAttacked) {
        setActionError("You have already attacked this turn.");
        return;
      }

      if (hasMoved) {
        setActionError("You cannot attack after moving.");
        return;
      }

      const target = getAttackTarget(direction);
      if (!target) {
        setActionError("No valid target to attack in that direction.");
        return;
      }

      setIsSubmitting(true);
      setActionError(null);
      setActionMessage(null);

      try {
        const isPlayerTarget = target.type === "player";
        const endpoint = isPlayerTarget ? "attack" : "bonusAttack";
        const response = await fetch(
          `${process.env.NEXT_PUBLIC_API_BASE_URL}api/game/${endpoint}?roomId=${roomId}`,
          {
            method: "POST",
            body: JSON.stringify(
              isPlayerTarget
                ? {
                  AttackerId: me.Id,
                  DefenderId: target.targetId,
                }
                : {
                  X: target.coords.x,
                  Y: target.coords.y,
                  Id: me.Id,
                }
            ),
          }
        );

        if (!response.ok) {
          const errorText = await response.text();
          throw new Error(errorText || `Failed to submit attack (${response.status})`);
        }

        setHasAttacked(true);
        setRemainingMovement(0);
        setActionMessage(
          target.type === "player"
            ? "Attack submitted against the enemy player."
            : "Attack submitted against the enemy treasure chest."
        );
      } catch (error) {
        setActionError(error instanceof Error ? error.message : "Failed to process attack");
      } finally {
        setIsSubmitting(false);
      }
    },
    [getAttackTarget, hasAttacked, hasMoved, hasSkipped, me, roomId]
  );

  return (
    <main className="mx-auto w-full max-w-4xl p-6">
      <Card className="w-full">
        <CardHeader className="space-y-2">
          <CardTitle>Game View</CardTitle>
          <div className="text-sm text-muted-foreground">
            <div>Room: {roomId}</div>
            <div>Player: {currentPlayer?.Name ?? playerId}</div>
          </div>
        </CardHeader>
        <CardContent className="space-y-6">
          {!connected ? (
            <p className="text-sm text-muted-foreground">
              Connecting to the game&hellip;
            </p>
          ) : (
            <div className="space-y-4">
              <div className="relative w-fit">
                <div className="inline-grid rounded-xl border border-border overflow-hidden"
                  style={{ gridTemplateColumns: `repeat(${cols ?? 0}, 5rem)` }}
                >
                  {grid.flatMap((row, ri) =>
                    row.map(([playerNames, itemLabels, flag, team], ci) => {
                      const isInaccessible = flag === Flag.INACCESSIBLE;
                      const isCenterTile = me && ri === centerPosition.r && ci === centerPosition.c;
                      const visibleNames = (() => {
                        if (!isCenterTile) return playerNames;
                        if (!playerNames) return "";
                        const parts = playerNames.split(",").map((p) => p.trim()).filter(Boolean);
                        return parts.filter((p) => p !== me?.Name && p !== me?.Name?.[0]).join(",");
                      })();
                      const tileClasses = [
                        "relative flex h-20 w-20 flex-col items-center justify-center gap-1 text-lg font-semibold",
                        "border border-border -ml-px -mt-px",
                        "transition-colors",
                        isInaccessible ? "bg-black text-white" : TeamTileClass(team),
                        isCenterTile ? "ring-2 ring-emerald-500 ring-offset-2 ring-offset-background" : "",
                      ].join(" ");

                      const playerNameClasses = [
                        "leading-none",
                        isInaccessible ? "text-white" : TeamTextClass(team),
                      ].join(" ");

                      const itemLabelClasses = [
                        "text-xs font-normal leading-none",
                        isInaccessible ? "text-white/80" : "text-muted-foreground",
                      ].join(" ");

                      return (
                        <div
                          key={`${ri}-${ci}`}
                          className={tileClasses}
                        >
                          {flag === Flag.TREASURE_CHEST ? (
                            <div className="absolute top-1 right-1">
                              <Package className="h-4 w-4 text-amber-600/80 dark:text-amber-400/80" />
                            </div>
                          ) : null}
                          {flag === Flag.SPAWN ? (
                            <div className="absolute top-1 right-1">
                              <FlagIcon className="h-4 w-4 text-amber-600/80 dark:text-amber-400/80" />
                            </div>
                          ) : null}
                          <span className={playerNameClasses}>
                            {visibleNames || "\u00a0"}
                          </span>
                          {itemLabels ? (
                            <span className={itemLabelClasses}>
                              {itemLabels}
                            </span>
                          ) : null}
                          {isCenterTile ? (
                            <span className="absolute bottom-1 right-1 rounded-full bg-emerald-500 px-2 py-0.5 text-xs font-semibold text-white shadow">
                              Self
                            </span>
                          ) : null}
                        </div>
                      );
                    })
                  )}
                </div>

                {me ? (
                  <DirectionalControls
                    gridSize={{ cols, rows: grid.length }}
                    onSelect={handleMove}
                    position={centerPosition}
                    showDirection={showDirection}
                  />
                ) : null}
              </div>

              {me ? (
                <div className="space-y-2">
                  <div className="flex items-center gap-2 text-sm font-semibold text-foreground">
                    <Sword className="h-4 w-4" />
                    <span>Attack</span>
                  </div>
                  <div className="grid w-fit grid-cols-3 gap-2">
                    <div />
                    <Button
                      size="sm"
                      variant="destructive"
                      disabled={isSubmitting || !canAttackDirection("up")}
                      onClick={() => handleAttack("up")}
                    >
                      Attack ↑
                    </Button>
                    <div />
                    <Button
                      size="sm"
                      variant="destructive"
                      disabled={isSubmitting || !canAttackDirection("left")}
                      onClick={() => handleAttack("left")}
                    >
                      Attack ←
                    </Button>
                    <div />
                    <Button
                      size="sm"
                      variant="destructive"
                      disabled={isSubmitting || !canAttackDirection("right")}
                      onClick={() => handleAttack("right")}
                    >
                      Attack →
                    </Button>
                    <div />
                    <Button
                      size="sm"
                      variant="destructive"
                      disabled={isSubmitting || !canAttackDirection("down")}
                      onClick={() => handleAttack("down")}
                    >
                      Attack ↓
                    </Button>
                    <div />
                  </div>
                </div>
              ) : null}

              <div className="flex flex-wrap items-center gap-3 text-sm">
                <Button
                  size="sm"
                  variant="outline"
                  disabled={isSubmitting}
                  onClick={handleSkip}
                >
                  Skip Turn
                </Button>

                {actionMessage ? (
                  <span className="text-emerald-600 dark:text-emerald-400">{actionMessage}</span>
                ) : null}
                {actionError ? (
                  <span className="text-destructive">{actionError}</span>
                ) : null}
              </div>

              <div className="space-y-1 text-xs text-muted-foreground">
                <div>
                  Movement remaining: {remainingMovement} / {movementCapacity}
                </div>
                <div>
                  Attack status: {hasSkipped
                    ? "Turn skipped"
                    : hasAttacked
                      ? "Already attacked"
                      : hasMoved
                        ? "Unavailable after moving"
                        : "Available"}
                </div>
              </div>
            </div>
          )}

          {connected ? (
            <div className="space-y-1 text-sm text-muted-foreground">
              <div className="font-semibold text-foreground">Players</div>
              <ul className="list-disc pl-4 space-y-1">
                {players.map((player) => (
                  <li key={player.Id}>
                    <span className={TeamTextClass(player.Team)}>
                      {player.Name}
                    </span>{" "}
                    ({player.Class}) at ({player.X}, {player.Y})
                  </li>
                ))}
              </ul>
            </div>
          ) : null}
        </CardContent>
      </Card>
    </main>
  );
}
