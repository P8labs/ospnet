import * as React from "react"
import { cn } from "@/lib/utils"

interface DialogContextType {
  open: boolean
  setOpen: (open: boolean) => void
}

const DialogContext = React.createContext<DialogContextType | undefined>(
  undefined
)

const useDialog = () => {
  const context = React.useContext(DialogContext)
  if (!context) {
    throw new Error("useDialog must be used within a Dialog provider")
  }
  return context
}

interface DialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  children: React.ReactNode
}

function Dialog({ open, onOpenChange, children }: DialogProps) {
  return (
    <DialogContext.Provider value={{ open, setOpen: onOpenChange }}>
      {children}
    </DialogContext.Provider>
  )
}

const DialogTrigger = React.forwardRef<
  HTMLButtonElement,
  React.ButtonHTMLAttributes<HTMLButtonElement> & { asChild?: boolean }
>(({ onClick, children, ...props }, ref) => {
  const { setOpen } = useDialog()
  return (
    <button
      ref={ref}
      onClick={(e) => {
        setOpen(true)
        onClick?.(e)
      }}
      {...props}
    >
      {children}
    </button>
  )
})
DialogTrigger.displayName = "DialogTrigger"

const DialogPortal = ({ children }: { children: React.ReactNode }) => {
  const { open } = useDialog()
  if (!open) return null
  return <>{children}</>
}

const DialogClose = React.forwardRef<
  HTMLButtonElement,
  React.ButtonHTMLAttributes<HTMLButtonElement>
>(({ onClick, children, ...props }, ref) => {
  const { setOpen } = useDialog()
  return (
    <button
      ref={ref}
      onClick={(e) => {
        setOpen(false)
        onClick?.(e)
      }}
      {...props}
    >
      {children}
    </button>
  )
})
DialogClose.displayName = "DialogClose"

const DialogOverlay = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, onClick, ...props }, ref) => {
  const { setOpen } = useDialog()
  return (
    <div
      ref={ref}
      className={cn("fixed inset-0 z-50 bg-black/80", className)}
      onClick={(e) => {
        setOpen(false)
        onClick?.(e)
      }}
      {...props}
    />
  )
})
DialogOverlay.displayName = "DialogOverlay"

const DialogContent = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, children, ...props }, ref) => {
  const { open } = useDialog()

  if (!open) return null

  return (
    <DialogPortal>
      <DialogOverlay />
      <div
        ref={ref}
        className={cn(
          "fixed top-[50%] left-[50%] z-50 w-full max-w-lg translate-x-[-50%] translate-y-[-50%] rounded-lg border border-border bg-background p-6 shadow-lg duration-200",
          open &&
            "animate-in fade-in-0 zoom-in-95 slide-in-from-left-1/2 slide-in-from-top-[48%]",
          className
        )}
        onClick={(e) => e.stopPropagation()}
        {...props}
      >
        {children}
      </div>
    </DialogPortal>
  )
})
DialogContent.displayName = "DialogContent"

const DialogHeader = ({
  className,
  ...props
}: React.HTMLAttributes<HTMLDivElement>) => (
  <div
    className={cn(
      "flex flex-col space-y-1.5 text-center sm:text-left",
      className
    )}
    {...props}
  />
)
DialogHeader.displayName = "DialogHeader"

const DialogFooter = ({
  className,
  ...props
}: React.HTMLAttributes<HTMLDivElement>) => (
  <div
    className={cn(
      "flex flex-col-reverse sm:flex-row sm:justify-end sm:space-x-2",
      className
    )}
    {...props}
  />
)
DialogFooter.displayName = "DialogFooter"

const DialogTitle = React.forwardRef<
  HTMLHeadingElement,
  React.HTMLAttributes<HTMLHeadingElement>
>(({ className, ...props }, ref) => (
  <h2
    ref={ref}
    className={cn(
      "text-lg leading-none font-semibold tracking-tight",
      className
    )}
    {...props}
  />
))
DialogTitle.displayName = "DialogTitle"

const DialogDescription = React.forwardRef<
  HTMLParagraphElement,
  React.HTMLAttributes<HTMLParagraphElement>
>(({ className, ...props }, ref) => (
  <p
    ref={ref}
    className={cn("text-sm text-muted-foreground", className)}
    {...props}
  />
))
DialogDescription.displayName = "DialogDescription"

export {
  Dialog,
  DialogPortal,
  DialogOverlay,
  DialogTrigger,
  DialogClose,
  DialogContent,
  DialogHeader,
  DialogFooter,
  DialogTitle,
  DialogDescription,
}
