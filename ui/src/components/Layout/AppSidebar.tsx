import { NavLink } from "react-router-dom";
import {
  LayoutDashboard,
  Coins,
  Users,
  Trophy,
  Gift,
  ChevronRight,
} from "lucide-react";
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarHeader,
  useSidebar,
} from "@/components/ui/sidebar";

const menuItems = [
  {
    title: "控制台",
    url: "/",
    icon: LayoutDashboard,
  },
  {
    title: "积分类型",
    url: "/points-types",
    icon: Coins,
  },
];

export function AppSidebar() {
  const { open } = useSidebar();

  return (
    <Sidebar collapsible="icon" className="border-r border-sidebar-border">
      <SidebarHeader className="border-b border-sidebar-border">
        <div className="flex items-center gap-2 px-4 py-4">
          <div className="w-8 h-8 rounded-lg bg-gradient-primary flex items-center justify-center">
            <Coins className="w-5 h-5 text-primary-foreground" />
          </div>
          {open && (
            <div className="flex flex-col">
              <span className="font-semibold text-sidebar-foreground">积分平台</span>
              <span className="text-xs text-sidebar-foreground/60">运营管理后台</span>
            </div>
          )}
        </div>
      </SidebarHeader>

      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupLabel>主菜单</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {menuItems.map((item) => (
                <SidebarMenuItem key={item.title}>
                  <SidebarMenuButton asChild>
                    <NavLink
                      to={item.url}
                      end={item.url === "/"}
                      className={({ isActive }) =>
                        `flex items-center gap-3 px-3 py-2 rounded-lg transition-all ${
                          isActive
                            ? "bg-sidebar-accent text-sidebar-primary font-medium"
                            : "text-sidebar-foreground/70 hover:bg-sidebar-accent/50 hover:text-sidebar-foreground"
                        }`
                      }
                    >
                      <item.icon className="w-5 h-5" />
                      <span>{item.title}</span>
                      <ChevronRight className="w-4 h-4 ml-auto opacity-0 group-hover:opacity-100 transition-opacity" />
                    </NavLink>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
    </Sidebar>
  );
}
