import { useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { AdminLayout } from "@/components/Layout/AdminLayout";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { ArrowLeft, Search, Plus, Minus, Eye, Edit, Trash2 } from "lucide-react";
import { AdjustPointsDialog } from "@/components/UserPoints/AdjustPointsDialog";
import { RewardDialog } from "@/components/Rewards/RewardDialog";
import { useToast } from "@/hooks/use-toast";

// Mock data - 实际项目中从 API 获取
const mockPointsType = {
  id: 1,
  name: "基础积分",
  code: "BASE_POINTS",
  description: "用户基础行为积分",
  status: "active",
  totalUsers: 1234,
};

const mockUsers = [
  { id: 1, name: "张三", userId: "U001", balance: 1250, lastUpdate: "2024-01-15" },
  { id: 2, name: "李四", userId: "U002", balance: 2100, lastUpdate: "2024-01-14" },
  { id: 3, name: "王五", userId: "U003", balance: 850, lastUpdate: "2024-01-13" },
  { id: 4, name: "赵六", userId: "U004", balance: 3200, lastUpdate: "2024-01-12" },
];

const mockLeaderboard = [
  { id: 1, rank: 1, name: "李四", userId: "U002", points: 2100, change: "+150" },
  { id: 2, rank: 2, name: "赵六", userId: "U004", points: 3200, change: "+80" },
  { id: 3, rank: 3, name: "张三", userId: "U001", points: 1250, change: "-20" },
  { id: 4, rank: 4, name: "王五", userId: "U003", points: 850, change: "+45" },
];

const mockRewards = [
  { id: 1, name: "10元优惠券", type: "兑换", cost: 100, stock: 500, status: "active" },
  { id: 2, name: "50元优惠券", type: "兑换", cost: 500, stock: 100, status: "active" },
  { id: 3, name: "排行榜第一名奖励", type: "排行榜", cost: 0, stock: 1, status: "active" },
];

const mockRewardRecords = [
  { id: 1, userName: "张三", userId: "U001", rewardName: "10元优惠券", type: "兑换", points: 100, status: "已领取", time: "2024-01-15 14:30" },
  { id: 2, userName: "李四", userId: "U002", rewardName: "50元优惠券", type: "兑换", points: 500, status: "已领取", time: "2024-01-15 10:20" },
  { id: 3, userName: "王五", userId: "U003", rewardName: "排行榜第一名奖励", type: "发放", points: 0, status: "已发放", time: "2024-01-14 18:00" },
  { id: 4, userName: "赵六", userId: "U004", rewardName: "10元优惠券", type: "兑换", points: 100, status: "已领取", time: "2024-01-14 16:45" },
  { id: 5, userName: "张三", userId: "U001", rewardName: "10元优惠券", type: "兑换", points: 100, status: "已领取", time: "2024-01-13 09:15" },
];

export default function PointsTypeDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const { toast } = useToast();
  
  const [searchQuery, setSearchQuery] = useState("");
  const [filterPeriod, setFilterPeriod] = useState("all");
  const [filterRewardType, setFilterRewardType] = useState("all");
  const [adjustDialogOpen, setAdjustDialogOpen] = useState(false);
  const [rewardDialogOpen, setRewardDialogOpen] = useState(false);
  const [selectedUser, setSelectedUser] = useState<typeof mockUsers[0] | undefined>();
  const [selectedReward, setSelectedReward] = useState<typeof mockRewards[0] | undefined>();

  const handleAdjustPoints = (user: typeof mockUsers[0]) => {
    setSelectedUser(user);
    setAdjustDialogOpen(true);
  };

  const handleViewDetails = (user: typeof mockUsers[0]) => {
    toast({
      title: "查看详情",
      description: `查看用户 ${user.name} 的积分详情`,
    });
  };

  const handleAddReward = () => {
    setSelectedReward(undefined);
    setRewardDialogOpen(true);
  };

  const handleEditReward = (reward: typeof mockRewards[0]) => {
    setSelectedReward(reward);
    setRewardDialogOpen(true);
  };

  const handleDeleteReward = (reward: typeof mockRewards[0]) => {
    toast({
      title: "删除成功",
      description: `奖励 "${reward.name}" 已删除`,
    });
  };

  const filteredUsers = mockUsers.filter((user) =>
    user.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    user.userId.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const filteredRewards = mockRewards.filter((reward) =>
    (filterRewardType === "all" || reward.type === filterRewardType) &&
    reward.name.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <AdminLayout>
      <div className="space-y-6">
        <div className="flex items-center gap-4">
          <Button variant="ghost" size="icon" onClick={() => navigate("/points-types")}>
            <ArrowLeft className="w-5 h-5" />
          </Button>
          <div className="flex-1">
            <div className="flex items-center gap-3 mb-2">
              <h1 className="text-3xl font-bold text-foreground">{mockPointsType.name}</h1>
              <Badge variant={mockPointsType.status === "active" ? "default" : "secondary"}>
                {mockPointsType.status === "active" ? "启用" : "禁用"}
              </Badge>
            </div>
            <p className="text-muted-foreground">{mockPointsType.description}</p>
          </div>
        </div>

        <Card className="shadow-card">
          <CardContent className="p-6">
            <div className="grid grid-cols-3 gap-6">
              <div>
                <p className="text-sm text-muted-foreground mb-1">积分编码</p>
                <p className="text-lg font-mono font-semibold">{mockPointsType.code}</p>
              </div>
              <div>
                <p className="text-sm text-muted-foreground mb-1">用户总数</p>
                <p className="text-lg font-semibold text-primary">{mockPointsType.totalUsers}</p>
              </div>
              <div>
                <p className="text-sm text-muted-foreground mb-1">状态</p>
                <p className="text-lg font-semibold">{mockPointsType.status === "active" ? "启用中" : "已禁用"}</p>
              </div>
            </div>
          </CardContent>
        </Card>

        <Tabs defaultValue="users" className="space-y-6">
          <TabsList className="bg-muted/50">
            <TabsTrigger value="users">用户积分</TabsTrigger>
            <TabsTrigger value="leaderboard">排行榜</TabsTrigger>
            <TabsTrigger value="rewards">奖励管理</TabsTrigger>
          </TabsList>

          {/* 用户积分 Tab */}
          <TabsContent value="users" className="space-y-4">
            <Card className="shadow-card">
              <CardContent className="p-4">
                <div className="relative">
                  <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
                  <Input
                    placeholder="搜索用户名或用户ID"
                    className="pl-10"
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                  />
                </div>
              </CardContent>
            </Card>

            <Card className="shadow-card">
              <CardContent className="p-0">
                <div className="overflow-x-auto">
                  <table className="w-full">
                    <thead className="bg-muted/50 border-b">
                      <tr>
                        <th className="text-left p-4 font-medium text-sm text-muted-foreground">用户</th>
                        <th className="text-left p-4 font-medium text-sm text-muted-foreground">用户ID</th>
                        <th className="text-left p-4 font-medium text-sm text-muted-foreground">余额</th>
                        <th className="text-left p-4 font-medium text-sm text-muted-foreground">最后更新</th>
                        <th className="text-left p-4 font-medium text-sm text-muted-foreground">操作</th>
                      </tr>
                    </thead>
                    <tbody>
                      {filteredUsers.length === 0 ? (
                        <tr>
                          <td colSpan={5} className="p-12 text-center text-muted-foreground">
                            未找到匹配的用户
                          </td>
                        </tr>
                      ) : (
                        filteredUsers.map((user) => (
                          <tr key={user.id} className="border-b last:border-0 hover:bg-muted/30 transition-colors">
                            <td className="p-4">
                              <div className="flex items-center gap-3">
                                <div className="w-10 h-10 rounded-full bg-gradient-primary flex items-center justify-center text-white font-semibold">
                                  {user.name.charAt(0)}
                                </div>
                                <span className="font-medium">{user.name}</span>
                              </div>
                            </td>
                            <td className="p-4">
                              <span className="font-mono text-sm">{user.userId}</span>
                            </td>
                            <td className="p-4">
                              <span className="font-semibold text-lg text-primary">{user.balance}</span>
                            </td>
                            <td className="p-4 text-sm text-muted-foreground">{user.lastUpdate}</td>
                            <td className="p-4">
                              <div className="flex items-center gap-2">
                                <Button variant="ghost" size="icon" title="查看详情" onClick={() => handleViewDetails(user)}>
                                  <Eye className="w-4 h-4" />
                                </Button>
                                <Button variant="ghost" size="icon" title="增加积分" onClick={() => handleAdjustPoints(user)}>
                                  <Plus className="w-4 h-4 text-success" />
                                </Button>
                                <Button variant="ghost" size="icon" title="扣减积分" onClick={() => handleAdjustPoints(user)}>
                                  <Minus className="w-4 h-4 text-destructive" />
                                </Button>
                              </div>
                            </td>
                          </tr>
                        ))
                      )}
                    </tbody>
                  </table>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          {/* 排行榜 Tab */}
          <TabsContent value="leaderboard" className="space-y-4">
            <Card className="shadow-card">
              <CardContent className="p-4">
                <Select value={filterPeriod} onValueChange={setFilterPeriod}>
                  <SelectTrigger className="w-[180px]">
                    <SelectValue placeholder="选择周期" />
                  </SelectTrigger>
                  <SelectContent className="bg-popover z-50">
                    <SelectItem value="all">全部</SelectItem>
                    <SelectItem value="daily">每日</SelectItem>
                    <SelectItem value="weekly">每周</SelectItem>
                    <SelectItem value="monthly">每月</SelectItem>
                  </SelectContent>
                </Select>
              </CardContent>
            </Card>

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
                        <th className="text-left p-4 font-medium text-sm text-muted-foreground">变化</th>
                      </tr>
                    </thead>
                    <tbody>
                      {mockLeaderboard.map((entry) => (
                        <tr key={entry.id} className="border-b last:border-0 hover:bg-muted/30 transition-colors">
                          <td className="p-4">
                            <div className={`w-8 h-8 rounded-full flex items-center justify-center font-bold ${
                              entry.rank === 1 ? "bg-yellow-500 text-white" :
                              entry.rank === 2 ? "bg-gray-400 text-white" :
                              entry.rank === 3 ? "bg-orange-600 text-white" :
                              "bg-muted text-foreground"
                            }`}>
                              {entry.rank}
                            </div>
                          </td>
                          <td className="p-4">
                            <div className="flex items-center gap-3">
                              <div className="w-10 h-10 rounded-full bg-gradient-primary flex items-center justify-center text-white font-semibold">
                                {entry.name.charAt(0)}
                              </div>
                              <span className="font-medium">{entry.name}</span>
                            </div>
                          </td>
                          <td className="p-4">
                            <span className="font-mono text-sm">{entry.userId}</span>
                          </td>
                          <td className="p-4">
                            <span className="font-semibold text-lg text-primary">{entry.points}</span>
                          </td>
                          <td className="p-4">
                            <Badge variant={entry.change.startsWith("+") ? "default" : "secondary"}>
                              {entry.change}
                            </Badge>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          {/* 奖励管理 Tab */}
          <TabsContent value="rewards" className="space-y-4">
            <div className="flex items-center justify-between">
              <Card className="shadow-card flex-1 mr-4">
                <CardContent className="p-4">
                  <div className="flex gap-4">
                    <div className="flex-1 relative">
                      <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
                      <Input
                        placeholder="搜索奖励名称"
                        className="pl-10"
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                      />
                    </div>
                    <Select value={filterRewardType} onValueChange={setFilterRewardType}>
                      <SelectTrigger className="w-[180px]">
                        <SelectValue placeholder="奖励类型" />
                      </SelectTrigger>
                      <SelectContent className="bg-popover z-50">
                        <SelectItem value="all">全部类型</SelectItem>
                        <SelectItem value="兑换">兑换奖励</SelectItem>
                        <SelectItem value="排行榜">排行榜奖励</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                </CardContent>
              </Card>
              <Button className="bg-gradient-primary hover:opacity-90" onClick={handleAddReward}>
                <Plus className="w-4 h-4 mr-2" />
                新增奖励
              </Button>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {filteredRewards.length === 0 ? (
                <Card className="shadow-card col-span-full">
                  <CardContent className="p-12 text-center">
                    <p className="text-muted-foreground">未找到匹配的奖励</p>
                  </CardContent>
                </Card>
              ) : (
                filteredRewards.map((reward) => (
                  <Card key={reward.id} className="shadow-card hover:shadow-card-hover transition-shadow">
                    <CardContent className="p-6">
                      <div className="flex items-start justify-between mb-4">
                        <Badge variant={reward.type === "兑换" ? "default" : "secondary"}>
                          {reward.type}
                        </Badge>
                        <Badge variant={reward.status === "active" ? "default" : "secondary"}>
                          {reward.status === "active" ? "上架" : "下架"}
                        </Badge>
                      </div>
                      <h3 className="text-lg font-semibold mb-2">{reward.name}</h3>
                      <div className="space-y-2 mb-4">
                        <div className="flex justify-between text-sm">
                          <span className="text-muted-foreground">所需积分</span>
                          <span className="font-semibold text-primary">{reward.cost || "无需积分"}</span>
                        </div>
                        <div className="flex justify-between text-sm">
                          <span className="text-muted-foreground">剩余库存</span>
                          <span className="font-semibold">{reward.stock}</span>
                        </div>
                      </div>
                      <div className="flex gap-2">
                        <Button variant="outline" size="sm" className="flex-1" onClick={() => handleEditReward(reward)}>
                          <Edit className="w-4 h-4 mr-1" />
                          编辑
                        </Button>
                        <Button variant="outline" size="sm" className="text-destructive" onClick={() => handleDeleteReward(reward)}>
                          <Trash2 className="w-4 h-4" />
                        </Button>
                      </div>
                    </CardContent>
                  </Card>
                ))
              )}
            </div>

            {/* 发放/领取记录 */}
            <Card className="shadow-card">
              <CardContent className="p-6">
                <h3 className="text-lg font-semibold mb-4">发放/领取记录</h3>
                <div className="overflow-x-auto">
                  <table className="w-full">
                    <thead className="bg-muted/50 border-b">
                      <tr>
                        <th className="text-left p-4 font-medium text-sm text-muted-foreground">用户</th>
                        <th className="text-left p-4 font-medium text-sm text-muted-foreground">用户ID</th>
                        <th className="text-left p-4 font-medium text-sm text-muted-foreground">奖励名称</th>
                        <th className="text-left p-4 font-medium text-sm text-muted-foreground">类型</th>
                        <th className="text-left p-4 font-medium text-sm text-muted-foreground">消耗积分</th>
                        <th className="text-left p-4 font-medium text-sm text-muted-foreground">状态</th>
                        <th className="text-left p-4 font-medium text-sm text-muted-foreground">时间</th>
                      </tr>
                    </thead>
                    <tbody>
                      {mockRewardRecords.length === 0 ? (
                        <tr>
                          <td colSpan={7} className="p-12 text-center text-muted-foreground">
                            暂无发放/领取记录
                          </td>
                        </tr>
                      ) : (
                        mockRewardRecords.map((record) => (
                          <tr key={record.id} className="border-b last:border-0 hover:bg-muted/30 transition-colors">
                            <td className="p-4">
                              <div className="flex items-center gap-3">
                                <div className="w-10 h-10 rounded-full bg-gradient-primary flex items-center justify-center text-white font-semibold">
                                  {record.userName.charAt(0)}
                                </div>
                                <span className="font-medium">{record.userName}</span>
                              </div>
                            </td>
                            <td className="p-4">
                              <span className="font-mono text-sm">{record.userId}</span>
                            </td>
                            <td className="p-4">
                              <span className="font-medium">{record.rewardName}</span>
                            </td>
                            <td className="p-4">
                              <Badge variant={record.type === "兑换" ? "default" : "secondary"}>
                                {record.type}
                              </Badge>
                            </td>
                            <td className="p-4">
                              <span className="font-semibold text-primary">
                                {record.points > 0 ? record.points : "无"}
                              </span>
                            </td>
                            <td className="p-4">
                              <Badge variant="outline">{record.status}</Badge>
                            </td>
                            <td className="p-4 text-sm text-muted-foreground">{record.time}</td>
                          </tr>
                        ))
                      )}
                    </tbody>
                  </table>
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>

        <AdjustPointsDialog
          open={adjustDialogOpen}
          onOpenChange={setAdjustDialogOpen}
          user={selectedUser}
        />

        <RewardDialog
          open={rewardDialogOpen}
          onOpenChange={setRewardDialogOpen}
          reward={selectedReward}
        />
      </div>
    </AdminLayout>
  );
}
