import { Teams } from "@/model/Tile";

export default function TeamTextClass(team: Teams) {
  switch (team) {
    case Teams.TEAM1:
      return "text-blue-700 dark:text-blue-300";
    case Teams.TEAM2:
      return "text-red-700 dark:text-red-300";
    default:
      return "text-foreground";
  }
}

