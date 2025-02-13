import * as React from "react"
import * as AvatarPrimitive from "@radix-ui/react-avatar"

import { cn } from "@/lib/utils"
import { User } from "@/subdomains/users/models/user"
import { SessionUserData } from "@/subdomains/users/models/session"

const Avatar = React.forwardRef<
  React.ElementRef<typeof AvatarPrimitive.Root>,
  React.ComponentPropsWithoutRef<typeof AvatarPrimitive.Root>
>(({ className, ...props }, ref) => (
  <AvatarPrimitive.Root
    ref={ref}
    className={cn(
      "relative flex h-10 w-10 shrink-0 overflow-hidden rounded-full",
      className
    )}
    {...props}
  />
))
Avatar.displayName = AvatarPrimitive.Root.displayName

const AvatarImage = React.forwardRef<
  React.ElementRef<typeof AvatarPrimitive.Image>,
  React.ComponentPropsWithoutRef<typeof AvatarPrimitive.Image>
>(({ className, ...props }, ref) => (
  <AvatarPrimitive.Image
    ref={ref}
    className={cn("aspect-square h-full w-full", className)}
    {...props}
  />
))
AvatarImage.displayName = AvatarPrimitive.Image.displayName

const AvatarFallback = React.forwardRef<
  React.ElementRef<typeof AvatarPrimitive.Fallback>,
  React.ComponentPropsWithoutRef<typeof AvatarPrimitive.Fallback>
>(({ className, ...props }, ref) => (
  <AvatarPrimitive.Fallback
    ref={ref}
    className={cn(
      "flex h-full w-full items-center justify-center rounded-full bg-muted",
      className
    )}
    {...props}
  />
))
AvatarFallback.displayName = AvatarPrimitive.Fallback.displayName

export { Avatar, AvatarImage, AvatarFallback }

export const UserAvatar: React.FC<{ user: User }> = ({ user }) => {
  const initials = user.UserName.Value.slice(0, 2).toUpperCase()
  return <Avatar className="border flex items-center justify-center">
    <AvatarImage src={user.Image} />
    <AvatarFallback>{initials}</AvatarFallback>
  </Avatar>
}

export const SessionAvatar: React.FC<{ session: SessionUserData }> = ({ session: user }) => {
  const initials = user.UserName.slice(0, 2).toUpperCase()
  return <Avatar className="border flex items-center justify-center">
    <AvatarImage src={user.UserImage} />
    <AvatarFallback>{initials}</AvatarFallback>
  </Avatar>
}