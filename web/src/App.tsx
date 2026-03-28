import { useState } from "react"
import { Toaster } from "sonner"
import { Sidebar, NavItem, Header, MainContent } from "@/components/layout"
import { Button } from "@/components/ui/button"
import { OnboardingDialog } from "@/features/onboard/OnboardingDialog"
import { NodesTable } from "@/features/nodes/NodesTable"
import { Zap } from "lucide-react"

export function App() {
  const [onboardingOpen, setOnboardingOpen] = useState(false)

  return (
    <div className="flex min-h-screen bg-background w-full">
      <Sidebar>
        <NavItem icon={<Zap className="h-5 w-5" />} label="Dashboard" active />
      </Sidebar>

      <MainContent>
        <Header
          title="Dashboard"
          subtitle="Monitor and manage your distributed system nodes"
          actions={
            <Button size="lg" onClick={() => setOnboardingOpen(true)}>
              + Add Node
            </Button>
          }
        />

        <main className="flex-1 px-8 py-6">
          <div className="space-y-6">
            <div>
              <h2 className="mb-4 text-lg font-semibold">Nodes</h2>
              <NodesTable />
            </div>
          </div>
        </main>
      </MainContent>

      <OnboardingDialog
        open={onboardingOpen}
        onOpenChange={setOnboardingOpen}
      />

      <Toaster position="bottom-right" />
    </div>
  )
}

export default App
