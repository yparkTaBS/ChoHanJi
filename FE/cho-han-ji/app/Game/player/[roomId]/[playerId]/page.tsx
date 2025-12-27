"use client";

import { use, useEffect, useMemo, useRef, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import TeamTileClass from "@/components/ui/TeamTileClass";
import TeamTextClass from "@/components/ui/TeamTextClass";
import Engine from "@/controller/Engine";
import Item from "@/model/Item";
import Player from "@/model/Player";
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

          setRenderedGrid(engine.RenderAll());
          setPlayers(msgBody.Players ?? []);
          setItems(msgBody.Items ?? []);
          setConnected(true);
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
                row.map(([playerNames, itemLabels, flag, team], ci) => (
                  <div
                    key={`${ri}-${ci}`}
                    className={[
                      "relative flex h-20 w-20 flex-col items-center justify-center gap-1 text-lg font-semibold",
                      "border border-border -ml-px -mt-px",
                      "transition-colors",
                      TeamTileClass(team),
                    ].join(" ")}
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
                    <span className={["leading-none", TeamTextClass(team)].join(" ")}>
                      {playerNames || "\u00a0"}
                    </span>
                    {itemLabels ? (
                      <span className="text-xs font-normal text-muted-foreground leading-none">
                        {itemLabels}
                      </span>
                    ) : null}
                  </div>
                ))
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

              {items.length > 0 ? (
                <>
                  <div className="font-semibold text-foreground pt-2">Items</div>
                  <ul className="list-disc pl-4 space-y-1">
                    {items.map((item) => (
                      <li key={item.Id} className="text-muted-foreground">
                        {item.Name} at ({item.X}, {item.Y})
                      </li>
                    ))}
                  </ul>
                </>
              ) : null}
            </div>
          ) : null}
        </CardContent>
      </Card>
    </main>
  );
}
