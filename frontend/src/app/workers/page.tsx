"use client";

import { useEffect, useState } from "react";
import { PageHeader, DataTable, Button } from "@/components/ui";
import { apiClient } from "@/lib/api/client";
import type { Worker } from "@/lib/types";

export default function WorkersPage() {
  const [workers, setWorkers] = useState<Worker[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    apiClient<Worker[]>("/api/v1/workers")
      .then(setWorkers)
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return <div className="text-sm text-gray-500">Loading workers...</div>;
  }

  if (error) {
    return (
      <div className="rounded-md bg-red-50 p-4 text-sm text-red-700">
        Failed to load workers: {error}
      </div>
    );
  }

  const columns = [
    {
      key: "name",
      header: "Name",
      render: (w: Worker) => `${w.firstName} ${w.lastName}`,
    },
    { key: "email", header: "Email" },
    {
      key: "phone",
      header: "Telephone",
      render: (w: Worker) => w.phone ?? "—",
    },
  ];

  return (
    <div>
      <PageHeader
        title="Workers"
        description="Manage security workers, their qualifications, and company memberships."
        action={<Button>Add Worker</Button>}
      />
      <DataTable
        columns={columns}
        data={workers}
        emptyMessage="No workers found."
      />
    </div>
  );
}
