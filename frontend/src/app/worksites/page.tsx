"use client";

import { useEffect, useState } from "react";
import { PageHeader, DataTable, Button } from "@/components/ui";
import { apiClient } from "@/lib/api/client";
import type { Worksite } from "@/lib/types";

export default function WorksitesPage() {
  const [worksites, setWorksites] = useState<Worksite[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    apiClient<Worksite[]>("/api/v1/worksites")
      .then(setWorksites)
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return <div className="text-sm text-gray-500">Loading worksites...</div>;
  }

  if (error) {
    return (
      <div className="rounded-md bg-red-50 p-4 text-sm text-red-700">
        Failed to load worksites: {error}
      </div>
    );
  }

  const columns = [
    { key: "name", header: "Name" },
    { key: "address", header: "Address" },
    {
      key: "latitude",
      header: "Coordinates",
      render: (w: Worksite) =>
        w.latitude && w.longitude
          ? `${w.latitude.toFixed(4)}, ${w.longitude.toFixed(4)}`
          : "—",
    },
  ];

  return (
    <div>
      <PageHeader
        title="Worksites"
        description="View and manage worksites across all companies."
        action={<Button>Add Worksite</Button>}
      />
      <DataTable
        columns={columns}
        data={worksites}
        emptyMessage="No worksites found."
      />
    </div>
  );
}
