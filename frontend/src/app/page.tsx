import Link from "next/link";

const quickLinks = [
  {
    title: "Companies",
    description: "Manage security companies and their details.",
    href: "/companies",
    colour: "bg-blue-50 border-blue-200 hover:bg-blue-100",
  },
  {
    title: "Worksites",
    description: "View and manage worksites across all companies.",
    href: "/worksites",
    colour: "bg-green-50 border-green-200 hover:bg-green-100",
  },
  {
    title: "Workers",
    description: "Manage workers, qualifications, and memberships.",
    href: "/workers",
    colour: "bg-purple-50 border-purple-200 hover:bg-purple-100",
  },
  {
    title: "Shifts",
    description: "Schedule, assign, and track security shifts.",
    href: "/shifts",
    colour: "bg-orange-50 border-orange-200 hover:bg-orange-100",
  },
  {
    title: "Alarms",
    description: "Monitor active alarms raised by workers.",
    href: "/alarms",
    colour: "bg-red-50 border-red-200 hover:bg-red-100",
  },
];

export default function Dashboard() {
  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900">Dashboard</h1>
      <p className="mt-1 text-sm text-gray-500">
        Welcome to SiteSecurity — your security workforce management platform.
      </p>

      <div className="mt-8 grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
        {quickLinks.map((link) => (
          <Link
            key={link.href}
            href={link.href}
            className={`rounded-lg border p-6 transition-colours ${link.colour}`}
          >
            <h2 className="text-lg font-semibold text-gray-900">
              {link.title}
            </h2>
            <p className="mt-2 text-sm text-gray-600">{link.description}</p>
          </Link>
        ))}
      </div>
    </div>
  );
}
