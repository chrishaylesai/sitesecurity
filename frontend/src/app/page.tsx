"use client";

import Link from "next/link";
import { useAuth } from "@/lib/auth/context";

const adminLinks = [
  {
    title: "Companies",
    description: "Manage security companies and their details.",
    href: "/companies",
  },
  {
    title: "Worksites",
    description: "View and manage worksites across all companies.",
    href: "/worksites",
  },
  {
    title: "Workers",
    description: "Manage workers, qualifications, and memberships.",
    href: "/workers",
  },
  {
    title: "Reports",
    description: "Manage report templates and view submitted shift reports.",
    href: "/reports",
  },
];

const workerLinks = [
  {
    title: "Shifts",
    description: "View your upcoming shifts and browse available ones.",
    href: "/shifts",
  },
  {
    title: "Check-ins",
    description: "Record your GPS check-in when arriving on site.",
    href: "/check-ins",
  },
  {
    title: "Alarms",
    description: "Raise an alarm or view active alerts.",
    href: "/alarms",
  },
];

function LinkCard({
  title,
  description,
  href,
}: {
  title: string;
  description: string;
  href: string;
}) {
  return (
    <Link
      href={href}
      className="rounded-lg border border-gray-200 bg-white p-6 transition-colors hover:bg-gray-50"
    >
      <h3 className="text-lg font-semibold text-gray-900">{title}</h3>
      <p className="mt-2 text-sm text-gray-600">{description}</p>
    </Link>
  );
}

export default function Dashboard() {
  const { isAuthenticated, isLoading, login } = useAuth();

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-24 text-gray-400">
        Loading…
      </div>
    );
  }

  if (!isAuthenticated) {
    return (
      <div className="flex flex-col items-center justify-center py-24 text-center">
        <h1 className="text-2xl font-bold text-gray-900">
          Welcome to SiteSecurity
        </h1>
        <p className="mt-2 text-sm text-gray-500">
          Sign in to access the workforce management platform.
        </p>
        <button
          onClick={login}
          className="mt-6 rounded-md bg-blue-600 px-5 py-2.5 text-sm font-semibold text-white hover:bg-blue-700"
        >
          Sign in
        </button>
      </div>
    );
  }

  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900">Dashboard</h1>
      <p className="mt-1 text-sm text-gray-500">
        Welcome to SiteSecurity — your security workforce management platform.
      </p>

      <section className="mt-8">
        <h2 className="mb-4 text-sm font-semibold uppercase tracking-wide text-gray-500">
          Administration
        </h2>
        <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
          {adminLinks.map((link) => (
            <LinkCard key={link.href} {...link} />
          ))}
        </div>
      </section>

      <section className="mt-10">
        <h2 className="mb-4 text-sm font-semibold uppercase tracking-wide text-gray-500">
          Field Operations
        </h2>
        <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {workerLinks.map((link) => (
            <LinkCard key={link.href} {...link} />
          ))}
        </div>
      </section>
    </div>
  );
}
