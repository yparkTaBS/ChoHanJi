"use client";

import { use, useEffect, useRef, useState, useMemo } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import Engine from "@/controller/Engine";
import { Message, TypedMessage } from "@/model/SSEMessage";
import { GameConnected } from "@/model/SSEMessages/GameConnected";
import TeamTileClass from "@/components/ui/TeamTileClass";
import TeamTextClass from "@/components/ui/TeamTextClass";
import { Package, Flag as FlagIcon } from "lucide-react";
import { Flag } from "@/model/Tile"

export default function Page({ params }: { params: Promise<{ roomId: string }> }) {
  const esRef = useRef<EventSource | null>(null);
  const engineRef = useRef<Engine | null>(null);
  type RenderedGrid = ReturnType<Engine["RenderAll"]>;
  type RenderedTile = RenderedGrid[number][number];
  const [renderedGrid, setRenderedGrid] = useState<RenderedGrid | null>();
  const { roomId } = use(params);

  const grid = renderedGrid ?? [];
  const cols = grid[0]?.length ?? 0;

  useEffect(() => {
    const es = new EventSource(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}api/game/admin?roomId=${roomId}`
    );
    esRef.current = es;

    es.onopen = () => console.log("SSE connected");
    es.onmessage = (event) => {
      console.log("SSE raw:", event.data);
      try {
        const data = JSON.parse(event.data) as { MessageType: string; Message?: unknown };

        const baseMessage = new Message(data.MessageType);

        if (baseMessage.MessageType === "Connection") {
          if (!data.Message) return;

          const msgBody = data.Message as GameConnected;

          const snapshot = msgBody;
          const engine = new Engine(snapshot.MapWidth, snapshot.MapHeight);
          engine.Initialize(msgBody.Tiles, msgBody.Players, msgBody.Items);
          engineRef.current = engine;

          setRenderedGrid(engine.RenderAll());
        }
      } catch (e) {
        console.log(e);
        return;
      }
    };
    es.onerror = (err) => console.log("SSE error", err);

    return () => {
      es.close();
      esRef.current = null;
    };
  }, [roomId]);

  return (
    <Card className="w-fit">
      <CardHeader>
        <CardTitle>Map</CardTitle>
      </CardHeader>
      <CardContent className="space-y-6">
        {/* Map */}
        <div className="inline-grid rounded-xl border border-border overflow-hidden"
          style={{ gridTemplateColumns: `repeat(${cols ?? 0}, 5rem)` }}
        >
          {grid.flatMap((row, ri) =>
            row.map(([players, items, flag, team], ci) => (
              <div key={`${ri}-${ci}`}
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
                <span className={["leading-none", TeamTextClass(team)].join(" ")}>{players || "\u00a0"}</span>
              </div>
            ))
          )}
        </div>
      </CardContent>
    </Card>
  );
}

