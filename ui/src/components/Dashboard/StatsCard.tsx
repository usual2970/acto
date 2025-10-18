import { LucideIcon } from "lucide-react";
import { Card, CardContent } from "@/components/ui/card";

interface StatsCardProps {
  title: string;
  value: string | number;
  icon: LucideIcon;
  trend?: {
    value: string;
    isPositive: boolean;
  };
  gradient?: string;
}

export function StatsCard({ title, value, icon: Icon, trend, gradient = "gradient-primary" }: StatsCardProps) {
  return (
    <Card className="shadow-card hover:shadow-card-hover transition-shadow duration-200">
      <CardContent className="p-6">
        <div className="flex items-start justify-between">
          <div className="flex-1">
            <p className="text-sm font-medium text-muted-foreground mb-1">{title}</p>
            <h3 className="text-3xl font-bold text-foreground mb-2">{value}</h3>
            {trend && (
              <p className={`text-sm font-medium ${trend.isPositive ? 'text-success' : 'text-destructive'}`}>
                {trend.isPositive ? '↑' : '↓'} {trend.value}
              </p>
            )}
          </div>
          <div className={`w-12 h-12 rounded-xl bg-${gradient} flex items-center justify-center shadow-glow`}>
            <Icon className="w-6 h-6 text-white" />
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
