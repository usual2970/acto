import { useState } from "react";
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Switch } from "@/components/ui/switch";
import { useToast } from "@/hooks/use-toast";

interface PointsTypeDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  pointsType?: {
    id: string;
    name: string;

    description: string;
    enabled: boolean;
  };
}

export function PointsTypeDialog({ open, onOpenChange, pointsType }: PointsTypeDialogProps) {
  const [name, setName] = useState(pointsType?.name || "");
  const [code, setCode] = useState(pointsType?.name || "");
  const [description, setDescription] = useState(pointsType?.description || "");
  const [isActive, setIsActive] = useState<boolean>(!pointsType || pointsType.enabled);
  const { toast } = useToast();

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!name || !code) {
      toast({
        title: "提交失败",
        description: "请填写必填字段",
        variant: "destructive",
      });
      return;
    }

    toast({
      title: pointsType ? "更新成功" : "创建成功",
      description: `积分类型 "${name}" 已${pointsType ? "更新" : "创建"}`,
    });

    onOpenChange(false);
    // 重置表单
    setName("");
    setCode("");
    setDescription("");
    setIsActive(true);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px] bg-card">
        <DialogHeader>
          <DialogTitle>{pointsType ? "编辑积分类型" : "新增积分类型"}</DialogTitle>
          <DialogDescription>
            {pointsType ? "修改积分类型信息" : "创建一个新的积分类型"}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="name">
              积分名称 <span className="text-destructive">*</span>
            </Label>
            <Input
              id="name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="例如：基础积分"
              required
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="code">
              积分编码 <span className="text-destructive">*</span>
            </Label>
            <Input
              id="code"
              value={code}
              onChange={(e) => setCode(e.target.value)}
              placeholder="例如：BASE_POINTS"
              className="font-mono"
              required
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="description">描述</Label>
            <Textarea
              id="description"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="描述积分类型的用途和规则"
              rows={3}
            />
          </div>

          <div className="flex items-center justify-between">
            <Label htmlFor="status">启用状态</Label>
            <Switch
              id="status"
              checked={isActive}
              onCheckedChange={setIsActive}
            />
          </div>

          <DialogFooter>
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              取消
            </Button>
            <Button type="submit" className="bg-gradient-primary hover:opacity-90">
              {pointsType ? "保存修改" : "创建"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
