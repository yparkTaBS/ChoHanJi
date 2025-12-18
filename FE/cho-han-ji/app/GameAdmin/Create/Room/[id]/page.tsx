import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";

type Player = {
  id: string;
  name: string;
}

type PageProps = {
  params: Promise<{ id: string }>;
}

export default async function AdminWaitingRoom({ params }: { params: { id: string } }) {
  const players: Player[] = [];
  const roomId = (await params).id;
  const playerCount = players.length;

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
              <span className="text-sm font-medium">{playerCount}</span>
            </div>
          </div>

          <Separator />

          <section className="space-y-2">
            <h2 className="text-sm font-semibold">Player list</h2>

            {playerCount === 0 ? (
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
