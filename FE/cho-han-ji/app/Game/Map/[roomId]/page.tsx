"use client";

import { use, useEffect, useRef, useState, useMemo } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import Engine from "@/controller/Engine";
import { Message } from "@/model/SSEMessage";
import { GameConnected } from "@/model/SSEMessages/GameConnected";
import TeamTileClass from "@/components/ui/TeamTileClass";
import TeamTextClass from "@/components/ui/TeamTextClass";
import { Package, Flag as FlagIcon } from "lucide-react";
import Player from "@/model/Player";
import { Flag, Teams } from "@/model/Tile";

export default function Page({ params }: { params: Promise<{ roomId: string }> }) {
  const esRef = useRef<EventSource | null>(null);
  type RenderedGrid = ReturnType<Engine["RenderAll"]>;
  const [renderedGrid, setRenderedGrid] = useState<RenderedGrid | null>();
  const [players, setPlayers] = useState<Player[]>([]);
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

          setRenderedGrid(engine.RenderAll());
          setPlayers(msgBody.Players ?? []);
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

  const teamPlayers = useMemo(() => {
    const team1 = players.filter((player) => player.Team === Teams.TEAM1);
    const team2 = players.filter((player) => player.Team === Teams.TEAM2);

    return { team1, team2 };
  }, [players]);

  const renderTeam = (team: Teams, title: string, teamMembers: Player[]) => (
    <div className="flex-1 rounded-lg border border-border bg-muted/30 p-3">
      <div className={["mb-2 text-sm font-semibold", TeamTextClass(team)].join(" ")}>{title}</div>
      <ul className="space-y-1 text-sm text-muted-foreground">
        {teamMembers.length > 0 ? (
          teamMembers.map((player) => <li key={player.Id}>{player.Name}</li>)
        ) : (
          <li className="italic text-muted-foreground">No players</li>
        )}
      </ul>
    </div>
  );

  return (
    <Card className="w-full max-w-full">
      <CardHeader>
        <CardTitle>Map</CardTitle>
      </CardHeader>
      <CardContent className="space-y-6">
        <div className="flex flex-col gap-6 xl:flex-row xl:items-start">
          {/* Map */}
          <div className="inline-grid flex-shrink-0 rounded-xl border border-border overflow-hidden"
            style={{ gridTemplateColumns: `repeat(${cols ?? 0}, 5rem)` }}
          >
            {grid.flatMap((row, ri) =>
              row.map(([players, _items, flag, team], ci) => (
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

          {/* Players */}
          <div className="flex min-w-[16rem] flex-1 flex-col gap-3 xl:max-w-xs">
            <div className="text-sm font-semibold text-foreground">Players</div>
            <div className="flex flex-col gap-3 sm:flex-row sm:gap-4 xl:flex-col">
              {renderTeam(Teams.TEAM1, "Team 1", teamPlayers.team1)}
              {renderTeam(Teams.TEAM2, "Team 2", teamPlayers.team2)}
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
