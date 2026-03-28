import { cn } from "@/lib/utils"

interface SidebarProps {
  children?: React.ReactNode
}

export function Sidebar({ children }: SidebarProps) {
  return (
    <aside className="fixed top-0 left-0 h-screen w-64 border-r border-border bg-background">
      <div className="flex h-full flex-col">
        <div className="border-b border-border px-6 py-6">
          <h1 className="text-lg font-bold tracking-tight">OSPNet</h1>
          <p className="mt-1 text-xs text-muted-foreground">
            Infrastructure Dashboard
          </p>
        </div>
        <nav className="flex-1 overflow-y-auto px-3 py-4">{children}</nav>
      </div>
    </aside>
  )
}

interface NavItemProps extends React.HTMLAttributes<HTMLDivElement> {
  active?: boolean
  icon?: React.ReactNode
  label: string
}

export function NavItem({
  active,
  icon,
  label,
  className,
  ...props
}: NavItemProps) {
  return (
    <div
      className={cn(
        "flex cursor-pointer items-center gap-3 rounded-md px-3 py-2 text-sm font-medium transition-colors",
        active
          ? "bg-primary text-primary-foreground"
          : "text-muted-foreground hover:bg-muted hover:text-foreground",
        className
      )}
      {...props}
    >
      {icon && <span className="h-5 w-5">{icon}</span>}
      <span>{label}</span>
    </div>
  )
}

interface HeaderProps {
  title: string
  subtitle?: string
  actions?: React.ReactNode
}

export function Header({ title, subtitle, actions }: HeaderProps) {
  return (
    <header className="border-b border-border bg-background px-8 py-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold tracking-tight">{title}</h2>
          {subtitle && (
            <p className="mt-1 text-sm text-muted-foreground">{subtitle}</p>
          )}
        </div>
        {actions && <div className="flex gap-3">{actions}</div>}
      </div>
    </header>
  )
}

interface MainContentProps {
  children: React.ReactNode
}

export function MainContent({ children }: MainContentProps) {
  return (
    <main className="ml-64 min-h-screen bg-background w-full">
      <div className="flex flex-col">{children}</div>
    </main>
  )
}
