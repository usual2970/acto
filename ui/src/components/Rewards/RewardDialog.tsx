import { useState } from "react";
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { useToast } from "@/hooks/use-toast";

interface RewardDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  reward?: {
    id: number;
    name: string;
    cost: number;
    stock: number;
    type: string;
  };
}

export function RewardDialog({ open, onOpenChange, reward }: RewardDialogProps) {
  const [name, setName] = useState(reward?.name || "");
  const [cost, setCost] = useState(reward?.cost?.toString() || "");
  const [stock, setStock] = useState(reward?.stock?.toString() || "");
  const [type, setType] = useState(reward?.type || "积分兑换");
  const [description, setDescription] = useState("");
  const { toast } = useToast();

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!name || !cost || !stock) {
      toast({
        title: "提交失败",
        description: "请填写所有必填字段",
        variant: "destructive",
      });
      return;
    }

    toast({
      title: reward ? "更新成功" : "创建成功",
      description: `奖励 "${name}" 已${reward ? "更新" : "创建"}`,
    });

    onOpenChange(false);
    // 重置表单
    setName("");
    setCost("");
    setStock("");
    setType("积分兑换");
    setDescription("");
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px] bg-card">
        <DialogHeader>
          <DialogTitle>{reward ? "编辑奖励" : "新增奖励"}</DialogTitle>
          <DialogDescription>
            {reward ? "修改奖励信息" : "创建一个新的奖励"}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="name">
              奖励名称 <span className="text-destructive">*</span>
            </Label>
            <Input
              id="name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="例如：iPhone 15 Pro"
              required
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="type">
              奖励类型 <span className="text-destructive">*</span>
            </Label>
            <Select value={type} onValueChange={setType}>
              <SelectTrigger id="type">
                <SelectValue placeholder="选择奖励类型" />
              </SelectTrigger>
              <SelectContent className="bg-popover z-50">
                <SelectItem value="积分兑换">积分兑换</SelectItem>
                <SelectItem value="排行榜奖励">排行榜奖励</SelectItem>
                <SelectItem value="活动奖励">活动奖励</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="cost">
                消耗积分 <span className="text-destructive">*</span>
              </Label>
              <Input
                id="cost"
                type="number"
                value={cost}
                onChange={(e) => setCost(e.target.value)}
                placeholder="0"
                min="0"
                required
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="stock">
                库存数量 <span className="text-destructive">*</span>
              </Label>
              <Input
                id="stock"
                type="number"
                value={stock}
                onChange={(e) => setStock(e.target.value)}
                placeholder="0"
                min="0"
                required
              />
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="description">奖励描述</Label>
            <Textarea
              id="description"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="描述奖励的详细信息"
              rows={3}
            />
          </div>

          <DialogFooter>
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              取消
            </Button>
            <Button type="submit" className="bg-gradient-primary hover:opacity-90">
              {reward ? "保存修改" : "创建"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
