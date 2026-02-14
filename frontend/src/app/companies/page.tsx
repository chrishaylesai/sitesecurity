"use client";

import { useEffect, useState } from "react";
import { PageHeader, DataTable, Button } from "@/components/ui";
import { apiClient } from "@/lib/api/client";
import type { Company } from "@/lib/types";

export default function CompaniesPage() {
  const [companies, setCompanies] = useState<Company[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    apiClient<Company[]>("/api/v1/companies")
      .then(setCompanies)
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return <div className="text-sm text-gray-500">Loading companies...</div>;
  }

  if (error) {
    return (
      <div className="rounded-md bg-red-50 p-4 text-sm text-red-700">
        Failed to load companies: {error}
      </div>
    );
  }

  const columns = [
    { key: "name", header: "Name" },
    { key: "email", header: "Email" },
    { key: "phone", header: "Telephone" },
    {
      key: "address",
      header: "Address",
      render: (c: Company) => (
        <span className="max-w-xs truncate">{c.address ?? "—"}</span>
      ),
    },
  ];

  return (
    <div>
      <PageHeader
        title="Companies"
        description="Manage security companies registered on the platform."
        action={<Button>Add Company</Button>}
      />
      <DataTable
        columns={columns}
        data={companies}
        emptyMessage="No companies found. Add one to get started."
      />
    </div>
  );
}
