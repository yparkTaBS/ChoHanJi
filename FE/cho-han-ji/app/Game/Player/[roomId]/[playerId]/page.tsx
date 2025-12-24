"use client";

import { use, useEffect, useRef, useState, useMemo } from "react";

export default function Page({
  params,
}: {
  params: Promise<{ roomId: string; characterId: string }>;
}) {
  const esRef = useRef<EventSource | null>(null);

  const { roomId, characterId } = use(params);

  useEffect(() => {
    const es = new EventSource(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}api/game/player?roomId=${roomId}&playerId=${characterId}`
    );
    esRef.current = es;

    es.onopen = () => console.log("SSE connected");
    es.onmessage = (event) => console.log("SSE message:", event.data);
    es.onerror = (err) => console.log("SSE error", err);

    return () => {
      es.close();
      esRef.current = null;
    };
  }, [roomId, characterId]);
  return (
    <main className="mx-auto w-full max-w-2xl p-6">
    </main>
  )
}
