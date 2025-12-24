"use client";

import { use, useEffect, useRef, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { Button } from "@/components/ui/button";
import { useRouter } from "next/navigation";
import { Message, TypedMessage } from "@/model/SSEMessage";

type Player = {
  id: string;
  name: string;
};

export default function AdminWaitingRoom({ params }: { params: Promise<{ id: string }> }) {
  const esRef = useRef<EventSource | null>(null);

  const { id: roomId } = use(params);
  const router = useRouter();
  const [players, setPlayers] = useState<Player[]>([]);
  const [isStarting, setIsStarting] = useState(false);
  const [isCancelling, setIsCancelling] = useState(false);

  useEffect(() => {
    const es = new EventSource(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}api/room/waiting/admin?roomId=${roomId}`
    );
    esRef.current = es;

    es.onopen = () => {
      console.log("SSE connected");
    };

    es.onmessage = (event) => {
      console.log("SSE raw:", event.data);
      try {
        const data = JSON.parse(event.data) as { MessageType: string; Message?: string };

        const baseMessage = new Message(data.MessageType);

        if (baseMessage.MessageType === "PlayerConnected") {
          if (!data.Message) return;

          // 2nd parse: Message is a JSON string
          const playerPartial = JSON.parse(data.Message) as { id: string; name: string };

          const player: Player = {
            id: playerPartial.id,
            name: playerPartial.name,
          };

          const typed = new TypedMessage<Player>(data.MessageType, player);

          // ✅ De-duplicate players by id
          setPlayers((prev) => {
            const idx = prev.findIndex((p) => p.id === typed.Message.id);
            if (idx !== -1) {
              const next = [...prev];
              next[idx] = typed.Message;
              return next;
            }
            return [...prev, typed.Message];
          });
        }
      } catch (e) {
        console.log(e);
        return;
      }
    };

    es.onerror = (err) => {
      console.log("SSE error", err);
      router.back();
    };

    return () => {
      es.close();
      esRef.current = null;
    };
  }, [roomId, router]);

  const handleCancelGame = async () => {
    try {
      setIsCancelling(true);

      // TODO: Implement later
      /*
      const res = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}api/room/cancel?roomId=${roomId}`,
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
        }
      );

      if (!res.ok) {
        throw new Error(`Cancel failed: ${res.status}`);
      }
      */

      // Optional: go back to previous page after cancel
      router.back();
    } catch (err) {
      console.log(err);
    } finally {
      setIsCancelling(false);
    }
  };

  const handleStartGame = async () => {
    try {
      setIsStarting(true);

      const res = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}api/game/start?roomId=${roomId}`,
        {
          method: "POST",
        }
      );

      if (!res.ok) {
        throw new Error(`Start failed: ${res.status}`);
      }

      // Optional: navigate to game page after start
      router.push(`/Game/Map/${roomId}`);
    } catch (err) {
      console.error(err);
    } finally {
      setIsStarting(false);
    }
  };

  return (
    <main className="mx-auto w-full max-w-2xl p-6">
      <Card className="rounded-2xl">
        <CardHeader>
          <CardTitle className="text-2xl">Room</CardTitle>
          <CardDescription>Details for a specific room</CardDescription>
        </CardHeader>

        <CardContent className="space-y-4">
          <div className="grid gap-2">
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">ID</span>
              <span className="font-mono text-sm">{roomId}</span>
            </div>

            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">Players</span>
              <span className="text-sm font-medium">{players.length}</span>
            </div>
          </div>

          {/* ✅ Buttons */}
          <div className="flex gap-2">
            <Button
              variant="destructive"
              className="flex-1"
              onClick={handleCancelGame}
              disabled={isCancelling || isStarting}
            >
              {isCancelling ? "Cancelling..." : "Cancel the game"}
            </Button>

            <Button
              className="flex-1"
              onClick={handleStartGame}
              disabled={isStarting || isCancelling || players.length === 0}
            >
              {isStarting ? "Starting..." : "Start the game"}
            </Button>
          </div>

          <Separator />

          <section className="space-y-2">
            <h2 className="text-sm font-semibold">Player list</h2>

            {players.length === 0 ? (
              <p className="text-sm text-muted-foreground">
                No players yet. (List will appear once provided.)
              </p>
            ) : (
              <ul className="space-y-2">
                {players.map((p) => (
                  <li
                    key={p.id}
                    className="flex items-center justify-between rounded-xl border p-3"
                  >
                    <span className="font-medium">{p.name}</span>
                    <span className="font-mono text-xs text-muted-foreground">{p.id}</span>
                  </li>
                ))}
              </ul>
            )}
          </section>
        </CardContent>
      </Card>
    </main>
  );
}

