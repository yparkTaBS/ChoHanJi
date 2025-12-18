"use client";

import React, { useState } from "react";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea"
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { useRouter } from "next/navigation"

type Payload = { MapWidth: number; MapHeight: number; Items: string };

export default function CreateRoomPage() {
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const router = useRouter();

  async function onSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    setError(null);
    setSubmitting(true);

    try {
      const form = new FormData(e.currentTarget);
      const mapWidth = Number(form.get("MapWidth"));
      const mapHeight = Number(form.get("MapHeight"));
      const items = String(form.get("items") ?? "");

      const payload: Payload = { MapWidth: mapWidth, MapHeight: mapHeight, Items: items };

      const res = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}api/room`, {
        method: "POST",
        body: JSON.stringify(payload)
      });

      if (!res.ok) {
        const text = await res.text();
        throw new Error(text || `Request failed (${res.status})`);
      }

      const data: { MapId: string } = await res.json();
      router.push(`/GameAdmin/Create/Room/${data.MapId}`)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unknown error");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <main className="min-h-screen flex items-center justify-center">
      <Card className="w-full max-w-sm">
        <CardHeader className="text-center">
          <CardTitle className="text-2xl">
            Create a Room
          </CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={onSubmit} className="space-y-6">
            <div className="grid gap-2">
              <Label htmlFor="MapWidth">MapWidth</Label>
              <Input
                id="MapWidth"
                name="MapWidth"
                type="number"
                min={1}
                step={1}
                required
                placeholder="e.g. 100"
              />
            </div>

            <div className="grid gap-2">
              <Label htmlFor="mapHeight">Map height</Label>
              <Input
                id="MapHeight"
                name="MapHeight"
                type="number"
                min={1}
                step={1}
                required
                placeholder="e.g. 80"
              />
            </div>

            <div className="grid gap-2">
              <Label htmlFor="items">Items</Label>
              <Textarea
                id="items"
                name="items"
                required
                placeholder='Separated by comma. e.g. sausage,ham,bacon'
                className="min-h-[140px]"
              />
            </div>

            {error ? <p className="text-sm text-destructive">{error}</p> : null}

            <Button type="submit" className="w-full" disabled={submitting}>
              {submitting ? "Submitting..." : "Submit"}
            </Button>
          </form>
        </CardContent>
      </Card>
    </main>
  )
}
