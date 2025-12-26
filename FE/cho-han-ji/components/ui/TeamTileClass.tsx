import { Teams } from "@/model/Tile";

export default function TeamTileClass(team: Teams) {
  switch (team) {
    case Teams.TEAM1:
      return "bg-blue-50/70 dark:bg-blue-950/25 border-blue-200/60 dark:border-blue-900/40";
    case Teams.TEAM2:
      return "bg-red-50/70 dark:bg-red-950/25 border-red-200/60 dark:border-red-900/40";
    default:
      return "bg-background border-border";
  }
}

