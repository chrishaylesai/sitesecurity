"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useAuth } from "@/lib/auth";

const navItems = [
  { href: "/", label: "Dashboard" },
  { href: "/companies", label: "Companies" },
  { href: "/worksites", label: "Worksites" },
  { href: "/workers", label: "Workers" },
  { href: "/shifts", label: "Shifts" },
  { href: "/reports", label: "Reports" },
  { href: "/check-ins", label: "Check-ins" },
  { href: "/alarms", label: "Alarms" },
];

function roleLabel(roles: string[]): string | null {
  if (roles.includes("company_admin")) return "Admin";
  if (roles.includes("site_admin")) return "Manager";
  if (roles.includes("worker")) return "Worker";
  return null;
}

export function Navigation() {
  const pathname = usePathname();
  const { user, isAuthenticated, login, logout } = useAuth();

  return (
    <nav className="bg-gray-900">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="flex h-16 items-center justify-between">
          <div className="flex items-center">
            <Link href="/" className="text-xl font-bold text-white">
              SiteSecurity
            </Link>
            <div className="ml-10 flex items-baseline space-x-4">
              {navItems.map((item) => {
                const isActive =
                  item.href === "/"
                    ? pathname === "/"
                    : pathname.startsWith(item.href);
                return (
                  <Link
                    key={item.href}
                    href={item.href}
                    className={`rounded-md px-3 py-2 text-sm font-medium ${
                      isActive
                        ? "bg-gray-800 text-white"
                        : "text-gray-300 hover:bg-gray-700 hover:text-white"
                    }`}
                  >
                    {item.label}
                  </Link>
                );
              })}
            </div>
          </div>
          <div className="flex items-center gap-4">
            {isAuthenticated && user ? (
              <>
                <div className="flex items-center gap-2 text-sm text-gray-300">
                  <span className="font-medium">
                    {user.name || user.email}
                  </span>
                  {roleLabel(user.roles) && (
                    <span className="rounded-full bg-gray-700 px-2 py-0.5 text-xs text-gray-300">
                      {roleLabel(user.roles)}
                    </span>
                  )}
                </div>
                <button
                  onClick={() => logout()}
                  className="rounded-md bg-gray-700 px-3 py-1.5 text-sm font-medium text-white hover:bg-gray-600"
                >
                  Sign out
                </button>
              </>
            ) : (
              <button
                onClick={() => login()}
                className="rounded-md bg-blue-600 px-3 py-1.5 text-sm font-medium text-white hover:bg-blue-500"
              >
                Sign in
              </button>
            )}
          </div>
        </div>
      </div>
    </nav>
  );
}
