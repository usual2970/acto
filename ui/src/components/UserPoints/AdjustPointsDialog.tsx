import { useState } from "react";
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { useToast } from "@/hooks/use-toast";
import { Plus, Minus } from "lucide-react";

interface AdjustPointsDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  user?: {
    id: number;
    name: string;
    userId: string;
    balance: number;
  };
}

export function AdjustPointsDialog({ open, onOpenChange, user }: AdjustPointsDialogProps) {
  const [type, setType] = useState<"add" | "subtract">("add");
  const [amount, setAmount] = useState("");
  const [reason, setReason] = useState("");
  const { toast } = useToast();

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    const pointsAmount = parseInt(amount);
    if (!amount || pointsAmount <= 0) {
      toast({
        title: "提交失败",
        description: "请输入有效的积分数量",
        variant: "destructive",
      });
      return;
    }

    const newBalance = type === "add" 
      ? (user?.balance || 0) + pointsAmount 
      : (user?.balance || 0) - pointsAmount;

    toast({
      title: "操作成功",
      description: `已为用户 ${user?.name} ${type === "add" ? "增加" : "扣减"} ${pointsAmount} 积分，当前余额：${newBalance}`,
    });

    onOpenChange(false);
    setAmount("");
    setReason("");
    setType("add");
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px] bg-card">
        <DialogHeader>
          <DialogTitle>调整用户积分</DialogTitle>
          <DialogDescription>
            为 <span className="font-semibold text-foreground">{user?.name}</span> ({user?.userId}) 调整积分
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="p-4 bg-muted rounded-lg">
            <p className="text-sm text-muted-foreground mb-1">当前余额</p>
            <p className="text-2xl font-bold text-primary">{user?.balance || 0}</p>
          </div>

          <div className="space-y-3">
            <Label>操作类型</Label>
            <RadioGroup value={type} onValueChange={(value) => setType(value as "add" | "subtract")} className="grid grid-cols-2 gap-4">
              <Label
                htmlFor="add"
                className={`flex items-center justify-center gap-2 p-4 border-2 rounded-lg cursor-pointer transition-colors ${
                  type === "add" ? "border-success bg-success/10" : "border-border hover:border-success/50"
                }`}
              >
                <RadioGroupItem value="add" id="add" className="sr-only" />
                <Plus className="w-5 h-5 text-success" />
                <span className="font-medium">增加积分</span>
              </Label>
              <Label
                htmlFor="subtract"
                className={`flex items-center justify-center gap-2 p-4 border-2 rounded-lg cursor-pointer transition-colors ${
                  type === "subtract" ? "border-destructive bg-destructive/10" : "border-border hover:border-destructive/50"
                }`}
              >
                <RadioGroupItem value="subtract" id="subtract" className="sr-only" />
                <Minus className="w-5 h-5 text-destructive" />
                <span className="font-medium">扣减积分</span>
              </Label>
            </RadioGroup>
          </div>

          <div className="space-y-2">
            <Label htmlFor="amount">
              积分数量 <span className="text-destructive">*</span>
            </Label>
            <Input
              id="amount"
              type="number"
              value={amount}
              onChange={(e) => setAmount(e.target.value)}
              placeholder="输入积分数量"
              min="1"
              required
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="reason">操作原因</Label>
            <Textarea
              id="reason"
              value={reason}
              onChange={(e) => setReason(e.target.value)}
              placeholder="记录调整积分的原因"
              rows={3}
            />
          </div>

          <DialogFooter>
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              取消
            </Button>
            <Button type="submit" className={type === "add" ? "bg-gradient-success" : "bg-destructive"}>
              确认{type === "add" ? "增加" : "扣减"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
