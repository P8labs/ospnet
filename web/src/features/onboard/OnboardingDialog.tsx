import { useState } from "react"
import { useCreateToken } from "@/hooks"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Loader2, Copy, Check } from "lucide-react"
import { toast } from "sonner"

interface OnboardingDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function OnboardingDialog({
  open,
  onOpenChange,
}: OnboardingDialogProps) {
  const [copiedField, setCopiedField] = useState<string | null>(null)
  const createTokenMutation = useCreateToken()

  const handleGenerateToken = async () => {
    try {
      await createTokenMutation.mutateAsync()
      toast.success("Token generated successfully!")
    } catch (error) {
      toast.error(
        error instanceof Error ? error.message : "Failed to generate token"
      )
    }
  }

  const handleCopy = (text: string, field: string) => {
    navigator.clipboard.writeText(text)
    setCopiedField(field)
    toast.success("Copied to clipboard!")
    setTimeout(() => setCopiedField(null), 2000)
  }

  const token = createTokenMutation.data?.token
  const expiresAt = createTokenMutation.data?.expires_at

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl">
        <DialogHeader>
          <DialogTitle>Add New Node</DialogTitle>
          <DialogDescription>
            Generate an onboarding token and follow the installation
            instructions
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-6">
          {/* Token Generation Section */}
          <div className="space-y-3">
            <h3 className="text-sm font-semibold">Step 1: Generate Token</h3>
            <Button
              onClick={handleGenerateToken}
              disabled={createTokenMutation.isPending}
              className="w-full"
              size="lg"
            >
              {createTokenMutation.isPending && (
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              )}
              Generate Onboarding Token
            </Button>

            {createTokenMutation.isError && (
              <div className="rounded-md bg-red-50 p-3 text-sm text-red-800">
                Failed to generate token. Please try again.
              </div>
            )}

            {token && (
              <div className="space-y-3 rounded-lg bg-slate-50 p-4">
                <div>
                  <p className="mb-2 text-xs font-medium text-muted-foreground">
                    Token
                  </p>
                  <div className="flex items-center gap-2">
                    <code className="flex-1 overflow-hidden rounded bg-white px-3 py-2 font-mono text-sm text-ellipsis">
                      {token}
                    </code>
                    <Button
                      variant="outline"
                      size="icon"
                      onClick={() => handleCopy(token, "token")}
                    >
                      {copiedField === "token" ? (
                        <Check className="h-4 w-4" />
                      ) : (
                        <Copy className="h-4 w-4" />
                      )}
                    </Button>
                  </div>
                </div>

                {expiresAt && (
                  <div>
                    <p className="mb-2 text-xs font-medium text-muted-foreground">
                      Expires At
                    </p>
                    <p className="font-mono text-sm text-foreground">
                      {new Date(expiresAt).toLocaleString()}
                    </p>
                  </div>
                )}
              </div>
            )}
          </div>

          {/* Installation Instructions */}
          <div className="space-y-3">
            <h3 className="text-sm font-semibold">Step 2: Install on Node</h3>
            <p className="text-sm text-muted-foreground">
              Run this command on the node you want to add:
            </p>

            <div className="flex items-center gap-2">
              <code className="flex-1 rounded bg-slate-900 px-4 py-3 font-mono text-sm text-slate-50">
                curl -sSL https://ospnet.run/install.sh | bash
              </code>
              <Button
                variant="outline"
                size="icon"
                onClick={() =>
                  handleCopy(
                    "curl -sSL https://ospnet.run/install.sh | bash",
                    "command"
                  )
                }
              >
                {copiedField === "command" ? (
                  <Check className="h-4 w-4" />
                ) : (
                  <Copy className="h-4 w-4" />
                )}
              </Button>
            </div>

            {token && (
              <div className="rounded-lg bg-blue-50 p-4">
                <p className="mb-2 text-sm font-medium text-blue-900">
                  When prompted, enter this token:
                </p>
                <div className="flex items-center gap-2">
                  <code className="flex-1 rounded bg-white px-3 py-2 font-mono text-sm text-blue-900">
                    {token}
                  </code>
                  <Button
                    variant="outline"
                    size="icon"
                    onClick={() => handleCopy(token, "token-prompt")}
                  >
                    {copiedField === "token-prompt" ? (
                      <Check className="h-4 w-4" />
                    ) : (
                      <Copy className="h-4 w-4" />
                    )}
                  </Button>
                </div>
              </div>
            )}
          </div>

          {/* Footer */}
          <div className="flex justify-end gap-3 border-t border-border pt-4">
            <Button variant="outline" onClick={() => onOpenChange(false)}>
              Close
            </Button>
            {token && <Button onClick={() => onOpenChange(false)}>Done</Button>}
          </div>
        </div>
      </DialogContent>
    </Dialog>
  )
}
