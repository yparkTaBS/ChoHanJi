"use client"

import React, { useState } from "react"
import { useRouter } from "next/navigation"

import { Button } from "@/components/ui/button"
import { Label } from "@/components/ui/label"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"

type CharacterClass = "fighter" | "ranger" | "thief"
type TeamNumber = 1 | 2

type Payload = {
  RoomId: string
  UserName: string
  Class: CharacterClass
  TeamNumber: TeamNumber
}

function isCharacterClass(v: string): v is CharacterClass {
  return v === "fighter" || v === "ranger" || v === "thief"
}

function isTeamNumberString(v: string): v is "1" | "2" {
  return v === "1" || v === "2"
}

export default function CreateCharacterPage() {
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [charClass, setCharClass] = useState<CharacterClass | "">("")
  const [teamNumber, setTeamNumber] = useState<TeamNumber>(1)
  const router = useRouter()

  async function onSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault()
    setError(null)
    setSubmitting(true)

    try {
      const form = new FormData(e.currentTarget)

      const roomId = String(form.get("RoomId") ?? "").trim()
      const userName = String(form.get("UserName") ?? "").trim()
      const playerClass = charClass
      const playerTeamNumber = teamNumber

      if (!roomId) throw new Error("Room ID is required.")
      if (!userName) throw new Error("User name is required.")
      if (!playerClass) throw new Error("Class is required.")
      if (!playerTeamNumber) throw new Error("TeamNumber is required")

      const payload: Payload = { RoomId: roomId, UserName: userName, Class: playerClass, TeamNumber: teamNumber }

      const res = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}api/character`,
        {
          method: "POST",
          // headers: { "Content-Type": "application/json" },
          body: JSON.stringify(payload),
        }
      )

      if (!res.ok) {
        const text = await res.text()
        throw new Error(text || `Request failed (${res.status})`)
      }

      const data: { CharacterId: string } = await res.json()
      console.log(data.CharacterId)
      router.push(`/Player/${roomId}/${data.CharacterId}`)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unknown error")
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <main className="min-h-screen flex items-center justify-center">
      <Card className="w-full max-w-sm">
        <CardHeader className="text-center">
          <CardTitle className="text-2xl">Create a Character</CardTitle>
        </CardHeader>

        <CardContent>
          <form onSubmit={onSubmit} className="space-y-6">
            <div className="grid gap-2">
              <Label htmlFor="RoomId">Room ID</Label>
              <Input id="RoomId" name="RoomId" required placeholder="e.g. 6f2c1f2a-..." />
            </div>

            <div className="grid gap-2">
              <Label htmlFor="UserName">User name</Label>
              <Input id="UserName" name="UserName" required placeholder="e.g. yeounjun" />
            </div>

            <div className="grid gap-2">
              <Label htmlFor="Class">Class</Label>

              <Select
                value={charClass}
                onValueChange={(value) => {
                  // Deduction: shadcn Select always gives string, so we must narrow it
                  if (isCharacterClass(value)) setCharClass(value)
                  else setCharClass("") // safety fallback
                }}
              >
                <SelectTrigger id="Class">
                  <SelectValue placeholder="Select a class" />
                </SelectTrigger>

                <SelectContent>
                  <SelectItem value="fighter">Fighter</SelectItem>
                  <SelectItem value="ranger">Ranger</SelectItem>
                  <SelectItem value="thief">Thief</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="grid gap-2">
              <Label htmlFor="TeamNumber">Team Number</Label>
              <Select
                // Deduction/work: Select expects string; state is number union => stringify it
                value={String(teamNumber)}
                onValueChange={(value) => {
                  // Deduction/work: convert "1"/"2" -> 1|2 safely
                  if (isTeamNumberString(value)) {
                    setTeamNumber(value === "1" ? 1 : 2)
                  }
                }}
              >
                <SelectTrigger id="TeamNumber">
                  <SelectValue placeholder="Select TeamNumber" />
                </SelectTrigger>

                <SelectContent>
                  <SelectItem value="1">1</SelectItem>
                  <SelectItem value="2">2</SelectItem>
                </SelectContent>
              </Select>

              {/* Optional: hidden input only needed if you were *not* using JSON payload;
                  harmless to keep, but it’s redundant since you send `teamNumber` in payload. */}
              <input type="hidden" name="TeamNumber" value={String(teamNumber)} />
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

