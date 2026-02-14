"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

const navItems = [
  { href: "/", label: "Dashboard" },
  { href: "/companies", label: "Companies" },
  { href: "/worksites", label: "Worksites" },
  { href: "/workers", label: "Workers" },
  { href: "/shifts", label: "Shifts" },
  { href: "/alarms", label: "Alarms" },
];

export function Navigation() {
  const pathname = usePathname();

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
          <div className="text-sm text-gray-400">
            {/* Auth status will go here */}
          </div>
        </div>
      </div>
    </nav>
  );
}
