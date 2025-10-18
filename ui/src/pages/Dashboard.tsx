import { AdminLayout } from "@/components/Layout/AdminLayout";
import { StatsCard } from "@/components/Dashboard/StatsCard";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Coins, Users, Trophy, Gift, TrendingUp, Activity } from "lucide-react";

export default function Dashboard() {
  return (
    <AdminLayout>
      <div className="space-y-6">
        <div>
          <h1 className="text-3xl font-bold text-foreground mb-2">控制台</h1>
          <p className="text-muted-foreground">欢迎回来，这是您的平台数据概览</p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <StatsCard
            title="总积分发放"
            value="1,234,567"
            icon={Coins}
            trend={{ value: "12.5%", isPositive: true }}
            gradient="gradient-primary"
          />
          <StatsCard
            title="活跃用户"
            value="8,432"
            icon={Users}
            trend={{ value: "8.2%", isPositive: true }}
            gradient="gradient-secondary"
          />
          <StatsCard
            title="奖励兑换"
            value="2,345"
            icon={Gift}
            trend={{ value: "3.1%", isPositive: false }}
            gradient="gradient-success"
          />
          <StatsCard
            title="积分类型"
            value="12"
            icon={Trophy}
            gradient="gradient-primary"
          />
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <Card className="shadow-card">
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <TrendingUp className="w-5 h-5 text-primary" />
                积分趋势
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="h-64 flex items-center justify-center text-muted-foreground">
                图表区域 - 可集成 recharts 组件
              </div>
            </CardContent>
          </Card>

          <Card className="shadow-card">
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Activity className="w-5 h-5 text-primary" />
                最近活动
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {[
                  { user: "用户A", action: "兑换了奖励", time: "5分钟前" },
                  { user: "用户B", action: "获得积分 +100", time: "10分钟前" },
                  { user: "用户C", action: "登上排行榜第1名", time: "15分钟前" },
                  { user: "用户D", action: "兑换了奖励", time: "20分钟前" },
                ].map((activity, index) => (
                  <div key={index} className="flex items-center justify-between py-2 border-b last:border-0">
                    <div>
                      <p className="font-medium text-foreground">{activity.user}</p>
                      <p className="text-sm text-muted-foreground">{activity.action}</p>
                    </div>
                    <span className="text-xs text-muted-foreground">{activity.time}</span>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </AdminLayout>
  );
}
