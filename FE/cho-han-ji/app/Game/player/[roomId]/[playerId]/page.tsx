"use client";

import { use, useEffect, useMemo, useRef, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import TeamTileClass from "@/components/ui/TeamTileClass";
import TeamTextClass from "@/components/ui/TeamTextClass";
import Engine from "@/controller/Engine";
import Item from "@/model/Item";
import Player, { PlayerClass, PlayerInstance } from "@/model/Player";
import { Message } from "@/model/SSEMessage";
import { GameConnected } from "@/model/SSEMessages/GameConnected";
import { Flag } from "@/model/Tile";
import { Flag as FlagIcon, Package } from "lucide-react";

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

          const sight = myInstance.Class === PlayerClass.Thief ? 3 : 2;

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
    () => players.find((player) => player.Id === playerId),
    [players, playerId]
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
            <div className="inline-grid rounded-xl border border-border overflow-hidden"
              style={{ gridTemplateColumns: `repeat(${cols ?? 0}, 5rem)` }}
            >
              {grid.flatMap((row, ri) =>
                row.map(([playerNames, itemLabels, flag, team], ci) => {
                  const isInaccessible = flag === Flag.INACCESSIBLE;
                  const tileClasses = [
                    "relative flex h-20 w-20 flex-col items-center justify-center gap-1 text-lg font-semibold",
                    "border border-border -ml-px -mt-px",
                    "transition-colors",
                    isInaccessible ? "bg-black text-white" : TeamTileClass(team),
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
                        {playerNames || "\u00a0"}
                      </span>
                      {itemLabels ? (
                        <span className={itemLabelClasses}>
                          {itemLabels}
                        </span>
                      ) : null}
                    </div>
                  );
                })
              )}
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
