import { AdminLayout } from "@/components/Layout/AdminLayout";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Trophy, Medal, Award } from "lucide-react";

const mockLeaderboard = [
  { rank: 1, name: "张三", userId: "U001", points: 12500, change: "+2" },
  { rank: 2, name: "李四", userId: "U002", points: 11800, change: "-1" },
  { rank: 3, name: "王五", userId: "U003", points: 10300, change: "+1" },
  { rank: 4, name: "赵六", userId: "U004", points: 9500, change: "0" },
  { rank: 5, name: "孙七", userId: "U005", points: 8900, change: "+3" },
  { rank: 6, name: "周八", userId: "U006", points: 8200, change: "-2" },
  { rank: 7, name: "吴九", userId: "U007", points: 7600, change: "0" },
  { rank: 8, name: "郑十", userId: "U008", points: 7100, change: "+1" },
];

const getRankIcon = (rank: number) => {
  switch (rank) {
    case 1:
      return <Trophy className="w-6 h-6 text-yellow-500" />;
    case 2:
      return <Medal className="w-6 h-6 text-gray-400" />;
    case 3:
      return <Award className="w-6 h-6 text-amber-700" />;
    default:
      return <span className="w-6 h-6 flex items-center justify-center font-bold text-muted-foreground">{rank}</span>;
  }
};

export default function Leaderboard() {
  return (
    <AdminLayout>
      <div className="space-y-6">
        <div>
          <h1 className="text-3xl font-bold text-foreground mb-2">排行榜</h1>
          <p className="text-muted-foreground">查看用户积分排名</p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {mockLeaderboard.slice(0, 3).map((user, index) => (
            <Card key={user.userId} className="shadow-card-hover">
              <CardHeader>
                <CardTitle className="flex items-center justify-between">
                  <span className="text-lg">第 {user.rank} 名</span>
                  {getRankIcon(user.rank)}
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="text-center space-y-3">
                  <div className="w-20 h-20 mx-auto rounded-full bg-gradient-primary flex items-center justify-center text-white text-2xl font-bold">
                    {user.name.charAt(0)}
                  </div>
                  <div>
                    <p className="font-semibold text-lg">{user.name}</p>
                    <p className="text-sm text-muted-foreground">{user.userId}</p>
                  </div>
                  <div className="pt-2 border-t">
                    <p className="text-3xl font-bold text-primary">{user.points}</p>
                    <p className="text-sm text-muted-foreground">积分</p>
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>

        <Card className="shadow-card">
          <CardContent className="p-0">
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead className="bg-muted/50 border-b">
                  <tr>
                    <th className="text-left p-4 font-medium text-sm text-muted-foreground">排名</th>
                    <th className="text-left p-4 font-medium text-sm text-muted-foreground">用户</th>
                    <th className="text-left p-4 font-medium text-sm text-muted-foreground">用户ID</th>
                    <th className="text-left p-4 font-medium text-sm text-muted-foreground">积分</th>
                    <th className="text-left p-4 font-medium text-sm text-muted-foreground">排名变化</th>
                  </tr>
                </thead>
                <tbody>
                  {mockLeaderboard.map((user) => (
                    <tr key={user.userId} className="border-b last:border-0 hover:bg-muted/30 transition-colors">
                      <td className="p-4">
                        <div className="flex items-center gap-2">
                          {getRankIcon(user.rank)}
                        </div>
                      </td>
                      <td className="p-4">
                        <div className="flex items-center gap-3">
                          <div className="w-10 h-10 rounded-full bg-gradient-secondary flex items-center justify-center text-white font-semibold">
                            {user.name.charAt(0)}
                          </div>
                          <span className="font-medium">{user.name}</span>
                        </div>
                      </td>
                      <td className="p-4">
                        <span className="font-mono text-sm">{user.userId}</span>
                      </td>
                      <td className="p-4">
                        <span className="font-semibold text-lg text-primary">{user.points}</span>
                      </td>
                      <td className="p-4">
                        <Badge
                          variant={
                            user.change.startsWith("+")
                              ? "default"
                              : user.change.startsWith("-")
                              ? "destructive"
                              : "secondary"
                          }
                        >
                          {user.change}
                        </Badge>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </CardContent>
        </Card>
      </div>
    </AdminLayout>
  );
}
