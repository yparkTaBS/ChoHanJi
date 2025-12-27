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
import { Flag } from "@/model/Tile";
import { Flag as FlagIcon, Package } from "lucide-react";
import Change from "@/model/Change";
import { Button } from "@/components/ui/button";

type RenderedGrid = ReturnType<Engine["RenderAll"]>;

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

      const deltas: Record<Direction, [number, number]> = {
        up: [0, -1],
        down: [0, 1],
        left: [-1, 0],
        right: [1, 0],
      };

      const [dx, dy] = deltas[direction];
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

      const deltas: Record<Direction, [number, number]> = {
        up: [0, -1],
        down: [0, 1],
        left: [-1, 0],
        right: [1, 0],
      };

      const [dx, dy] = deltas[direction];
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
