import Link from "next/link";

const adminLinks = [
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
    title: "Reports",
    description: "Manage report templates and view submitted shift reports.",
    href: "/reports",
    colour: "bg-teal-50 border-teal-200 hover:bg-teal-100",
  },
];

const workerLinks = [
  {
    title: "Shifts",
    description: "View your upcoming shifts and browse available ones.",
    href: "/shifts",
    colour: "bg-orange-50 border-orange-200 hover:bg-orange-100",
  },
  {
    title: "Check-ins",
    description: "Record your GPS check-in when arriving on site.",
    href: "/check-ins",
    colour: "bg-cyan-50 border-cyan-200 hover:bg-cyan-100",
  },
  {
    title: "Alarms",
    description: "Raise an alarm or view active alerts.",
    href: "/alarms",
    colour: "bg-red-50 border-red-200 hover:bg-red-100",
  },
];

function LinkCard({
  title,
  description,
  href,
  colour,
}: {
  title: string;
  description: string;
  href: string;
  colour: string;
}) {
  return (
    <Link
      href={href}
      className={`rounded-lg border p-6 transition-colors ${colour}`}
    >
      <h3 className="text-lg font-semibold text-gray-900">{title}</h3>
      <p className="mt-2 text-sm text-gray-600">{description}</p>
    </Link>
  );
}

export default function Dashboard() {
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
