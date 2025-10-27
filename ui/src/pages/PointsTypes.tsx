import { useState } from "react";
import { AdminLayout } from "@/components/Layout/AdminLayout";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Plus, Edit, Trash2, ToggleLeft, ToggleRight, Search } from "lucide-react";
import { useNavigate } from "react-router-dom";
import { PointsTypeDialog } from "@/components/PointsTypes/PointsTypeDialog";
import { useToast } from "@/hooks/use-toast";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { useQuery } from "@tanstack/react-query";
import { pointsTypeApi } from "@/services/api";

export default function PointsTypes() {
  const navigate = useNavigate();
  const [dialogOpen, setDialogOpen] = useState(false);
  const [selectedType, setSelectedType] = useState<typeof pointTypes[0] | undefined>();
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [typeToDelete, setTypeToDelete] = useState<typeof pointTypes[0] | null>(null);
  const [searchQuery, setSearchQuery] = useState("");
  const { toast } = useToast();

  const { isPending, error, data: pointTypes = [] } = useQuery({
    queryKey: ['pointsTypeStats'],
    queryFn: () =>
      pointsTypeApi.getList()
  })

  const handleEdit = (type: typeof pointTypes[0]) => {
    setSelectedType(type);
    setDialogOpen(true);
  };

  const handleAdd = () => {
    setSelectedType(undefined);
    setDialogOpen(true);
  };

  const handleToggleStatus = (type: typeof pointTypes[0]) => {
    const newStatus = type.enabled ? "启用" : "禁用";
    toast({
      title: "状态已更新",
      description: `积分类型 "${type.displayName}" 已${newStatus}`,
    });
  };

  const handleDeleteClick = (type: typeof pointTypes[0]) => {
    setTypeToDelete(type);
    setDeleteDialogOpen(true);
  };

  const handleDeleteConfirm = () => {
    if (typeToDelete) {
      toast({
        title: "删除成功",
        description: `积分类型 "${typeToDelete.displayName}" 已删除`,
      });
      setDeleteDialogOpen(false);
      setTypeToDelete(null);
    }
  };

  const filteredTypes = pointTypes.filter(
    (type) =>
      type.displayName.toLowerCase().includes(searchQuery.toLowerCase())
  );


  return (
    <AdminLayout>
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-foreground mb-2">积分类型</h1>
            <p className="text-muted-foreground">管理所有积分类型和规则</p>
          </div>
          <Button className="bg-gradient-primary hover:opacity-90" onClick={handleAdd}>
            <Plus className="w-4 h-4 mr-2" />
            新增积分类型
          </Button>
        </div>

        <Card className="shadow-card">
          <CardContent className="p-4">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
              <Input
                placeholder="搜索积分类型名称或编码"
                className="pl-10"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
              />
            </div>
          </CardContent>
        </Card>

        <div className="grid grid-cols-1 gap-4">
          {filteredTypes.length === 0 ? (
            <Card className="shadow-card">
              <CardContent className="p-12 text-center">
                <p className="text-muted-foreground">未找到匹配的积分类型</p>
              </CardContent>
            </Card>
          ) : (
            filteredTypes.map((type) => (
              <Card
                key={type.id}
                className="shadow-card hover:shadow-card-hover transition-shadow cursor-pointer"
                onClick={() => navigate(`/points-types/${type.id}`)}
              >
                <CardContent className="p-6">
                  <div className="flex items-center justify-between">
                    <div className="flex-1">
                      <div className="flex items-center gap-3 mb-2">
                        <h3 className="text-xl font-semibold text-foreground">{type.displayName}</h3>
                        <Badge variant={type.enabled ? "default" : "secondary"}>
                          {type.enabled ? "启用" : "禁用"}
                        </Badge>
                      </div>
                      <p className="text-sm text-muted-foreground mb-3">{type.description}</p>
                      <div className="flex items-center gap-6 text-sm">
                        <span className="text-muted-foreground">
                          编码: <span className="font-mono text-foreground">{type.uri}</span>
                        </span>
                        <span className="text-muted-foreground">
                          用户数: <span className="font-semibold text-foreground">{100}</span>
                        </span>
                      </div>
                    </div>
                    <div className="flex items-center gap-2" onClick={(e) => e.stopPropagation()}>
                      <Button variant="ghost" size="icon" onClick={() => handleToggleStatus(type)} title={type.enabled ? "启用" : "禁用"}>
                        {type.enabled ? (
                          <ToggleRight className="w-5 h-5 text-success" />
                        ) : (
                          <ToggleLeft className="w-5 h-5 text-muted-foreground" />
                        )}
                      </Button>
                      <Button variant="ghost" size="icon" onClick={() => handleEdit(type)} title="编辑">
                        <Edit className="w-5 h-5" />
                      </Button>
                      <Button variant="ghost" size="icon" className="text-destructive" onClick={() => handleDeleteClick(type)} title="删除">
                        <Trash2 className="w-5 h-5" />
                      </Button>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))
          )}
        </div>

        <PointsTypeDialog
          open={dialogOpen}
          onOpenChange={setDialogOpen}
          pointsType={selectedType}
        />

        <AlertDialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
          <AlertDialogContent className="bg-card">
            <AlertDialogHeader>
              <AlertDialogTitle>确认删除</AlertDialogTitle>
              <AlertDialogDescription>
                确定要删除积分类型 "<span className="font-semibold text-foreground">{typeToDelete?.displayName}</span>" 吗？
                此操作不可撤销，相关的用户积分数据也将受到影响。
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel>取消</AlertDialogCancel>
              <AlertDialogAction onClick={handleDeleteConfirm} className="bg-destructive hover:bg-destructive/90">
                删除
              </AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      </div>
    </AdminLayout>
  );
}
