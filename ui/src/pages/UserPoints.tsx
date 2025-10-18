import { useState } from "react";
import { AdminLayout } from "@/components/Layout/AdminLayout";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Search, Plus, Minus, Eye } from "lucide-react";
import { AdjustPointsDialog } from "@/components/UserPoints/AdjustPointsDialog";
import { useToast } from "@/hooks/use-toast";

const mockUsers = [
  { id: 1, name: "张三", userId: "U001", pointsType: "基础积分", balance: 1250, lastUpdate: "2024-01-15" },
  { id: 2, name: "李四", userId: "U002", pointsType: "基础积分", balance: 2100, lastUpdate: "2024-01-14" },
  { id: 3, name: "王五", userId: "U003", pointsType: "活动积分", balance: 850, lastUpdate: "2024-01-13" },
  { id: 4, name: "赵六", userId: "U004", pointsType: "消费积分", balance: 3200, lastUpdate: "2024-01-12" },
];

export default function UserPoints() {
  const [searchQuery, setSearchQuery] = useState("");
  const [filterType, setFilterType] = useState("all");
  const [dialogOpen, setDialogOpen] = useState(false);
  const [selectedUser, setSelectedUser] = useState<typeof mockUsers[0] | undefined>();
  const { toast } = useToast();

  const handleAdjustPoints = (user: typeof mockUsers[0]) => {
    setSelectedUser(user);
    setDialogOpen(true);
  };

  const handleViewDetails = (user: typeof mockUsers[0]) => {
    toast({
      title: "查看详情",
      description: `查看用户 ${user.name} 的积分详情`,
    });
  };

  const filteredUsers = mockUsers.filter((user) => {
    const matchesSearch = 
      user.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      user.userId.toLowerCase().includes(searchQuery.toLowerCase());
    const matchesFilter = filterType === "all" || user.pointsType === filterType;
    return matchesSearch && matchesFilter;
  });

  return (
    <AdminLayout>
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-foreground mb-2">用户积分</h1>
            <p className="text-muted-foreground">查看和管理用户积分账户</p>
          </div>
        </div>

        <Card className="shadow-card">
          <CardContent className="p-4">
            <div className="flex gap-4">
              <div className="flex-1 relative">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
                <Input 
                  placeholder="搜索用户名或用户ID" 
                  className="pl-10"
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                />
              </div>
              <Select value={filterType} onValueChange={setFilterType}>
                <SelectTrigger className="w-[180px]">
                  <SelectValue placeholder="积分类型" />
                </SelectTrigger>
                <SelectContent className="bg-popover z-50">
                  <SelectItem value="all">全部类型</SelectItem>
                  <SelectItem value="基础积分">基础积分</SelectItem>
                  <SelectItem value="活动积分">活动积分</SelectItem>
                  <SelectItem value="消费积分">消费积分</SelectItem>
                </SelectContent>
              </Select>
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
                    <th className="text-left p-4 font-medium text-sm text-muted-foreground">积分类型</th>
                    <th className="text-left p-4 font-medium text-sm text-muted-foreground">余额</th>
                    <th className="text-left p-4 font-medium text-sm text-muted-foreground">最后更新</th>
                    <th className="text-left p-4 font-medium text-sm text-muted-foreground">操作</th>
                  </tr>
                </thead>
                <tbody>
                  {filteredUsers.length === 0 ? (
                    <tr>
                      <td colSpan={6} className="p-12 text-center text-muted-foreground">
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
                        <Badge variant="outline">{user.pointsType}</Badge>
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
                          <Button variant="ghost" size="icon" title="调整积分" onClick={() => handleAdjustPoints(user)}>
                            <Plus className="w-4 h-4 text-success" />
                          </Button>
                          <Button variant="ghost" size="icon" title="调整积分" onClick={() => handleAdjustPoints(user)}>
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

        <AdjustPointsDialog
          open={dialogOpen}
          onOpenChange={setDialogOpen}
          user={selectedUser}
        />
      </div>
    </AdminLayout>
  );
}
