"use client";

import { use, useEffect, useRef, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { useRouter } from "next/navigation"
import { Message, TypedMessage } from "@/model/SSEMessage";

type Player = {
  id: string;
  name: string;
}

export default function AdminWaitingRoom({ params }: { params: Promise<{ id: string }> }) {
  const esRef = useRef<EventSource | null>(null);

  const { id: roomId } = use(params);
  const router = useRouter()
  const [players, setPlayers] = useState<Player[]>([])

  useEffect(() => {
    const es = new EventSource(`${process.env.NEXT_PUBLIC_API_BASE_URL}api/room/waiting/admin?roomId=${roomId}`);
    esRef.current = es;

    es.onopen = () => {
      console.log("SSE connected");
    }

    es.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data) as { MessageType: string; Message?: string }

        const baseMessage = new Message(data.MessageType)

        if (baseMessage.MessageType === "PlayerConnected") {
          if (!data.Message) return

          // 2nd parse: Message is a JSON string
          const playerPartial = JSON.parse(data.Message) as { id: string; name: string }

          const player: Player = {
            id: playerPartial.id,
            name: playerPartial.name,
          }

          const typed = new TypedMessage<Player>(data.MessageType, player)

          setPlayers((prev) => [...prev, typed.Message])
        }
      } catch (e) {
        console.log(e)
        return
      }

    }

    es.onerror = (err) => {
      console.error("SSE error", err);
      router.back()
    }

    return () => {
      es.close();
      esRef.current = null;
    }
  }, [roomId]);

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
                    <span className="font-mono text-xs text-muted-foreground">
                      {p.id}
                    </span>
                  </li>
                ))}
              </ul>
            )}
          </section>
        </CardContent>

      </Card>
    </main>
  )
}
