import Link from "next/link"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"

export default function Home() {
  return (
    <main className="min-h-screen flex items-center justify-center">
      <Card className="w-full max-w-sm">
        <CardHeader className="w-full text-center">
          <CardTitle className="whitespace-normal break-words leading-snug">
            Budaejjigae Craft
          </CardTitle>
        </CardHeader>
        <CardContent className="flex flex-col gap-4">
          <Button size="lg" className="w-full" asChild>
            <Link href="GameAdmin/Create">
              Create a Room
            </Link>
          </Button>
          <Button size="lg">
            Join a Room
          </Button>
        </CardContent>
      </Card>
    </main>
  );
}
