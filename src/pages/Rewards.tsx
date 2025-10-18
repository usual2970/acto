import { useState } from "react";
import { AdminLayout } from "@/components/Layout/AdminLayout";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Plus, Edit, Eye, Package } from "lucide-react";
import { RewardDialog } from "@/components/Rewards/RewardDialog";
import { useToast } from "@/hooks/use-toast";

const mockRewards = [
  { id: 1, name: "iPhone 15 Pro", cost: 50000, stock: 5, claimed: 12, type: "积分兑换", status: "active" },
  { id: 2, name: "小米手环", cost: 5000, stock: 50, claimed: 234, type: "积分兑换", status: "active" },
  { id: 3, name: "现金红包 ¥10", cost: 1000, stock: 1000, claimed: 567, type: "积分兑换", status: "active" },
  { id: 4, name: "冠军奖杯", cost: 0, stock: 1, claimed: 0, type: "排行榜奖励", status: "active" },
];

export default function Rewards() {
  const [dialogOpen, setDialogOpen] = useState(false);
  const [selectedReward, setSelectedReward] = useState<typeof mockRewards[0] | undefined>();
  const [filterType, setFilterType] = useState("all");
  const { toast } = useToast();

  const handleEdit = (reward: typeof mockRewards[0]) => {
    setSelectedReward(reward);
    setDialogOpen(true);
  };

  const handleAdd = () => {
    setSelectedReward(undefined);
    setDialogOpen(true);
  };

  const handleView = (reward: typeof mockRewards[0]) => {
    toast({
      title: "查看详情",
      description: `查看奖励 "${reward.name}" 的详细信息和兑换记录`,
    });
  };

  const filteredRewards = mockRewards.filter(
    (reward) => filterType === "all" || reward.type === filterType
  );

  return (
    <AdminLayout>
      <div className="space-y-6">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <h1 className="text-3xl font-bold text-foreground mb-2">奖励管理</h1>
            <p className="text-muted-foreground">配置和管理所有奖励</p>
          </div>
          <div className="flex items-center gap-3">
            <Select value={filterType} onValueChange={setFilterType}>
              <SelectTrigger className="w-[150px]">
                <SelectValue placeholder="奖励类型" />
              </SelectTrigger>
              <SelectContent className="bg-popover z-50">
                <SelectItem value="all">全部类型</SelectItem>
                <SelectItem value="积分兑换">积分兑换</SelectItem>
                <SelectItem value="排行榜奖励">排行榜奖励</SelectItem>
                <SelectItem value="活动奖励">活动奖励</SelectItem>
              </SelectContent>
            </Select>
            <Button className="bg-gradient-primary hover:opacity-90" onClick={handleAdd}>
              <Plus className="w-4 h-4 mr-2" />
              新增奖励
            </Button>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
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
                  <div className="w-12 h-12 rounded-xl bg-gradient-success flex items-center justify-center">
                    <Package className="w-6 h-6 text-white" />
                  </div>
                  <Badge variant={reward.type === "积分兑换" ? "default" : "secondary"}>
                    {reward.type}
                  </Badge>
                </div>

                <h3 className="text-lg font-semibold text-foreground mb-2">{reward.name}</h3>

                <div className="space-y-2 mb-4 text-sm">
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">消耗积分</span>
                    <span className="font-semibold text-primary">{reward.cost}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">库存</span>
                    <span className="font-semibold">{reward.stock}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">已兑换</span>
                    <span className="font-semibold">{reward.claimed}</span>
                  </div>
                </div>

                <div className="flex items-center gap-2">
                  <Button variant="outline" size="sm" className="flex-1" onClick={() => handleView(reward)}>
                    <Eye className="w-4 h-4 mr-1" />
                    查看
                  </Button>
                  <Button variant="outline" size="sm" className="flex-1" onClick={() => handleEdit(reward)}>
                    <Edit className="w-4 h-4 mr-1" />
                    编辑
                  </Button>
                </div>
              </CardContent>
            </Card>
            ))
          )}
        </div>

        <RewardDialog
          open={dialogOpen}
          onOpenChange={setDialogOpen}
          reward={selectedReward}
        />
      </div>
    </AdminLayout>
  );
}
